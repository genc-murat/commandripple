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
	config, err := getSSHConfig(user)
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

func getSSHConfig(user string) (*ssh.ClientConfig, error) {
	authMethods, err := getAuthMethods()
	if err != nil {
		return nil, fmt.Errorf("failed to get authentication methods: %v", err)
	}

	return &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: This is insecure for production use
	}, nil
}

func getAuthMethods() ([]ssh.AuthMethod, error) {
	var authMethods []ssh.AuthMethod

	// Try SSH key authentication first
	if sshKeyAuth, err := getSSHKeyAuth(); err == nil {
		authMethods = append(authMethods, sshKeyAuth)
	} else {
		fmt.Printf("SSH key authentication failed: %v\n", err)
	}

	// Try SSH Agent if key auth failed
	if len(authMethods) == 0 {
		if sshAgentAuth, err := getSSHAgentAuth(); err == nil {
			authMethods = append(authMethods, sshAgentAuth)
		} else {
			fmt.Printf("SSH agent authentication failed: %v\n", err)
		}
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no SSH authentication methods available")
	}

	return authMethods, nil
}

func getSSHAgentAuth() (ssh.AuthMethod, error) {
	if runtime.GOOS == "windows" {
		return nil, fmt.Errorf("SSH Agent auth not implemented for Windows")
	}

	socket := os.Getenv("SSH_AUTH_SOCK")
	if socket == "" {
		return nil, fmt.Errorf("SSH_AUTH_SOCK not set")
	}

	conn, err := net.Dial("unix", socket)
	if err != nil {
		return nil, fmt.Errorf("failed to open SSH_AUTH_SOCK: %v", err)
	}

	agentClient := agent.NewClient(conn)
	return ssh.PublicKeysCallback(agentClient.Signers), nil
}

func getSSHKeyAuth() (ssh.AuthMethod, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %v", err)
	}

	keyFile := filepath.Join(home, ".ssh", "id_rsa")
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key from %s: %v", keyFile, err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	return ssh.PublicKeys(signer), nil
}
