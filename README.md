# Link Shortner

A simple Webapp + API to short long links. Written in Go Lang and SQL.

## Features

- Web UI
- API for Devs
- Admin Panel
- Using SQL

## Installation

Install Link Shortnet

- Install [go 1.17](https://go.dev/doc/install) (latest)
- Clone repo
- `go run main.go`
- Edit "passwd.txt" with a strong password. It will admin login password

## Screenshots

![App Screenshot](https://telegra.ph/file/059c3e361f1c8fa0bc9fd.png)

![App Screenshot](https://telegra.ph/file/38aed6b5a23ce07fc97a0.png)


## API Reference

#### Short a Long Link

```http
  POST /api
```

| Parameter | Type     | Description                      |
| :-------- | :------- | :------------------------------- |
|   `key`   | `string` | **Required**. Unique key for URL |
|   `url`   | `string` | **Required**. Long URL           |

## License

[GPL v3](https://www.gnu.org/licenses/gpl-3.0.en.html)
