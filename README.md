# FlClash - Linux 无图形终端代理客户端

> **部署包下载**: [github.com/SZzyw/FLCLASH-LINUX-project](https://github.com/SZzyw/FLCLASH-LINUX-project)


![Screenshot](https://zhaoyingwei.dpdns.org/unimportant/20260528181800301.png)

FlClash 是一款基于 [Clash.Meta](https://github.com/MetaCubeX/Clash.Meta) 内核的 **Linux 纯终端代理客户端**，采用双进程架构，无需桌面环境即可运行。适用于服务器、开发机、ARM 设备等无 GUI 场景。


## 架构概览

```
┌──────────────────────────────────────────┐
│  flclash-headless                        │
│  ├── daemon  (后台守护进程)               │
│  └── tui / status / start / stop ...     │
│           │                              │
│           │ Unix Socket (JSON-RPC)       │
│           ▼                              │
│  ┌──────────────────────────────────┐    │
│  │  FlClashCore (Clash.Meta 引擎)   │    │
│  │  混合端口 :7890                  │    │
│  │  控制端口 :9090                  │    │
│  └──────────────────────────────────┘    │
└──────────────────────────────────────────┘
```

- **flclash-headless** (Go 编写) — 管理守护进程 + CLI/TUI 客户端
- **FlClashCore** (预编译) — 静态链接的 Clash.Meta 内核，无需外部运行时

## 功能特性

- **TUI 终端界面** — 仪表盘、代理节点切换、延迟测试、日志查看、配置管理
- **订阅管理** — URL 在线导入、本地文件导入、多配置切换、订阅更新
- **双进程架构** — 后台 daemon 管理内核生命周期，TUI/CLI 通过 Unix Socket 通信
- **TUN 虚拟网卡** — 支持全局透明代理（需 root）
- **三种模式** — 规则(Rule) / 全局(Global) / 直连(Direct)
- **代理节点组** — 查看组详情、切换节点、批量延迟测试
- **systemd 集成** — 支持开机自启
- **持久化存储** — 配置、状态、首选项自动保存

## 快速开始

### 1. 环境要求

- Linux (x86_64 / arm64)
- Go 1.22+（仅编译 headless 时需要）

### 2. 编译 headless 程序

```bash
cd headless
go build -o flclash-headless .
```

### 3. 给核心加执行权限

```bash
chmod +x FlClashCore
```

### 4. 启动

```bash
# 启动后台 daemon（需 root 才能使用 TUN）
sudo ./flclash-headless daemon &

# 进入 TUI 交互界面
sudo ./flclash-headless tui
```

### 5. TUI 首次配置

进入 TUI 后：
1. 按 `c` 进入配置管理
2. 按 `a` 输入订阅 URL 或本地文件路径
3. 按 `b` 返回仪表盘
4. 按 `r` 启动代理核心

## TUI 快捷键

| 按键 | 功能 |
|------|------|
| `1` | 仪表盘（流量、速度、内存） |
| `2` / `g` | 代理节点组 |
| `3` / `x` | 全局出口选择 |
| `4` / `m` | 切换运行模式 |
| `5` / `l` | 查看日志 |
| `6` / `c` | 配置/订阅管理 |
| `r` | 启动/停止核心 |
| `n` | 切换 TUN |
| `q` | 退出 TUI |

## CLI 命令

```bash
sudo ./flclash-headless status            # 查看运行状态
sudo ./flclash-headless start             # 启动代理核心
sudo ./flclash-headless stop              # 停止代理核心
sudo ./flclash-headless restart           # 重启代理核心
sudo ./flclash-headless global <出口名>   # 切换全局出口
sudo ./flclash-headless tun on|off        # 切换 TUN
sudo ./flclash-headless logs              # 查看日志
```

## systemd 开机自启

```bash
sudo cp flclash-headless /usr/local/bin/
sudo cp FlClashCore /usr/local/bin/
sudo chmod +x /usr/local/bin/FlClashCore

sudo tee /etc/systemd/system/flclash-headless.service << 'EOF'
[Unit]
Description=FlClash Headless Daemon
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/flclash-headless daemon
Restart=on-failure
RestartSec=3
User=root
WorkingDirectory=/root

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now flclash-headless
```

## 目录结构

```
flcc/
├── assets/data/           # GEOIP / GEOSITE 路由数据库
├── core/                  # FlClashCore Go 源码
├── headless/              # flclash-headless Go 源码
│   ├── action/            # 业务操作（启动/停止/导入/切换...）
│   ├── app/               # 应用状态管理
│   ├── configbuilder/     # YAML 配置构建
│   ├── coreclient/        # 核心进程通信
│   ├── input/             # 终端输入
│   ├── model/             # 数据模型
│   ├── renderer/          # TUI 渲染
│   ├── storage/           # JSON 持久化
│   └── util/              # 工具函数
├── FlClashCore            # 预编译核心引擎（53MB 静态链接）
├── 正确.md                 # Ubuntu 22.04 部署指南
└── LICENSE                # GPL v3
```

## 注意事项

- **执行权限**: `FlClashCore` 必须有 `chmod +x`
- **工作目录**: 命令在 `headless/` 目录下执行，相对路径 `../` 指向项目根目录
- **TUN 需 root**: TUN 模式需要 root 权限
- **数据目录**: 自动使用 `headless/../assets/data/`（GEOIP/GEOSITE 文件所在目录）
- 默认混合端口 `7890`，外部控制 `127.0.0.1:9090`
- 更换订阅后需重启核心

## 许可证

[GNU General Public License v3](LICENSE)
