FROM alpine
COPY kraken-dca /
ENTRYPOINT ["/kraken-dca"]