#!/bin/bash
APP_NAME=geoip
GROUP=$APP_NAME
BIN_DIR=/usr/bin
CONF_DIR=/etc/$APP_NAME
CONF_FILE=config.toml
DATA_DIR=/var/lib/$APP_NAME
DEFAULTS_FILE=/etc/default/$APP_NAME
LOG_DIR=/var/log/$APP_NAME
LOGROTATE_DIR=/etc/logrotate.d
SYSTEMD_DIR=/lib/systemd/system

set -x
rm -rf $CONF_DIR
rm -rf $DATA_DIR
rm -f $DEFAULTS_FILE
rm -rf $LOG_DIR
rm -f $BIN_DIR/$APP_NAME

systemctl stop $APP_NAME
systemctl disable $APP_NAME
rm -f $SYSTEMD_DIR/$APP_NAME.service
systemctl daemon-reload

userdel $APP_NAME
set +x