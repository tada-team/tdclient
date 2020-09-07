package tdclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

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
		inbox:   make(chan event, 100),
		outbox:  make(chan event, 100),
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

type params map[string]interface{}

type event struct {
	Name      string `json:"event"`
	Params    params `json:"params"`
	ConfirmId string `json:"confirm_id,omitempty"`
	raw       []byte
}

type wsClient struct {
	*Session
	team   string
	conn   *websocket.Conn
	closed bool
	inbox  chan event
	outbox chan event
	fail   chan error
}

func (w *wsClient) SendPlainMessage(to, text string) string {
	type messageContent struct {
		Text string `json:"text"`
		Type string `json:"type"`
	}

	uid := uuid.New().String()
	w.send("client.message.updated", params{
		"message_id": uid,
		"to":         to,
		"content": messageContent{
			Type: "plain",
			Text: text,
		},
	})

	return uid
}

func (w *wsClient) DeleteMessage(uid string) {
	w.send("client.message.delete", params{
		"message_id": uid,
	})
}

func (w *wsClient) Ping() string {
	return w.send("client.ping", params{})
}

var wsTimeout = errors.New("Timeout")

type Message struct {
	PushText  string `json:"push_text,omitempty"`
	MessageId string `json:"message_id"`
}

type serverMessageUpdated struct {
	Name   string `json:"event"`
	Params struct {
		Messages []Message `json:"messages"`
		Delayed  bool      `json:"delayed"`
	} `json:"params"`
}

type serverConfirm struct {
	Name   string `json:"event"`
	Params struct {
		ConfirmId string `json:"confirm_id"`
	} `json:"params"`
}

func (w *wsClient) waitForMessage(timeout time.Duration) (Message, bool, error) {
	v := serverMessageUpdated{}
	err := w.waitFor("server.message.updated", timeout, &v)
	if err != nil {
		return Message{}, false, err
	}
	return v.Params.Messages[0], v.Params.Delayed, nil
}

func (w *wsClient) waitForConfirm(timeout time.Duration) (string, error) {
	v := serverConfirm{}
	err := w.waitFor("server.confirm", timeout, &v)
	if err != nil {
		return "", err
	}
	return v.Params.ConfirmId, nil
}

func (w *wsClient) waitFor(name string, timeout time.Duration, v interface{}) error {
	for {
		select {
		case event := <-w.inbox:
			w.Session.logger.Println("got:", string(event.raw))
			if event.Name == name {
				if err := json.Unmarshal(event.raw, &v); err != nil {
					w.fail <- errors.Wrap(err, "json fail")
					return nil
				}
				return nil
			}
		case <-time.After(timeout):
			return wsTimeout
		}
	}
}

func (w *wsClient) send(name string, params params) string {
	uid := uuid.New().String()
	w.outbox <- event{
		Name:      name,
		Params:    params,
		ConfirmId: uid,
	}
	return uid
}

func (w *wsClient) outboxLoop() {
	for !w.closed {
		data := <-w.outbox

		b, err := json.Marshal(data)
		if err != nil {
			w.fail <- errors.Wrap(err, "json marshal fail")
			return
		}

		w.Session.logger.Println("send:", string(b))
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

		v := event{}
		if err := json.Unmarshal(data, &v); err != nil {
			w.fail <- errors.Wrap(err, "json fail")
			return
		}

		if v.ConfirmId != "" {
			w.send("client.confirm", params{
				"confirm_id": v.ConfirmId,
			})
		}

		v.raw = data
		w.inbox <- v
	}
}
