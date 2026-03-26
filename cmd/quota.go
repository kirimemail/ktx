package cmd

import (
	"fmt"

	smtpsdk "github.com/kirimemail/kirimemail-smtp-go-sdk"
)

func QuotaCmd(client *smtpsdk.Client) error {
	result, err := client.User().GetQuota()
	if err != nil {
		return err
	}

	fmt.Printf("Current quota: %d\n", result.Data.CurrentQuota)
	fmt.Printf("Max quota: %d\n", result.Data.MaxQuota)
	fmt.Printf("Usage: %.1f%%\n", result.Data.UsagePercentage)
	return nil
}
