FROM ubuntu:24.04

ENV TZ="America/Los_Angeles"

# Install basic development tools and iptables/ipset
RUN apt update && apt install -y \
  less git procps sudo fzf zsh man-db unzip \
  gnupg2 gh iptables ipset iproute2 dnsutils \
  aggregate jq nodejs npm wget ca-certificates \
  curl pkg-config build-essential libssl-dev \
  locales golang-go

RUN locale-gen

# Ensure default user has access to /usr/local/share
RUN mkdir -p /usr/local/share/npm-global && \
  chown -R ubuntu:ubuntu /usr/local/share

# Persist bash history.
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R ubuntu /commandhistory

# Set `DEVCONTAINER` environment variable to help with AI agent orientation
ENV DEVCONTAINER=true

# Create workspace and config directories and set permissions
RUN mkdir -p /workspace /home/ubuntu/.claude && \
  chown -R ubuntu:ubuntu /workspace /home/ubuntu/.claude

WORKDIR /workspace

# Install git-delta for better git diffs
RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://github.com/dandavison/delta/releases/download/0.18.2/git-delta_0.18.2_${ARCH}.deb" && \
  sudo dpkg -i "git-delta_0.18.2_${ARCH}.deb" && \
  rm "git-delta_0.18.2_${ARCH}.deb"

# Install fzf shell completion and key bindings
RUN curl https://raw.githubusercontent.com/junegunn/fzf/master/shell/completion.zsh -o /usr/share/doc/fzf/examples/completion.zsh
RUN curl https://raw.githubusercontent.com/junegunn/fzf/master/shell/key-bindings.zsh -o /usr/share/doc/fzf/examples/key-bindings.zsh

# Copy and set up firewall script
COPY assets/init-firewall.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/init-firewall.sh && \
  echo "ubuntu ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/ubuntu-firewall && \
  chmod 0440 /etc/sudoers.d/ubuntu-firewall

# Set up non-root user
USER ubuntu

# Install global packages
ENV NPM_CONFIG_PREFIX=/usr/local/share/npm-global
ENV PATH=$PATH:/usr/local/share/npm-global/bin

# Set the default shell to zsh rather than sh
ENV SHELL=/bin/zsh

# Default powerline10k theme
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.2.0/zsh-in-docker.sh)" -- \
  -p git \
  -p fzf \
  -a "source /usr/share/doc/fzf/examples/key-bindings.zsh" \
  -a "source /usr/share/doc/fzf/examples/completion.zsh" \
  -a "export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  -x

# Install Claude
RUN npm install -g @anthropic-ai/claude-code

# Configuration vars
ENV NODE_OPTIONS=--max-old-space-size=4096
ENV CLAUDE_CONFIG_DIR=/home/ubuntu/.claude
ENV POWERLEVEL9K_DISABLE_GITSTATUS=true
ENV CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1

# initialize Claude settings
COPY assets/claude-settings.default.json /home/ubuntu/.claude/settings.json

COPY assets/startup.sh /usr/local/bin
ENTRYPOINT [ "/usr/local/bin/startup.sh" ]
