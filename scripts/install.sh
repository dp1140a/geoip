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

function install_systemd {
    cp -f $APP_NAME.service $1
    systemctl enable $APP_NAME || true
    systemctl daemon-reload || true
}

function install_update_rcd {
    update-rc.d $APP_NAME defaults
}

function install_chkconfig {
    chkconfig --add $APP_NAME
}

# Create user isf not exists
if ! id $APP_NAME &>/dev/null; then
    echo -e "User $APP_NAME not found.  Creating user $APP_NAME. "
    useradd --system -U -M $APP_NAME -s /bin/false -d $DATA_DIR
fi

# check if DATA_DIR exists
if [ ! -d "$DATA_DIR" ]; then
    echo -e "Creating data dir at $DATA_DIR"
    mkdir -p $DATA_DIR
    chown $APP_NAME:$GROUP $DATA_DIR
fi

# check if LOG_DIR exists
if [ ! -d "$LOG_DIR" ]; then
    echo -e "Creating log dir at $LOG_DIR"
    mkdir -p $LOG_DIR
    chown $APP_NAME:$GROUP $LOG_DIR
fi

#create base conf file
if [[ ! -f $CONF_FILE ]]; then
  echo -e "Creating sample config at $(pwd)/$CONF_FILE"
  ./$APP_NAME config >> $CONF_FILE
fi

# Add defaults file, if it doesn't exist
if [[ ! -f "$DEFAULTS_FILE" ]]; then
    echo -e "Creating defaults file at $DEFAULTS_FILE"
    touch $DEFAULTS_FILE
    echo "CONF_FILE=$CONF_DIR/$CONF_FILE" >> $DEFAULTS_FILE
fi

#Check conf_dir exists
if [ ! -d "$CONF_DIR" ]; then
  mkdir $CONF_DIR
fi

# If 'config.toml' is not present use package's sample (fresh install)
if [[ ! -f $CONF_DIR/$CONF_FILE ]]; then
    echo -e "No existing conf found. Copying sample conf to $CONF_DIR/$CONF_FILE"
   cp $CONF_FILE $CONF_DIR/$CONF_FILE
fi

#install
chown $APP_NAME:$GROUP $APP_NAME
#Copy binary to BIN_DIR
echo -e "Installing binary at $BIN_DIR"
cp $APP_NAME $BIN_DIR

#Install as systemd service
if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
  echo -e "Installing systemd service at /lib/systemd/system/geoip.service"
	install_systemd $SYSTEMD_DIR/$APP_NAME.service
	deb-systemd-invoke restart $APP_NAME.service || echo "WARNING: systemd not running."
fi

echo -e "Installation of $APP_NAME complete"