package hosts

import (
	"os"
	"strings"
)

const HostsPath = "/etc/hosts"

type Hosts struct {
	oldContent string
	newContent string
}

func (h *Hosts) Save() error {
	return os.WriteFile(HostsPath, []byte(h.newContent), 0644)
}

func ReadHostsFile() (*Hosts, error) {
	content, err := os.ReadFile(HostsPath)
	if err != nil {
		return nil, err
	}
	return &Hosts{oldContent: string(content), newContent: ""}, nil
}

func (h *Hosts) AddOrReplaceHost(ip, name string) *Hosts {
	if h.newContent == "" {
		h.newContent = h.oldContent
	}

	newEntry := ip + " " + name
	buffer := ""
	found := false
	for _, line := range strings.Split(h.oldContent, "\n") {
		if strings.Contains(line, name) {
			found = true
			buffer += newEntry
		} else {
			buffer += line + "\n"
		}
	}
	h.newContent = buffer
	if !found {
		h.newContent += newEntry + "\n"
	}
	return h
}
