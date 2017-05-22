#!/bin/bash

# This script creates downloadable releases of the CLI app.

APP_FOLDER="../cmd/screp"
echo Using app folder: $APP_FOLDER

# Acquire app name
# Expected format: `    appName = "thename"`
APP_NAME=$(more $APP_FOLDER/screp.go | grep "appName\s*=" | cut -d '"' -f 2)
if [ -z "$APP_NAME" ]; then
    echo Could not detect app name!
    exit 1
fi
echo Detected app name: $APP_NAME

# Acquire app version
# Expected format: `    appVersion = "theversion"`
APP_VERSION=$(more $APP_FOLDER/screp.go | grep "appVersion\s*=" | cut -d '"' -f 2)
if [ -z "$APP_VERSION" ]; then
    echo Could not detect app version!
    exit 2
fi
echo Detected app version: $APP_VERSION

START=$(date +%s)

for REL_OS in linux windows darwin
do
    for REL_ARCH in amd64 386
    do
        REL_NAME=$APP_NAME-$APP_VERSION-$REL_OS-$REL_ARCH
        if [ $REL_OS = "windows" ]; then
            REL_NAME=$REL_NAME.zip
        else
            REL_NAME=$REL_NAME.tar.gz
        fi
        echo Creating release $REL_NAME...
        rm $REL_NAME 2> /dev/null
        EXEC_NAME=$APP_NAME
        if [ $REL_OS = "windows" ]; then
            EXEC_NAME=$EXEC_NAME.exe
        fi

        GOOS=$REL_OS GOARCH=$REL_ARCH go build -o $EXEC_NAME $APP_FOLDER || exit 3

        if [ $REL_OS = "windows" ]; then
            zip -q $REL_NAME $EXEC_NAME
        else
            tar -zcf $REL_NAME $EXEC_NAME
        fi
        rm $EXEC_NAME
    done
done

END=$(date +%s)
DIFF=$(echo "$END - $START" | bc)
echo Done in $DIFF sec.
