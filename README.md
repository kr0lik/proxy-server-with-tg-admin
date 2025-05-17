# proxy-server-with-tg-admin

Socks5 Proxy server with remote control via telegram bot.


### Supported commands:

/top

/create {username} [password] [ttl]

/activate {username} [ttl]

/deactivate {username}

/rename {username} {username}

/password {username} [password]

/ttl {username} [ttl]

/stat {username}

/clear {username}

/delete {username}

/users


## Using

Build bin file

```
docker run -v .:/app -w /app --rm golang:1.24-alpine go build -o socks5-proxy-server-with-tg cmd/server/main.go
```

Run bin file

```
socks5-proxy-server-with-tg --port-socks5=1080 --sqlite-path=./.data --telegram-bot-token= --telegram-admin-id=
```

Where:
 - port-socks5 - Socks5 proxy server port
 - sqlite-path - Path where will be storing server data
 - telegram-bot-token - Your bot token
 - telegram-admin-id - Your user id

## Develop

```
docker run -v .:/app -w /app --rm golang:1.24-alpine go run cmd/server/main.go --env=dev --port-socks5=1080 --telegram-bot-token= --telegram-admin-id=
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v2.1.6-alpine golangci-lint run
```