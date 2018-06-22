
FROM golang:1.10.3 as builder
RUN go get -v github.com/euforia/pseudo
WORKDIR /go/src/github.com/euforia/pseudo/
RUN make test
RUN make dist

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/euforia/pseudo/dist/pseudo-linux /bin/pseudo
CMD ["pseudo"]
