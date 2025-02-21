package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/relex/aini"
)

var (
	inventory string
	output    string
	backup    bool
)

func init() {
	flag.StringVar(&inventory, "i", "/etc/ansible/hosts", "path to the inventory file")
	flag.StringVar(&output, "o", "/etc/ssh/ssh_config.d/ansible.conf", "path to the output file")
	flag.BoolVar(&backup, "b", false, "backup the old one")
}

func main() {
	flag.Parse()

	data, err := aini.ParseFile(inventory)
	if err != nil {
		log.Fatal(err)
	}

	cfgs := make([]SSHConfig, 0)
	for _, item := range data.Hosts {
		host := item.Vars["ansible_host"]
		port := item.Vars["ansible_ssh_port"]
		user := item.Vars["ansible_ssh_user"]
		identityFile := item.Vars["ansible_ssh_private_key_file"]

		if host == "" {
			slog.Warn(fmt.Sprintf("[%s] empty host, skip", item.Name))
			continue
		}

		cfgs = append(cfgs, SSHConfig{
			Host:         item.Name,
			HostName:     host,
			User:         user,
			IdentityFile: identityFile,
			Port:         port,
		})
	}

	buf, err := GenConfig(cfgs)
	if err != nil {
		log.Fatal(err)
	}

	if backup && FileExists(output) {
		if err := os.Rename(output, output+".bak"); err != nil {
			log.Fatal(err)
		}
		log.Printf("backup old ssh config file to %s", output+".bak")

	}

	if err := os.WriteFile(output, buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
	log.Printf("write new ssh config file to %s", output)
}
