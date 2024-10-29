package android

type IpInfo struct {
	Ip           string `json:"ip"`
	ProxyCountry string `json:"proxy_country"`
	Org          string `json:"org"`
}

type VpnInfo struct {
	VpnName string `json:"vpn_name"`
	VpnId   string `json:"vpn_id"`
}

type DeviceInfo struct {
	DevicePlatform string `json:"device_platform"`
	DeviceId       string `json:"device_id"`
	DeviceType     string `json:"device_type"`
}

type Session struct {
	IpInfo       IpInfo     `json:"ip_Info"`
	VpnInfo      VpnInfo    `json:"vpn_info"`
	DeviceInfo   DeviceInfo `json:"device_info"`
	DeviceParams Device     `json:"device_params"`
}

type Account struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Passwd    string `json:"passwd"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}
