FROM golang:1.14.2-alpine3.11
RUN apk update
RUN apk add --no-cache git
RUN apk --update add ca-certificates
WORKDIR /home
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o gothanks .

FROM scratch
WORKDIR /home/
COPY --from=0 /home/gothanks /usr/bin/
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["gothanks"]
