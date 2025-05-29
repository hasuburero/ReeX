package ssh

import (
	"errors"
	"fmt"
	"github.com/hasuburero/ReeX/lib/controller/config/confssh"
	"io"
	"os"
)

import (
	"golang.org/x/crypto/ssh"
)

// section 1
type SSH struct {
	Hosts []Host
}

type Host struct {
	NodeName string
	IP       string
	User     string
	SSHconf  *ssh.ClientConfig
}

// section 2
const (
	Auth_Pubkey  = "pubkey"
	Auth_Passkey = "passkey"

	ErrorPubkeyPath = "Invalid Pubkey path\n"
)

// section3
func Init(arg []confssh.Node) (map[string]*SSH, error) {
	var SSHhosts map[string]*Host = make(map[string]*Host)
	for _, node := range arg {
		if node.Nodename == "" {
			return nil, errors.New("Invalid node name\n")
		}
		_, exists := SSHhosts[node.Nodename]
		if exists {
			return nil, errors.New("Nodename is already exists\n")
		}

		addr := node.IP
		user := node.User
		if user != "root" {
			fmt.Println("not a root user")
		}
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
				return SSH{}, err
			}

			signer, err := ssh.ParsePrivateKey(key)
			if err != nil {
				return SSH{}, err
			}

			config := &ssh.ClientConfig{
				User: node.User,
				Auth: []ssh.AuthMethod{
					ssh.PublicKeys(signer),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}
			new_host := new(Host)
			new_host.NodeName = node.Nodename
			new_host.IP = addr
			new_host.User = user
		}
	}

	return SSHhosts, nil
}
