FROM golang:latest as builder

WORKDIR /go/src/github.com/storefinder/indexer

COPY . .

RUN go get -u golang.org/x/vgo

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }

RUN vgo test $(vgo list ./... | grep -v /vendor/ | grep -v /build/) -cover \
 && CGO_ENABLED=0 GOOS=linux vgo build -a -installsuffix cgo -o indexer 

 FROM alpine:latest

 RUN apk --no-cache add ca-certificates curl bash

 COPY --from=builder /go/src/github.com/storefinder/indexer .

 CMD ["./indexer"]
 