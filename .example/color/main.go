package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	for i := 0; i < 256; i++ {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("%d", i))).Render(fmt.Sprintf("hello world %d", i)))
	}
}
