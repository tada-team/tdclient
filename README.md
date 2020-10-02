# tdclient
Simple client library.

## Feedback
Join team: https://web.tada.team/api/v3/join/nWqYT1DgHnsS08pck1l9eq4ELVgVCm7q6xTxtHEVnu

### Example
```go
package main

import (
	"fmt"
	"log"
	"sync"
	"github.com/tada-team/tdclient"
)

func main() {
	session, err := tdclient.NewSession("https://web.tada.team")
	if err != nil {
		panic(err)
	}

	session.SetToken("YOUR_TOKEN")
	session.SetVerbose(true)
	if err := session.Ping(); err != nil {
		panic(err)
	}

	_, err := session.AddContact("TEAM_UID", "+70000000000")
	if err != nil {
		panic(err)
	}
}

```

## Snippets

### Get server version

```bash
go run examples/serverinfo/main.go
go run examples/serverinfo/main.go -verbose
go run examples/serverinfo/main.go -verbose -server https://demo.tada.team
```

Any example have `-verbose` and `-server` options. Default server value is `https://web.tada.team`.

### Authorization token

Some examples need authorization token. You have two options:

#### Bot token
 
Create bot by typing `/newbot <BOTNAME>` command in @TadaBot direct chat. 
Allowed for team admin only.

#### Token from your personal account

For server with sms authorization:
```
go run examples/smsauth/main.go
```

For server with Active Directory authorization:
```
go run examples/passwordauth/main.go
```

### Contacts

```go run examples/contacts/main.go -team <team uid> -token <token>```

### Messaging

Using websockets:
```
go run examples/send-message/main.go -team <team uid> -token <token> -chat <chat jid> -message <message text>
```
or

Using http API only:
```
go run examples/send-message-http/main.go -team <team uid> -token <token> -chat <chat jid> -message <message text>
```

How to get team uid and chat jid:

```https://demo.tada.team/dbd248d7-25c2-4e8f-a23a-99baf63223e9/chats/g-dce6f5fd-b741-40a6-aa9c-c0e928d9dac5```
 
 * dbd248d7-25c2-4e8f-a23a-99baf63223e9 – team uid
 * g-dce6f5fd-b741-40a6-aa9c-c0e928d9dac5 – chat jid

### Simple bot

Echobot: make response to every direct message.

```
go run examples/echobot/main.go -team <team uid> -token <token>
```
