FROM golang:1.11-alpine3.8 AS builder

WORKDIR /go/src/github.com/woraphol-j/url-shortener

# ======================================================
# (1) Install dep dependency tool
# - debug tool(not in use now)
# ======================================================
RUN apk add --no-cache git
RUN go get github.com/derekparker/delve/cmd/dlv

CMD ["go", "run", "main.go"]
