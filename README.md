[![Publish Archivar Docker](https://github.com/rwese/archivar/actions/workflows/release-package.yml/badge.svg?branch=main&event=release)](https://github.com/rwese/archivar/actions/workflows/release-package.yml)

# Archivar

Is a small golang tool to archive/move resources from one location to another,
for example archiving mails sent to a specific mailadress on a webdav storage.
Previously I used IFTTT or YahooPipelines for these tasks but I wanted to keep
working with GO and tinker around.

**THIS - IS - (TO BE) - HEAVILY - OVERENGINEERED** but clean _enough_

## DIAGRAM

```
┌─────────────┐  ┌──────────────┐  ┌─────────────┐  ┌────────────┐
│ GATHERER    │  │ FILTERS      │  │ PROCESSOR   │  │ ARCHIVER   │
│ │           │  │ │            │  │ │           │  │ │          │
│ └►IMAP      ├─►│ │►Filename   ├─►│ └►SANATIZER ├─►│ │►WEBDAV   │
│             │  │ └►Filesize   │  │             │  │ └►GDRIVE   │
│             │  │              │  │             │  │            │
└─────────────┘  └───┬──────────┤  └─┬───────────┤  └────────────┘
                     │ FILTER   │    │ PROCESS   │
                     │ │        │    └───────────┘
                     │ │►ACCEPT │
                     │ │►REJECT │
                     │ └►MISS   │
                     │          │
                     └──────────┘
```

## TODO

- General
  - [ ] Solution to remember already archived things
  - [ ] How to handle Processors adding files
  - [ ] Tests now that things shape up
    - [ ] Archiver
    - [x] Filters
    - [ ] Gatherer
    - [x] Processors
  - [ ] Prometheus Instrumenting
  - [ ] More docs
  - [x] Use Factories to reduce "new" logic
  - [x] Middleware-like for Processors and Filters
  - [x] Use Github-Actions
  - [x] Use Github-Packages for DockerImages
  - [x] MultiStaged DockerFile
- Gatherers
  - [x] IMAP
  - [ ] POP3
  - [ ] Webdav
  - [ ] Dropbox
  - [ ] Google Drive
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
  - [ ] OCR
  - [ ] Anti Virus (rly?)
- Filters
  - [x] Filename
  - [x] Filesize
  - [ ] Image (Size, Dpi)
- [x] cleanup logging
  - [x] properly apply log levels to output
- [x] deamonize - let's call it "daemonized"
  - [x] graceful shutdown
- [x] global service structgen to hold logger and other global stuff
