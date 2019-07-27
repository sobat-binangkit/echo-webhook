############################
# STEP 1
############################
FROM frolvlad/alpine-glibc

# Maintainer
LABEL maintainer="Ahmad R. Djarkasih<djarkasih@gmail.com>"

# Build Time Variables
ARG appname=webhook
ARG httpport=8080
ARG httpsport=8585

# Environment Variables
ENV APP_NAME=$appname \
    APP_DIR=/apps \
    HTTP_PORT=$httpport \
    HTTPS_PORT=$httpsport 

WORKDIR /apps

COPY ./$APP_NAME main

# Run the app
CMD ["/apps/main"]
