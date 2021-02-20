# builder image
FROM golang:alpine as builder
RUN mkdir /build
ADD bot.go questions.yml /build/
WORKDIR /build
RUN apk add --update --no-cache ca-certificates git
ENV GO111MODULE=off
RUN go get github.com/go-telegram-bot-api/telegram-bot-api && go get io/ioutil && go get gopkg.in/yaml.v2
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o island_bot bot.go


# generate clean, final image for end users
FROM scratch
COPY --from=builder /build/island_bot .
COPY --from=builder /etc/ssl/cert.pem /etc/ssl/cert.pem

# executable
ENTRYPOINT [ "./island_bot" ]

