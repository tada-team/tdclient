# tdclient

Simple client library.

## Examples

### Server version

```
go run examples/serverinfo/main.go
go run examples/serverinfo/main.go -verbose
go run examples/serverinfo/main.go -verbose -server https://demo.tada.team
```

### Messaging

```
go run examples/message/main.go -team <team uid> -token <bot token> -chat <chat jid> -message <message text>
```

For example:

```https://demo.tada.team/dbd248d7-25c2-4e8f-a23a-99baf63223e9/chats/g-dce6f5fd-b741-40a6-aa9c-c0e928d9dac5```

 * https://demo.tada.team — server url
 * dbd248d7-25c2-4e8f-a23a-99baf63223e9 – team uid
 * g-dce6f5fd-b741-40a6-aa9c-c0e928d9dac5 – chat jid
 
How to get bot token: type "/newbot <NAME>" command in @TadaBot direct chat
