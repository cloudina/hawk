FROM golang:1.22-alpine3.20 AS builder

RUN apk update upgrade;

ENV YARA 4.5.1

# Install Yara
RUN apk --update add --no-cache openssl file bison jansson ca-certificates zlib
RUN apk --update add --no-cache  \
  pkgconfig \
  openssl-dev \
  jansson-dev \
  build-base \
  libc-dev \
  file-dev \
  automake \
  autoconf \
  libtool \
  flex \
  git \
  gcc \
  libcrypto3 \
  libmagic-static \
  linux-headers \
  && echo "===> Install Yara from source..." 

RUN cd /tmp \
  && git clone --recursive --branch v${YARA} https://github.com/VirusTotal/yara.git \
  && cd /tmp/yara \
  && ./bootstrap.sh \
  && sync \
  && ./configure --enable-magic \
  --enable-crypto \
  --enable-cuckoo \
  && make \
  && make install \
  && rm -rf /tmp/* 

RUN mkdir -p /uploads
RUN mkdir -p /rules
RUN mkdir -p /go/src/app
WORKDIR /go/src/app

COPY *.go /go/src/app/

RUN go mod init hawk
RUN go get -d -v ./...
RUN go build -o /go/bin/hawk

RUN git clone https://github.com/Yara-Rules/rules.git /rules

FROM alpine:3.20

# Update
RUN apk update upgrade

# Install git
RUN apk --update add --no-cache  \
    git \
    libc6-compat \
    tzdata

# Set timezone to Europe/Zurich
RUN ln -s /usr/share/zoneinfo/Europe/London /etc/localtime

# Install ClamAV
RUN apk --update add --no-cache clamav clamav-libunrar \
    && mkdir /run/clamav \
    && chown clamav:clamav /run/clamav \
    && chown -R clamav:clamav /var/lib/clamav/

RUN apk --update add --no-cache  \
    jansson \
    libmagic

COPY config/clamd.conf /etc/clamav/clamd.conf
COPY config/freshclam.conf /etc/clamav/freshclam.conf

RUN freshclam  --no-dns

COPY entrypoint.sh /usr/bin
COPY --from=builder /go/bin/hawk /usr/bin/
COPY --from=builder /usr/local/lib/libyara* /usr/local/lib/

COPY --from=builder /rules /rules

RUN ldconfig /etc/ld.so.conf.d

ENV IPADDR 0.0.0.0
ENV PORT 9999
EXPOSE ${PORT}
ENV INDEXES -i /rules/malware_index.yar

ENV NO_OF_CHECKS_FOR_DB_UPDATE=24

RUN addgroup --gid 10001 --system hawkgroup \
 && adduser  --uid 10000 --system --ingroup hawkgroup --home /home/hawkuser hawkuser

RUN chmod -R +r /rules
RUN chown -R  hawkuser:hawkgroup /usr/bin/hawk

RUN chown -R  hawkuser:hawkgroup /usr/bin/freshclam
RUN chown -R  hawkuser:hawkgroup /usr/sbin/clamd
RUN chmod +x  /usr/bin/entrypoint.sh
RUN chown -R  hawkuser:hawkgroup /usr/bin/entrypoint.sh
RUN chown -R  hawkuser:hawkgroup /etc/clamav/clamd.conf
RUN chown -R  hawkuser:hawkgroup /var/log/clamav/
RUN chown -R  hawkuser:hawkgroup /run/clamav/
RUN chown -R  hawkuser:hawkgroup /var/lib/clamav/

RUN apk add --no-cache tini

ENTRYPOINT ["/sbin/tini", "--", "entrypoint.sh"]

RUN apk add --no-cache bind-tools

USER hawkuser
