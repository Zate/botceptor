FROM golang:latest as builder
WORKDIR /go/src/github.com/zate/botceptor/
RUN go get -u github.com/golang/dep/cmd/dep
ADD main.go .
ADD followers_training.t .
ADD followers.mod .
RUN dep init && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o botceptor main.go

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/zate/botceptor/botceptor .
COPY .secrets.yaml .

CMD ["/botceptor"]