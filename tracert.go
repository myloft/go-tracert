package main

import (
	"errors"
	"fmt"
	"net"
	"os"
)

// LookupHostIP 解析 Host 对应的 IP 地址
func LookupHostIP(host string) ([]net.IPAddr, error) {
	ip46 := make([]net.IPAddr, 2)
	ips, err := net.LookupIP(host)
	if err != nil {
		return ip46, err
	}

	// 获取 ipv4 地址
	for _, ip := range ips {
		if ip.To4() != nil {
			ip46[0].IP = ip
			break
		}
	}

	// 获取 ipv6 地址
	for _, ip := range ips {
		if ip.To16() != nil && ip.To4() == nil {
			ip46[1].IP = ip
			break
		}
	}

	// 不存在 A AAAA 记录
	if ip46[0].IP == nil && ip46[1].IP == nil {
		return ip46, errors.New("No A & AAAA Records")
	}

	return ip46, nil
}

// Tracert 路由追踪
func Tracert(host string, maxhoop int, ttype int) {
	ips, err := LookupHostIP(host)
	if err != nil {
		fmt.Printf("无法解析目标系统名称 %s。\n\n", host)
		os.Exit(0)
	}
	// 根据选项进行追踪 优先 IPv6
	if ttype == 4 && ips[0].IP != nil {
		Tracert4(host, ips[0], maxhoop)
	} else if ttype == 6 && ips[1].IP != nil {
		Tracert6(host, ips[1], maxhoop)
	} else if ttype == 0 && ips[1].IP != nil {
		Tracert6(host, ips[1], maxhoop)
	} else if ttype == 0 && ips[0].IP != nil {
		Tracert4(host, ips[0], maxhoop)
	} else {
		fmt.Printf("无法解析目标系统名称 %s。\n\n", host)
	}
}
