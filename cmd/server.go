package main

import (
	"flag"
	"fmt"
	"joerx/minecraft-cli/internal/server"
	"log"
)

var Version = "development"

func main() {
	cfg := server.Config{}

	flag.StringVar(&cfg.Addr, "addr", ":8080", "Server address")
	flag.StringVar(&cfg.RCONHostPort, "rcon-addr", "127.0.0.1:25575", "Address of Minecraft RCON server")
	flag.StringVar(&cfg.RCONPasswd, "rcon-passwd", "passwd", "Password for Minecraft RCON server")
	flag.StringVar(&cfg.MCWorldDir, "world-dir", "./server/world", "Directory with Minecraft world data")
	flag.StringVar(&cfg.UnitName, "unit-name", "minecraft.service", "Systemd unit name used for the minecraft server")
	flag.StringVar(&cfg.S3Bucket, "s3-bucket", "", "S3 bucket to upload backup files to")
	flag.StringVar(&cfg.S3Region, "s3-region", "", "AWS region the bucket is located in")

	flag.Parse()

	log.Printf("Server config: %#v", cfg)

	cmd := flag.Arg(0)
	switch cmd {
	case "server":
		if err := server.Run(cfg); err != nil {
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
