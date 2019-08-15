package network

import (
	"net"
	"net/http"
	"strings"
)

var localNetworks []*net.IPNet

func init() {
	localNetworks = make([]*net.IPNet, 19)
	for i, sNetwork := range []string{
		"10.0.0.0/8",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"172.17.0.0/12",
		"172.18.0.0/12",
		"172.19.0.0/12",
		"172.20.0.0/12",
		"172.21.0.0/12",
		"172.22.0.0/12",
		"172.23.0.0/12",
		"172.24.0.0/12",
		"172.25.0.0/12",
		"172.26.0.0/12",
		"172.27.0.0/12",
		"172.28.0.0/12",
		"172.29.0.0/12",
		"172.30.0.0/12",
		"172.31.0.0/12",
		"192.168.0.0/16",
	} {
		_, network, _ := net.ParseCIDR(sNetwork)
		localNetworks[i] = network
	}
}

// 获取请求的ip地址
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// HasLocalIPAddr 检测 IP 地址字符串是否是内网地址
func HasLocalIPAddr(ip string) bool {
	return HasLocalIP(net.ParseIP(ip))
}

// HasLocalIP 检测 IP 地址是否是内网地址
func HasLocalIP(ip net.IP) bool {
	for _, network := range localNetworks {
		if network.Contains(ip) {
			return true
		}
	}

	return ip.IsLoopback()
}

// ClientPublicIP 尽最大努力实现获取客户端公网 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" && !HasLocalIPAddr(ip) {
			return ip
		}
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" && !HasLocalIPAddr(ip) {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		if !HasLocalIPAddr(ip) {
			return ip
		}
	}

	return ""
}