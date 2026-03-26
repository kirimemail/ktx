package cmd

import (
	"fmt"
	"strings"

	smtpsdk "github.com/kirimemail/kirimemail-smtp-go-sdk"
)

func SendCmd(client *smtpsdk.Client, defaultDomain string, args []string) error {
	var from, to, subject, text, html string
	var toList []string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-from":
			if i+1 < len(args) {
				from = args[i+1]
				i++
			}
		case "-to":
			if i+1 < len(args) {
				to = args[i+1]
				toList = strings.Split(to, ",")
				i++
			}
		case "-subject":
			if i+1 < len(args) {
				subject = args[i+1]
				i++
			}
		case "-text":
			if i+1 < len(args) {
				text = args[i+1]
				i++
			}
		case "-html":
			if i+1 < len(args) {
				html = args[i+1]
				i++
			}
		case "-domain":
			if i+1 < len(args) {
				defaultDomain = args[i+1]
				i++
			}
		}
	}

	if from == "" || len(toList) == 0 || subject == "" {
		return fmt.Errorf("usage: ktx send -domain <domain> -from <from> -to <to> -subject <subject> [-text <text>] [-html <html>]")
	}

	if defaultDomain == "" {
		return fmt.Errorf("domain is required (use -domain flag or KIRIM_DOMAIN env)")
	}

	result, err := client.Messages().Send(defaultDomain, smtpsdk.MessageSendRequest{
		From:    from,
		To:      toList,
		Subject: subject,
		Text:    text,
		HTML:    html,
	}, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Message sent: %s\n", result.Message)
	return nil
}
