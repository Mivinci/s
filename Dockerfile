FROM alpine:3.12

RUN mkdir /app
ADD bin/gc-sso-linux /app
ADD html /app/html
ADD static /app/static
WORKDIR /app

ENV MODE="production"
ENV DB="/data/db"
ENV KEY="wfs1101"

ENTRYPOINT [ "./gc-sso-linux" ]
