FROM debian:trixie

RUN apt-get update; \
    apt-get install -y --no-install-recommends g++ libc6-dev curl; \
    apt clean; \
    rm -rf /var/lib/apt/lists/* ; 

COPY worker/bin /codeForge

CMD ["/bin/bash"]
