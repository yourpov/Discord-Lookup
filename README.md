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

**Discord-Lookup** is a Go API for retrieving Discord user data by ID.  
It returns avatars, banners, badges, and account creation dates

## Tech Stack

- [Go](https://golang.org/)  
- [Discord API](https://discord.com/developers/docs/intro)  

## Setup

```bash
# Clone & enter project
git clone https://github.com/yourpov/Discord-Lookup
cd Discord-Lookup

# Run server
go run main.go
```

> The server will start at <http://localhost:8080> by default

---

## API Endpoint

```html
http://localhost:<port>/lookup?id=<DiscordID>
```


## Examples

<img width="939" height="437" alt="image" src="https://github.com/user-attachments/assets/7254cb49-cf9f-45ca-bc67-fc0624de02f7" />

<img width="1259" height="963" alt="image" src="https://github.com/user-attachments/assets/7fc4b143-8fa5-46a9-9a63-13f84e5ac650" />
