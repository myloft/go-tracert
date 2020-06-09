package main

import (
	"flag"
	"fmt"
)

func main() {
	var v4, v6 bool
	var host string
	var maxhoop int

	// 解析命令行参数
	flag.IntVar(&maxhoop, "h", 30, "搜索目标的最大跃点数。")
	flag.BoolVar(&v4, "4", false, "强制使用 IPv4。")
	flag.BoolVar(&v6, "6", false, "强制使用 IPv6。")
	flag.Parse()

	if flag.Arg(0) != "" {
		host = flag.Arg(0)
	} else {
		fmt.Printf("\n用法: tracert [-h maximum_hops] [-4] [-6] target_name\n")
		fmt.Printf("\n选项:\n\t-h maximum_hops\t搜索目标的最大跃点数。\n\t-4\t强制使用 IPv4。\n\t-6\t强制使用 IPV6\n\n")
		return
	}

	// 根据参数进行 Tracert
	if !v4 && !v6 {
		Tracert(host, maxhoop, 0)
	} else if v4 && !v6 {
		Tracert(host, maxhoop, 4)
	} else if !v4 && v6 {
		Tracert(host, maxhoop, 6)
	} else {
		fmt.Printf("\n用法: tracert [-h maximum_hops] [-4] [-6] target_name\n")
		fmt.Printf("\n选项:\n\t-h maximum_hops\t搜索目标的最大跃点数。\n\t-4\t使用 IPv4。\n\t-6\t使用 IPV6\n\n")
	}
}
