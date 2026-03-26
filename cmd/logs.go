package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	smtpsdk "github.com/kirimemail/kirimemail-smtp-go-sdk"
)

func LogsCmd(client *smtpsdk.Client, defaultDomain string, args []string) error {
	if defaultDomain == "" {
		return fmt.Errorf("domain is required (use -domain flag or KIRIM_DOMAIN env)")
	}

	var startDate, endDate, sender, recipient *string
	limit := 100
	offset := 0

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-start":
			if i+1 < len(args) {
				startDate = smtpsdk.StringPtr(args[i+1])
				i++
			}
		case "-end":
			if i+1 < len(args) {
				endDate = smtpsdk.StringPtr(args[i+1])
				i++
			}
		case "-sender":
			if i+1 < len(args) {
				sender = smtpsdk.StringPtr(args[i+1])
				i++
			}
		case "-recipient":
			if i+1 < len(args) {
				recipient = smtpsdk.StringPtr(args[i+1])
				i++
			}
		}
	}

	result, err := client.Logs().Get(defaultDomain, startDate, endDate, sender, recipient, &limit, &offset)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tEVENT\tMESSAGE GUID")
	for _, log := range result.Data {
		t := time.Unix(log.Timestamp, 0)
		fmt.Fprintf(w, "%s\t%s\t%s\n", t.Format(time.RFC3339), log.EventType, log.MessageGUID)
	}
	return w.Flush()
}
