package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	smtpsdk "github.com/kirimemail/kirimemail-smtp-go-sdk"
)

func DomainsCmd(client *smtpsdk.Client, defaultDomain string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("domains command requires subcommand: list, create, get, delete, verify, setup-auth-domain, verify-auth-domain, delete-auth-domain, setup-tracklink, verify-tracklink, delete-tracklink")
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
	case "setup-auth-domain":
		return domainsSetupAuthDomain(client, args[1:])
	case "verify-auth-domain":
		return domainsVerifyAuthDomain(client, args[1:])
	case "delete-auth-domain":
		return domainsDeleteAuthDomain(client, args[1:])
	case "setup-tracklink":
		return domainsSetupTracklink(client, args[1:])
	case "verify-tracklink":
		return domainsVerifyTracklink(client, args[1:])
	case "delete-tracklink":
		return domainsDeleteTracklink(client, args[1:])
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

	err := client.Domains().Delete(domain)
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

func domainsSetupAuthDomain(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx domains setup-auth-domain <domain> [auth_domain]")
	}
	domain := args[0]
	authDomain := domain
	if len(args) > 1 {
		authDomain = args[1]
	}

	result, err := client.Domains().SetupAuthDomain(domain, smtpsdk.AuthDomainSetupRequest{
		AuthDomain:    authDomain,
		DKIMKeyLength: 2048,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Auth domain setup: %s\n", result.Data.AuthDomain)
	return nil
}

func domainsVerifyAuthDomain(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx domains verify-auth-domain <domain>")
	}
	domain := args[0]

	result, err := client.Domains().VerifyAuthDomain(domain)
	if err != nil {
		return err
	}

	fmt.Printf("Auth DKIM: %v\nAuth SPF: %v\nAuth MX: %v\n",
		result.Records.AuthDKIM, result.Records.AuthSPF, result.Records.AuthMX)
	return nil
}

func domainsDeleteAuthDomain(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx domains delete-auth-domain <domain>")
	}
	domain := args[0]

	err := client.Domains().DeleteAuthDomain(domain)
	if err != nil {
		return err
	}

	fmt.Printf("Auth domain deleted: %s\n", domain)
	return nil
}

func domainsSetupTracklink(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx domains setup-tracklink <domain> [tracking_domain]")
	}
	domain := args[0]
	trackingDomain := domain
	if len(args) > 1 {
		trackingDomain = args[1]
	}

	result, err := client.Domains().SetupTracklink(domain, smtpsdk.TracklinkSetupRequest{
		TrackingDomain: trackingDomain,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Tracklink setup: %s\n", result.Data.TrackingDomain)
	return nil
}

func domainsVerifyTracklink(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx domains verify-tracklink <domain>")
	}
	domain := args[0]

	result, err := client.Domains().VerifyTracklink(domain)
	if err != nil {
		return err
	}

	fmt.Printf("CNAME: %v\nTracking Domain: %v\n",
		result.Records.CNAME, result.Records.TrackingDomain)
	return nil
}

func domainsDeleteTracklink(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx domains delete-tracklink <domain>")
	}
	domain := args[0]

	err := client.Domains().DeleteTracklink(domain)
	if err != nil {
		return err
	}

	fmt.Printf("Tracklink deleted: %s\n", domain)
	return nil
}
