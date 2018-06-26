FROM golang:1.10.3 as build
WORKDIR /
RUN go get github.com/golang/dep/cmd/dep

RUN ping -c 1 consul.test > consul-ping.log
RUN ping -c 1 vault.test > vault-ping.log

FROM alpine
WORKDIR /
COPY --from=build /consul-ping.log /
COPY --from=build /vault-ping.log /
