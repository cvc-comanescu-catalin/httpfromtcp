package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.OpenFile("messages.txt", os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	linesChan := getLinesChannel(f)
	for line := range linesChan {
		fmt.Println("read:", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)

		buffer := make([]byte, 8, 8)
		text := ""

		for {
			n, err := f.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					lines <- text
				} else {
					log.Fatal(err)
				}
				break
			}

			text += string(buffer[:n])
			textParts := strings.FieldsFunc(text, func(c rune) bool { return c == '\n' || c == '\r' })
			if len(textParts) > 1 {
				for _, part := range textParts[:len(textParts)-1] {
					lines <- part
				}
				text = textParts[len(textParts)-1]
			}
		}
	}()
	return lines
}
