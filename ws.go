package tdclient

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/tada-team/tdproto"
)

var (
	Timeout     = errors.New("Timeout")
	defaultSize = 20
)

func (s *Session) Ws(team string, onfail func(error)) (*WsSession, error) {
	if s.token == "" {
		return nil, errors.New("empty token")
	}

	u := s.server
	u.Path = "/messaging/" + team
	u.Scheme = strings.Replace(u.Scheme, "http", "ws", 1)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{
		"token": []string{s.token},
	})

	if err != nil {
		return nil, err
	}

	w := &WsSession{
		session:        s,
		team:           team,
		websocket:      conn,
		eventListeners: make([]eventListener, 0),
	}

	w.ctx, w.cancel = context.WithCancel(context.Background())

	go w.inboxLoop()

	return w, nil
}

type serverEvent struct {
	name string
	raw  []byte
}

type eventListener struct {
	eventChannel    chan serverEvent
	finishedChannel chan struct{}
}

type WsSession struct {
	session             *Session
	currentError        error
	eventListeners      []eventListener
	eventListenerMutext sync.Mutex
	team                string
	websocket           *websocket.Conn
	ctx                 context.Context
	cancel              context.CancelFunc
	sendMutex           sync.Mutex
}

func (w *WsSession) Ping() string {
	confirmId := tdproto.ConfirmId()
	w.SendRaw(tdproto.XClientPing(confirmId))
	return confirmId
}

func (w *WsSession) SendPlainMessage(to tdproto.JID, text string) string {
	uid := uuid.New().String()
	w.SendEvent(tdproto.NewClientMessageUpdated(tdproto.ClientMessageUpdatedParams{
		MessageId: uid,
		To:        to,
		Content: tdproto.MessageContent{
			Type: tdproto.MediatypePlain,
			Text: text,
		},
	}))
	return uid
}

func (w *WsSession) DeleteMessage(uid string) error {
	return w.SendEvent(tdproto.NewClientMessageDeleted(uid))
}

func (w *WsSession) WaitForMessage() (tdproto.Message, bool, error) {
	v := new(tdproto.ServerMessageUpdated)
	if err := w.WaitFor(v); err != nil {
		return tdproto.Message{}, false, err
	}
	return v.Params.Messages[0], v.Params.Delayed, nil
}

func (w *WsSession) WaitForConfirm() (string, error) {
	v := getServerConfirm()
	defer releaseServerConfirm(v)
	if err := w.WaitFor(v); err != nil {
		return "", err
	}
	return v.Params.ConfirmId, nil
}

func (w *WsSession) createListener(eventName string) (*eventListener, error) {
	listener := eventListener{
		eventChannel:    make(chan serverEvent),
		finishedChannel: make(chan struct{}, 1),
	}

	func() {
		w.eventListenerMutext.Lock()
		defer w.eventListenerMutext.Unlock()

		w.eventListeners = append(w.eventListeners, listener)
	}()

	return &listener, nil
}

func (w *WsSession) removeLisener(listenerData *eventListener) {
	listenerData.finishedChannel <- struct{}{}
}

func (w *WsSession) WaitFor(v tdproto.Event) error {
	name := v.GetName()

	listener, err := w.createListener(name)
	if err != nil {
		return err
	}
	defer w.removeLisener(listener)

	for {
		select {
		case ev := <-(*listener).eventChannel:
			tdclientGlgLogger.Debug("recieved event: ", string(ev.raw))
			switch ev.name {
			case name:
				if err := JSON.Unmarshal(ev.raw, &v); err != nil {
					return errors.Wrapf(err, "json fail on %v", string(ev.raw))
				}
				return nil
			}
		case <-time.After(httpClient.Timeout):
			return Timeout
		}
	}
}

func (w *WsSession) SendRaw(b []byte) error {
	w.sendMutex.Lock()
	defer w.sendMutex.Unlock()

	err := w.websocket.SetWriteDeadline(time.Now().Add(httpClient.Timeout))
	if err != nil {
		return err
	}

	tdclientGlgLogger.Debug("raw sent:", string(b))
	if err := w.websocket.WriteMessage(websocket.BinaryMessage, b); err != nil {
		tdclientGlgLogger.Warn(errors.Wrap(err, ""))
		return err
	}

	err = w.websocket.SetWriteDeadline(time.Time{})
	if err != nil {
		return err
	}

	return nil
}

func (w *WsSession) Close() error {
	closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "tdclient closing")

	if err := w.SendRaw(closeMessage); err != nil {
		return err
	}
	tdclientGlgLogger.Info("sent closing message")

	w.cancel()
	return w.websocket.Close()
}

func (w *WsSession) SendEvent(event tdproto.Event) error {
	b, err := JSON.Marshal(event)
	if err != nil {
		tdclientGlgLogger.Warn(errors.Wrap(err, ""))
		return err
	}
	tdclientGlgLogger.Info("sending event:", event)
	return w.SendRaw(b)
}

func (w *WsSession) inboxLoop() {
	for {

		_, data, err := w.websocket.ReadMessage()
		if err != nil {
			tdclientGlgLogger.Error("websocket reading error: ", err)
			w.currentError = err
			return
		}

		tdclientGlgLogger.Debugf("received websocket data %q", data)

		var receivedEvent map[string]interface{}
		err = json.Unmarshal(data, &receivedEvent)
		if err != nil {
			tdclientGlgLogger.Warn("failed to unmarshal json event: ", err)
			continue
		}

		// Try to get confirm_id and resend it back
		confirmIdInterface := receivedEvent["confirm_id"]
		confirmId, ok := confirmIdInterface.(string)
		if ok {
			if confirmId != "" {
				w.SendRaw(tdproto.XClientConfirm(confirmId))
			}
		}
		eventNameInterface := receivedEvent["event"]
		eventName, ok := eventNameInterface.(string)
		if !ok {
			tdclientGlgLogger.Warn("failed to get event name of event, got: ", eventNameInterface)
			continue
		}

		if eventName == "server.warning" {
			tdclientGlgLogger.Warnf("recieved server warning: %q", receivedEvent["params"])
		}

		ev := serverEvent{
			name: eventName,
			raw:  data,
		}
		func() {
			w.eventListenerMutext.Lock()
			defer w.eventListenerMutext.Unlock()
			futureListeners := make([]eventListener, 0)

			for _, listener := range w.eventListeners {
				select {
				case listener.eventChannel <- ev:
					{

					}
				case <-listener.finishedChannel:
					{
						continue
					}
				}

				futureListeners = append(futureListeners, listener)
			}

			w.eventListeners = futureListeners
		}()
		select {
		case <-w.ctx.Done():
			return
		default:
		}

	}
}

func (w *WsSession) SendCallOffer(jid tdproto.JID, sdp string) {
	callOffer := new(tdproto.ClientCallOffer)
	callOffer.Name = callOffer.GetName()
	callOffer.Params.Jid = jid
	callOffer.Params.Trickle = false
	callOffer.Params.Sdp = sdp
	w.SendEvent(callOffer)
}

func (w *WsSession) SendCallLeave(jid tdproto.JID) {
	callLeave := new(tdproto.ClientCallLeave)
	callLeave.Name = callLeave.GetName()
	callLeave.Params.Jid = jid
	callLeave.Params.Reason = ""
	w.SendEvent(callLeave)
}

func (w *WsSession) ForeachMessage(messageHandler func(chan tdproto.Message, chan error)) error {

	eventName := tdproto.ServerMessageUpdated{}.GetName()

	listener, err := w.createListener(eventName)
	if err != nil {
		return err
	}
	defer w.removeLisener(listener)

	messages := make(chan tdproto.Message)
	errorsChan := make(chan error, 1)

	go messageHandler(messages, errorsChan)

	for {
		event := new(tdproto.ServerMessageUpdated)
		select {
		case ev := <-listener.eventChannel:
			tdclientGlgLogger.Debug("recieved event: ", string(ev.raw))
			switch ev.name {
			case eventName:
				if err := JSON.Unmarshal(ev.raw, &event); err != nil {
					return errors.Wrapf(err, "json fail on %v", string(ev.raw))
				}
				select {
				case err := <-errorsChan:
					{
						return err
					}
				case messages <- event.Params.Messages[0]:
					{

					}
				}
			}
		case err := <-errorsChan:
			return err
		}
	}
}
