FROM alpine:latest

WORKDIR /tally
COPY tally /tally/tally

RUN mkdir /tally/config

ENTRYPOINT ["/tally/tally"]
CMD ["--port", "8080", "--config", "/tally/config/tally.conf"]
