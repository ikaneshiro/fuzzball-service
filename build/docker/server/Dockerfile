ARG GOLANG_TAG=1.14-alpine

FROM golang:$GOLANG_TAG as build-stage
RUN apk add --no-cache git ca-certificates && update-ca-certificates
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go run ./mage-bootstrap.go install

FROM scratch
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-stage /go/bin/server /usr/local/bin/
USER 1000
EXPOSE 8080/tcp
ENTRYPOINT ["server"]
