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

## Example

See `etc/archivar.yaml.dist` for a full example.

Minimal example config

```yaml
# etc/archivar.yaml
Jobs:
  imap_to_webdav:
Gatherers:
  imap_mail_account:
    Type: imap
    Config:
      Server: server:993
      Username:
      Password:
      # DeleteDownloaded: False
Archivers:
  webdav_nextcloud:
    Type: webdav
    Config:
      Username:
      Password:
      Server: https://server/remote.php/dav/files/username/
      UploadDirectory: /upload/
```

```yaml
# docker-compose.yml
# with bind-mount
version: "2.3"

services:
  archivar:
    image: docker.pkg.github.com/rwese/archivar/archivar
    restart: unless-stopped
    volumes:
      - "./etc:/etc/go-archivar"
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

## Jobs

### Definition

Each job consists of `1:Gatherer -> [x:Filters -> x:Processors] -> 1:Archiver`

```yaml
Jobs:
  imap_to_webdav:
    Interval: 600
    Gatherer: imap_mail_account
    Archiver: webdav_nextcloud
    Filters:
      - pdf_only
      - filesize_filter
    Processors:
      - only_nice_chars_and_trim
```

## Gatherers

### IMAP

```yaml
Gatherers:
  <gatherer_name>:
    Type: imap
    Config:
      Server: server:993
      Username:
      Password:
      # DeleteDownloaded: False
      # AllowInsecureSSL: False
```

## Archivers

### Webdav

```yaml
Archivers:
  <archiver_name>:
    Type: webdav
    Config:
      Username:
      Password:
      Server: https://server/remote.php/dav/files/username/
      UploadDirectory: /upload/
```

## Filters

Filters are used to filter files provided by the gatherer.

Possible results:

- Allow
- Reject
- Miss

### Filesize

Verify if file is at least `MinSizeBytes` and/or is below `MaxSizeBytes`.

```yaml
Filters:
  min_1B_max_100MB:
    Type: filesize
        Config:
          MinSizeBytes: 100
          MaxSizeBytes: 100000000
```

### Filename

Test the Filename against defined regex's.

Tests are always Allow -> Reject -> Miss (Allow).

- If you wish to not allow missed regex's have a . regex.
- All regexes are partial by default if you wish full match use ^abcd$
- Case-insensitive is (?i)

```yaml
Filters:
  pdf_only:
    Type: filename
        Config:
          Allow:
            - (?i).pdf$
          Reject:
```

## Processors

### Sanatizer

Is used to perform manipulation on the gathered files.

```yaml
Processors:
  only_nice_chars_and_trim:
    Type: sanatizer
    Config:
      TrimWhitespaces: True
      CharacterBlacklistRegexs:
        - "[^[:word:]-_. ]"
```

#### Notes

- Golang Regexp uses [RE2, a regular expression library](https://github.com/google/re2/wiki/Syntax)
