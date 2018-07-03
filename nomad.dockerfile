
FROM frolvlad/alpine-glibc

VOLUME /opt/nomad
RUN apk update && apk add curl && cd /tmp/ && \
    curl -O https://releases.hashicorp.com/nomad/0.8.4/nomad_0.8.4_linux_amd64.zip && \
    unzip nomad_0.8.4_linux_amd64.zip && mv nomad /bin/

CMD ["nomad", "agent", "-server", "-data-dir=/opt/nomad", "-bootstrap-expect=1"]
