// Public domain. No rights reserved.
// Original developer: Luis Forra <luis.forra@gmail.com>

package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type HostEntry struct {
	Alias        string
	HostName     string
	User         string
	Port         string
	IdentityFile string
	FileName     string
}

func parseSSHConfig(path string) ([]HostEntry, error) {
	visited := make(map[string]bool)
	entries, err := parseFile(path, visited)
	if err != nil {
		return nil, err
	}

	var result []HostEntry
	seen := make(map[string]bool)
	for _, e := range entries {
		if e.Alias == "" || isWildcard(e.Alias) {
			continue
		}
		key := strings.ToLower(e.Alias)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, e)
	}

	return result, nil
}

func parseFile(path string, visited map[string]bool) ([]HostEntry, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if visited[absPath] {
		return nil, nil
	}
	visited[absPath] = true

	f, err := os.Open(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	dir := filepath.Dir(absPath)

	var entries []HostEntry
	var curPatterns []string
	var curHost *HostEntry

	finalizeHost := func() {
		if curHost == nil {
			return
		}
		for _, p := range curPatterns {
			if !isWildcard(p) {
				entry := *curHost
				entry.Alias = p
				entries = append(entries, entry)
			}
		}
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		for strings.HasSuffix(line, "\\") && scanner.Scan() {
			line = strings.TrimSuffix(line, "\\") + strings.TrimSpace(scanner.Text())
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.Fields(trimmed)
		keyword := strings.ToLower(parts[0])
		rest := strings.TrimSpace(trimmed[len(parts[0]):])

		switch keyword {
		case "include":
			finalizeHost()
			files := expandIncludes(rest, dir)
			for _, fpath := range files {
				sub, err := parseFile(fpath, visited)
				if err != nil {
					continue
				}
				entries = append(entries, sub...)
			}
			curHost = nil
			curPatterns = nil

		case "host":
			finalizeHost()
			curPatterns = parts[1:]
			curHost = &HostEntry{FileName: absPath}

		case "match":
			finalizeHost()
			curHost = nil
			curPatterns = nil

		default:
			if curHost == nil {
				continue
			}
			switch keyword {
			case "hostname":
				curHost.HostName = rest
			case "user":
				curHost.User = rest
			case "port":
				curHost.Port = rest
			case "identityfile":
				curHost.IdentityFile = expandUser(rest)
			}
		}
	}

	finalizeHost()

	if err := scanner.Err(); err != nil {
		return entries, err
	}
	return entries, nil
}

func expandIncludes(pattern string, baseDir string) []string {
	var result []string
	for _, p := range strings.Fields(pattern) {
		p = expandUser(p)
		if !filepath.IsAbs(p) {
			p = filepath.Join(baseDir, p)
		}
		matches, err := filepath.Glob(p)
		if err != nil {
			continue
		}
		result = append(result, matches...)
	}
	return result
}

func expandUser(path string) string {
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[1:])
	}
	return path
}

func isWildcard(pattern string) bool {
	return strings.ContainsAny(pattern, "*?")
}
