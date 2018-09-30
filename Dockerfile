FROM golang:1.11-alpine3.8 AS builder

# ======================================================
# (1) Install dep
# ======================================================
RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

# ======================================================
# (2) Create workspace folder
# ======================================================
RUN mkdir -p /go/src/git.linecorp.com/LINE-Thailand/bp-message-server
WORKDIR /go/src/git.linecorp.com/LINE-Thailand/bp-message-server

# ================================================================
# (3) Copy Gopkg.toml Gopkg.lock and install library dependencies
# These layers will only be re-built when Gopkg files are updated
# ================================================================
COPY Gopkg.lock Gopkg.toml /go/src/git.linecorp.com/LINE-Thailand/bp-message-server/
RUN dep ensure -vendor-only

# ======================================================
# (4) Copy all project and build it
# This layer will be rebuilt when ever a file
# has changed in the project directory
# ======================================================
COPY . /go/src/git.linecorp.com/LINE-Thailand/bp-message-server
RUN go build -o /bin/bp-message-server

# ==================================================================
# (5) Build the final image based on executable from build image
# ==================================================================
FROM alpine3.8
COPY --from=builder /bin/bp-message-server /bin/bp-message-server
EXPOSE 8080
CMD ["/bin/bp-message-server"]
