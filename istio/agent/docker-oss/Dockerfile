FROM nginx:1.13.7

RUN apt-get update && apt-get install -y curl vim-tiny net-tools procps libcurl3-gnutls

RUN rm /etc/nginx/nginx.conf \
	&& rm -rf /etc/nginx/conf.d

ADD build/modules /etc/nginx/modules/
ADD build/libs /usr/local/lib/
RUN ldconfig

# Debug options
ARG DISABLE_MIXER_REPORT=1
ARG DISABLE_MIXER_CHECK=1
ARG DISABLE_TRACING=0

# Debug env variables
ENV DISABLE_MIXER_REPORT $DISABLE_MIXER_REPORT
ENV DISABLE_MIXER_CHECK $DISABLE_MIXER_CHECK
ENV DISABLE_TRACING $DISABLE_TRACING

ADD build/agent /usr/local/bin

RUN mkdir -p /etc/istio/proxy \
	&& chmod -R g=u /etc/istio 

ENTRYPOINT ["/usr/local/bin/agent"]
