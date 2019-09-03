FROM golang:1.12.6

ENV path /go/src/micro_framework
WORKDIR ${path}
COPY . ${path}

RUN go build -i -v -o server/micro_framework server/server.go \
  && cp server/micro_framework /usr/bin/ 

ENTRYPOINT ["/usr/bin/micro_framework"]
CMD ["--help"]
