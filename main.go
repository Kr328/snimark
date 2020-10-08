package main

import (
	"flag"
	"fmt"

	"github.com/kr328/snimark/config"
	"github.com/kr328/snimark/server"
)

func main() {
	var path string

	flag.StringVar(&path, "c", "", "Configuration file path")

	flag.Parse()

	if path == "" {
		println("Invalid configuration file path")

		return
	}

	cfg, err := config.Parse(path)
	if err != nil {
		fmt.Printf("Parse config: %s\n", err.Error())
		return
	}

	srv, err := server.New(cfg)
	if err != nil {
		fmt.Printf("Start service: %s\n", err.Error())
		return
	}

	if err := srv.Exec(); err != nil {
		fmt.Printf("Execute service: %s\n", err.Error())
		return
	}
}
