package mi

import (
	"bufio"
	"net"
	"time"
)

//Client 用于连接到OpenVPN管理界面
type Client struct {
	MINetwork string
	MIAddress string
}

//NewClient 初始化管理接口客户端结构
func NewClient(network, address string) *Client {
	c := &Client{
		MINetwork: network, //Management Interface Network
		MIAddress: address, //Management Interface Address
	}

	return c
}

//GetPid 获取OpenVPN服务器的进程ID
func (c *Client) GetPid() (int64, error) {
	str, err := c.Execute("pid")
	if err != nil {
		return -1, err
	}
	return ParsePid(str)
}

//GetVersion 获取OpenVPN服务器的版本
func (c *Client) GetVersion() (*Version, error) {
	str, err := c.Execute("version")
	if err != nil {
		return nil, err
	}
	return ParseVersion(str)
}

//GetStatus 获取连接的客户端列表和路由表
func (c *Client) GetStatus() (*Status, error) {
	str, err := c.Execute("status 2")
	if err != nil {
		return nil, err
	}
	return ParseStatus(str)
}

//GetStatus 获取连接的客户端列表和路由表
func (c *Client) GetLogs() (string, error) {
	str, err := c.Execute("log all")
	if err != nil {
		return "", err
	}
	return str, err
}

//返回已连接客户端的数量和网络流量的总数
func (c *Client) GetLoadStats() (*LoadStats, error) {
	str, err := c.Execute("load-stats")
	if err != nil {
		return nil, err
	}
	return ParseStats(str)
}

//杀死OpenVPN连接
func (c *Client) KillSession(cname string) (string, error) {
	str, err := c.Execute("kill " + cname)
	if err != nil {
		return "", err
	}
	return ParseKillSession(str)
}

//Signal 向守护程序发送信号
func (c *Client) Signal(signal string) error {
	str, err := c.Execute("signal " + signal)
	if err != nil {
		return err
	}
	return ParseSignal(str)
}

//Execute 连接到OpenVPN服务器，发送命令并读取响应
func (c *Client) Execute(cmd string) (string, error) {
	conn, err := net.DialTimeout(c.MINetwork, c.MIAddress, 3*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	buf := bufio.NewReader(conn)
	buf.ReadString('\n') //read welcome message
	err = SendCommand(conn, cmd)
	if err != nil {
		return "", err
	}

	return ReadResponse(buf)
}
