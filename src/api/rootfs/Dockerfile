FROM alpine:3.6

RUN apk -U add ca-certificates && \
    rm -Rf /var/cache/apk/*
EXPOSE 8081
ENV PORT 8081
ENV HOST 0.0.0.0
COPY usr/bin/monocular /usr/local/bin/monocular
CMD ["monocular"]
