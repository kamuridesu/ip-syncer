package hosts

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Hosts struct {
	hostPath string
	content  string
}

func (h *Hosts) Save() error {
	return os.WriteFile(h.hostPath, []byte(h.content), 0644)
}

func ReadHostsFile(hostPath string) (*Hosts, error) {
	content, err := os.ReadFile(hostPath)
	if err != nil {
		return nil, err
	}
	return &Hosts{hostPath: hostPath, content: string(content)}, nil
}

func (h *Hosts) AddOrReplaceHost(ip, name string) error {

	newEntry := ip + " " + name
	buffer := ""
	found := false
	oldContent := h.content
	for _, line := range strings.Split(h.content, "\n") {
		if strings.Contains(line, name) {
			if strings.Contains(line, ip) {
				return nil
			}
			slog.Info(fmt.Sprintf("IP assign to client %s already exists and will be replaced\n", name))
			found = true
			buffer += newEntry
		} else {
			buffer += line + "\n"
		}
	}
	h.content = buffer
	if !found {
		h.content += newEntry + "\n"
	}
	err := h.Save()
	if err != nil {
		h.content = oldContent
	}
	return err
}
