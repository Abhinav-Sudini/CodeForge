FROM debian:trixie

RUN apt-get update; \
    apt-get install -y --no-install-recommends curl; \
    apt clean; \
    rm -rf /var/lib/apt/lists/* ; 

COPY master/bin /codeForge

CMD ["/codeForge/master.bin"]
