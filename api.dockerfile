FROM golang:1.10.3 as build
WORKDIR /go/src/github.com/euforia/thrap
#  Move this out
COPY  . .
RUN make deps
RUN make test
RUN make dist

FROM alpine
VOLUME /secrets.hcl
WORKDIR /
COPY --from=build /go/src/github.com/euforia/thrap/dist/thrap-linux /usr/bin/thrap
CMD ["thrap"]
