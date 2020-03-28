# go-websocket-sample

## practice
https://github.com/oreilly-japan/go-programming-blueprints


## Build Setup

```bash
# build
cd chat
go build -o chat

# run
env GITHUB_ID=hogehoge GITHUB_SECRET=fugafua GITHUB_REDIRECT=http://localhost:8080/auth/callback/github  ./chat -addr=":8080"
```

GITHUB_ID, GITHUB_SECRET: for OAuth

see

https://developer.github.com/apps/building-oauth-apps/authorizing-oauth-apps/
