package io

import (
	"io"
	"bufio"
	"os"
	"os/exec"
	"strings"
	"fmt"
	"bytes"
)

type Console struct {
	reader io.Reader
	sttySettings string
}

func NewConsole() *Console {
	c := new(Console)
	c.reader = bufio.NewReader(os.Stdin)
	
	return c
}

func (c *Console) Read() {
	output, err := exec.Command("stty", "-F", "/dev/tty", "-g").CombinedOutput()
	if err != nil {
		panic(err)	
	}
	c.sttySettings = strings.TrimSpace(string(output[:]))
	defer exec.Command("stty", "-F", "/dev/tty", c.sttySettings).Run()
	exec.Command("stty", "-F", "/dev/tty", "-icanon", "-echo", "min", "1", "time", "0").Run()
	
	buff := bytes.Buffer{};
	char := make([]byte, 1)
	for {
		_, err = os.Stdin.Read(char)
		if err != nil {
			panic(err)
		}
		
		if char[0] == 'e' {
			fmt.Println("")
			fmt.Println(buff.String())
			break;
		} else if char[0] == 127 {
			if buff.Len() > 0 {
				fmt.Print("\b \b")
				buff.Truncate(buff.Len() - 1)
			}
		} else if char[0] == '\t' {
			c.handleKeyTab()
		} else if char[0] == '\x1b' {
			_, err = os.Stdin.Read(char)
			if char[0] == '[' {
				_, err = os.Stdin.Read(char)
				switch char[0] {
					case 'A': c.handleKeyUp()
					case 'B': c.handleKeyDown()
				}
			}
		} else {
			fmt.Printf("%c", char[0])
			buff.WriteByte(char[0])
		}
	}
}

func (c *Console) handleKeyUp() {
	fmt.Println("UP")
}

func (c *Console) handleKeyDown() {
	fmt.Println("DOWN")
}

func (c *Console) handleKeyTab() {
	fmt.Println("TAB")
}

