# tdclient

Simple client library.

```go
import "github.com/tada-team/tdclient"

func main() {
    verbose := true
    token := "secret"
    team := "uid"
    chat := "chat jid"
    
    client, err := tdclient.NewSession("https://web.tada.team", verbose)
    if err != nil {
        panic(err)
    }
    
    client.SetToken("secret")
   
    ws, err := client.Ws(team, nil)
    if err != nil {
        panic(err)
    }
    
    ws.SendPlainMessage(chat, "hello, world") 
}
```
