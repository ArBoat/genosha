FROM bash
WORKDIR $GOPATH/src/genosha

COPY ./genosha             $GOPATH/src/genosha/
COPY ./auth/               $GOPATH/src/genosha/auth/
COPY ./configure_docker/   $GOPATH/src/genosha/configure_docker/
RUN  mkdir -p              $GOPATH/src/genosha/log/
RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*
RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
    ca-certificates \
    && update-ca-certificates 2>/dev/null || true

EXPOSE 8080
CMD ["./genosha"]