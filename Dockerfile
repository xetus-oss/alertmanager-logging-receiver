FROM golang:1.14 AS builder
RUN go version

WORKDIR /go/src/app
ENV GOARCH amd64
ENV GOOS linux
ENV CGO_ENABLED 0

COPY . /go/src/app
RUN go build -o bin/receiver ./receiver


FROM scratch

WORKDIR /receiver
COPY --from=builder /go/src/app/bin/receiver .

EXPOSE 8080
ENTRYPOINT ["./receiver"]