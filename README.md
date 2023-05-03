# Minecraft Server Control Panel

Run a Minecraft server:

```sh
java -Xmx1024M -Xms1024M -jar server.jar nogui
```

Run this application:

```sh
go run cmd/server.go --rcon-addr 127.0.0.1:25575 server
```

Hot reload using [nodemon](https://www.npmjs.com/package/nodemon):

```sh
nodemon -e go,html --exec go run cmd/server.go --rcon-addr 127.0.0.1:25575 server --signal SIGTERM
```

## Terraform

An S3 bucket is required to store backups. Apply the config in `./terraform` to create one. **Warning**: The resulting setup is _not safe for production use!_

```sh
terraform -chdir=terraform apply
```

Don't forget to destroy any resources created when you're done.

## Development VM

Make sure Terraform has been applied. Then start the VM (NB, the `-f` flag will forcibly destroy any existing VM of the same name):

```sh
./start-vm.sh -f steve
```

If output shows no network address, run `virsh domifaddr steve` again after a few secs:

```sh
$ virsh domifaddr steve
 Name       MAC address          Protocol     Address
-------------------------------------------------------------------------------
 vnet2      52:54:00:51:88:b5    ipv4         192.168.124.73/24
```

SSH into the VM using the IP obtained from previous step:

```sh
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no minecraft@<VM_IP>
```

Start the development server:

```sh
cd /mnt/code
go run cmd/server.go --rcon-passwd password --s3-bucket '<BUCKET_NAME_FROM_TERRAFORM>' --s3-region ap-southeast-1 server
```

The web interface will be available at `<VM_IP>:8080`.
