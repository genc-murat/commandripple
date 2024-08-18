package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

// FileTransfer transfers a file to or from a remote machine using SSH
func FileTransfer(args []string) error {
	if len(args) < 5 {
		return fmt.Errorf("usage: file_transfer [user] [host] [port] [source] [destination]")
	}
	user := args[0]
	host := args[1]
	port := args[2]
	source := args[3]
	destination := args[4]

	// Determine if we're uploading or downloading
	isUpload := !strings.HasPrefix(source, fmt.Sprintf("%s@%s:", user, host))

	// Setup SSH client configuration
	config, err := getSSHConfig(user)
	if err != nil {
		return fmt.Errorf("failed to configure SSH client: %v", err)
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)
	if err != nil {
		return fmt.Errorf("failed to connect to SSH server: %v", err)
	}
	defer client.Close()

	if isUpload {
		return uploadFile(client, source, destination)
	}
	return downloadFile(client, source, destination)
}

func uploadFile(client *ssh.Client, source, destination string) error {
	// Open the source file
	srcFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer srcFile.Close()

	// Get file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %v", err)
	}

	// Create an SSH session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer session.Close()

	// Setup the remote command to receive the file
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintf(w, "C0644 %d %s\n", srcInfo.Size(), filepath.Base(destination))
		io.Copy(w, srcFile)
		fmt.Fprint(w, "\x00")
	}()

	// Run the command to receive the file
	if err := session.Run(fmt.Sprintf("scp -t %s", destination)); err != nil {
		return fmt.Errorf("failed to run remote scp command: %v", err)
	}

	fmt.Printf("File uploaded successfully to %s\n", destination)
	return nil
}

func downloadFile(client *ssh.Client, source, destination string) error {
	// Create an SSH session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer session.Close()

	// Setup the remote command to send the file
	remoteCmd := fmt.Sprintf("scp -f %s", source)
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to setup stdout pipe: %v", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to setup stdin pipe: %v", err)
	}

	if err := session.Start(remoteCmd); err != nil {
		return fmt.Errorf("failed to start remote scp command: %v", err)
	}

	// Open the destination file
	dstFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dstFile.Close()

	// Perform the file transfer
	_, err = fmt.Fprint(stdin, "\x00")
	if err != nil {
		return fmt.Errorf("failed to send null byte: %v", err)
	}

	var fileSize int64
	_, err = fmt.Fscanf(stdout, "C0644 %d %s", &fileSize, &source)
	if err != nil {
		return fmt.Errorf("failed to read file info: %v", err)
	}

	_, err = io.CopyN(dstFile, stdout, fileSize)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	fmt.Printf("File downloaded successfully to %s\n", destination)
	return nil
}
