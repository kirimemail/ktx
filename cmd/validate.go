package cmd

import (
	"fmt"
	"strings"

	smtpsdk "github.com/kirimemail/kirimemail-smtp-go-sdk"
)

func ValidateCmd(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("validate command requires subcommand: email, batch")
	}

	subCmd := args[0]
	switch subCmd {
	case "email":
		return validateEmail(client, args[1:])
	case "batch":
		return validateBatch(client, args[1:])
	default:
		return fmt.Errorf("unknown validate subcommand: %s", subCmd)
	}
}

func validateEmail(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx validate email <email>")
	}
	email := args[0]

	result, err := client.Validation().ValidateEmail(email)
	if err != nil {
		return err
	}

	fmt.Printf("Email: %s\nValid: %v\nIs spamtrap: %v\nSpamtrap score: %.2f\n",
		result.Data.Email, result.Data.IsValid, result.Data.IsSpamtrap, result.Data.SpamtrapScore)
	return nil
}

func validateBatch(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx validate batch <email1,email2,...>")
	}

	emails := strings.Split(args[0], ",")

	result, err := client.Validation().ValidateEmailsBatch(emails)
	if err != nil {
		return err
	}

	fmt.Printf("Total: %d\nValid: %d\nInvalid: %d\n",
		result.Data.Summary.Total, result.Data.Summary.Valid, result.Data.Summary.Invalid)
	return nil
}
