package tdclient

import "bytes"

func XNewClientPing(confirmId string) []byte {
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
