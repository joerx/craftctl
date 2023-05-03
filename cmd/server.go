package main

import (
	"flag"
	"fmt"
	"joerx/minecraft-cli/internal/server"
	"log"
)

var Version = "development"

func main() {
	srvOpts := server.Opts{}

	flag.StringVar(&srvOpts.Addr, "addr", ":8080", "Server address")
	flag.StringVar(&srvOpts.RCONHostPort, "rcon-addr", "127.0.0.1:25575", "Address of Minecraft RCON server")
	flag.StringVar(&srvOpts.RCONPasswd, "rcon-passwd", "passwd", "Password for Minecraft RCON server")
	flag.StringVar(&srvOpts.MCWorldDir, "world-dir", "./server/world", "Directory with Minecraft world data")
	flag.StringVar(&srvOpts.UnitName, "unit-name", "minecraft.service", "Systemd unit name used for the minecraft server")
	flag.StringVar(&srvOpts.S3Bucket, "s3-bucket", "", "S3 bucket to upload backup files to")
	flag.StringVar(&srvOpts.S3Region, "s3-region", "", "AWS region the bucket is located in")

	flag.Parse()

	log.Printf("server opts: %#v", srvOpts)

	cmd := flag.Arg(0)
	switch cmd {
	case "server":
		if err := server.Run(srvOpts); err != nil {
			log.Fatal(err)
		}
	case "version":
		version()
	default:
		log.Fatalf("Invalid command: '%s'", cmd)
	}
}

func version() {
	fmt.Println(Version)
}
