package commands

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// RemoteExecute executes a command on a remote machine via SSH
func RemoteExecute(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("usage: remote_execute [user] [host] [port] [command]")
	}
	user := args[0]
	host := args[1]
	port := args[2]
	command := strings.Join(args[3:], " ") // Allow multi-word commands

	// Setup SSH client configuration
	config, err := getSSHForRemoteConfig(user)
	if err != nil {
		return fmt.Errorf("failed to configure SSH client: %v", err)
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)
	if err != nil {
		return fmt.Errorf("unable to connect: %v", err)
	}
	defer client.Close()

	// Create a new SSH session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("unable to create session: %v", err)
	}
	defer session.Close()

	// Execute the command
	output, err := session.CombinedOutput(command)
	if err != nil {
		return fmt.Errorf("command execution failed: %v\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}

func getSSHForRemoteConfig(user string) (*ssh.ClientConfig, error) {
	authMethods, err := getAuthMethods()
	if err != nil {
		return nil, err
	}

	return &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: This is insecure for production use
	}, nil
}

func getAuthMethods() ([]ssh.AuthMethod, error) {
	var authMethods []ssh.AuthMethod

	// Try SSH Agent first
	if sshAgentAuth, err := getSSHAgentAuth(); err == nil {
		authMethods = append(authMethods, sshAgentAuth)
	}

	// Then try SSH key authentication
	if sshKeyAuth, err := getSSHKeyAuth(); err == nil {
		authMethods = append(authMethods, sshKeyAuth)
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no SSH authentication methods available")
	}

	return authMethods, nil
}

func getSSHAgentAuth() (ssh.AuthMethod, error) {
	if sshAgentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgentConn).Signers), nil
	}
	return nil, fmt.Errorf("failed to connect to SSH agent")
}

func getSSHKeyAuth() (ssh.AuthMethod, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %v", err)
	}

	keyFile := filepath.Join(home, ".ssh", "id_rsa")
	if runtime.GOOS == "windows" {
		keyFile = filepath.Join(home, ".ssh", "id_rsa")
	}

	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	return ssh.PublicKeys(signer), nil
}
