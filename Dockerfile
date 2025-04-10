# builder image
FROM golang:alpine as builder
RUN mkdir /build
COPY bot.go /build/
COPY questions.yml /build/
COPY go.mod /build/
COPY go.sum /build/
WORKDIR /build
RUN apk add --update --no-cache ca-certificates git
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -o island_bot bot.go


# generate clean, final image for end users
FROM scratch
COPY --from=builder /build/island_bot .
COPY --from=builder /build/questions.yml .
COPY --from=builder /etc/ssl/cert.pem /etc/ssl/cert.pem

# executable
ENTRYPOINT [ "./island_bot" ]
