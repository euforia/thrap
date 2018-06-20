FROM alpine
#
# Dummy
#
RUN apk update && \
    apk add curl

CMD ["cmd", "arg"]

ENV key00=value key01="value two"
ENV key10="some thing" key11=7

EXPOSE 9090
EXPOSE 9091/tcp
EXPOSE 9092/udp

ENTRYPOINT ["foobar"]
