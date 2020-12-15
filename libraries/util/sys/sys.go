package sys

import (
	"bytes"
	"errors"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"

	utli_err "gin-api/libraries/util/error"
)

//执行系统命令
func ExecCommand(name string, arg ...string) string{
	cmd := exec.Command(name, arg...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	_ = cmd.Run()

	return buf.String()
}

//执行系统pipe命令
//ps := exec.Command("ps", "-ef")
//grep := exec.Command("grep", "-i", "php-fpm")
func ExecCommandGrep(command *exec.Cmd, grep *exec.Cmd) string {
	r, w := io.Pipe() // 创建一个管道
	defer r.Close()
	defer w.Close()
	command.Stdout = w // command向管道的一端写
	grep.Stdin = r     // grep从管道的一端读

	var buffer bytes.Buffer
	grep.Stdout = &buffer // grep的输出为buffer

	_ = command.Start()
	_ = grep.Start()
	command.Wait()
	w.Close()
	grep.Wait()

	return buffer.String()
}

func ExternalIP() (string, error) {
	iFaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iFace := range iFaces {
		if iFace.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iFace.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iFace.Addrs()
		if err != nil {
			return "", err
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
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

//获得本机名
func HostName() string {
	hostNamePrefix := ""
	host, err := os.Hostname()
	utli_err.Must(err)
	if err == nil {
		parts := strings.SplitN(host, ".", 2)
		if len(parts) > 0 {
			hostNamePrefix = parts[0]
		}
	}
	return hostNamePrefix
}

func GetInternalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	panic("cannot get internal ip")
}
