FROM debian:trixie-slim

RUN apt-get update -y && \
    apt-get install wget p7zip -y && \
    apt-get clean

WORKDIR /app

ARG MONERO_WALLET_VERSION=0.18.4.5

RUN wget https://downloads.getmonero.org/cli/monero-linux-x64-v${MONERO_WALLET_VERSION}.tar.bz2 && \
    7z x monero-linux-x64-v${MONERO_WALLET_VERSION}.tar.bz2 && \
    7z x monero-linux-x64-v${MONERO_WALLET_VERSION}.tar && \
    rm -f monero-linux-x64-v${MONERO_WALLET_VERSION}.tar.bz2 && \
    rm -f monero-linux-x64-v${MONERO_WALLET_VERSION}.tar

EXPOSE 18080
EXPOSE 18081
EXPOSE 18083
EXPOSE 18089

COPY ./monerod.conf /etc/monerod.conf
WORKDIR /app/monero-x86_64-linux-gnu-v${MONERO_WALLET_VERSION}

CMD [ "./monerod", "--config-file", "/etc/monerod.conf" ]
