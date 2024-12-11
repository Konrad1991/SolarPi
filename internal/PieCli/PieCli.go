package main

import (
	"flag"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Underline(true)

	nameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	ageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Italic(true)

	name := flag.String("name", "world", "a name to greet")
	age := flag.Int("age", 0, "your age")

	flag.Parse()

	fmt.Println(titleStyle.Render("👋 Welcome to the Styled CLI!"))
	fmt.Printf("Hello, %s! You are %s years old.\n",
		nameStyle.Render(*name),
		ageStyle.Render(fmt.Sprintf("%d", *age)),
	)
}
