package wizard

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
)

type AuthCredentials struct {
	Username string
	Password string
}

func RunAuthWizard() (*AuthCredentials, error) {
	var username, password string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Username").
				Value(&username).
				Validate(func(s string) error {
					if len(s) < 1 {
						return errors.New("username cannot be empty")
					}
					if len(s) < 3 {
						return errors.New("username must be at least 3 characters")
					}
					return nil
				}),

			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&password).
				Validate(func(s string) error {
					if len(s) < 1 {
						return errors.New("password cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("wizard cancelled: %w", err)
	}

	return &AuthCredentials{
		Username: username,
		Password: password,
	}, nil
}
