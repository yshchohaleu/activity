FROM golang:1.23-bullseye

# Create a non-root user
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Install necessary packages first
RUN apt-get update && apt-get install -y \
    sudo \
    git \
    curl \
    vim \
    && rm -rf /var/lib/apt/lists/*

# Set up Go environment
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
ENV CGO_ENABLED=1
ENV GO111MODULE=on

# Create the user with a home directory
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

# Set up workspace and Go directories
RUN mkdir -p /app /go && \
    chown -R $USERNAME:$USERNAME /app /go

# Install Go tools as root before switching to vscode user
RUN go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install github.com/99designs/gqlgen@latest 

RUN chown -R $USERNAME:$USERNAME $GOPATH

WORKDIR /app

# Switch to non-root user
USER $USERNAME 