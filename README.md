# go-tracert
计算机网路课程项目：基于 ICMP 使用 Golang 实现的 Tracert 应用。

## 说明
### 主要功能
- IPv4 Tracert
- IPv6 Tracert
- DNS Lookup
- Reverse DNS Lookup
- Set Maxhoop

### 依赖库
- flag
- fmt
- net
- golang.org/x/net/icmp
- golang.org/x/net/ipv4
- golang.org/x/net/ipv6

由于使用了 net 库中的 `ListenPacket()` 方法，在 Unix 和类 Unix 系统下运行需要 root 权限，否则会提示 `socket: operation not permitted`。

## 项目结构和功能
### main.go
接受用户从命令行参数传来的目标地址、最大跃点数、IPv4/IPv6 选项，调用 tracert.go 中的 `tracert` 方法。

### tracert.go
使用 `LookupIP` 方法解析目标的 IPv4 和 IPv6 地址。根据传入 IPv4/IPv6 参数，调用 tracert4/6.go 中的 `tracert4/6` 方法。如未指定 IPv4/IPv6 则优先使用 IPv6 追踪。

### tracert4/6.go
使用 `ListenPacket` 监听接口，使用 `ipv4/6` 和 `ICMP` 提供的方法封装请求和设置参数。每个跃点值发送三个包，根据 ICMP 的响应类型，输出结果，如果类型为 `ICMPTypeEchoReply` 提前结束循环，完成追踪。


## 使用示例
### 编译
```
$ go build
```
### 帮助
```
$ sudo ./tracert

用法: tracert [-h maximum_hops] [-4] [-6] target_name

选项:
        -h maximum_hops 搜索目标的最大跃点数。
        -4      强制使用 IPv4。
        -6      强制使用 IPV6
```
### 路由追踪
#### 默认
```
$ sudo ./tracert www.cloudflare.com

通过最多 30 个跃点跟踪
到 www.cloudflare.com [2606:4700::6811:d209] 的路由:

1       0 ms    0 ms    0 ms    Hidden for privacy
2       6 ms    4 ms    3 ms    Hidden for privacy
3       5 ms    3 ms    3 ms    Hidden for privacy
4       14 ms   15 ms   15 ms   2408:8000:a004:1::14
5       11 ms   12 ms   12 ms   2408:8000:2:8::
6       11 ms   11 ms   11 ms   2408:8000:2:511::
7       *       *       *       请求超时。
8       *       *       *       请求超时。
9       *       *       *       请求超时。
10      *       *       *       请求超时。
11      *       163 ms  158 ms  ae-3.r00.tokyjp08.jp.bb.gin.ntt.net. [2001:218:0:2000::2d7]
12      160 ms  157 ms  153 ms  as7515.ntt.net. [2001:218:2000:5000::26]
13      169 ms  169 ms  170 ms  2606:4700::6811:d209

跟踪完成。
```

#### 强制使用 IPv4
```
$ sudo ./tracert -4 www.cloudflare.com

通过最多 30 个跃点跟踪
到 www.cloudflare.com [104.17.209.9] 的路由:

1       0 ms    0 ms    0 ms    _gateway [172.16.0.1]
2       *       *       *       请求超时。
3       4 ms    4 ms    4 ms    Hidden for privacy
4       15 ms   15 ms   15 ms   Hidden for privacy
5       10 ms   10 ms   10 ms   219.158.18.1
6       16 ms   15 ms   15 ms   219.158.19.82
7       10 ms   15 ms   15 ms   219.158.19.85
8       195 ms  199 ms  199 ms  219.158.116.234
9       174 ms  173 ms  174 ms  219.158.39.106
10      255 ms  247 ms  266 ms  ip4.gtt.net. [173.205.51.142]
11      245 ms  226 ms  235 ms  104.17.209.9

跟踪完成。
```

#### 强制使用 IPv6
```
$ sudo ./tracert -6 www.cloudflare.com

通过最多 30 个跃点跟踪
到 www.cloudflare.com [2606:4700::6811:d109] 的路由:

1       0 ms    0 ms    0 ms    Hidden for privacy
2       4 ms    3 ms    3 ms    Hidden for privacy
3       3 ms    3 ms    4 ms    Hidden for privacy
4       9 ms    15 ms   16 ms   2408:8000:a004:1::14
5       12 ms   12 ms   12 ms   2408:8000:2:90c::
6       22 ms   14 ms   14 ms   2408:8000:2:50d::
7       *       116 ms  *       2408:8000:2:686::1
8       *       *       118 ms  2408:8000:2:7a2::
9       *       *       *       请求超时。
10      *       *       *       请求超时。
11      *       158 ms  158 ms  ae-3.r00.tokyjp08.jp.bb.gin.ntt.net. [2001:218:0:2000::2d7]
12      162 ms  *       164 ms  as7515.ntt.net. [2001:218:2000:5000::26]
13      193 ms  193 ms  201 ms  2606:4700::6811:d109

跟踪完成。
```

#### 指定最大跃点数
```
$ sudo ./tracert -h 10 www.cloudflare.com

通过最多 10 个跃点跟踪
到 www.cloudflare.com [2606:4700::6811:d209] 的路由:

1       0 ms    0 ms    0 ms    Hidden for privacy
2       4 ms    3 ms    4 ms    Hidden for privacy
3       5 ms    5 ms    4 ms    Hidden for privacy
4       13 ms   15 ms   15 ms   2408:8000:a004:1::14
5       12 ms   13 ms   10 ms   2408:8000:2:8::
6       11 ms   11 ms   11 ms   2408:8000:2:511::
7       *       *       72 ms   2408:8000:2:686::1
8       *       117 ms  119 ms  2408:8000:2:7a3::
9       *       *       *       请求超时。
10      *       *       *       请求超时。

跟踪完成。
```

## 已知问题
1. 由于使用了 `ListenPacket()` 方法，Windows 下无法使用。

2. 目标地址只能放在其它参数后，否则其它参数不生效。