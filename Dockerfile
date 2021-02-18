# builder image
FROM golang:latest as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go get github.com/go-telegram-bot-api/telegram-bot-api
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o island_bot .


# generate clean, final image for end users
FROM scratch
COPY --from=builder /build/island_bot .

# executable
ENTRYPOINT [ "./island_bot" ]
