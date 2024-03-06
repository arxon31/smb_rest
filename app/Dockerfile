FROM golang:1.21-alpine as builder
LABEL authors="arxon31"

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN go build -o ./bin/server ./cmd/server/main.go


FROM alpine as runner

WORKDIR /usr/local/bin

COPY --from=builder /usr/local/src/bin/server /usr/local/bin
COPY --from=builder /usr/local/src/cmd/server/config/config.yaml ./config/

ENTRYPOINT ["./server"]