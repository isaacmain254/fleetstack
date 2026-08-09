FROM ubuntu:22.04

RUN apt-get update && apt-get install -y curl ca-certificates gnupg lsb-release
# Install Docker CLI
RUN curl -fsSL https://get.docker.com | sh

# Install Nixpacks
RUN curl -sSL https://nixpacks.com/install.sh | bash

# Install Railpack
RUN curl -sSL https://railpack.com/install.sh | bash

# Install Buildpacks
COPY --from=buildpacksio/pack:0.39.1 /usr/local/bin/pack /usr/local/bin/pack