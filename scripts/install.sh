#!/bin/bash
USER=geoip
GROUP=geoip
BIN_DIR=/usr/bin
DATA_DIR=/var/lib/$USER
LOG_DIR=/var/log/$USER
LOGROTATE_DIR=/etc/logrotate.d

function install_systemd {
    cp -f geoip.service $1
    systemctl enable geoip|| true
    systemctl daemon-reload || true
}

function install_update_rcd {
    update-rc.d geoip defaults
}

function install_chkconfig {
    chkconfig --add geoip
}

if ! id geoip &>/dev/null; then
    useradd --system -U -M $USER -s /bin/false -d $DATA_DIR
fi

# check if DATA_DIR exists
if [ ! -d "$DATA_DIR" ]; then
    mkdir -p $DATA_DIR
    chown $USER:$GROUP $DATA_DIR
fi

# check if LOG_DIR exists
if [ ! -d "$LOG_DIR" ]; then
    mkdir -p $LOG_DIR
    chown $USER:$GROUP $DATA_DIR
fi

#create base conf file
if [[ ! -f geoip.toml ]]; then
  ./$USER config >> config.toml
fi

#install
chown $USER:$GROUP $USER
#Copy binary to BIN_DIR
cp $USER $BIN_DIR

# Add defaults file, if it doesn't exist
if [[ ! -f /etc/default/geoip ]]; then
    touch /etc/default/geoip
fi

# If 'geoip.conf' is not present use package's sample (fresh install)
if [[ ! -f /etc/geoip/geoip.toml ]]; then
   cp /etc/geoip/geoip.toml /etc/geoip/geoip.toml
fi

#Install as systemd service
if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
	install_systemd /lib/systemd/system/geoip.service
	deb-systemd-invoke restart geoip.service || echo "WARNING: systemd not running."
fi