package main

import (
	"bytes"
	"text/template"
)

type SSHConfig struct {
	Host, HostName, User, IdentityFile, ProxyCommand string
	Port                                             string
}

const tpl = `
Host {{.Host}}
    HostName {{.HostName}}
    {{if .User}}User {{.User}}{{end}}
    {{if .IdentityFile}}IdentityFile {{.IdentityFile}}{{end}}
    {{if .Port}}Port {{.Port}}{{end}}
    {{if .ProxyCommand}}ProxyCommand {{.ProxyCommand}}{{end}}
`

func GenConfig(cfgs []SSHConfig) (*bytes.Buffer, error) {
	tmpl, _ := template.New("tpl").Parse(tpl)
	buf := bytes.NewBufferString("")

	for _, cfg := range cfgs {
		err := tmpl.Execute(buf, cfg)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}
