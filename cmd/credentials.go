package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	smtpsdk "github.com/kirimemail/kirimemail-smtp-go-sdk"
)

func CredentialsCmd(client *smtpsdk.Client, defaultDomain string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("credentials command requires subcommand: list, create, get, delete, reset-password")
	}

	subCmd := args[0]
	switch subCmd {
	case "list":
		return credentialsList(client, args[1:])
	case "create":
		return credentialsCreate(client, args[1:])
	case "get":
		return credentialsGet(client, args[1:])
	case "delete":
		return credentialsDelete(client, args[1:])
	case "reset-password":
		return credentialsResetPassword(client, args[1:])
	default:
		return fmt.Errorf("unknown credentials subcommand: %s", subCmd)
	}
}

func credentialsList(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx credentials list <domain>")
	}
	domain := args[0]

	result, err := client.Credentials().List(domain, nil, nil)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tDOMAIN")
	for _, c := range result.Data {
		fmt.Fprintf(w, "%s\t%s\n", c.Username, result.Domain)
	}
	return w.Flush()
}

func credentialsCreate(client *smtpsdk.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: ktx credentials create <domain> <username>")
	}
	domain := args[0]
	username := args[1]

	result, err := client.Credentials().Create(domain, smtpsdk.CredentialCreateRequest{
		Username: username,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Credential created:\n")
	fmt.Printf("  Username: %s\n", result.Data.Credential.Username)
	fmt.Printf("  Password: %s\n", result.Data.Password)
	return nil
}

func credentialsGet(client *smtpsdk.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: ktx credentials get <domain> <guid>")
	}
	domain := args[0]
	guid := args[1]

	result, err := client.Credentials().Get(domain, guid)
	if err != nil {
		return err
	}

	fmt.Printf("Username: %s\n", result.Username)
	return nil
}

func credentialsDelete(client *smtpsdk.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: ktx credentials delete <domain> <guid>")
	}
	domain := args[0]
	guid := args[1]

	_, err := client.Credentials().Delete(domain, guid)
	if err != nil {
		return err
	}

	fmt.Printf("Credential deleted: %s\n", guid)
	return nil
}

func credentialsResetPassword(client *smtpsdk.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: ktx credentials reset-password <domain> <guid>")
	}
	domain := args[0]
	guid := args[1]

	result, err := client.Credentials().ResetPassword(domain, guid)
	if err != nil {
		return err
	}

	fmt.Printf("New password: %s\n", result.Data.NewPassword)
	return nil
}
