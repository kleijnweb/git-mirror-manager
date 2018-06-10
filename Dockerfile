FROM golang:1.10.1-alpine as golang
RUN apk add -U --no-cache ca-certificates git alpine-sdk
ADD . /go/src/github.com/kleijnweb/git-mirror-manager
WORKDIR /go/src/github.com/kleijnweb/git-mirror-manager
RUN make build

FROM alpine/git
WORKDIR /
COPY LICENSE /
COPY --from=golang /go/src/github.com/kleijnweb/git-mirror-manager/git-mirror-manager .

EXPOSE 8080

ENTRYPOINT ["./git-mirror-manager"]
