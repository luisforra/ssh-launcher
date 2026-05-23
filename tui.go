// Public domain. No rights reserved.
// Original developer: Luis Forra <luis.forra@gmail.com>

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

type listState struct {
	names     []string
	filter    string
	filtered  []int
	cursor    int
	lastLines int
}

func selectHost(hosts []HostEntry) string {
	if len(hosts) == 0 {
		return ""
	}

	names := make([]string, len(hosts))
	for i, h := range hosts {
		if h.User != "" {
			names[i] = h.User + "@" + h.Alias
		} else {
			names[i] = h.Alias
		}
	}

	filtered := make([]int, len(hosts))
	for i := range filtered {
		filtered[i] = i
	}

	state := &listState{
		names:    names,
		filtered: filtered,
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return selectHostSimple(hosts, names)
	}
	defer term.Restore(fd, oldState)

	fmt.Fprint(os.Stderr, "\033[?25l")
	defer fmt.Fprint(os.Stderr, "\033[?25h")

	renderList(state)

	for {
		key, err := readKey()
		if err != nil {
			break
		}

		switch {
		case key == "enter":
			if len(state.filtered) > 0 {
				clearList(state)
				return hosts[state.filtered[state.cursor]].Alias
			}
		case key == "esc" || key == "ctrl+c":
			clearList(state)
			return ""
		case key == "backspace":
			if len(state.filter) > 0 {
				state.filter = state.filter[:len(state.filter)-1]
				applyFilter(state)
			}
		case key == "up":
			if state.cursor > 0 {
				state.cursor--
			}
		case key == "down":
			if state.cursor < len(state.filtered)-1 {
				state.cursor++
			}
		case len(key) == 1 && key[0] >= ' ' && key[0] <= '~':
			state.filter += key
			applyFilter(state)
		}

		renderList(state)
	}

	clearList(state)
	return ""
}

func applyFilter(state *listState) {
	filter := strings.ToLower(state.filter)
	filtered := make([]int, 0, len(state.names))
	for i, name := range state.names {
		if filter == "" || strings.Contains(strings.ToLower(name), filter) {
			filtered = append(filtered, i)
		}
	}
	state.filtered = filtered
	if state.cursor >= len(state.filtered) {
		state.cursor = max(0, len(state.filtered)-1)
	}
	if state.cursor < 0 {
		state.cursor = 0
	}
}

func renderList(state *listState) {
	_, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || h < 5 {
		h = 24
	}

	if state.lastLines > 0 {
		fmt.Fprintf(os.Stderr, "\033[%dA", state.lastLines)
		fmt.Fprint(os.Stderr, "\033[J")
	}

	total := len(state.filtered)
	visible := h - 5
	if visible < 1 {
		visible = 1
	}

	lines := 0

	fmt.Fprintf(os.Stderr, "SSH Launcher  \033[2m(%d host%s)\033[0m\r\n", total, plural(total))
	lines++

	fmt.Fprint(os.Stderr, "\r\n")
	lines++

	start := 0
	if state.cursor >= visible {
		start = state.cursor - visible + 1
	}
	end := start + visible
	if end > total {
		end = total
	}

	for i := start; i < end; i++ {
		idx := state.filtered[i]
		name := state.names[idx]
		if i == state.cursor {
			fmt.Fprintf(os.Stderr, "\033[7m \033[1m>\033[0m\033[7m %s \033[0m\r\n", name)
		} else {
			fmt.Fprintf(os.Stderr, "  %s\r\n", name)
		}
		lines++
	}

	if end < total {
		fmt.Fprintf(os.Stderr, "  \033[2m... %d more\033[0m\r\n", total-end)
		lines++
	}

	fmt.Fprint(os.Stderr, "\r\n")
	lines++

	filterLine := state.filter
	if total == 0 {
		filterLine = "\033[31m(no match)\033[0m " + state.filter
	}
	fmt.Fprintf(os.Stderr, "filter: %s\u258c\r\n", filterLine)
	lines++

	state.lastLines = lines
}

func clearList(state *listState) {
	if state.lastLines > 0 {
		fmt.Fprintf(os.Stderr, "\033[%dA", state.lastLines)
		fmt.Fprint(os.Stderr, "\033[J")
		state.lastLines = 0
	}
}

func selectHostSimple(hosts []HostEntry, names []string) string {
	fmt.Fprintln(os.Stderr, "Available hosts:")
	for i, n := range names {
		fmt.Fprintf(os.Stderr, "  %2d. %s\n", i+1, n)
	}
	fmt.Fprint(os.Stderr, "\nEnter number (or 0 to cancel): ")

	var n int
	_, err := fmt.Scanf("%d", &n)
	if err != nil || n < 1 || n > len(hosts) {
		return ""
	}
	return hosts[n-1].Alias
}

func readKey() (string, error) {
	buf := make([]byte, 8)
	n, err := os.Stdin.Read(buf[:1])
	if err != nil || n == 0 {
		return "", err
	}

	b := buf[0]
	if b != 0x1b {
		switch b {
		case 0x0d:
			return "enter", nil
		case 0x7f, 0x08:
			return "backspace", nil
		case 0x03:
			return "ctrl+c", nil
		default:
			if b >= ' ' && b <= '~' {
				return string(b), nil
			}
			return "", nil
		}
	}

	done := make(chan int, 1)
	go func() {
		n, _ = os.Stdin.Read(buf[1:])
		done <- n
	}()

	select {
	case n := <-done:
		if n >= 2 && buf[1] == '[' {
			switch buf[2] {
			case 'A':
				return "up", nil
			case 'B':
				return "down", nil
			case 'C':
				return "right", nil
			case 'D':
				return "left", nil
			}
		}
		return "esc", nil
	case <-time.After(50 * time.Millisecond):
		return "esc", nil
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
