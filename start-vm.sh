#!/bin/bash

# Script to start a development VM using libvirt. Metadata for the VM will be stored in a new folder `.vm` in 
# the current working dir. Cloud-init will be used to install a Minecraft server running as systemd unit for the primary user 
# (`minecraft` by default). Current working directory will be shared with the VM using virtiofs (mounted automatically).
# What exactly will be installed in the VM can be managed in `install-server.sh`
# Base image for the VM can be customized, for now this script support Ubuntu only and the install script depends
# heavily on the Ubuntu utilities and repositories. Other Debian-based distros may or may not work.

#./start-vm -f -r jammy my-project

set -e -o pipefail

# Image will be downloaded if not exists
VM_IMG_DIR=${VM_IMG_DIR:-$HOME/Devel/isos}

# Cleanup existing VM if exists, yes/no
DO_FORCE="no"

# Only Ubuntu releases are supported for now
RELEASE="jammy"

# Name of the libvirt domain to create
NAME=$(basename $PWD)

# Local cache dir
VMDIR=.vm

parse_opts() {
  while getopts 'fr:h' opt; do
    case "$opt" in
      f)
        DO_FORCE=yes
        ;;
      r)
        arg="$OPTARG"
        RELEASE=$arg
        ;;
      ?|h)
        echo "Usage: $(basename $0) [-f] [-r RELEASE] NAME"
        exit 1
        ;;
    esac
  done
  shift "$(($OPTIND -1))"

  if [[ "$1" != "" ]]; then
    NAME=$1
  fi
}

create_vm() {
  local NAME=$1
  local IMG_PATH=$2

  local SSH_RSA=$(cat ~/.ssh/id_rsa.pub)

  local AWS_ACCESS_KEY_ID=$(terraform -chdir=terraform output -raw aws_access_key_id)
  local AWS_SECRET_ACCESS_KEY=$(terraform -chdir=terraform output -raw aws_secret_access_key)

  if [[ -z "$AWS_ACCESS_KEY_ID" ]]; then
    >&2 echo "No valid aws_access_key_id in terraform output"
    return 1
  fi

  # Local cache dir
  mkdir -p $VMDIR

  # Generate meta-data and user-data files
  cat << EOF > $VMDIR/meta-data
instance-id: $NAME
local-hostname: $NAME
EOF

  cat << EOF > $VMDIR/cloudinit.yml
#cloud-config

users:
- name: minecraft
  ssh_authorized_keys:
  - $SSH_RSA
  sudo: ['ALL=(ALL) NOPASSWD:ALL']
  groups: sudo
  shell: /bin/bash

packages:
- golang-go
- curl
- gnupg

mounts:
- [ code, /mnt/code, virtiofs, "rw,relatime,nofail", "0", "0"]

write_files:
- path: /home/minecraft/.aws/credentials
  content: |
    [profile default]
    aws_access_key_id = $AWS_ACCESS_KEY_ID
    aws_secret_access_key = $AWS_SECRET_ACCESS_KEY
- path: /home/minecraft/.aws/config
  content: |
    [profile default]
    region = ap-southeast-1

EOF

  # Generate multi-part cloudinit mime archive
  # See https://cloudinit.readthedocs.io/en/latest/explanation/format.html#mime-multi-part-archive
  cloud-init devel make-mime -a $VMDIR/cloudinit.yml:cloud-config -a install-server.sh:x-shellscript > $VMDIR/user-data

  # Generate file system from base image
  qemu-img create -b $IMG_PATH -f qcow2 -F qcow2 $VMDIR/$NAME.qcow2 10G

  # Generate ISO image for cloudinit
  genisoimage -output $VMDIR/cidata.iso -V cidata -r -J $VMDIR/user-data $VMDIR/meta-data

  # Virt-install
  virt-install \
    --name $NAME \
    --ram 2048 \
    --vcpus 2 \
    --import \
    --disk path=$VMDIR/$NAME.qcow2,format=qcow2 \
    --disk path=$VMDIR/cidata.iso,device=cdrom \
    --os-variant=ubuntu20.04 \
    --memorybacking access.mode=shared \
    --filesystem source=$PWD,target=code,accessmode=passthrough,driver.type=virtiofs \
    --noautoconsole

}

cleanup_vm() {
  local NAME=$1
  echo "Cleaning up existing VM '${NAME}'"

  if ! virsh domstate $NAME > /dev/null 2>&1; then
    echo "Domain ${NAME} does not seem to exist"
    return 0
  fi


  if [[ $(virsh domstate $NAME) == "running" ]]; then
    echo "Shutting down domain $NAME"
    virsh shutdown ${NAME}
    sleep 3
  fi

  while :; do
    local STATE=$(virsh domstate $NAME)
    if [[ $STATE == "shut off" ]]; then
      break
    else
      echo "'$NAME' is in state state $STATE, waiting for 'shut off'"
      sleep 1
    fi
  done

  echo "Undefining domain ${NAME}"
  virsh undefine ${NAME}
}

download_image() {
  local RELEASE=$1
  local IMG_URL="~"
  local IMG_NAME="~"

  case $RELEASE in
  jammy | focal | bionic)
    IMG_URL=https://cloud-images.ubuntu.com/$RELEASE/current/$RELEASE-server-cloudimg-amd64.img
    IMG_NAME=$RELEASE-server-cloudimg-amd64.img
    ;;
  *)
    >&2 echo "Unsupported release: '$RELEASE'"
    exit 1
    ;;
  esac

  # If base image does not exist on disk, download it
  local IMG_PATH="${VM_IMG_DIR}/${IMG_NAME}"
  if [[ ! -f $IMG_PATH ]]; then
    >&2 echo "Downloading $IMG_URL"
    mkdir -p $VM_IMG_DIR
    curl -o $IMG_PATH $IMG_URL
  else
    >&2 echo "Using existing image $IMG_PATH for '$RELEASE'"
  fi

  echo $IMG_PATH
}


parse_opts $@
echo "RELEASE=$RELEASE; NAME=$NAME; DO_FORCE=$DO_FORCE"

if [[ -d $VMDIR ]]; then
  if [[ "$DO_FORCE" != "yes" ]]; then
    echo "VM already exists, to re-create, run vm-destroy.sh first"
    exit 1
  fi
  cleanup_vm $NAME
  rm -rf $VMDIR
fi

IMG_PATH=$(download_image $RELEASE)

create_vm $NAME $IMG_PATH

while :; do
  STATE=$(virsh domstate $NAME)
  if [[ $STATE == "running" ]]; then
    break
  else
    echo "'$NAME' is in state state $STATE, waiting for 'running'"
    sleep 1
  fi
  sleep 1
done

# FIXME: find a way to get network status dynamically
echo "Waiting for network"
sleep 10
virsh domifaddr $NAME
