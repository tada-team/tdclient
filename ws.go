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
		Session: s,
		team:    team,
		conn:    conn,
		inbox:   make(chan serverEvent, 100),
		outbox:  make(chan tdproto.Event, 100),
		fail:    make(chan error),
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
	*Session
	team   string
	conn   *websocket.Conn
	closed bool
	inbox  chan serverEvent
	outbox chan tdproto.Event
	fail   chan error
}

func (w *WsSession) Ping() string {
	return w.Send(tdproto.NewClientPing())
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
	if err := w.WaitFor("server.message.updated", &v); err != nil {
		return tdproto.Message{}, false, err
	}
	return v.Params.Messages[0], v.Params.Delayed, nil
}

func (w *WsSession) WaitForConfirm() (string, error) {
	v := new(tdproto.ServerConfirm)
	if err := w.WaitFor("server.confirm", v); err != nil {
		return "", err
	}
	return v.Params.ConfirmId, nil
}

func (w *WsSession) WaitFor(name string, v interface{}) error {
	for {
		select {
		case ev := <-w.inbox:
			w.logger.Println("got:", string(ev.raw))
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
		case <-time.After(w.Timeout):
			return Timeout
		}
	}
}

func (w *WsSession) Send(event tdproto.Event) string {
	w.outbox <- event
	return event.GetConfirmId()
}

func (w *WsSession) Close() error {
	w.closed = true
	return w.conn.Close()
}

func (w *WsSession) outboxLoop() {
	for !w.closed {
		event := <-w.outbox

		b, err := JSON.Marshal(event)
		if err != nil {
			w.fail <- errors.Wrap(err, "json marshal fail")
			return
		}

		w.logger.Println("send:", string(b))
		if err := w.conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
			w.fail <- errors.Wrap(err, "ws client fail")
			return
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
		if err := JSON.Unmarshal(data, &v); err != nil {
			w.fail <- errors.Wrap(err, "json fail")
			return
		}

		if v.ConfirmId != "" {
			w.Send(tdproto.NewClientConfirm(v.ConfirmId))
		}

		w.inbox <- serverEvent{
			name: v.Name,
			raw:  data,
		}
	}
}
