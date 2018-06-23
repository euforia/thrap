FROM golang:1.10.3

RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/euforia/thrap

COPY Gopkg.* ./

RUN dep ensure -v -vendor-only

COPY . .

RUN make test
RUN make dist

FROM alpine
VOLUME /secrets.hcl
WORKDIR /
COPY --from=build /go/src/github.com/euforia/thrap/dist/thrap-linux /usr/bin/thrap
CMD ["thrap"]
