# builder image
FROM golang:latest as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN GO111MODULE=on go get github.com/go-telegram-bot-api/telegram-bot-api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o island_bot bot.go


# generate clean, final image for end users
FROM scratch
COPY --from=builder /build/island_bot .

# executable
ENTRYPOINT [ "./island_bot" ]