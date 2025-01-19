package util

import (
	"net"

	"github.com/spf13/viper"
)

// 判断是否是开发环境
func IsDev() bool {
	return viper.GetString("env") == "development" || viper.GetString("env") == "dev"
}

// 辅助函数
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// GetAllIPAddresses 获取所有网络接口的 IP 地址
func GetAllIPAddresses() ([]string, error) {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		// 跳过 down 状态的接口
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// 跳过回环接口
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			// 只获取 IPv4 地址
			ip = ip.To4()
			if ip == nil {
				continue
			}
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}
