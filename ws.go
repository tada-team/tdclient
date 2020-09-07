package tdclient

import (
	"encoding/json"
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

var WsTimeout = errors.New("Timeout")

func (s *Session) WsClient(team string, onfail func(error)) (*wsClient, error) {
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

	w := &wsClient{
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

type wsClient struct {
	*Session
	team   string
	conn   *websocket.Conn
	closed bool
	inbox  chan serverEvent
	outbox chan tdproto.Event
	fail   chan error
}

func (w *wsClient) Ping() string {
	return w.send(tdproto.NewClientPing())
}

func (w *wsClient) SendPlainMessage(to tdproto.JID, text string) string {
	uid := uuid.New().String()
	w.send(tdproto.NewClientMessageUpdated(tdproto.ClientMessageUpdatedParams{
		MessageId: uid,
		To:        to,
		Content: tdproto.MessageContent{
			Type: tdproto.MediatypePlain,
			Text: text,
		},
	}))
	return uid
}

func (w *wsClient) DeleteMessage(uid string) string {
	return w.send(tdproto.NewClientMessageDeleted(uid))
}

func (w *wsClient) WaitForMessage() (tdproto.Message, bool, error) {
	v := new(tdproto.ServerMessageUpdated)
	err := w.waitFor("server.message.updated", &v)
	if err != nil {
		return tdproto.Message{}, false, err
	}
	return v.Params.Messages[0], v.Params.Delayed, nil
}

func (w *wsClient) WaitForConfirm() (string, error) {
	v := new(tdproto.ServerConfirm)
	err := w.waitFor("server.confirm", v)
	if err != nil {
		return "", err
	}
	return v.Params.ConfirmId, nil
}

func (w *wsClient) waitFor(name string, v interface{}) error {
	for {
		select {
		case ev := <-w.inbox:
			w.logger.Println("got:", string(ev.raw))
			switch ev.name {
			case name:
				if err := json.Unmarshal(ev.raw, &v); err != nil {
					w.fail <- errors.Wrapf(err, "json fail on %v", string(ev.raw))
					return nil
				}
			case "server.warning":
				t := new(tdproto.ServerWarning)
				if err := json.Unmarshal(ev.raw, &t); err != nil {
					w.fail <- errors.Wrapf(err, "json fail on %v", string(ev.raw))
					return nil
				}
				log.Println("tdclient: warn:", t.Params.Message)
			case "server.panic":
				t := new(tdproto.ServerPanic)
				if err := json.Unmarshal(ev.raw, &t); err != nil {
					w.fail <- errors.Wrapf(err, "json fail on %v", string(ev.raw))
					return nil
				}
				w.fail <- fmt.Errorf("server panic: %s", t.Params.Code)
				return nil
			}
			return nil
		case <-time.After(w.Timeout):
			return WsTimeout
		}
	}
}

func (w *wsClient) send(e tdproto.Event) string {
	w.outbox <- e
	return e.GetConfirmId()
}

func (w *wsClient) outboxLoop() {
	for !w.closed {
		data := <-w.outbox

		b, err := json.Marshal(data)
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

func (w wsClient) inboxLoop() {
	for !w.closed {
		_, data, err := w.conn.ReadMessage()
		if err != nil {
			w.fail <- errors.Wrap(err, "conn read fail")
			return
		}

		v := new(tdproto.BaseEvent)
		if err := json.Unmarshal(data, &v); err != nil {
			w.fail <- errors.Wrap(err, "json fail")
			return
		}

		if v.ConfirmId != "" {
			w.send(tdproto.NewClientConfirm(v.ConfirmId))
		}

		w.inbox <- serverEvent{
			name: v.Name,
			raw:  data,
		}
	}
}
