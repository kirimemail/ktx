package cmd

import (
	"testing"

	smtpsdk "github.com/kirimemail/kirimemail-smtp-go-sdk"
)

func TestDomainsCmd_Subcommands(t *testing.T) {
	client := smtpsdk.NewClient("testuser", "testtoken")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "no subcommand",
			args:    []string{},
			wantErr: true,
			errMsg:  "domains command requires subcommand",
		},
		{
			name:    "unknown subcommand",
			args:    []string{"unknown"},
			wantErr: true,
			errMsg:  "unknown domains subcommand",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DomainsCmd(client, "example.com", tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainsCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if tt.wantErr && err.Error() != tt.errMsg && !contains(err.Error(), tt.errMsg) {
					t.Errorf("DomainsCmd() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

func TestCredentialsCmd_Subcommands(t *testing.T) {
	client := smtpsdk.NewClient("testuser", "testtoken")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no subcommand",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "unknown subcommand",
			args:    []string{"unknown"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CredentialsCmd(client, "example.com", tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("CredentialsCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendCmd_Validation(t *testing.T) {
	client := smtpsdk.NewClient("testuser", "testtoken")

	tests := []struct {
		name    string
		args    []string
		domain  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "missing required flags",
			args:    []string{},
			domain:  "",
			wantErr: true,
			errMsg:  "usage: ktx send",
		},
		{
			name:    "missing domain",
			args:    []string{"-from", "a@b.com", "-to", "c@d.com", "-subject", "test"},
			domain:  "",
			wantErr: true,
			errMsg:  "domain is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendCmd(client, tt.domain, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("SendCmd() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateCmd_Subcommands(t *testing.T) {
	client := smtpsdk.NewClient("testuser", "testtoken")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no subcommand",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "unknown subcommand",
			args:    []string{"unknown"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCmd(client, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSuppressionsCmd_Subcommands(t *testing.T) {
	client := smtpsdk.NewClient("testuser", "testtoken")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no subcommand",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "unknown subcommand",
			args:    []string{"unknown"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SuppressionsCmd(client, "example.com", tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("SuppressionsCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebhooksCmd_Subcommands(t *testing.T) {
	client := smtpsdk.NewClient("testuser", "testtoken")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no subcommand",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "unknown subcommand",
			args:    []string{"unknown"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WebhooksCmd(client, "example.com", tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("WebhooksCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogsCmd_MissingDomain(t *testing.T) {
	client := smtpsdk.NewClient("testuser", "testtoken")

	err := LogsCmd(client, "", []string{})
	if err == nil {
		t.Error("LogsCmd() expected error for missing domain")
	}
	if !contains(err.Error(), "domain is required") {
		t.Errorf("LogsCmd() error = %v, want error containing 'domain is required'", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
