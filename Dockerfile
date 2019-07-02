############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Maintainer
LABEL maintainer="Sobat Binangkit<sobat.binangkit@gmail.com>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /go/src/github.com/sobat-binangkit/echo-setup

# Get the source from github.com.
RUN git clone https://github.com/sobat-binangkit/echo-setup.git .

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -o echobase .

############################
# STEP 2 build a small image
############################

FROM scratch

# Environment Variables
ENV DATA_PATH=/app/data
ENV HTTP_PORT=8080
ENV HTTPS_PORT=8585 

WORKDIR /app

# Copy our static executable.
COPY --from=builder /go/src/github.com/sobat-binangkit/echo-setup/echobase .

# Run the hello binary.
ENTRYPOINT ["/app/echobase"]