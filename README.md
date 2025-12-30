<div align="center" id="top">

# Discord-Lookup
</div>
<p align="center">
  <img alt="Top language" src="https://img.shields.io/github/languages/top/yourpov/Discord-Lookup?color=56BEB8">
  <img alt="Language count" src="https://img.shields.io/github/languages/count/yourpov/Discord-Lookup?color=56BEB8">
  <img alt="Repository size" src="https://img.shields.io/github/repo-size/yourpov/Discord-Lookup?color=56BEB8">
  <img alt="License" src="https://img.shields.io/github/license/yourpov/Discord-Lookup?color=56BEB8">
</p>

---

## About

**Discord-Lookup** is a Go API for fetching user data from a discord ID.  
It returns avatars, banners, badges, and account creation dates

## Tech Stack

- [Go](https://golang.org/)  
- [Discord API v10](https://discord.com/developers/docs/intro)

## Setup

```bash
git clone https://github.com/yourpov/Discord-Lookup
cd Discord-Lookup

go run main.go
```

> The server will start at <http://localhost:8080> by default

---

## API Endpoint

```html
http://localhost:<port>/lookup?id=<DiscordID>
```

## Examples

> API Response: (http://localhost:port/lookup?id=DiscordID)
<img alt="image" src="https://imgur.com/zAjCUQh.png" />


HTML: (http://localhost:port)
<img alt="image" src="https://i.imgur.com/AeGm4as.png" />
