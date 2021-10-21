package tdclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/tada-team/tdproto"
	"github.com/tada-team/timerpool"
	"github.com/valyala/fastjson"
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
		session:   s,
		team:      team,
		conn:      conn,
		inbox:     make(chan serverEvent, defaultSize),
		outBytes:  make(chan []byte, defaultSize),
		listeners: make(map[string]chan []byte),
		fail:      make(chan error),
	}

	w.ctx, w.cancel = context.WithCancel(context.Background())

	go func() {
		err := <-w.fail
		if err != nil {
			if onfail == nil {
				s.logger.Panic("ws client fail:", err)
			}
			onfail(err)
		}
	}()

	go w.outboxLoop()
	go w.inboxLoop()

	return w, nil
}

type serverEvent struct {
	name string
	raw  []byte
}

type WsSession struct {
	session   *Session
	team      string
	conn      *websocket.Conn
	inbox     chan serverEvent
	outBytes  chan []byte
	fail      chan error
	listeners map[string]chan []byte
	ctx       context.Context
	cancel    context.CancelFunc
}

func (w *WsSession) Ping() string {
	confirmId := tdproto.ConfirmId()
	w.SendRaw(tdproto.XClientPing(confirmId))
	return confirmId
}

func (w *WsSession) SendPlainMessage(to tdproto.JID, text string) string {
	uid := uuid.New().String()
	w.Send(tdproto.NewClientMessageUpdated(tdproto.ClientMessageUpdatedParams{
		MessageId: uid,
		To:        to,
		Content: tdproto.MessageContent{
			Type: tdproto.MediatypePlain,
			Text: text,
		},
	}))
	return uid
}

func (w *WsSession) DeleteMessage(uid string) string {
	return w.Send(tdproto.NewClientMessageDeleted(uid))
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

func (w *WsSession) ListenFor(v tdproto.Event) chan []byte {
	ch := make(chan []byte, defaultSize)
	w.listeners[v.GetName()] = ch
	return ch
}

func (w *WsSession) WaitFor(v tdproto.Event) error {
	name := v.GetName()

	timer := timerpool.Get(httpClient.Timeout)
	defer timerpool.Release(timer)

	for {
		select {
		case ev := <-w.inbox:
			w.session.logger.Println("got:", string(ev.raw))
			switch ev.name {
			case name:
				if err := JSON.Unmarshal(ev.raw, &v); err != nil {
					w.fail <- errors.Wrapf(err, "json fail on %v", string(ev.raw))
					return nil
				}
				return nil
			case "server.warning":
				t := new(tdproto.ServerWarning)
				if err := JSON.Unmarshal(ev.raw, &t); err != nil {
					w.fail <- errors.Wrapf(err, "json fail on %v", string(ev.raw))
					return nil
				}
				log.Println("tdclient: warn:", t.Params.Message)
			}
		case <-timer.C:
			w.fail <- Timeout
			return Timeout
		}
	}
}

func (w *WsSession) Send(event tdproto.Event) string {
	b, err := JSON.Marshal(event)
	if err != nil {
		w.fail <- errors.Wrap(err, "json marshal fail")
	}
	w.SendRaw(b)
	return event.GetConfirmId()
}

func (w *WsSession) SendRaw(b []byte) {
	w.outBytes <- b
}

func (w *WsSession) Close() error {
	w.cancel()
	return w.conn.Close()
}

func (w *WsSession) outboxLoop() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case b := <-w.outBytes:
			w.session.logger.Println("send:", string(b))
			if err := w.conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
				w.fail <- errors.Wrap(err, "ws client fail")
				return
			}
		}
	}
}

func (w *WsSession) inboxLoop() {
	var parser fastjson.Parser
	for {
		_, data, err := w.conn.ReadMessage()
		if err != nil {
			if w.ctx.Err() == nil {
				w.fail <- errors.Wrap(err, "conn read fail")
			}
			return
		}

		w.session.logger.Println("got:", string(data))
		v, err := parser.ParseBytes(data)
		if err != nil {
			w.fail <- errors.Wrapf(err, "invalid json: `%s`", string(data))
			return
		}

		confirmId := string(v.GetStringBytes("confirm_id"))
		if confirmId != "" {
			w.SendRaw(tdproto.XClientConfirm(confirmId))
		}

		ev := serverEvent{
			name: string(v.GetStringBytes("event")),
			raw:  data,
		}

		ch := w.listeners[ev.name]
		if ch != nil {
			select {
			case ch <- ev.raw:
			default:
				w.fail <- fmt.Errorf("listener %s chan is full", ev.name)
			}
			continue
		}

		select {
		case w.inbox <- ev:
		case <-w.ctx.Done():
			return
		default:
			w.fail <- errors.Wrapf(err, "full inbox")
		}
	}
}

func (w *WsSession) SendCallOffer(jid tdproto.JID, sdp string) {
	callOffer := new(tdproto.ClientCallOffer)
	callOffer.Name = callOffer.GetName()
	callOffer.Params.Jid = jid
	callOffer.Params.Trickle = false
	callOffer.Params.Sdp = sdp
	w.Send(callOffer)
}

func (w *WsSession) SendCallLeave(jid tdproto.JID) {
	callLeave := new(tdproto.ClientCallLeave)
	callLeave.Name = callLeave.GetName()
	callLeave.Params.Jid = jid
	callLeave.Params.Reason = ""
	w.Send(callLeave)
}
