# builder image
FROM golang:latest as builder
RUN mkdir /build
ADD bot.go /build/
WORKDIR /build
ENV GO111MODULE=off
RUN go get github.com/go-telegram-bot-api/telegram-bot-api 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o island_bot bot.go


# generate clean, final image for end users
FROM scratch
COPY --from=builder /build/island_bot .

# executable
CMD apk add --update --no-cache ca-certificates git
ENTRYPOINT [ "./island_bot" ]

