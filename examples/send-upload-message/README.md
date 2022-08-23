## Message sending without websockets

POST `https://web.tada.team/api/v4/teams/<team-uid>/chats/<chat-jid>/messages`
```http request
Content-Type: multipart/form-data;

form-data
file: sample.txt
```

Sample response:
```json
{
    "_time": "70ms",
    "ok": true,
    "result": {
        "content": {
            "text": "File: https://web.tada.team/u/1ed22cd5effd618486440242ac120006/sample.txt",
            "type": "file",
            "upload": "1ed22cd5-effd-6184-8644-0242ac120006",
            "mediaURL": "https://web.tada.team/u/1ed22cd5effd618486440242ac120006/sample.txt",
            "name": "sample.txt"
        },
        "push_text": "File",
        "from": "d-b09ad656-4cba-d004-957d-74c6fe7785a7",
        "to": "t-3d2c7835-2fb9-aed7-4d25-abbb056d12a3",
        "message_id": "be20fba3-478e-e81d-8eee-b40ce68a2a22",
        "created": "2020-09-29T15:32:13.939083Z",
        "gentime": 1601393533939084058,
        "chat_type": "task",
        "chat": "t-3d2c7835-4d25-2fb9-aed7-abbb056d12a3",
        "prev": "a43a869b-4a41-d734-87b5-709a8bf94f72",
        "is_last": true,
        "silently": true,
        "editable_until": "2020-09-30T15:32:13.939083Z",
        "links": [
          {
            "pattern": "/u/1ed22cd5effd618486440242ac120006/sample.txt",
            "url": "https://web.tada.team/u/1ed22cd5effd618486440242ac120006/sample.txt",
            "text": "sample.txt"
          }
        ],
        "markup": [
          {
              "op": 27,
              "cl": 73,
              "typ": "link",
              "url": "https://web.tada.team/u/1ed22cd5effd618486440242ac120006/sample.txt",
              "repl": "sample.txt"
          }
        ],
        "num": 7
    }
}
```

Message fields description: https://github.com/tada-team/tdproto/blob/master/message.go
