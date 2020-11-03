package tdclient

import (
	"os"
	"testing"
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

	c.SetToken(tokenResp.Token)

	anyTeam := tokenResp.Me.Teams[0]
	ws, err := c.Ws(anyTeam.Uid, func(err error) {
		t.Fatal(err)
	})
	if err != nil {
		t.Fatal(err)
	}

	confirmId := ws.Ping()
	if confirmId == "" {
		t.Error("invalid confirm id")
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(key + " variable not set")
	}
	return v
}
