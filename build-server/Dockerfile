# Use a lightweight Ubuntu image
FROM ubuntu:focal

# Update and install necessary tools
RUN apt-get update && apt-get install -y curl tar git

# Install Node.js
RUN curl -sL https://deb.nodesource.com/setup_20.x | bash -
RUN apt-get install -y nodejs

# Download and install Go
RUN curl -fsSL https://go.dev/dl/go1.21.1.linux-amd64.tar.gz -o go.tar.gz \
    && tar -C /usr/local -xzf go.tar.gz \
    && rm go.tar.gz

# Add Go to PATH
ENV PATH="/usr/local/go/bin:${PATH}"

# Upgrade packages
RUN apt-get upgrade -y

# Set the working directory
WORKDIR /home/app

# Copy application files
COPY . .

# Ensure the script is executable
RUN chmod +x script.sh

# Define the entry point
ENTRYPOINT ["/home/app/script.sh"]
