FROM alpine
COPY kraken-dca /usr/bin/kraken-dca
ENTRYPOINT ["/usr/bin/kraken-dca"]