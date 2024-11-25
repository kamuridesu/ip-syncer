package shared

import "os"

type EnvInfo struct {
	AuthKey        string
	ServerEndpoint string
	Name           string
}

var (
	Info = EnvInfo{
		AuthKey:        os.Getenv("AUTH_KEY"),
		ServerEndpoint: os.Getenv("SERVER_ENDPOINT"),
		Name:           os.Getenv("CLIENT"),
	}
)

type IPInfo struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
}

func NewIPInfo(ip, name string) *IPInfo {
	return &IPInfo{
		IP:   ip,
		Name: name,
	}
}

func (i *IPInfo) Equals(other *IPInfo) bool {
	return i.IP == other.IP && i.Name == other.Name
}
