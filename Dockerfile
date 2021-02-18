# builder image
FROM golang:latest as builder
RUN mkdir /build
ADD *.go /build/
WORKDIR /build
RUN go get github.com/go-telegram-bot-api/telegram-bot-api
RUN GOOS=linux go build bot.go


# generate clean, final image for end users
FROM scratch
COPY --from=builder /build/bot .

# executable
ENTRYPOINT [ "./bot" ]
