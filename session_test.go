package tdclient

import (
	"os"
	"testing"

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

	anyTeam := me.Teams[0]
	contacts, err := c.Contacts(anyTeam.Uid)
	if err != nil {
		t.Fatal(err)
	}

	var anyCoworker tdproto.Contact
	for _, contact := range contacts {
		if contact.CanSendMessage != nil && *contact.CanSendMessage {
			anyCoworker = contact
			break
		}
	}
	if anyCoworker.Jid.Empty() {
		t.Error("coworker not fouind in contacts")
	}

	ws, err := c.Ws(anyTeam.Uid, func(err error) {
		t.Fatal(err)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()

	t.Run("ping", func(t *testing.T) {
		confirmId := ws.Ping()
		ev := new(tdproto.ServerConfirm)
		if err := ws.WaitFor(ev); err != nil {
			t.Fatal(err)
		}
		if ev.ConfirmId != confirmId {
			t.Error("confirmId mismatched")
		}
	})

	t.Run("create message", func(t *testing.T) {
		messageUid := ws.SendPlainMessage(anyCoworker.Jid, kozma.Say())
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
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(key + " variable not set")
	}
	return v
}
