FROM scratch

ADD playground-protoactor.amd64 playground-protoactor

ENTRYPOINT ["./playground-protoactor"]
