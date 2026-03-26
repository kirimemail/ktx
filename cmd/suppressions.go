package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	smtpsdk "github.com/kirimemail/kirimemail-smtp-go-sdk"
)

func SuppressionsCmd(client *smtpsdk.Client, defaultDomain string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("suppressions command requires subcommand: list, create-whitelist, delete")
	}

	subCmd := args[0]
	switch subCmd {
	case "list":
		return suppressionsList(client, args[1:])
	case "create-whitelist":
		return suppressionsCreateWhitelist(client, args[1:])
	case "delete":
		return suppressionsDelete(client, args[1:])
	default:
		return fmt.Errorf("unknown suppressions subcommand: %s", subCmd)
	}
}

func suppressionsList(client *smtpsdk.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ktx suppressions list <domain> [type]")
	}
	domain := args[0]

	var result *smtpsdk.SuppressionListResponse
	var err error

	if len(args) > 1 {
		supType := args[1]
		switch supType {
		case "unsubscribe":
			result, err = client.Suppressions().ListUnsubscribes(domain, nil, nil, nil)
		case "bounce":
			result, err = client.Suppressions().ListBounces(domain, nil, nil, nil)
		case "whitelist":
			result, err = client.Suppressions().ListWhitelists(domain, nil, nil, nil)
		default:
			result, err = client.Suppressions().List(domain, smtpsdk.StringPtr(supType), nil, nil, nil)
		}
	} else {
		result, err = client.Suppressions().List(domain, nil, nil, nil, nil)
	}

	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "RECIPIENT\tTYPE\tDESCRIPTION")
	for _, s := range result.Data {
		desc := ""
		if s.Description != nil {
			desc = *s.Description
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", s.Recipient, s.Type, desc)
	}
	return w.Flush()
}

func suppressionsCreateWhitelist(client *smtpsdk.Client, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: ktx suppressions create-whitelist <domain> <recipient> <type> [description]")
	}
	domain := args[0]
	recipient := args[1]
	recipientType := args[2]
	description := ""
	if len(args) > 3 {
		description = args[3]
	}

	result, err := client.Suppressions().CreateWhitelist(domain, smtpsdk.WhitelistCreateRequest{
		Recipient:     recipient,
		RecipientType: recipientType,
		Description:   description,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Whitelist entry created: %s\n", result.Data.Recipient)
	return nil
}

func suppressionsDelete(client *smtpsdk.Client, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: ktx suppressions delete <domain> <type> <id1,id2,...>")
	}
	domain := args[0]
	supType := args[1]
	idStr := args[2]

	idParts := strings.Split(idStr, ",")
	ids := make([]int, 0, len(idParts))
	for _, p := range idParts {
		if id, err := strconv.Atoi(p); err == nil {
			ids = append(ids, id)
		}
	}

	var result *smtpsdk.SuppressionDeleteResponse
	var err error

	switch supType {
	case "unsubscribe":
		result, err = client.Suppressions().DeleteUnsubscribes(domain, ids)
	case "bounce":
		result, err = client.Suppressions().DeleteBounces(domain, ids)
	case "whitelist":
		result, err = client.Suppressions().DeleteWhitelists(domain, ids)
	default:
		return fmt.Errorf("unknown suppression type: %s", supType)
	}

	if err != nil {
		return err
	}

	fmt.Printf("Deleted %d suppressions\n", result.DeletedCount)
	return nil
}
