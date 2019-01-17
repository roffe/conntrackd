FROM golang:alpine3.8 as builder
RUN apk add --no-cache git
COPY . $GOPATH/src/github.com/shopgun/conntrackd/
WORKDIR $GOPATH/src/github.com/shopgun/conntrackd/
RUN go get -d -v
RUN CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -o /go/bin/conntrackd *.go

FROM alpine:3.8
RUN adduser -D conntrackd
RUN apk add --no-cache conntrack-tools
COPY --from=builder /go/bin/conntrackd /opt/conntrackd/conntrackd
COPY dummy.sh /opt/conntrackd/dummy.sh
RUN chmod +x /opt/conntrackd/dummy.sh
WORKDIR /opt/conntrackd
EXPOSE 2112
USER conntrackd
ENTRYPOINT ["/opt/conntrackd/conntrackd"]