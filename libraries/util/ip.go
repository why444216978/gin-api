package util

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func GetRealIP(r *http.Request) string {
	remoteAddr := ""
	forward := r.Header.Get("X-Forwarded-For")
	var err error
	//从左到右找到第一个非内网地址为止
	ips := strings.Split(forward, ",")
	for _, ip := range ips {
		if !IsInnerIp(ip) {
			remoteAddr = ip
			break
		}
	}

	if remoteAddr == "" {
		remoteAddr = r.Header.Get("X-Real-Ip")
	}

	if remoteAddr == "" {
		remoteAddr, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Printf("parse remote failed,error is %s", err)
			remoteAddr = "127.0.0.1" //give a ip address to avoid return nil.
		}
		if remoteAddr == "::1" {
			remoteAddr = "127.0.0.1"
		}
	}

	return remoteAddr
}

//判断一个ip是否是内网ip
func IsInnerIp(ip_str string) bool {
	/**
	A类地址：10.0.0.0--10.255.255.255 (10/8 prefix)
	B类地址：172.16.0.0--172.31.255.255 (172.16/12 prefix)
	C类地址：192.168.0.0--192.168.255.255 (192.168/16 prefix)
	*/
	ip := net.ParseIP(ip_str)

	ip4 := ip.To4()
	if ip4 != nil {
		if ip4[0] == 10 {
			return true
		} else if ip4[0] == 172 {
			if ip4[1] >= 16 && ip4[1] <= 31 {
				return true
			}
		} else if ip4[0] == 192 && ip4[1] == 168 {
			return true
		}
	} else {
		//解析不了默认为内网ip
		return true
	}

	//排除回环地址
	return ip4.IsLoopback()
}
