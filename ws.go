package tdclient

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/tada-team/tdproto"
)

var Timeout = errors.New("Timeout")

func (s *Session) Ws(team string, onfail func(error)) (*WsSession, error) {
	if s.token == "" {
		return nil, errors.New("empty token")
	}

	u := s.server
	u.Path = fmt.Sprintf("/messaging/%s", team)
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
		inbox:     make(chan serverEvent, 100),
		outEvents: make(chan tdproto.Event, 100),
		outBytes:  make(chan []byte, 100),
		fail:      make(chan error),
	}

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
	closed    bool
	inbox     chan serverEvent
	outEvents chan tdproto.Event
	outBytes  chan []byte
	fail      chan error
}

func (w *WsSession) Ping() string {
	confirmId := tdproto.ConfirmId()
	w.SendRaw(XNewClientPing(confirmId))
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

func (w *WsSession) WaitFor(v tdproto.Event) error {
	name := v.GetName()
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
			case "server.panic":
				t := new(tdproto.ServerPanic)
				if err := JSON.Unmarshal(ev.raw, &t); err != nil {
					w.fail <- errors.Wrapf(err, "json fail on %v", string(ev.raw))
					return nil
				}
				w.fail <- fmt.Errorf("server panic: %s", t.Params.Code)
				return nil
			}
		case <-time.After(w.session.Timeout):
			return Timeout
		}
	}
}

func (w *WsSession) Send(event tdproto.Event) string {
	w.outEvents <- event
	return event.GetConfirmId()
}

func (w *WsSession) SendRaw(b []byte) {
	w.outBytes <- b
}

func (w *WsSession) Close() error {
	w.closed = true
	return w.conn.Close()
}

func (w *WsSession) outboxLoop() {
	for !w.closed {
		select {
		case b := <-w.outBytes:
			w.session.logger.Println("send:", string(b))
			if err := w.conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
				w.fail <- errors.Wrap(err, "ws client fail")
				return
			}
		case e := <-w.outEvents:
			b, err := JSON.Marshal(e)
			if err != nil {
				w.fail <- errors.Wrap(err, "json marshal fail")
				return
			}
			w.outBytes <- b
		}
	}
}

func (w WsSession) inboxLoop() {
	for !w.closed {
		_, data, err := w.conn.ReadMessage()
		if err != nil {
			w.fail <- errors.Wrap(err, "conn read fail")
			return
		}

		v := new(tdproto.BaseEvent)
		if err := JSON.Unmarshal(data, v); err != nil {
			w.fail <- errors.Wrap(err, "invalid json")
			return
		}

		eventName := v.Name
		confirmId := v.GetConfirmId()

		//v, err := fastjson.ParseBytes(data)
		//if err != nil {
		//	w.fail <- errors.Wrap(err, "invalid json")
		//	return
		//}
		//
		//eventName := v.GetStringBytes("event")
		//confirmId := v.GetStringBytes("confirm_id")

		if len(confirmId) > 0 {
			w.Send(tdproto.NewClientConfirm(string(confirmId)))
		}

		w.inbox <- serverEvent{
			name: string(eventName),
			raw:  data,
		}
	}
}
