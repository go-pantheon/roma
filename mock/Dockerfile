FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

ARG BIN

WORKDIR /app

COPY ${BIN} .

VOLUME /app/conf
VOLUME /app/gamedata

CMD ["./client", "-conf", "/app/conf", "-gamedata", "/app/gamedata/json/"]
