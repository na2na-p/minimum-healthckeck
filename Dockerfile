FROM golang:1.26.6-alpine3.24@sha256:af8d6740070b8906d12eae1c3e3ea0957fb63f492051ea05e354c38ef9fe88df AS builder

WORKDIR /go/src/app

COPY . .

RUN CGO_ENABLED=0 go build -o app .

FROM gcr.io/distroless/static:nonroot@sha256:f7f8f729987ad0fdf6b05eeeae94b26e6a0f613bdf46feea7fc40f7bd72953e6

WORKDIR /go/bin

USER nonroot

COPY --chown=nonroot:nonroot --from=builder /go/src/app/app ./app

ENTRYPOINT ["./app"]
