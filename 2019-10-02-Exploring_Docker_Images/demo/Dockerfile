FROM alpine:3.10

RUN apk add --no-cache figlet

COPY --from=hairyhenderson/gomplate:slim /gomplate /bin/gomplate

ENV CITY=Ottawa

CMD gomplate -i 'Hello, {{ getenv `CITY` }}' | figlet
