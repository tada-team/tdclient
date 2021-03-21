[![Codecov coverage build](https://github.com/tada-team/tdclient/actions/workflows/codecov.yml/badge.svg)](https://github.com/tada-team/tdclient/actions/workflows/codecov.yml) [![codecov](https://codecov.io/gh/tada-team/tdclient/branch/master/graph/badge.svg)](https://codecov.io/gh/tada-team/tdclient)

# tdclient
Simple client library.

## Feedback
Join team: http://tada.team/apitalks

### Example
```go
package main

import (
	"fmt"
	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdproto/tdapi"
)

func main() {
    session, err := tdclient.NewSession("https://web.tada.team")
    if err != nil {
       panic(err)
    }

    // Create new bot (/newbot command) or use examples/smsauth to get own account token
    session.SetToken("YOUR_TOKEN")

    // How to see team_uid: https://web.tada.team/{team_uid}/chats/{chat_jid}
    teamUid := "YOUR_TEAM_UID" 
	
    // show all requests/responses
    session.SetVerbose(true)
 
    // Check connection
    if err := session.Ping(); err != nil {
        panic(err)
    }
    
    // Invite new member to your team
    phone := "+70001234567"
    contact, err := session.AddContact(teamUid, phone)
    if err != nil {
        panic(err)
    }
    fmt.Println("contact created:", contact.Jid)
    
    // Send hello to direct
    msg, err := session.SendPlaintextMessage(teamUid, contact.Jid, "Hi there!") 
    fmt.Println("message sent at:", msg.Created)
    
    // Create new task. All Fields: https://github.com/tada-team/tdproto/blob/master/tdapi/task.go
    taskChat, err := session.CreateTask(teamUid, tdapi.Task{
        Description: "do it, do it now",
        Assignee:    contact.Jid,
        Public:      true, // task visible for all team members
    }) 
    fmt.Println("task created:", taskChat.Jid)
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
```bash
go run examples/smsauth/main.go
```

For server with Active Directory authorization:
```bash
go run examples/passwordauth/main.go
```

### Contacts

```bash
go run examples/contacts/main.go -team <team uid> -token <token>
```

### Messaging

Using websockets:
```bash
go run examples/send-message-ws/main.go -team <team uid> -token <token> -chat <chat jid> -message <message text>
```
or

Using http API only:
```bash
go run examples/send-message/main.go -team <team uid> -token <token> -chat <chat jid> -message <message text>
```

How to get team uid and chat jid:

```https://demo.tada.team/dbd248d7-25c2-4e8f-a23a-99baf63223e9/chats/g-dce6f5fd-b741-40a6-aa9c-c0e928d9dac5```
 
 * dbd248d7-25c2-4e8f-a23a-99baf63223e9 – team uid
 * g-dce6f5fd-b741-40a6-aa9c-c0e928d9dac5 – chat jid

### Simple bot

Echobot: make response to every direct message.

```bash
go run examples/echobot/main.go -team <team uid> -token <token>
```

### Delete system alerts

If you not need cleans alerts now, set `-dryrun` options. Query ruturns you list message delete prepareing

```bash
go run examples/delete-messages/main.go -token <token> -team <team uid> -chat <chat jid> -date <2017-12-31 or 2017-01-01 00:10:00>
```
