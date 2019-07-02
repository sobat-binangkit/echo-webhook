############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Maintainer
LABEL maintainer="Sobat Binangkit<sobat.binangkit@gmail.com>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/sobat-binangkit/echo-setup

# Get the source from github.com.
RUN git clone https://github.com/sobat-binangkit/echo-setup.git .

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN go build -o $GOPATH/bin/main

############################
# STEP 2 build a small image
############################

FROM scratch

ARG domain_name

# Environment Variables
ENV DOMAIN_NAME=domain_name
ENV DATA_PATH=/app/data
ENV HTTP_PORT=8080
ENV HTTPS_PORT=8585 

WORKDIR /app

# Copy our static executable.
COPY --from=builder $GOPATH/bin/main .

# Run the hello binary.
ENTRYPOINT ["/app/main"]