package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/user"

	"ktx/cmd"

	smtpsdk "github.com/kirimemail/kirimemail-smtp-go-sdk"
)

const defaultBaseURL = "https://smtp-app.kirim.email"

type Config struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	BaseURL  string `json:"base_url,omitempty"`
	Domain   string `json:"domain,omitempty"`
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 2 {
		return printUsage()
	}

	username := flag.String("username", "", "KirimEmail username")
	token := flag.String("token", "", "KirimEmail API token")
	baseURL := flag.String("baseurl", "", "API base URL")
	domain := flag.String("domain", "", "Default domain to use")
	flag.Parse()

	cfg, err := loadConfig()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if *username == "" {
		*username = cfg.Username
	}
	if *username == "" {
		*username = os.Getenv("KIRIM_USERNAME")
	}
	if *token == "" {
		*token = cfg.Token
	}
	if *token == "" {
		*token = os.Getenv("KIRIM_TOKEN")
	}
	if *baseURL == "" {
		*baseURL = cfg.BaseURL
	}
	if *baseURL == "" {
		*baseURL = defaultBaseURL
	}
	if *domain == "" {
		*domain = cfg.Domain
	}
	if *domain == "" {
		*domain = os.Getenv("KIRIM_DOMAIN")
	}

	if *username == "" || *token == "" {
		return fmt.Errorf("username and token are required")
	}

	client := smtpsdk.NewClient(*username, *token, smtpsdk.WithBaseURL(*baseURL))

	cmdName := args[1]
	subArgs := args[2:]

	switch cmdName {
	case "domains":
		return cmd.DomainsCmd(client, *domain, subArgs)
	case "credentials":
		return cmd.CredentialsCmd(client, *domain, subArgs)
	case "send":
		return cmd.SendCmd(client, *domain, subArgs)
	case "validate":
		return cmd.ValidateCmd(client, subArgs)
	case "logs":
		return cmd.LogsCmd(client, *domain, subArgs)
	case "suppressions":
		return cmd.SuppressionsCmd(client, *domain, subArgs)
	case "webhooks":
		return cmd.WebhooksCmd(client, *domain, subArgs)
	case "quota":
		return cmd.QuotaCmd(client)
	case "config":
		return configCmd(*username, *baseURL, *domain)
	default:
		return printUsage()
	}
}

func loadConfig() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(usr.HomeDir + "/.ktxrc")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func configCmd(username, baseURL, domain string) error {
	cfg := Config{
		Username: username,
		Token:    "",
		BaseURL:  baseURL,
		Domain:   domain,
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("Config saved to ~/.ktxrc:\n%s\n", string(data))
	return nil
}

func printUsage() error {
	fmt.Println(`ktx - KirimEmail SMTP CLI client

Usage:
  ktx [command] [flags]

Commands:
  domains      Manage domains (list, create, get, delete, verify)
  credentials  Manage SMTP credentials (list, create, get, delete, reset-password)
  send         Send an email
  validate     Validate email addresses (email, batch)
  logs         Retrieve email logs
  suppressions Manage suppressions (list, create-whitelist, delete)
  webhooks     Manage webhooks (list, create, get, delete, test)
  quota        Get user quota information

Flags:
  -username string   KirimEmail username (or KIRIM_USERNAME env)
  -token string     KirimEmail API token (or KIRIM_TOKEN env)
  -baseurl string   API base URL (default "https://smtp-app.kirim.email")
  -domain string    Default domain to use (or KIRIM_DOMAIN env)

Config File (~/.ktxrc):
  {"username":"user","token":"token","base_url":"https://smtp-app.kirim.email","domain":"example.com"}

Examples:
  ktx domains list -username user -token token
  ktx send -domain example.com -from noreply@example.com -to user@example.com -subject "Hello" -text "Message"
  ktx validate email user@example.com
  ktx quota`)
	return nil
}
