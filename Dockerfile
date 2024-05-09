FROM golang:1.22 as builder
WORKDIR /app
COPY go.* .
COPY *.go .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build 


FROM scratch
COPY --from=builder /app/simple-app .
CMD [ "./simple-app" ]
