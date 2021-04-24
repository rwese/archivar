# Archivar

Is a small golang tool to archive/move resources from one location to another.

For example archiving mails sent to a specific mailadress on a webdav storage.

## DIAGRAM

```
┌─────────────┐                           ┌──────────────┐
│  GATHERER   │                           │  ARCHIVER    │
│  │          │       ┌────────────┐      │  │           │
│  └►IMAP     │       │            │      │  └►WEBDAV    │
│             ├──────►│  ARCHIVAR  ├─────►│              │
│             │       │            │      │              │
└─────────────┘       └────────────┘      └──────────────┘
```

## TODO

- [x] cleanup logging
  - [ ] properly apply log levels to output
- [x] deamonize
  - [ ] graceful shutdown
- [x] global service structgen to hold logger and other global stuff
