# Build stage
FROM golang:1.24.4 as builder

WORKDIR /go/src/github.com/ccamel/playground-protoactor.go

COPY . .

RUN make build

# Run stage
FROM scratch

WORKDIR /root/

COPY --from=builder /go/src/github.com/ccamel/playground-protoactor.go/playground-protoactor .

ENTRYPOINT ["./playground-protoactor"]
