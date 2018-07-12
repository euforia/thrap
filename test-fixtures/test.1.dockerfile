FROM golang:1.10.3
WORKDIR /
RUN go get github.com/golang/dep/cmd/dep

RUN ping -c 1 consul.test > consul-ping.log
RUN ping -c 1 vault.test > vault-ping.log

FROM alpine
WORKDIR /
COPY --from=0 /consul-ping.log /
COPY --from=0 /vault-ping.log /
