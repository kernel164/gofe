package settings

import (
	"gopkg.in/ini.v1"
	"log"
	"net/url"
)

type BackendType string
type SchemeType string

const (
	SSH BackendType = "ssh"
)

const (
	HTTP SchemeType = "http"
	UNIX SchemeType = "unix"
)

var (
	Backend    BackendType
	BackendSsh struct {
		Host string
	}
	Server struct {
		Scheme  SchemeType
		Bind    string
		Statics []string
	}
)

func Load() {
	Cfg, err := ini.Load("gofe.ini")
	if err != nil {
		log.Println(err)
		return
	}

	// Global Section
	g := Cfg.Section("")
	Backend = SSH
	if g.Key("BACKEND").String() == "ssh" {
		ssh := Cfg.Section("backend.ssh")
		BackendSsh.Host = ssh.Key("HOST").MustString("localhost:22")
	}

	// Server section
	server := Cfg.Section("server")
	bind := server.Key("BIND").MustString("http://localhost:4000")
	u, err := url.Parse(bind)
	if err != nil {
		panic(err)
	}
	if u.Scheme == "http" {
		Server.Scheme = HTTP
	}
	if u.Scheme == "unix" {
		Server.Scheme = UNIX
	}
	Server.Bind = u.Host
	Server.Statics = server.Key("STATICS").Strings(",")
}
