FROM nixos/nix
WORKDIR /tmp
COPY . .
RUN nix build --extra-experimental-features "nix-command flakes" --show-trace

FROM ubuntu:jammy
# explicitly don't install 3d packages - they're huge and we don't need them
RUN export DEBIAN_FRONTEND=noninteractive \
    && apt-get update \
    && apt-get install -y gpgv software-properties-common curl xvfb unzip git\
    && add-apt-repository --yes ppa:kicad/kicad-6.0-releases \
    && apt update && apt-get install -y --install-recommends kicad kicad-packages3d-
RUN curl -Lo /opt/ibom.zip "https://github.com/openscopeproject/InteractiveHtmlBom/archive/refs/heads/master.zip" \
    && unzip /opt/ibom.zip -d /opt
RUN curl -Lo /opt/kiauto.deb "https://github.com/INTI-CMNB/KiAuto/releases/download/v1.6.15/kiauto_1.6.15-1_all.deb" \
    && DEBIAN_FRONTEND=noninteractive apt install -y /opt/kiauto.deb
COPY --from=0 /tmp/result/bin/kicadsitegenerator /opt
COPY run.sh /opt
RUN mkdir /opt/project
WORKDIR /opt/project
ENTRYPOINT ["/opt/run.sh"]
