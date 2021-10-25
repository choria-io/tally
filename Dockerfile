FROM alpine:latest

WORKDIR /tally
COPY tally /tally/tally

RUN apk update && \
    apk add shadow && \
    mkdir -p /tally/config /tally/ssl && \
    groupadd --gid 2048 choria && \
    useradd -c "Choria Orchestrator - choria.io" -m --uid 2048 --gid 2048 choria && \
    chown -R choria:choria /tally/config /tally/ssl

USER choria

ENTRYPOINT ["/tally/tally"]

CMD ["--port", "8080", "--config", "/tally/config/tally.conf"]
