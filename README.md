# ghost

Tools for network proxy, management and protection.

## 组成

+ 系统管理程序 cmd - CLI & Web API & GUI
    + [ ] ghc (ghost-client) 内网访问客户端
    + [ ] ghm (ghost-master) 内网路由主节点（需要公网IP、MongoDB）
    + [ ] ghn (ghost-node)   内网服务从节点
+ 内部网络模块 net (for ghost-master & ghost-node & ghost-client)
    + [ ] p2p 点对点连接
    + [ ] robot 爬虫（用于爬取代理资源站）、节点间共享代理资源
    + [ ] opt 网络优化、可用性、测速
    + [ ] cache 缓存
    + [ ] stat 流量统计、分析
    + [ ] conn 连接模式 (tcp/kcp)
+ 系统接入模块 gateway (for ghost-client)
    + [ ] iptables (for linux)
    + [ ] proxy (common)
    + [ ] nict (NIC teaming 网卡聚合)
    + [ ] vpn
+ 内网服务模块 service (for ghost-node & ghost-client)
    + [ ] webdav
    + [ ] fs 文件系统
    + [ ] staticfs 静态文件
    + [ ] vhost 虚拟主机
    + [ ] webserver Web服务器
    + [ ] proxy 反向代理
    + [ ] http 代理
    + [ ] socks 代理
    + [ ] ssr
    + [ ] v2ray
    + [ ] dns
+ 系统调试模块 debug
    + [ ] protocol 网络协议 (用于支持报文解析修改等)
        + [ ] http(s) (协议升级、隧道解析)
        + [ ] http2 考虑是否合并到上面
    + [ ] hook 调试钩子、抓包
+ 服务防御模块 defense // TODO: 伪装、防御、留证 组件
    + [ ] scan 端口扫描防御、信息欺骗、留证
~~+ 外部调用模块 (考虑之后独立出去)~~
~~    + [ ] sh 脚本解析~~
~~        > // TO-DO: 加入数组类型~~
~~        > // TO-DO: 当前语法分析器是暂时的，之后要重写~~
~~    + [ ] exec 程序调用封装~~
+ 外部调用模块（使用成熟的脚本）
    + https://github.com/yuin/gopher-lua
    + https://github.com/aarzilli/golua
    + https://www.youbbs.org/t/2851
    + https://github.com/Azure/golua
    + https://www.jianshu.com/p/c9ab2ae410b0
