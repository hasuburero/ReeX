package ssh

import (
	"errors"
	"fmt"
	"github.com/hasuburero/ReeX/lib/controller/config/confssh"
	"io"
	"os"
	"path/filepath"
)

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// section 1
type Host struct {
	NodeName string
	IP       string
	User     string
	WorkDir  string
	SSHconf  *ssh.ClientConfig
}

// section 2
const (
	Auth_Pubkey     = "pubkey"
	Auth_Passkey    = "passkey"
	ErrorPubkeyPath = "Invalid Pubkey path\n"

	Port    = ":22"
	CD      = "cd "
	AndList = " && "
)

// section3
func (self *Host) ExecAsync(cmd string) error {
	conn, err := ssh.Dial("tcp", self.IP+Port, self.SSHconf)
	if err != nil {
	}
}
func (self *Host) Exec(cmd string) ([]byte, error) {
	conn, err := ssh.Dial("tcp", self.IP+Port, self.SSHconf)
	if err != nil {
		return []byte{}, err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return []byte{}, err
	}
	defer session.Close()

	output, err := session.CombinedOutput(CD + self.WorkDir + AndList + cmd)
	if err != nil {
		return []byte{}, err
	}

	return output, nil
}

func (self *Host) Delete(target string) error {
	conn, err := ssh.Dial("tcp", self.IP+Port, self.SSHconf)
	if err != nil {
		return err
	}
	defer conn.Close()

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	err = sftpClient.Remove(target)
	if err != nil {
		return err
	}

	return nil
}

func (self *Host) CopyR2L(src, dest string) error {
	conn, err := ssh.Dial("tcp", self.IP+Port, self.SSHconf)
	if err != nil {
		return err
	}
	defer conn.Close()

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	localfd, err := os.Create(dest + filepath.Base(dest))
	if err != nil {
		return err
	}
	defer localfd.Close()

	remotefd, err := sftpClient.Open(dest)
	if err != nil {
		return err
	}
	defer remotefd.Close()

	_, err = io.Copy(localfd, remotefd)
	if err != nil {
		return err
	}

	return nil
}

func (self *Host) CopyL2R(src, dest string) error {
	conn, err := ssh.Dial("tcp", self.IP+Port, self.SSHconf)
	if err != nil {
		return err
	}
	defer conn.Close()

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	localfd, err := os.Open(src)
	if err != nil {
		return err
	}
	defer localfd.Close()

	remotefd, err := sftpClient.Create(dest + filepath.Base(src))
	if err != nil {
		return err
	}
	defer remotefd.Close()

	_, err = io.Copy(remotefd, localfd)
	if err != nil {
		return err
	}

	return nil
}

func Init(arg []confssh.Node) (map[string]*Host, error) {
	var SSHhosts map[string]*Host = make(map[string]*Host)
	for _, node := range arg {
		if node.Nodename == "" {
			return nil, errors.New("Invalid node name\n")
		}
		_, exists := SSHhosts[node.Nodename]
		if exists {
			return nil, errors.New("Nodename is already exists\n")
		}

		addr := node.IP + ":22"
		user := node.User
		if user != "root" {
			fmt.Println("not a root user")
		}
		auth := []ssh.AuthMethod{}
		ctx, exists := node.AuthType[Auth_Pubkey]
		if exists {
			if ctx == "" {
				return nil, errors.New(ErrorPubkeyPath)
			}

			fd, err := os.Open(ctx)
			if err != nil {
				return nil, errors.New(ErrorPubkeyPath)
			}
			defer fd.Close()

			key, err := io.ReadAll(fd)
			if err != nil {
				return nil, err
			}

			signer, err := ssh.ParsePrivateKey(key)
			if err != nil {
				return nil, err
			}

			auth = append(auth, ssh.PublicKeys(signer))
		}

		ctx, exists = node.AuthType[Auth_Passkey]
		if exists {
			auth = append(auth, ssh.Password(ctx))
		}
		config := &ssh.ClientConfig{
			User:            node.User,
			Auth:            auth,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		new_host := new(Host)
		new_host.NodeName = node.Nodename
		new_host.IP = addr
		new_host.User = user
		new_host.SSHconf = config
	}

	return SSHhosts, nil
}
