package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

type Action struct {
	ID     string      `json:"id"`
	Type   int         `json:"type"`
	Data   interface{} `json:"data"`
}

type ActionResult struct {
	ID     string      `json:"id"`
	Method string      `json:"method"`
	Data   interface{} `json:"data"`
	Code   int         `json:"code"`
}

const (
	ActionInitClash = 1003
	ActionSetupConfig = 1004
	ActionGetProxies = 1009
)

func main() {
	// Listen for core connection
	listener, err := net.Listen("tcp", "127.0.0.1:8899")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Listen failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Listening on 127.0.0.1:8899")

	// Start the core (it will connect to us)
	cmd := exec.Command("/usr/local/bin/FlClashCore", "8899")
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Start core failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Core process started, waiting for connection...")

	// Accept core connection
	conn, err := listener.Accept()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Accept failed: %v\n", err)
		cmd.Process.Kill()
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Core connected")
	reader := bufio.NewReader(conn)

	sendAction := func(a Action) (*ActionResult, error) {
		data, _ := json.Marshal(a)
		_, err := conn.Write(append(data, '\n'))
		if err != nil {
			return nil, fmt.Errorf("write: %w", err)
		}
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		var result ActionResult
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			return nil, fmt.Errorf("unmarshal: %w", err)
		}
		return &result, nil
	}

	// InitClash
	fmt.Println("Sending InitClash...")
	params, _ := json.Marshal(map[string]interface{}{
		"home-dir": os.Getenv("HOME") + "/.local/share/flclash-headless",
		"version":  1,
	})
	result, err := sendAction(Action{
		ID:   fmt.Sprintf("%d", time.Now().UnixNano()),
		Type: ActionInitClash,
		Data: string(params),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "InitClash failed: %v\n", err)
		cmd.Process.Kill()
		os.Exit(1)
	}
	fmt.Printf("InitClash result: code=%d data=%v\n", result.Code, result.Data)

	// SetupConfig (this is the one that times out)
	fmt.Println("Sending SetupConfig...")
	sp, _ := json.Marshal(map[string]interface{}{
		"selected-map": map[string]string{},
		"test-url":     "https://www.gstatic.com/generate_204",
	})
	result, err = sendAction(Action{
		ID:   fmt.Sprintf("%d", time.Now().UnixNano()),
		Type: ActionSetupConfig,
		Data: string(sp),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "SetupConfig ERROR: %v\n", err)
	} else {
		fmt.Printf("SetupConfig result: code=%d data=%v\n", result.Code, result.Data)
	}

	// If SetupConfig succeeded, get proxies
	if err == nil && result.Code == 0 {
		fmt.Println("Sending GetProxies...")
		result, err = sendAction(Action{
			ID:   fmt.Sprintf("%d", time.Now().UnixNano()),
			Type: ActionGetProxies,
			Data: "",
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "GetProxies ERROR: %v\n", err)
		} else {
			fmt.Printf("GetProxies result: code=%d\n", result.Code)
		}
	}

	fmt.Println("DONE")
	cmd.Process.Kill()
	time.Sleep(500 * time.Millisecond)
}
