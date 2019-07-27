############################
# STEP 1
############################
FROM alpine:latest

# Maintainer
LABEL maintainer="Ahmad R. Djarkasih<djarkasih@gmail.com>"

# Build Time Variables
ARG appname=webhook
ARG appdir=/app
ARG httpport=8080
ARG httpsport=8585

# Environment Variables
ENV APP_NAME=$appname
ENV APP_DIR=$appdir
ENV DATA_DIR=$appdir/data
ENV PLUGIN_DIR=$appdir/handlers
ENV HTTP_PORT=$httpport
ENV HTTPS_PORT=$httpsport 

COPY ./$APP_NAME $APP_DIR/$APP_NAME

# Run the app
CMD ["sh","-c","$APP_DIR/$APP_NAME"]
