package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		runTUI()
		return
	}

	cmd := args[0]
	switch cmd {
	case "tui":
		runTUI()
	case "daemon":
		runDaemon()
	case "status":
		cliStatus()
	case "start":
		cliStart()
	case "stop":
		cliStop()
	case "restart":
		cliRestart()
	case "global":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "用法: flclash-headless global <出口名>")
			os.Exit(1)
		}
		cliGlobal(strings.Join(args[1:], " "))
	case "tun":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "用法: flclash-headless tun on|off")
			os.Exit(1)
		}
		switch args[1] {
		case "on":
			cliTun(true)
		case "off":
			cliTun(false)
		default:
			fmt.Fprintln(os.Stderr, "用法: flclash-headless tun on|off")
			os.Exit(1)
		}
	case "logs":
		cliLogs()
	default:
		fmt.Fprintf(os.Stderr, "未知命令: %s\n", cmd)
		fmt.Fprintln(os.Stderr, "可用命令: tui, daemon, status, start, stop, restart, global, tun, logs")
		os.Exit(1)
	}
}

func runTUI() {
	RunDaemonTUI()
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
