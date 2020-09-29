## Getting user token (not recommended, use bot tokens)

Part 1. Request for SMS:

POST `https://web.tada.team/api/v4/auth/sms/send-code`
```json
{
    "phone": "+79091234567"
}
```

Sample response:

```json
{
    "_time": "95ms",
    "ok": true,
    "result": {
        "phone": "+79091234567",
        "code_valid_until": "2020-09-29T17:32:51.623856Z",
        "next_code_at": "2020-09-29T17:23:50.623856Z",
        "code_length": 4
    }
}
```

Part 2. Confirm code from SMS:

POST `https://web.tada.team/api/v4/auth/sms/get-token`
```json
 {
    "code": "4321",
    "phone": "+79091234567"
}
```

Sample response:
```json5
{
    "_time": "63ms",
    "ok": true,
    "result": {
        "token": "*************",
        "me": {
            // ...some data...
        }
    }
}
```
