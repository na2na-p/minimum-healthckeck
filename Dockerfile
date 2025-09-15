FROM golang:1.25.1-alpine3.22 AS builder

WORKDIR /go/src/app

COPY . .

RUN CGO_ENABLED=0 go build -o app .

FROM gcr.io/distroless/static:nonroot

WORKDIR /go/bin

USER nonroot

COPY --chown=nonroot:nonroot --from=builder /go/src/app/app ./app

ENTRYPOINT ["./app"]
