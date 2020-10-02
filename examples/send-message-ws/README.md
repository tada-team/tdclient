## Message sending

### Create websocket connection
Connect to `wss://web.tada.team/messaging/<team-uid>` with header:
```
Token: <bot-token-value>
```

### Send composing event (optional)

```json
{
  "event": "client.chat.composing",
  "params": {
    "jid": "<chat-jid>",
    "composing": true
  }
}
```

All `client.chat.composing` arguments: https://github.com/tada-team/tdproto/blob/master/client_chat_composing.go

### Send message itself

```json
{
  "event": "client.message.updated",
  "params": {
    "to": "<chat-jid>",
    "content": {
      "text": "hello world!",
      "type":"plain"
    }
  }
}
```

All `client.message.updated` arguments: https://github.com/tada-team/tdproto/blob/master/client_message_updated.go
