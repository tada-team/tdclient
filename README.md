# tdclient

Simple client library.

```go
import "github.com/tada-team/tdclient"

func main() {
    team := "uid"
    chat := "chat jid"
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
    
    ws.SendPlainMessage(chat, "hello, world") 
}
```
