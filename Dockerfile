FROM ubuntu:hirsute
RUN export DEBIAN_FRONTEND=noninteractive && apt-get update && apt-get install -y gpgv software-properties-common curl xvfb && add-apt-repository --yes ppa:kicad/kicad-dev-nightly && apt update && apt-get install -y --install-recommends kicad-nightly
RUN curl -Lo /opt/ibom.zip "https://github.com/openscopeproject/InteractiveHtmlBom/archive/refs/heads/master.zip" && unzip /opt/ibom.zip -d /opt
COPY run.sh /opt
ENTRYPOINT ["/opt/run.sh"]
