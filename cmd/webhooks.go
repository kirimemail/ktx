package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	smtpsdk "github.com/kirimemail/kirimemail-smtp-go-sdk"
)

func WebhooksCmd(client *smtpsdk.Client, defaultDomain string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("webhooks command requires subcommand: list, create, get, delete, test")
	}

	subCmd := args[0]
	switch subCmd {
	case "list":
		return webhooksList(client, args[1:])
	case "create":
		return webhooksCreate(client, args[1:])
	case "get":
		return webhooksGet(client, args[1:])
	case "delete":
		return webhooksDelete(client, args[1:])
	case "test":
		return webhooksTest(client, args[1:])
	default:
		return fmt.Errorf("unknown webhooks subcommand: %s", subCmd)
	}
}

func webhooksList(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx webhooks list <domain>")
	}
	domain := args[0]

	result, err := client.Webhooks().List(domain, nil)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "GUID\tTYPE\tURL")
	for _, webhook := range result.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\n", webhook.WebhookGUID, webhook.Type, webhook.URL)
	}
	return w.Flush()
}

func webhooksCreate(client *smtpsdk.Client, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: ktx webhooks create <domain> <type> <url>")
	}
	domain := args[0]
	webhookType := args[1]
	url := args[2]

	result, err := client.Webhooks().Create(domain, smtpsdk.WebhookCreateRequest{
		Type: webhookType,
		URL:  url,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Webhook created: %s\n", result.Data.WebhookGUID)
	return nil
}

func webhooksGet(client *smtpsdk.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: ktx webhooks get <domain> <guid>")
	}
	domain := args[0]
	guid := args[1]

	result, err := client.Webhooks().Get(domain, guid)
	if err != nil {
		return err
	}

	fmt.Printf("GUID: %s\nType: %s\nURL: %s\n", result.WebhookGUID, result.Type, result.URL)
	return nil
}

func webhooksDelete(client *smtpsdk.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: ktx webhooks delete <domain> <guid>")
	}
	domain := args[0]
	guid := args[1]

	_, err := client.Webhooks().Delete(domain, guid)
	if err != nil {
		return err
	}

	fmt.Printf("Webhook deleted: %s\n", guid)
	return nil
}

func webhooksTest(client *smtpsdk.Client, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: ktx webhooks test <domain> <url> <event_type>")
	}
	domain := args[0]
	url := args[1]
	eventType := args[2]

	result, err := client.Webhooks().Test(domain, smtpsdk.WebhookTestRequest{
		URL:       url,
		EventType: eventType,
	})
	if err != nil {
		return err
	}

	if result.Data.ResponseStatus == 200 {
		fmt.Printf("Test successful! Response time: %dms\n", result.Data.ResponseTime)
	} else {
		fmt.Printf("Test failed with status: %d\n", result.Data.ResponseStatus)
	}
	return nil
}
