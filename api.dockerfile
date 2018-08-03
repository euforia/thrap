FROM golang:1.10.3 AS build
WORKDIR /go/src/github.com/euforia/thrap
RUN go get github.com/golang/dep/cmd/dep
COPY  Gopkg.* ./
RUN dep ensure -v -vendor-only

ARG STACK_VERSION
ARG NOMAD_ADDR
ARG VAULT_ADDR

ENV STACK_VERSION=${STACK_VERSION}
ENV NOMAD_ADDR=${NOMAD_ADDR}
ENV VAULT_ADDR=${VAULT_ADDR}

COPY  . .
RUN make test
RUN make dist/thrap-linux

# Publishable artifact
FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
VOLUME /secrets.hcl
WORKDIR /
EXPOSE 10000
COPY --from=build /go/src/github.com/euforia/thrap/dist/thrap-linux /usr/bin/thrap
RUN thrap configure --no-prompt
CMD ["thrap", "agent"]