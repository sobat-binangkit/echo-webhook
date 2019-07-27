############################
# STEP 1
############################
FROM alpine:latest

# Maintainer
LABEL maintainer="Ahmad R. Djarkasih<djarkasih@gmail.com>"

# Build Time Variables
ARG appname=webhook
ARG rootpath=/app
ARG httpport=8080
ARG httpsport=8585

# Environment Variables
ENV APP_NAME=$appname
ENV ROOT_DIR=$rootpath
ENV DATA_DIR=$rootpath/data
ENV PLUGIN_DIR=$rootpath/handlers
ENV HTTP_PORT=$httpport
ENV HTTPS_PORT=$httpsport 

ADD ./$APP_NAME $ROOT_DIR/$APP_NAME

# Run the app
ENTRYPOINT $ROOT_DIR/$APP_NAME