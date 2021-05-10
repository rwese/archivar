# Archivar

Is a small golang tool to archive/move resources from one location to another.

For example archiving mails sent to a specific mailadress on a webdav storage.

## DIAGRAM

```
┌─────────────┐  ┌──────────────┐  ┌────────────┐  ┌────────────┐
│ GATHERER    │  │ FILTERS      │  │ PROCESSOR  │  │ ARCHIVER   │
│ │           │  │ │            │  │ │          │  │ │          │
│ └►IMAP      ├─►│ │►Filename   ├─►│ └►SANATIZE ├─►│ │►WEBDAV   │
│             │  │ └►Filesize   │  │            │  │ └►GDRIVE   │
│             │  │              │  │            │  │            │
└─────────────┘  └───┬──────────┤  └─┬──────────┤  └────────────┘
                     │ FILTER   │    │ PROCESS  │
                     │ │        │    └──────────┘
                     │ │►ACCEPT │
                     │ │►REJECT │
                     │ └►MISS   │
                     │          │
                     └──────────┘
```

## TODO

- Gatherers
  - [x] IMAP
  - [ ] POP3
  - Reddit
    - [ ] Saved Posts
    - [ ] Top/New/Hot of Subreddit
    - [ ] w/o Post Comments?
- Archivers
  - [x] Webdav
  - [ ] Dropbox
  - [x] Google Drive
- Processors
  - [x] Sanatizer (Filename)
- Filters
  - [x] Filename
  - [x] Filesize
  - [ ] Image (Size, Dpi)
- [x] cleanup logging
  - [ ] properly apply log levels to output
- [x] deamonize
  - [x] graceful shutdown
- [x] global service structgen to hold logger and other global stuff
