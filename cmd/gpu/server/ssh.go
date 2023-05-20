package gpu_server

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Cli struct {
	user       string
	pwd        string
	host       string
	port       string
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func NewSSHClient(user, pwd, host, port string) *Cli {
	return &Cli{
		user: user,
		pwd:  pwd,
		host: host,
		port: port,
	}
}

func (c *Cli) getConfig_nokey() *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.pwd),
		},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config
}

func (c *Cli) Connect() error {
	config := c.getConfig_nokey()
	client, err := ssh.Dial("tcp", c.host+":"+c.port, config)
	if err != nil {
		return fmt.Errorf("connect server error: %w", err)
	}
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("new sftp client error: %w", err)
	}

	c.sshClient = client
	c.sftpClient = sftp
	return nil
}

func (c Cli) Run(cmd string) (string, error) {
	if c.sshClient == nil {
		if err := c.Connect(); err != nil {
			return "", err
		}
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("create new session error: %w", err)
	}
	defer session.Close()

	buf, err := session.CombinedOutput(cmd)
	return string(buf), err
}

func (c Cli) UploadFile(localFile, remoteFileName string) (int, error) {
	if c.sshClient == nil {
		if err := c.Connect(); err != nil {
			return -1, err
		}
	}
	file, err := os.Open(localFile)
	if nil != err {
		return -1, fmt.Errorf("open local file failed: %w", err)
	}
	defer file.Close()

	ftpFile, err := c.sftpClient.Create(remoteFileName)
	if nil != err {
		return -1, fmt.Errorf("Create remote path failed: %w", err)
	}
	defer ftpFile.Close()

	fileByte, err := ioutil.ReadAll(file)
	if nil != err {
		return -1, fmt.Errorf("read local file failed: %w", err)
	}

	ftpFile.Write(fileByte)

	return 0, nil
}

func (c Cli) DownloadFile(remoteFile, localFile string) (int, error) {
	if c.sshClient == nil {
		if err := c.Connect(); err != nil {
			return -1, err
		}
	}
	source, err := c.sftpClient.Open(remoteFile)
	if err != nil {
		return -1, fmt.Errorf("sftp client open file error: %w", err)
	}
	defer source.Close()

	target, err := os.OpenFile(localFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return -1, fmt.Errorf("open local file error: %w", err)
	}
	defer target.Close()

	n, err := io.Copy(target, source)
	if err != nil {
		return -1, fmt.Errorf("write file error: %w", err)
	}
	return int(n), nil
}
