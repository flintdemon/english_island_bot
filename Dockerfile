# builder image
FROM golang:alpine as builder
RUN mkdir /build
COPY bot.go /build/
COPY questions.yml /build/
WORKDIR /build
RUN apk add --update --no-cache ca-certificates git
ENV GO111MODULE=off
RUN go get github.com/go-telegram-bot-api/telegram-bot-api && go get io/ioutil && go get gopkg.in/yaml.v2
RUN CGO_ENABLED=0 GOOS=linux go build -a -o island_bot bot.go


# generate clean, final image for end users
FROM scratch
ENV TELETOKEN=$ISLAND_BOT_TOKEN
COPY --from=builder /build/island_bot .
COPY --from=builder /build/questions.yml .
COPY --from=builder /etc/ssl/cert.pem /etc/ssl/cert.pem

# executable
ENTRYPOINT [ "./island_bot" ]
