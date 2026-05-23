// Public domain. No rights reserved.
// Original developer: Luis Forra <luis.forra@gmail.com>

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	fmt.Fprintf(os.Stderr, "Public domain. Original developer: Luis Forra <luis.forra@gmail.com>. Built with deepseek-v4-flash-free (opencode.ai)\n")

	configPath := filepath.Join(userHomeDir(), ".ssh", "config")

	hosts, err := parseSSHConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: reading SSH config: %v\n", err)
		os.Exit(1)
	}
	if len(hosts) == 0 {
		fmt.Fprintln(os.Stderr, "No SSH hosts found.")
		os.Exit(1)
	}

	if len(os.Args) == 2 && os.Args[1] != "--help" && os.Args[1] != "-h" {
		launchSSH(os.Args[1])
		return
	}

	if len(os.Args) > 1 {
		fmt.Printf("Usage: sshl [host]\n\n")
		fmt.Println("Available hosts:")
		for _, h := range hosts {
			extra := ""
			if h.HostName != "" && h.HostName != h.Alias {
				extra = " (" + h.HostName + ")"
			}
			fmt.Printf("  %s%s\n", h.Alias, extra)
		}
		return
	}

	selected := selectHost(hosts)
	if selected == "" {
		os.Exit(0)
	}
	launchSSH(selected)
}

func userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

func launchSSH(host string) {
	path, err := exec.LookPath("ssh")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ssh.exe not found. Install OpenSSH Client (Windows Settings > Apps > Optional Features > OpenSSH Client) and try again.\n")
		os.Exit(1)
	}

	cmd := exec.Command(path, host)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}
