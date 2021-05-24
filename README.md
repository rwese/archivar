[![Publish Archivar Docker](https://github.com/rwese/archivar/actions/workflows/release-package.yml/badge.svg)](https://github.com/rwese/archivar/actions/workflows/release-package.yml)

# Archivar

Is a small golang tool to archive/move resources from one location to another,
for example archiving mails sent to a specific mailadress on a webdav storage.
Previously I used IFTTT or YahooPipelines for these tasks but I wanted to keep
working with GO and tinker around.

**THIS - IS - (TO BE) - HEAVILY - OVERENGINEERED** but clean _enough_

## DIAGRAM

```
┌──────────────┐  ┌──────────────┐  ┌─────────────┐  ┌──────────────┐
│ GATHERER     │  │ FILTERS      │  │ PROCESSOR   │  │ ARCHIVER     │
│ │            │  │ │            │  │ │           │  │ │            │
│ │►IMAP       ├─►│ │►Filename   ├─►│ │►SANATIZER ├─►│ │►WEBDAV     │
│ │►WEBDAV     │  │ └►Filesize   │  │ └►ENCRYPTER │  │ │►GDRIVE     │
│ └►FILESYSTEM │  │              │  │             │  │ └►FILESYSTEM │
└──┬───────────┤  └───┬──────────┤  └─┬───────────┤  └──┬───────────┤
   │ DOWNLOAD  │      │ FILTER   │    │ PROCESS   │     │ UPLOAD    │
   │ DELETE    │      │ │►ACCEPT │    └───────────┘     └───────────┘
   └───────────┘      │ │►REJECT │
                      │ └►MISS   │
                      └──────────┘
```

## Example

See `etc/archivar.yaml.dist` for a full example.

Minimal example config

```yaml
# etc/archivar.yaml
Jobs:
  imap_to_webdav:
    Gatherer: imap_mail_account
    Archiver: webdav_nextcloud
Gatherers:
  imap_mail_account:
    Type: imap
    Config:
      Server: server:993
      Username:
      Password:
      # DeleteDownloaded: False
      # AllowInsecureSSL: False
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
    image: docker.pkg.github.com/rwese/archivar/archivar:latest
    ## semver images
    # image: docker.pkg.github.com/rwese/archivar/archivar:0
    # image: docker.pkg.github.com/rwese/archivar/archivar:0.2
    # image: docker.pkg.github.com/rwese/archivar/archivar:0.2.1
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
  - [x] added PublicKey for file encryption/decryption and key-generation
  - [x] Use Factories to reduce "new" logic
  - [x] Middleware-like for Processors and Filters
  - [x] Use Github-Actions
  - [x] Use Github-Packages for DockerImages
  - [x] MultiStaged DockerFile
  - [x] MUST find a better naming for ´/internal` packages, no named imports
- Gatherers
  - [ ] POP3
  - [ ] Dropbox
  - [ ] Google Drive
  - [x] FileSystem
  - [x] IMAP
  - [x] Webdav
  - Reddit
    - [ ] Saved Posts
    - [ ] Top/New/Hot of Subreddit
    - [ ] w/o Post Comments?
- Archivers
  - [ ] Dropbox
  - [x] FileSystem
  - [x] Google Drive
  - [x] Webdav
- Processors
  - [x] Sanatizer (Filename)
  - [x] Encryption
    - [ ] Passphrase Support for decryption
    - [ ] Encrypt Metadata
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

### FILESYSTEM

```yaml
Gatherers:
  <gatherer_name>:
    Type: filesystem
    Config:
      Directory: /home/user/input_directory
      # DeleteDownloaded: true
```

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

### WEBDAV

```yaml
Gatherers:
  <gatherer_name>:
    Type: webdav
    Config:
      Username:
      Password:
      Server: https://server/remote.php/dav/files/username/
      UploadDirectory: /input_directory/
```

## Archivers

### Filesystem

Store files directly on the filesystem, write access required, works with
nfs or similar.

```yaml
Archivers:
  <archiver_name>:
    Type: filesystem
    Config:
      Directory: /home/user/archivar/
```

### Google-Drive

The setup of this archiver is a little complicated as it requires a client
registration etc etc. I will extend the documentation when I get to it. (i try)

```yaml
Archivers:
  <archiver_name>:
  google_drive:
    Type: gdrive
    OAuthToken: >
      {
        "access_token":"<your token here>",
        "token_type":"Bearer",
        "refresh_token":"<your refresh token here>",
        "expiry":"some-date"
      }
    ClientSecrets: >
      {"installed":
        {
          "client_id":"some_client_id",
          "project_id":"archivar",
          "auth_uri":"https://accounts.google.com/o/oauth2/auth",
          "token_uri":"https://oauth2.googleapis.com/token",
          "auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
          "client_secret":"some_client_secret",
          "redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
        }
      }
    UploadDirectory: /archivar/
```

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

Is used to perform manipulation filenames on the gathered files.

```yaml
Processors:
  only_nice_chars_and_trim:
    Type: sanatizer
    Config:
      TrimWhitespaces: True
      CharacterBlacklistRegexs:
        - "[^[:word:]-_. ]"
```

### Encrypter

Encrypter will encrypt, duh, files with the PublicKey given.

The Process works using EncryptOAEP(RSA-OAEP) to encrypt the AES passphrase
which is prepended to the encrypted fileBody both are then base64 encoded.

For decryption you can use the cli command:

```bash
./archivar encrypter decrypt \
--privateKey myKey.sec \
--srcFile /archive/somefile.txt.encrypted \
--destFile /archive/somefile.txt
```

I will work on the cli commands when the need arises, for now this works.

Passphrase support is TODO.

To split the encrypted-key from the encrypted-body you can use split:

```bash
./archivar encrypter split --srcFile /archive/somefile.txt.encrypted

```

```yaml
Processors:
  basic_encrypter:
    Type: encrypter
    Config:
      AddExtension: .thisIsEncryptedForMe
      DontRename: false # Default false
      PublicKey: |
        -----BEGIN RSA PUBLIC KEY-----
        << Enter your public key to encrypt files for >>
        -----END RSA PUBLIC KEY-----
```

#### Notes

- Golang Regexp uses [RE2](https://github.com/google/re2/wiki/Syntax)
- encryption snippets from @stupidbodo
  https://gist.github.com/stupidbodo/601b68bfef3449d1b8d9*
