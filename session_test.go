package tdclient

import (
	"os"
	"testing"
	"time"

	"github.com/tada-team/tdproto/tdapi"

	"github.com/tada-team/kozma"

	"github.com/tada-team/tdproto"
)

func TestSession(t *testing.T) {
	testServer := mustEnv("TEST_SERVER")
	testAccountPhone := mustEnv("TEST_ACCOUNT_PHONE")
	testAccountCode := mustEnv("TEST_ACCOUNT_CODE")

	c, err := NewSession(testServer)
	if err != nil {
		t.Fatal(err)
	}

	codeResp, err := c.AuthBySmsSendCode(testAccountPhone)
	if err != nil {
		t.Fatal(err)
	}

	if codeResp.CodeLength != len(testAccountCode) {
		t.Fatalf("invalid code length: %+v", codeResp)
	}

	tokenResp, err := c.AuthBySmsGetToken(testAccountPhone, testAccountCode)
	if err != nil {
		t.Fatal(err)
	}

	if len(tokenResp.Me.Teams) == 0 {
		t.Fatalf("invalid teams number: %d", len(tokenResp.Me.Teams))
	}

	me := tokenResp.Me
	c.SetToken(tokenResp.Token)

	team := me.Teams[0]
	contacts, err := c.Contacts(team.Uid)
	if err != nil {
		t.Fatal(err)
	}

	var coworker tdproto.Contact
	for _, contact := range contacts {
		if contact.CanSendMessage != nil && *contact.CanSendMessage {
			coworker = contact
			break
		}
	}
	if coworker.Jid.Empty() {
		t.Error("coworker not fouind in contacts")
	}

	ws, err := c.Ws(team.Uid, func(err error) {
		t.Fatal(err)
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("ping", func(t *testing.T) {
		confirmId := ws.Ping()
		ev := new(tdproto.ServerConfirm)
		if err := ws.WaitFor(ev); err != nil {
			t.Fatal(err)
		}
		if ev.Params.ConfirmId != confirmId {
			t.Error("confirmId mismatched: got:", ev.ConfirmId, "want:", confirmId)
		}
	})

	t.Run("create message", func(t *testing.T) {
		messageUid := ws.SendPlainMessage(coworker.Jid, kozma.Say())
		msg, _, err := ws.WaitForMessage()
		if err != nil {
			t.Fatal(err)
		}
		if msg.MessageId != messageUid {
			t.Fatal("invalid message uid")
		}

		t.Run("delete message", func(t *testing.T) {
			ws.DeleteMessage(messageUid)
			msg, _, err := ws.WaitForMessage()
			if err != nil {
				t.Fatal(err)
			}
			if msg.MessageId != messageUid {
				t.Fatal("invalid message uid")
			}
		})
	})

	t.Run("create task", func(t *testing.T) {
		text := kozma.Say()
		chat, err := c.CreateTask(team.Uid, tdapi.Task{
			Description: text,
			Tags:        []string{"autotest"},
			Assignee:    coworker.Jid,
			Deadline:    tdproto.IsoDatetime(time.Now().Add(time.Hour)),
			Public:      false,
			RemindAt:    tdproto.IsoDatetime(time.Now().Add(time.Minute)),
		})
		if err != nil {
			t.Error(err)
		}
		if chat.Description != text {
			t.Error("task description mismatched: want:", text, "got:", chat.Description)
		}
	})
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(key + " variable not set")
	}
	return v
}
