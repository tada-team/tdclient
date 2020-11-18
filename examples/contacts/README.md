## Server information
GET `https://web.tada.team/api/v4/teams/<team-id>/contacts/`
```json5
{
    "_time": "9ms",
    "ok": true,
    "result": [
        {
            "jid": "d-6bd5b0f5-254f-4f57-98f3-767a0af4c612",
            "display_name": "+7 (000) 000-0010",
            "short_name": "+7 (000) 000-0010",
            "contact_email": "",
            "contact_phone": "+70000000010",
            "icons": {
                "stub": "https://web.tada.team/u/8e88eb:10/256.png",
                "letters": "10",
                "color": "#8e88eb"
            },
            "role": "",
            "status": "member",
            "last_activity": null,
            "sections": [],
            "can_send_message": true,
            "can_call": true,
            "can_create_task": true,
            "can_add_to_group": true,
            "changeable_fields": []
        },
        {
            "jid": "d-c2e2c240-2630-4056-9973-40d367dcaf1f",
            "display_name": "+7 (987) 000-0000",
            "short_name": "+7 (987) 000-0000",
            "contact_email": "",
            "contact_phone": "+79870000000",
            "icons": {
                "stub": "https://web.tada.team/u/578e3c:00/256.png",
                "letters": "00",
                "color": "#578e3c"
            },
            "role": "",
            "status": "member",
            "last_activity": null,
            "sections": [],
            "can_send_message": true,
            "can_call": true,
            "can_create_task": true,
            "can_add_to_group": true,
            "changeable_fields": []
        },
        {
            "jid": "d-b5704414-f78c-4a77-8e00-c1593d398bc0",
            "display_name": "test test",
            "short_name": "Вы",
            "contact_email": "",
            "contact_phone": "+79999999999",
            "icons": {
                "stub": "https://web.tada.team/u/e36659:23/256.png",
                "letters": "23",
                "color": "#e36659"
            },
            "role": "",
            "status": "member",
            "last_activity": "2020-11-18T13:39:14.976754Z",
            "sections": [],
            "cant_send_message_reason": "Этот чат только для чтения",
            "can_create_task": true,
            "can_add_to_group": true,
            "can_delete": true,
            "changeable_fields": [
                "alt_send",
                "always_send_pushes",
                "asterisk_mention",
                "contact_email",
                "contact_mshort_view",
                "contact_phone",
                "contact_short_view",
                "contact_show_archived",
                "debug_show_activity",
                "default_lang",
                "family_name",
                "given_name",
                "group_mshort_view",
                "group_notifications_enabled",
                "group_short_view",
                "icons",
                "mood",
                "munread_first",
                "phone",
                "quiet_time_finish",
                "quiet_time_start",
                "role",
                "task_mshort_view",
                "task_notifications_enabled",
                "task_short_view",
                "timezone",
                "unread_first"
            ],
            "family_name": "test",
            "given_name": "test",
            "default_lang": "ru",
            "debug_show_activity": false,
            "alt_send": false,
            "asterisk_mention": false,
            "always_send_pushes": false,
            "timezone": "Europe/Moscow",
            "quiet_time_start": "00:30",
            "quiet_time_finish": "09:00",
            "group_notifications_enabled": true,
            "task_notifications_enabled": true,
            "contact_short_view": false,
            "group_short_view": false,
            "task_short_view": false,
            "contact_mshort_view": false,
            "group_mshort_view": false,
            "task_mshort_view": false,
            "contact_show_archived": false,
            "unread_first": false,
            "munread_first": false,
            "can_manage_sections": true,
            "can_manage_tags": true,
            "can_manage_integrations": true,
            "can_create_group": true,
            "can_join_public_groups": true,
            "can_join_public_tasks": true
        },
        {
            "jid": "d-b21aab48-29a7-4482-bdc4-67fab2908f96",
            "display_name": "Мои заметки",
            "short_name": "Мои заметки",
            "contact_email": "",
            "contact_phone": "",
            "icons": {
                "sm": {
                    "url": "https://web.tada.team/static/tada-bots/notes256.png",
                    "width": 256,
                    "height": 256
                },
                "lg": {
                    "url": "https://web.tada.team/static/tada-bots/notes512.png",
                    "width": 512,
                    "height": 512
                }
            },
            "role": "",
            "status": "member",
            "last_activity": "2020-11-18T13:39:32.799678Z",
            "botname": "notebot",
            "sections": [],
            "can_send_message": true,
            "changeable_fields": []
        },
    ]
}
```
