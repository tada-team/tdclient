## Task creation

POST `https://web.tada.team/api/v4/teams/<team-uid>/tasks`
```json
{
    "description": "task title\ntask details",
    "assignee": "d-e0c6ee1e-3db9-4e48-3842-7dcb98cdc5e8",
    "public": false
}
``` 

Sample response:
```json
 {
    "_time": "220ms",
    "ok": true,
    "result": {
        "jid": "t-04203bbd-ee5f-4d3c-b1b9-2126d88c10a5",
        "chat_type": "task",
        "base_gentime": 1603317463480886837,
        "gentime": 1603317463610019671,
        "created": "2020-10-21T21:57:43.453972Z",
        "display_name": "1931. task title",
        "icons": {
            "stub": "https://web.tada.team/u/579dda:02/256.png",
            "letters": "02",
            "color": "#579dda"
        },
        "counters_enabled": true,
        "can_call": true,
        "can_send_message": true,
        "notifications_enabled": true,
        "last_message": {
            "content": {
                "text": "Created task for @username: task title",
                "type": "change",
                "subtype": "newtask",
                "title": "Created task for @username: task title",
                "actor": "d-2788a2c0-be47-4b95-8a7d-889b599b7bfd"
            },
            "push_text": "Создана задача для @username: task title",
            "from": "d-2788a2c0-be47-4b95-8a7d-889b599b7bfd",
            "to": "t-04203bbd-ee5f-4d3c-b1b9-2126d88c10a5",
            "message_id": "ea7dda0e-bb05-4904-a3ff-12148ce197c4",
            "created": "2020-10-21T21:57:43.610018Z",
            "gentime": 1603317463610019671,
            "chat_type": "task",
            "chat": "t-04203bbd-ee5f-4d3c-b1b9-2126d88c10a5",
            "links": [
                {
                    "pattern": "@username",
                    "url": "tadateam://d-e0c6ee1e-3db9-4e48-abb3-7dcb98cdc5e8",
                    "text": "@username"
                }
            ],
            "is_first": true,
            "is_last": true,
            "silently": true,
            "editable_until": "2020-10-21T21:57:43.610018Z",
            "num": 0
        },
        "last_read_message_id": "ea7dda0e-bb05-4904-a3ff-12148ce197c4",
        "changeable_fields": [
            "assignee",
            "collapsed",
            "counters_enabled",
            "custom_color_index",
            "deadline",
            "description",
            "hidden",
            "items",
            "notifications_enabled",
            "observers",
            "pinned",
            "pinned_sort_ordering",
            "public",
            "section",
            "tags",
            "task_status"
        ],
        "num_members": 2,
        "can_delete": true,
        "description": "task title\ntask description",
        "feed": true,
        "num_items": 0,
        "num_checked_items": 0,
        "assignee": "d-e0c6ee1e-3db9-4e48-abb3-7dcb98cdc5e8",
        "num": 1931,
        "observers": [],
        "owner": "d-2788a2c0-be47-4b95-8a7d-889b599b7bfd",
        "task_status": "new",
        "title": "task title",
        "tabs": [
            "my",
            "out",
            "active"
        ],
        "can_delete_any_message": true,
        "can_set_important_any_message": true
    }
}
```
