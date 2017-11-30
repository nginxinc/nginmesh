FROM ubuntu:16.04


RUN apt-get update && apt-get install --no-install-recommends --no-install-suggests -y \
                         build-essential \
                         cmake \
                         git \
                         ca-certificates \
                         wget  \
                         libpcre3 \
                         libpcre3-dev  \
                         libz-dev \
                         libcurl4-gnutls-dev

COPY build_modules.sh /usr/local/bin/

ENTRYPOINT [ "/usr/local/bin/build_modules.sh" ]