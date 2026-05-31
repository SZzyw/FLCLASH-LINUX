package coreclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"

	"flclash-headless/model"
	"flclash-headless/storage"
)

type CoreStatus string

const (
	StatusStopped      CoreStatus = "stopped"
	StatusStarting     CoreStatus = "starting"
	StatusRunning      CoreStatus = "running"
	StatusError        CoreStatus = "error"
)

type CoreEventListener func(event CoreEvent)

type Client struct {
	mu          sync.RWMutex
	process     *exec.Cmd
	conn        net.Conn
	reader      *bufio.Reader
	status      CoreStatus
	corePath    string
	dataDir     string
	connectedAt time.Time
	listeners   []CoreEventListener

	pendingMu sync.Mutex
	pending   map[string]chan *ActionResult
}

func NewClient(corePath, dataDir string) *Client {
	return &Client{
		status:    StatusStopped,
		corePath:  corePath,
		dataDir:   dataDir,
		pending:   make(map[string]chan *ActionResult),
		listeners: make([]CoreEventListener, 0),
	}
}

func (c *Client) AddListener(listener CoreEventListener) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.listeners = append(c.listeners, listener)
}

func (c *Client) fireEvent(event CoreEvent) {
	c.mu.RLock()
	listeners := make([]CoreEventListener, len(c.listeners))
	copy(listeners, c.listeners)
	c.mu.RUnlock()
	for _, l := range listeners {
		l(event)
	}
}

func (c *Client) Status() CoreStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

func (c *Client) setStatus(s CoreStatus) {
	c.mu.Lock()
	c.status = s
	c.mu.Unlock()
}

func (c *Client) Uptime() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.status != StatusRunning {
		return 0
	}
	return time.Since(c.connectedAt)
}

func (c *Client) Start() error {
	c.setStatus(StatusStarting)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		c.setStatus(StatusError)
		return fmt.Errorf("listen: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	cmd := exec.Command(c.corePath, fmt.Sprintf("%d", port))
	if err := cmd.Start(); err != nil {
		listener.Close()
		c.setStatus(StatusError)
		return fmt.Errorf("start core: %w", err)
	}
	c.process = cmd

	acceptCh := make(chan net.Conn, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			c.setStatus(StatusError)
			return
		}
		acceptCh <- conn
		listener.Close()
	}()

	select {
	case conn := <-acceptCh:
		c.conn = conn
		c.reader = bufio.NewReader(conn)
		c.connectedAt = time.Now()
		c.setStatus(StatusRunning)
		go c.readLoop()
		return nil
	case <-time.After(10 * time.Second):
		listener.Close()
		cmd.Process.Kill()
		c.setStatus(StatusError)
		return fmt.Errorf("core connect timeout")
	}
}

func (c *Client) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	if c.process != nil && c.process.Process != nil {
		c.process.Process.Kill()
		c.process.Wait()
		c.process = nil
	}
	c.status = StatusStopped
	return nil
}

func (c *Client) readLoop() {
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			c.setStatus(StatusStopped)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var result ActionResult
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			continue
		}

		if result.ID == "" || result.ID == "message" {
			var event CoreEvent
			if dataBytes, err := json.Marshal(result.Data); err == nil {
				if err := json.Unmarshal(dataBytes, &event); err == nil {
					c.fireEvent(event)
				}
			}
			continue
		}

		c.pendingMu.Lock()
		ch, ok := c.pending[result.ID]
		if ok {
			delete(c.pending, result.ID)
		}
		c.pendingMu.Unlock()
		if ok {
			ch <- &result
		}
	}
}

func (c *Client) SendAction(action Action, timeout time.Duration) (*ActionResult, error) {
	data, err := json.Marshal(action)
	if err != nil {
		return nil, fmt.Errorf("marshal action: %w", err)
	}

	ch := make(chan *ActionResult, 1)
	c.pendingMu.Lock()
	c.pending[action.ID] = ch
	c.pendingMu.Unlock()

	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	if conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	if _, err := fmt.Fprintln(conn, string(data)); err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	select {
	case result := <-ch:
		return result, nil
	case <-time.After(timeout):
		c.pendingMu.Lock()
		delete(c.pending, action.ID)
		c.pendingMu.Unlock()
		return nil, fmt.Errorf("timeout")
	}
}

func (c *Client) InitClash(version int) error {
	params := InitParams{
		HomeDir: storage.GetDataDir(),
		Version: version,
	}
	paramsBytes, _ := json.Marshal(params)
	action := NewAction(ActionInitClash, string(paramsBytes))
	result, err := c.SendAction(action, 30*time.Second)
	if err != nil {
		return err
	}
	if result.Code == ResultError {
		return fmt.Errorf("init failed: %v", result.Data)
	}
	return nil
}

func (c *Client) SetupConfig(selectedMap map[string]string, testURL string) error {
	if testURL == "" {
		testURL = "https://www.gstatic.com/generate_204"
	}
	params := SetupParams{
		SelectedMap: selectedMap,
		TestUrl:     testURL,
	}
	paramsBytes, _ := json.Marshal(params)
	action := NewAction(ActionSetupConfig, string(paramsBytes))
	result, err := c.SendAction(action, 30*time.Second)
	if err != nil {
		return err
	}
	if result.Code == ResultError {
		return fmt.Errorf("setup config failed: %v", result.Data)
	}
	return nil
}

func (c *Client) UpdateConfig(params UpdateParams) error {
	paramsBytes, _ := json.Marshal(params)
	action := NewAction(ActionUpdateConfig, string(paramsBytes))
	result, err := c.SendAction(action, 10*time.Second)
	if err != nil {
		return err
	}
	if result.Code == ResultError {
		return fmt.Errorf("update config failed: %v", result.Data)
	}
	if msg, ok := result.Data.(string); ok && msg != "" {
		return fmt.Errorf("update config failed: %s", msg)
	}
	return nil
}

func (c *Client) ValidateConfig(path string) (string, error) {
	action := NewAction(ActionValidateConfig, path)
	result, err := c.SendAction(action, 10*time.Second)
	if err != nil {
		return "", err
	}
	if result.Code == ResultError {
		return "", fmt.Errorf("validate failed: %v", result.Data)
	}
	return "", nil
}

func (c *Client) GetProxies() (*ProxiesData, error) {
	action := NewAction(ActionGetProxies, nil)
	result, err := c.SendAction(action, 10*time.Second)
	if err != nil {
		return nil, err
	}
	if result.Code == ResultError {
		return nil, fmt.Errorf("get proxies failed: %v", result.Data)
	}
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, err
	}
	var proxies ProxiesData
	if err := json.Unmarshal(dataBytes, &proxies); err != nil {
		return nil, err
	}
	return &proxies, nil
}

func (c *Client) ChangeProxy(groupName, proxyName string) error {
	params := ChangeProxyParams{
		GroupName: groupName,
		ProxyName: proxyName,
	}
	paramsBytes, _ := json.Marshal(params)
	action := NewAction(ActionChangeProxy, string(paramsBytes))
	result, err := c.SendAction(action, 10*time.Second)
	if err != nil {
		return err
	}
	if result.Code == ResultError {
		return fmt.Errorf("change proxy failed: %v", result.Data)
	}
	return nil
}

func decodeResultJSON(data interface{}, v interface{}) error {
	switch x := data.(type) {
	case string:
		if x == "" {
			return nil
		}
		return json.Unmarshal([]byte(x), v)
	default:
		b, err := json.Marshal(x)
		if err != nil {
			return err
		}
		if string(b) == "null" {
			return nil
		}
		return json.Unmarshal(b, v)
	}
}

func (c *Client) GetTraffic() (*model.TrafficSnapshot, error) {
	action := NewAction(ActionGetTraffic, false)
	result, err := c.SendAction(action, 5*time.Second)
	if err != nil {
		return nil, err
	}
	if result.Code == ResultError {
		return &model.TrafficSnapshot{}, nil
	}
	var traffic model.TrafficSnapshot
	if err := decodeResultJSON(result.Data, &traffic); err != nil {
		return nil, err
	}
	return &traffic, nil
}

func (c *Client) GetTotalTraffic() (*model.TotalTraffic, error) {
	action := NewAction(ActionGetTotalTraffic, false)
	result, err := c.SendAction(action, 5*time.Second)
	if err != nil {
		return nil, err
	}
	var traffic model.TotalTraffic
	if err := decodeResultJSON(result.Data, &traffic); err != nil {
		return nil, err
	}
	return &traffic, nil
}

func (c *Client) StartListener() error {
	action := NewAction(ActionStartListener, nil)
	result, err := c.SendAction(action, 5*time.Second)
	if err != nil {
		return err
	}
	if result.Code == ResultError {
		return fmt.Errorf("start listener failed: %v", result.Data)
	}
	return nil
}

func (c *Client) StopListener() error {
	action := NewAction(ActionStopListener, nil)
	_, err := c.SendAction(action, 5*time.Second)
	return err
}

func (c *Client) StartLog() error {
	action := NewAction(ActionStartLog, nil)
	_, err := c.SendAction(action, 5*time.Second)
	return err
}

func (c *Client) StopLog() error {
	action := NewAction(ActionStopLog, nil)
	_, err := c.SendAction(action, 5*time.Second)
	return err
}

func (c *Client) Shutdown() error {
	action := NewAction(ActionShutdown, true)
	_, err := c.SendAction(action, 10*time.Second)
	return err
}

func (c *Client) CloseConnections() error {
	action := NewAction(ActionCloseConnections, nil)
	_, err := c.SendAction(action, 5*time.Second)
	return err
}

func (c *Client) ResetTraffic() error {
	action := NewAction(ActionResetTraffic, nil)
	_, err := c.SendAction(action, 5*time.Second)
	return err
}

func (c *Client) GetConnections() (interface{}, error) {
	action := NewAction(ActionGetConnections, nil)
	result, err := c.SendAction(action, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetConfig(path string) (interface{}, error) {
	action := NewAction(ActionGetConfig, path)
	result, err := c.SendAction(action, 10*time.Second)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) TestDelay(url, proxyName string) (int, error) {
	params := DelayParams{
		ProxyName: proxyName,
		Timeout:   5000,
		TestURL:   url,
	}
	paramsBytes, _ := json.Marshal(params)
	action := NewAction(ActionAsyncTestDelay, string(paramsBytes))
	result, err := c.SendAction(action, 10*time.Second)
	if err != nil {
		return -1, err
	}
	if result.Code == ResultError {
		return -1, nil
	}
	delayStr, ok := result.Data.(string)
	if !ok {
		return -1, nil
	}
	var delay DelayEvent
	if err := json.Unmarshal([]byte(delayStr), &delay); err != nil {
		return -1, nil
	}
	return delay.Value, nil
}

