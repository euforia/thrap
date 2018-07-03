FROM golang:1.10.3 AS build
WORKDIR /go/src/github.com/euforia/thrap
RUN go get github.com/golang/dep/cmd/dep
COPY  Gopkg.* ./
RUN dep ensure -v -vendor-only

ARG NOMAD_ADDR
ARG VAULT_ADDR
ENV NOMAD_ADDR=${NOMAD_ADDR}
ENV VAULT_ADDR=${VAULT_ADDR}

COPY  . .
RUN make test
RUN make dist

FROM alpine
VOLUME /secrets.hcl
WORKDIR /
COPY --from=build /go/src/github.com/euforia/thrap/dist/thrap-linux /usr/bin/thrap
CMD ["thrap"]