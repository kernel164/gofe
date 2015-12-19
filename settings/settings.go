package settings

import (
	"gopkg.in/ini.v1"
	"log"
)

type BackendType string

const (
	SSH BackendType = "ssh"
)

var (
	Backend BackendType
	SshHost string
)

func Load() {
	Cfg, err := ini.Load("gofe.ini")
	if err != nil {
		log.Println(err)
		return
	}
	g := Cfg.Section("")
	Backend = SSH
	if g.Key("BACKEND").String() == "ssh" {
		ssh := Cfg.Section("ssh")
		SshHost = ssh.Key("HOST").MustString("localhost:22")
	}
}
