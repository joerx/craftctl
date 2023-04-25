#!/bin/bash

set -e -o pipefail

exec > /var/log/minecraft-setup.log
exec 2>&1

USERNAME=minecraft
GROUPNAME=minecraft

INSTALL_DIR=/home/$USERNAME/server
TOOLS_DIR=/home/$USERNAME/tools
MCRCON_VERSION=0.7.2
MINECRAFT_DOWNLOAD_URL="https://piston-data.mojang.com/v1/objects/8f3112a1049751cc472ec13e397eade5336ca7ae/server.jar"

export DEBIAN_FRONTEND=noninteractive

# Update and upgrade system packages
apt-get update
apt-get \
  -o Dpkg::Options::=--force-confold \
  -o Dpkg::Options::=--force-confdef \
  -y --allow-downgrades --allow-remove-essential --allow-change-held-packages \
  dist-upgrade

apt-get -y install pwgen dbus-user-session

# Install Zulu SDK from their repos
# See https://docs.azul.com/core/zulu-openjdk/install/debian#install-from-azul-apt-repository
apt-get -y install gnupg curl
apt-key adv \
  --keyserver hkp://keyserver.ubuntu.com:80 \
  --recv-keys 0xB1998361219BD9C9

curl -O https://cdn.azul.com/zulu/bin/zulu-repo_1.0.0-3_all.deb
apt-get -y install ./zulu-repo_1.0.0-3_all.deb
apt-get update && apt-get -y install zulu19-jre

# Create user
loginctl enable-linger $USERNAME
mkdir -p $INSTALL_DIR

# Download mcrcon
echo "Downloading mcrcon version ${MCRCON_VERSION} "
mkdir -p $TOOLS_DIR/mcrcon
cd $TOOLS_DIR/mcrcon
wget -q https://github.com/Tiiffi/mcrcon/releases/download/v${MCRCON_VERSION}/mcrcon-${MCRCON_VERSION}-linux-x86-64.tar.gz
tar xzf mcrcon-${MCRCON_VERSION}-linux-x86-64.tar.gz
rm mcrcon-${MCRCON_VERSION}-linux-x86-64.tar.gz

RCON_PASSWORD=$(pwgen 20)

# Download minecraft
cd $INSTALL_DIR
echo "Download from ${MINECRAFT_DOWNLOAD_URL}"
wget -q "${MINECRAFT_DOWNLOAD_URL}"
echo "Download complete"

# Setup systemd service
cat <<- EOF > /usr/lib/systemd/user/minecraft.service
[Unit]
Description=Minecraft Server
Requires=dbus.socket

[Service]
Nice=1
SuccessExitStatus=0 1
WorkingDirectory=$INSTALL_DIR
ReadWriteDirectories=$INSTALL_DIR
ExecStart=/usr/bin/java -Xmx1024M -Xms1024M -jar server.jar nogui
ExecStop=$TOOLS_DIR/mcrcon/mcrcon -H 127.0.0.1 -P 25575 -p password stop

[Install]
WantedBy=default.target
EOF

chmod 664 /usr/lib/systemd/user/minecraft.service

# Basic configuration
echo "eula=true" > $INSTALL_DIR/eula.txt

cat << EOF > $INSTALL_DIR/server.properties
#Minecraft server properties
#(File modification date and time)
enable-jmx-monitoring=false
level-seed=
gamemode=survival
enable-command-block=false
enable-query=false
generator-settings={}
enforce-secure-profile=true
level-name=world
motd=A Minecraft Server
query.port=25565
pvp=true
generate-structures=true
max-chained-neighbor-updates=1000000
difficulty=easy
network-compression-threshold=256
max-tick-time=60000
require-resource-pack=false
use-native-transport=true
max-players=20
online-mode=true
enable-status=true
allow-flight=false
initial-disabled-packs=
broadcast-rcon-to-ops=true
view-distance=10
server-ip=
resource-pack-prompt=
allow-nether=true
server-port=25565
enable-rcon=true
rcon.password=password
rcon.port=25575
sync-chunk-writes=true
op-permission-level=4
prevent-proxy-connections=false
hide-online-players=false
resource-pack=
entity-broadcast-range-percentage=100
simulation-distance=10
player-idle-timeout=0
force-gamemode=false
rate-limit=0
hardcore=false
white-list=false
broadcast-console-to-ops=true
spawn-npcs=true
spawn-animals=true
function-permission-level=2
initial-enabled-packs=vanilla
level-type=minecraft\:normal
text-filtering-config=
spawn-monsters=true
enforce-whitelist=false
spawn-protection=16
resource-pack-sha1=
max-world-size=29999984
EOF

chown -R $USERNAME:$GROUPNAME $TOOLS_DIR
chown -R $USERNAME:$GROUPNAME $INSTALL_DIR

# Start minecraft
echo "Starting minecraft"

systemctl --user -M $USERNAME@ daemon-reload
systemctl --user -M $USERNAME@ start minecraft
