FROM golang:latest as builder

ENV GO111MODULE=on

WORKDIR /go/src
ADD . /go/src

COPY ./data/ /data/

RUN go get -v
RUN go build -ldflags "-linkmode external -extldflags -static" -a -o /go/bin/app

FROM scratch
COPY --from=builder /data /data
COPY --from=builder /go/bin/app /

EXPOSE 4041
ENTRYPOINT ["/app" ]
