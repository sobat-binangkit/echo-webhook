############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Maintainer Name
LABEL maintainer="Ahmad R. Djarkasih<djarkasih@gmail.com>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Build Time Variables
ARG projname=myproj
ARG currdir=/go/src/${projname}
ARG gitaddr=https://github.com/sobat-binangkit/echo-setup.git

WORKDIR currdir

# Get the source from github.com.
RUN git clone ${gitaddr} .

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -o echosvr .

############################
# STEP 2 build a small image
############################

FROM scratch

# Build Time Variables
ARG rootpath=/app
ARG httpport=8080
ARG httpsport=8585

# Environment Variables
ENV DATA_PATH=${rootpath}
ENV HTTP_PORT=${httpport}
ENV HTTPS_PORT=${httpsport} 

WORKDIR /app

# Copy our static executable.
COPY --from=builder ${currdir}/echosvr .

# Run the hello binary.
ENTRYPOINT ["/app/echosvr"]