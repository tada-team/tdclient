# tdclient

Simple client library.

### First example: get server version

```
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

### Messaging

```
go run examples/message/main.go -team <team uid> -token <bot token> -chat <chat jid> -message <message text>
```

How to get team uid and chat jid:

```https://demo.tada.team/dbd248d7-25c2-4e8f-a23a-99baf63223e9/chats/g-dce6f5fd-b741-40a6-aa9c-c0e928d9dac5```
 
 * dbd248d7-25c2-4e8f-a23a-99baf63223e9 – team uid
 * g-dce6f5fd-b741-40a6-aa9c-c0e928d9dac5 – chat jid

### Contacts

```go run examples/contacts/main.go -team <team uid> -token <bot token>```
