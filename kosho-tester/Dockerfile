FROM ubuntu:24.04

ENV TZ="America/Los_Angeles"

# Install basic development tools and iptables/ipset
RUN apt update && apt install -y \
  less git procps sudo man-db unzip \
  iptables curl iproute2 dnsutils \
  jq nodejs npm ca-certificates locales

# Configure locale
RUN locale-gen

# Ensure default user has access to /usr/local/share
RUN mkdir -p /usr/local/share/npm-global && \
  chown -R ubuntu:ubuntu /usr/local/share

# Create workspace and config directories and set permissions
RUN mkdir -p /workspace /home/ubuntu/.claude && \
  chown -R ubuntu:ubuntu /workspace /home/ubuntu/.claude

WORKDIR /workspace

# Copy and set up firewall script
COPY assets/init-firewall.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/init-firewall.sh && \
  echo "ubuntu ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/ubuntu-firewall && \
  chmod 0440 /etc/sudoers.d/ubuntu-firewall

# Set up non-root user
USER ubuntu

# Install Claude
ENV NPM_CONFIG_PREFIX=/usr/local/share/npm-global
ENV PATH=$PATH:/usr/local/share/npm-global/bin
RUN npm install -g @anthropic-ai/claude-code

# Configuration vars
ENV NODE_OPTIONS=--max-old-space-size=4096
ENV CLAUDE_CONFIG_DIR=/home/ubuntu/.claude
ENV CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1

# Setup seed repo for testing Kosho
COPY assets/seed.sh /usr/local/bin/seed.sh
RUN /usr/local/bin/seed.sh

COPY assets/claude.sh /usr/local/bin
COPY assets/startup.sh /usr/local/bin
ENTRYPOINT [ "/usr/local/bin/startup.sh" ]
