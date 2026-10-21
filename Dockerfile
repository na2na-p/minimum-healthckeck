FROM golang:1.27.3-alpine3.24@sha256:f2c872bcee949b7be9c490efd02b48a8a2fbf947cdf533f1c58525ec9d8170fd AS builder

WORKDIR /go/src/app

COPY . .

RUN CGO_ENABLED=0 go build -o app .

FROM gcr.io/distroless/static:nonroot@sha256:516757106de35205efb70b91ba03f6f1298424a0678c603255cf2852a7000087

WORKDIR /go/bin

USER nonroot

COPY --chown=nonroot:nonroot --from=builder /go/src/app/app ./app

ENTRYPOINT ["./app"]
