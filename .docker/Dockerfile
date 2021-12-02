FROM busybox:1.31.1-glibc

COPY skysight /
COPY buildInfo.json /
COPY i18n /
ENV GIN_MODE release
ENTRYPOINT ["/skysight"]