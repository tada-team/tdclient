package tdclient

import (
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
