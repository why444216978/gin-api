package util

import (
	"os"
	"strings"
)

var (
	hostNamePrefix string
)

func init() {
	host, err := os.Hostname()
	if err == nil {
		parts := strings.SplitN(host, ".", 2)
		if len(parts) > 0 {
			hostNamePrefix = parts[0]
		}
	}
}

//返回主机名前缀,如zf-test
func HostNamePrefix() string {
	return hostNamePrefix
}
