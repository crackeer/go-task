package util

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *ssh.Client
}

// NewSSHClient
//
//	@param server
//	@return *SSHClient
//	@return error
func NewSSHClient(host, port, user, password string) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return nil, fmt.Errorf("dial error: %v", err)
	}

	return &SSHClient{
		client: client,
	}, nil
}

func (s *SSHClient) Close() {
	s.client.Close()
}

func (s *SSHClient) Mkdir(tempDir string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("new ssh session error: %s", err.Error())
	}
	defer session.Close()
	cmd := fmt.Sprintf("mkdir -p %s", tempDir)
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("exec ssh command error: %s", err.Error())
	}
	return nil
}

func (s *SSHClient) RemoveDir(tempDir string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("new ssh session error: %s", err.Error())
	}
	defer session.Close()
	cmd := fmt.Sprintf("rm -rf %s", tempDir)
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("exec ssh command error: %s", err.Error())
	}
	return nil
}

func (s *SSHClient) Run(cmd string) error {
	session, err := s.client.NewSession()

	if err != nil {
		return fmt.Errorf("new ssh session error: %s", err.Error())
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	defer session.Close()
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("exec ssh command error: %s", err.Error())
	}
	return nil
}

func (s *SSHClient) K3sImport(remoteFile string) error {
	return s.Run(fmt.Sprintf("k3s ctr image import %s", remoteFile))
}

func (s *SSHClient) DockerLoad(remoteFile string) error {
	return s.Run(fmt.Sprintf("docker load -i %s", remoteFile))
}

// UploadTo
//
//	@param localFile
//	@param remoteFile
//	@return error
func (s *SSHClient) UploadTo(localFile, remoteFile string) error {
	sftpClient, err := sftp.NewClient(s.client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("open local file error: %s", err.Error())
	}
	defer srcFile.Close()

	_, err = os.Stat(localFile)
	if err != nil {
		return fmt.Errorf("file stat error: %v", err)
	}

	dstFile, err := sftpClient.Create(remoteFile)
	if err != nil {
		return fmt.Errorf("create remote file error: %s", err.Error())
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)

	return err
}

// DownloadTo
//
//	@param remoteFile
//	@param localFile
//	@return error
func (s *SSHClient) DownloadTo(remoteFile, localFile string) error {
	sftpClient, err := sftp.NewClient(s.client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := sftpClient.Open(remoteFile)
	if err != nil {
		return fmt.Errorf("open remote file error: %s", err.Error())
	}
	defer srcFile.Close()

	dstFile, err := os.Create(localFile)
	if err != nil {
		return fmt.Errorf("create local file error: %s", err.Error())
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)

	return err
}

// RunCommand
//
//	@param command
//	@return error
func (s *SSHClient) RunCommand(command string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("new ssh session error: %s", err.Error())
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	defer session.Close()
	if err := session.Run(command); err != nil {
		return fmt.Errorf("exec ssh command error: %s", err.Error())
	}
	return nil
}
