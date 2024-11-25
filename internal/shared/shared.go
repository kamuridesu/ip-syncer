package shared

type IPInfo struct {
	IP      string `json:"ip"`
	Name    string `json:"name"`
	Authkey string `json:"authkey"`
	Changed bool   `json:"changed"`
}

func NewIPInfo(ip string, authkey string, changed bool) *IPInfo {
	return &IPInfo{
		IP:      ip,
		Authkey: authkey,
		Changed: changed,
	}
}
