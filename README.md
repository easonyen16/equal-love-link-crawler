# equal-love-link-crawler

A CLI tool to back up chat messages and media from the [=LOVE](https://equal-love.link.cosm.jp/) and [≠ME](https://not-equal-me.link.cosm.jp/) fan platforms.

## Features

- Supports both =LOVE and ≠ME platforms
- Login with your account
- List all talk rooms with subscription status
- Back up chat messages and media (images/videos) from subscribed rooms
- Incremental backup — stops when it encounters already-saved content
- Files are named and timestamped using Japan Standard Time (JST)

## Requirements

- Go 1.26

## Usage

```bash
go run ./cmd/backup
```

You will be prompted to select a group, then enter your email and password. The password input is hidden.

```
グループ選択:
  1) =LOVE
  2) ≠ME
グループを選んでください (1/2): 1
Email: your@email.com
Password:

=== 購読済み (2/10) ===
  佐々木 舞香    Maika Sasaki
  齋藤 樹愛羅    Kiara Saito

=== 未購読 (8/10) ===
  ...

佐々木 舞香 のバックアップ開始...
  佐々木 舞香: 100 件保存 (page 1)
  佐々木 舞香: 87 件保存 (page 2)
```

## Output Structure

```
download/
├── 佐々木 舞香/
│   ├── 20260101120000.txt        # text message
│   ├── 20260101120000.jpeg       # single media attachment
│   ├── 20260102093000.txt
│   ├── 20260102093000-1.jpeg     # multiple media attachments
│   └── 20260102093000-2.jpeg
└── 齋藤 樹愛羅/
    └── ...
```

- File names are formatted as `YYYYMMDDHHmmss` in JST
- File modification times are set to match the original message timestamp
- Only artist-posted messages are saved (user replies are skipped)
- Media types supported: `jpeg`, `png`, `mp4`

## Project Structure

```
.
├── api/message/       # API client (auth, talk rooms, chat)
├── internal/backup/   # Backup logic (download, save, pagination)
└── cmd/backup/        # CLI entry point
```

## Build

```bash
go build -o backup ./cmd/backup
./backup
```
