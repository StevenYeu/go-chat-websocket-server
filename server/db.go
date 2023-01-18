package server

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Dial(ctx context.Context, addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}

type DBClient struct {
	sshClient *ssh.Client
}

type SSHInfo struct {
	Host     string
	Username string
	Port     int32
}

type DBInfo struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

func NewDBClient(sshInfo SSHInfo, dbInfo DBInfo) *DBClient {
	var agentClient agent.Agent

	// Establish a connection to the local ssh-agent
	if conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		defer conn.Close()

		// Create a new instance of the ssh agent
		agentClient = agent.NewClient(conn)
	}

	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User: sshInfo.Username,
		Auth: []ssh.AuthMethod{},
	}
	// When the agentClient connection succeeded, add them as AuthMethod
	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}

	if sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshInfo.Host, sshInfo.Host), sshConfig); err == nil {
		//defer sshcon.Close()
		// Now we register the ViaSSHDialer with the ssh connection as a parameter
		mysql.RegisterDialContext("mysql+tcp", (&ViaSSHDialer{sshcon}).Dial)

		// And now we can use our new driver with the regular mysql connection string tunneled through the SSH connection
		if db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@mysql+tcp(%s)/%s", dbInfo.Username, dbInfo.Password, dbInfo.Host, dbInfo.Name)); err == nil {

			fmt.Printf("Successfully connected to the db\n")

			//db.Close()

		} else {

			fmt.Printf("Failed to connect to the db: %s\n", err.Error())
		}

	}
}
