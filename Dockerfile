FROM golang:1.9.2

ENV path /go/src/route_guide
WORKDIR ${path}
COPY . ${path}

RUN go build -i -v -o server/micro_framework server/server.go \
  && cp server/micro_framework /usr/bin/ 

ENTRYPOINT ["/usr/bin/micro_framework"]
CMD ["--help"]
