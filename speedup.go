package tdclient

import (
	"bytes"
	"sync"

	"github.com/tada-team/tdproto"
)

var serverConfirmsPool = sync.Pool{New: func() interface{} {
	return new(tdproto.ServerConfirm)
}}

func getServerConfirm() *tdproto.ServerConfirm {
	return serverConfirmsPool.Get().(*tdproto.ServerConfirm)
}

func releaseServerConfirm(c *tdproto.ServerConfirm) {
	c.Params.ConfirmId = ""
	serverConfirmsPool.Put(c)
}

func xNewClientPing(confirmId string) []byte {
	return xConcat( `{"event":"client.ping","confirm_id":"`, confirmId,`"}`)
}

func xNewClientConfirm(confirmId string) []byte {
	return xConcat(`{"event":"client.confirm","params":{"confirm_id":"`, confirmId, `"}}`)
}

func xConcat(begin, mid, end string) []byte {
	var b bytes.Buffer
	b.Grow(len(begin) + len(mid) + len(end))
	b.WriteString(begin)
	b.WriteString(mid)
	b.WriteString(end)
	return b.Bytes()
}
