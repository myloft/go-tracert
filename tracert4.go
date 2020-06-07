package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// SendTracertMsg 向目标地址发送指定的 ttl 的 Tracert 包
func SendTracertMsg(dst net.IPAddr, ttl int) (int64, icmp.Type, net.Addr) {
	// 监听操作系统 raw socket 接口
	c, err := net.ListenPacket("ip4:1", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	// 设置 ipv4 头
	p := ipv4.NewPacketConn(c)
	if err := p.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
		log.Fatal(err)
	}
	// 构建 ICMP 消息
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
	// 创建 IP 数据报
	rb := make([]byte, 1500)
	// 将 ttl 作为
	wm.Body.(*icmp.Echo).Seq = ttl
	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}
	// 设置 TTL 值
	if err := p.SetTTL(ttl); err != nil {
		log.Fatal(err)
	}
	// 初始时间
	begin := time.Now()

	// 封装 ICMP 消息
	if _, err := p.WriteTo(wb, nil, &dst); err != nil {
		log.Fatal(err)
	}
	// 设置消息超时时间
	if err := p.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		log.Fatal(err)
	}
	// 读取返回的 IP 数据报信息
	n, _, peer, err := p.ReadFrom(rb)
	if err != nil {
		// 如果超时返回不可达
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return 0, ipv4.ICMPTypeDestinationUnreachable, peer
		}
		log.Fatal(err)
	}
	// 取出 ICMP 消息
	rm, err := icmp.ParseMessage(1, rb[:n])
	if err != nil {
		log.Fatal(err)
	}
	// 往返时间
	rtt := time.Since(begin).Milliseconds()
	return rtt, rm.Type, peer
}

// Tracert4 ipv4 路由追踪
func Tracert4(host string, dst net.IPAddr, maxhoop int) {
	// 反查 IP
	names, _ := net.LookupAddr(dst.IP.String())
	if names == nil {
		names = append(names, host)
	}
	fmt.Printf("\n通过最多 %v 个跃点跟踪\n到 %v [%s] 的路由:\n\n", maxhoop, names[0], dst.IP)

	// 发送 ICMP 包
ICMP:
	for i := 1; i <= maxhoop; i++ {
		rtts := make([]int64, 3)
		icmptypes := make([]icmp.Type, 3)
		peers := make([]net.Addr, 3)
		// 输出序号
		fmt.Printf("%d\t", i)

		// 第一组请求
		rtts[0], icmptypes[0], peers[0] = SendTracertMsg(dst, i)
		switch icmptypes[0] {
		case ipv4.ICMPTypeTimeExceeded:
			fmt.Printf("%d ms\t", rtts[0])
		case ipv4.ICMPTypeEchoReply:
			fmt.Printf("%d ms\t", rtts[0])
		default:
			fmt.Printf("*\t")
		}

		// 第二组请求
		rtts[1], icmptypes[1], peers[1] = SendTracertMsg(dst, i)
		switch icmptypes[1] {
		case ipv4.ICMPTypeTimeExceeded:
			fmt.Printf("%d ms\t", rtts[1])
		case ipv4.ICMPTypeEchoReply:
			fmt.Printf("%d ms\t", rtts[1])
		default:
			fmt.Printf("*\t")
		}

		// 第三组请求
		rtts[2], icmptypes[2], peers[2] = SendTracertMsg(dst, i)
		switch icmptypes[2] {
		case ipv4.ICMPTypeTimeExceeded:
			fmt.Printf("%d ms\t", rtts[2])
		case ipv4.ICMPTypeEchoReply:
			fmt.Printf("%d ms\t", rtts[2])
		default:
			fmt.Printf("*\t")
		}

		// 判断返回的 ICMP 状态
		for i, icmptype := range icmptypes {
			switch icmptype {
			case ipv4.ICMPTypeTimeExceeded:
				// 反查 IP
				names, _ := net.LookupAddr(peers[i].String())
				if names != nil {
					fmt.Printf("%s [%s]\n", names[0], peers[i])
				} else {
					fmt.Printf("%s\n", peers[i])
				}
				continue ICMP
			case ipv4.ICMPTypeEchoReply:
				// 反查 IP
				names, _ := net.LookupAddr(peers[i].String())
				if names != nil {
					fmt.Printf("%s [%s]\n", names[0], peers[i])
				} else {
					fmt.Printf("%s\n", peers[i])
				}
				break ICMP
			}
		}
		fmt.Printf("请求超时。\n")
	}
	fmt.Printf("\n跟踪完成。\n\n")
}
