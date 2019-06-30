FROM golang:latest as builder

# This is needed in order to avoid an installation outside of $GOPATH error
ENV GOBIN /go/bin

WORKDIR /go/src
ADD . /go/src

# No mkdir available in the scratch image, so we need to shoehorn the data file into place in the builder
RUN mkdir /data
COPY ./data/polygons.properties /data

RUN go get -v
RUN go build -ldflags "-linkmode external -extldflags -static" -a -o /go/bin/app

FROM scratch
COPY --from=builder /data /data
COPY --from=builder /go/bin/app /

EXPOSE 4041
ENTRYPOINT ["/app" ]
