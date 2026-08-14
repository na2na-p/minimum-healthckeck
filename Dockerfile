FROM golang:1.26.5-alpine3.24 AS builder

WORKDIR /go/src/app

COPY . .

RUN CGO_ENABLED=0 go build -o app .

FROM gcr.io/distroless/static:nonroot

WORKDIR /go/bin

USER nonroot

COPY --chown=nonroot:nonroot --from=builder /go/src/app/app ./app

ENTRYPOINT ["./app"]
