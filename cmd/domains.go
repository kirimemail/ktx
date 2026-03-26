package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	smtpsdk "github.com/kirimemail/kirimemail-smtp-go-sdk"
)

func DomainsCmd(client *smtpsdk.Client, defaultDomain string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("domains command requires subcommand: list, create, get, delete, verify")
	}

	subCmd := args[0]
	switch subCmd {
	case "list":
		return domainsList(client)
	case "create":
		return domainsCreate(client, args[1:])
	case "get":
		return domainsGet(client, args[1:])
	case "delete":
		return domainsDelete(client, args[1:])
	case "verify":
		return domainsVerify(client, args[1:])
	default:
		return fmt.Errorf("unknown domains subcommand: %s", subCmd)
	}
}

func domainsList(client *smtpsdk.Client) error {
	domains, err := client.Domains().List(nil, nil, nil)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DOMAIN\tSTATUS\tVERIFIED")
	for _, d := range domains.Data {
		fmt.Fprintf(w, "%s\t%v\t%v\n", d.Domain, d.Status, d.IsVerified)
	}
	return w.Flush()
}

func domainsCreate(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx domains create <domain>")
	}
	domain := args[0]

	result, err := client.Domains().Create(smtpsdk.DomainCreateRequest{
		Domain:        domain,
		DKIMKeyLength: 2048,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Domain created: %s\n", result.Data.Domain)
	return nil
}

func domainsGet(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx domains get <domain>")
	}
	domain := args[0]

	result, err := client.Domains().Get(domain)
	if err != nil {
		return err
	}

	fmt.Printf("Domain: %s\nStatus: %v\nVerified: %v\n", result.Domain, result.Status, result.IsVerified)
	return nil
}

func domainsDelete(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx domains delete <domain>")
	}
	domain := args[0]

	_, err := client.Domains().Delete(domain)
	if err != nil {
		return err
	}

	fmt.Printf("Domain deleted: %s\n", domain)
	return nil
}

func domainsVerify(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx domains verify <domain>")
	}
	domain := args[0]

	result, err := client.Domains().VerifyMandatoryRecords(domain)
	if err != nil {
		return err
	}

	fmt.Printf("DKIM: %v\nSPF: %v\nMX: %v\n",
		result.Records.DKIM, result.Records.SPF, result.Records.MX)
	return nil
}
