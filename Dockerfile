FROM golang:1.20.1-alpine3.17
RUN apk update
RUN apk add --no-cache git
RUN apk --update add ca-certificates

FROM scratch
WORKDIR /home/
COPY gothanks /usr/bin/
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["gothanks"]
