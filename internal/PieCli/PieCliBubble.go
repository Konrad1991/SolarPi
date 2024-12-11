/*
I. create user
 1. add user table in database
    id, username,
    password_hash, public_key (public ssh key),
    created_at, updated_at,
    user_root_directory (hashed?)
 2. add root create user
 3. add root update user
 4. add root delete user
 5. list all users

II. pie auth --> first as pure cli using flags (but also use bubbletea)
 1. Ask for name; Check valid and unique name (list_all_users)
 2. Choose between password and ssh
 3. enter password or public ssh key
 4. save password/key using charm
 5. run add_user

III. same as II. but interactie using bubbletea

... Pie down --> directly path
Pie down --> interactive explorer: cd through backup
*/
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	name         string
	age          int
	step         int
	input        string
	errorMessage string
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Underline(true)

	nameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	ageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Italic(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

const (
	stepName = 0
	stepAge  = 1
	stepDone = 2
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		if key == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.step {
		case stepName:
			if key == "enter" {
				if m.input == "" {
					m.errorMessage = "Name cannot be empty. Please enter your name."
					return m, nil
				}
				m.name = m.input
				m.input = ""
				m.step = stepAge
				m.errorMessage = ""
				return m, nil
			}
		case stepAge:
			if key == "enter" {
				var age int
				_, err := fmt.Sscanf(m.input, "%d", &age)
				if err != nil || age <= 0 {
					m.errorMessage = "Invalid age. Please enter a valid positive number."
					return m, nil
				}
				m.age = age
				m.input = ""
				m.step = stepDone
				return m, nil
			}
		}

		// Update the input
		if key == "backspace" {
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		} else {
			m.input += key
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.step {
	case stepName:
		return fmt.Sprintf(
			"%s\n\nEnter your name: %s\n\n%s",
			titleStyle.Render("👋 Welcome to the Styled CLI!"),
			nameStyle.Render(m.input),
			helpStyle.Render(m.errorMessage),
		)
	case stepAge:
		return fmt.Sprintf(
			"%s\n\nHello, %s! Now enter your age: %s\n\n%s",
			titleStyle.Render("Step 2: Age Input"),
			nameStyle.Render(m.name),
			ageStyle.Render(m.input),
			helpStyle.Render(m.errorMessage),
		)
	case stepDone:
		return fmt.Sprintf(
			"%s\n\nHello, %s! You are %s years old.\n\nPress Ctrl+C to exit.",
			titleStyle.Render("🎉 All Done!"),
			nameStyle.Render(m.name),
			ageStyle.Render(fmt.Sprintf("%d", m.age)),
		)
	}
	return ""
}

func main() {
	switch os.Args[1] {
	case "auth":
		p := tea.NewProgram(model{})
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting application: %v\n", err)
			os.Exit(1)
		}
	case "-man":
		fmt.Println("Manual authentication")
	case "--help":
		fmt.Println("Manual authentication")
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}
