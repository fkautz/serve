FROM golang:1-alpine as build
WORKDIR /build
COPY . .
RUN mkdir /www
RUN CGO_ENABLED=0 go install -ldflags="-extldflags=-static"

FROM scratch
COPY --from=build /go/bin/serve /serve
COPY --from=build /www /www
VOLUME ["/www"]
EXPOSE 8080
CMD ["/serve", "-d", "/www"]
