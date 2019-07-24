# Request Bodies for CloudWatchEvents

## Create / Update Schedule Event | Tower -> Builder

```json
{
  "operation": "UPDATE_SCHEDULE",
  "repository": "my-app",
  "branch": "feat/test",
  "shutdownSchedules": [
    {
      "cron": "(0 12 * * ? *)"
    }
  ],
  "startupSchedules": [
    {
      "cron": "(0 11 * * ? *)"
    }
  ]
}
```

## Delete Schedule Event | Tower -> Builder

```json
{
  "operation": "DELETE_SCHEDULE",
  "repository": "my-app",
  "branch": "feat/test"
}
```