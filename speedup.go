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
	const (
		begin = `{"event":"client.ping","confirm_id":"`
		end   = `"}`
	)
	var b bytes.Buffer
	b.Grow(len(begin) + len(confirmId) + len(end))
	b.WriteString(begin)
	b.WriteString(confirmId)
	b.WriteString(end)
	return b.Bytes()
}
