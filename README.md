# tdclient

Simple client library.

```go
import (
    "time"
    "github.com/tada-team/tdproto"
    "github.com/tada-team/tdclient"
)

func main() {
    // your team uid
    team := "uid"
    
    // group / direct / task chat identifier
    chat := "chat jid"
    
    // bot token. Type "/newbot <NAME>" command in @TadaBot direct chat
    token := "secret"
    
    client, err := tdclient.NewSession("https://web.tada.team")
    if err != nil {
        panic(err)
    }
    
    client.SetToken(token)
    client.SetVerbose(true) // enable logging
   
    features := client.Features()
    log.Prinln("server version:", features.Build)
   
    ws, err := client.Ws(team, nil)
    if err != nil {
        panic(err)
    }
    
    // composing like human. Full events list at https://github.com/tada-team/tdproto
    ws.Send(tdproto.NewClientChatComposing(jid, true, nil))
    time.Sleep(3*time.Second)
    
    // shortcut for simple messaging
    ws.SendPlainMessage(chat, "hello, world") 
}
```
