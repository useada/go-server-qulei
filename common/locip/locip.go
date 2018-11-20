package locip

import (
	"errors"
	"net"
)

func GetLocalIP() (string, error) {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		if ipnet.IP.To4() == nil {
			continue
		}
		return ipnet.IP.String(), nil
	}
	return "", errors.New("empty ip")
}
