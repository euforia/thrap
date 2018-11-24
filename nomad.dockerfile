
FROM frolvlad/alpine-glibc

ARG CONSUL_ADDR
ARG NOMAD_VERSION=0.8.4
ARG BOOTSTRAP_EXPECT
ENV CONSUL_ADDR="${CONSUL_ADDR}"
ENV NOMAD_VERSION="${NOMAD_VERSION}"
ENV BOOTSTRAP_EXPECT="${BOOTSTRAP_EXPECT}"

RUN apk update && apk add curl
RUN cd /tmp/ && \
    curl -O https://releases.hashicorp.com/nomad/${NOMAD_VERSION}/nomad_${NOMAD_VERSION}_linux_amd64.zip && \
    unzip nomad_${NOMAD_VERSION}_linux_amd64.zip


FROM frolvlad/alpine-glibc

EXPOSE 4646
EXPOSE 4647
EXPOSE 4648

ENV CONSUL_ADDR=""
ENV BOOTSTRAP_EXPECT=""

VOLUME /opt/nomad

COPY --from=0 /tmp/nomad /bin/
COPY nomad-entrypoint.sh /entrypoint.sh
CMD ["/entrypoint.sh"]
