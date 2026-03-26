# ktx - KirimEmail SMTP CLI Client

CLI tool for managing KirimEmail SMTP API - domains, credentials, email sending, validation, logs, suppressions, webhooks, and quota.

## Installation

```bash
go install
```

## Configuration

### Config File (`~/.ktxrc`)

Create `~/.ktxrc` with JSON format:

```json
{
  "username": "your-username",
  "token": "your-api-token",
  "base_url": "https://smtp-app.kirim.email",
  "domain": "default-domain.com"
}
```

### Environment Variables

```bash
export KIRIM_USERNAME=your-username
export KIRIM_TOKEN=your-api-token
export KIRIM_DOMAIN=default-domain.com
```

### Flags

All options can be passed as flags:

```bash
ktx [command] -username user -token token -domain example.com
```

## Usage

### Domains

```bash
# List all domains
ktx domains list

# With explicit domain
ktx domains list -domain example.com

# Create a domain
ktx domains create example.com

# Get domain details
ktx domains get example.com

# Delete a domain
ktx domains delete example.com

# Verify DNS records
ktx domains verify example.com
```

### Credentials

```bash
# List credentials (domain from config/env)
ktx credentials list

# With explicit domain
ktx credentials list -domain example.com

# Create a credential
ktx credentials create myuser
ktx credentials create -domain example.com myuser

# Get credential details
ktx credentials get <guid>
ktx credentials get -domain example.com <guid>

# Delete a credential
ktx credentials delete <guid>
ktx credentials delete -domain example.com <guid>

# Reset password
ktx credentials reset-password <guid>
ktx credentials reset-password -domain example.com <guid>
```

### Send Email

```bash
# Send plain text email (domain from config/env)
ktx send -from noreply@example.com -to user@example.com -subject "Hello" -text "Message body"

# With explicit domain
ktx send -domain example.com -from noreply@example.com -to user@example.com -subject "Hello" -text "Message body"

# Send HTML email
ktx send -from noreply@example.com -to user@example.com -subject "Hello" -html "<h1>Hello</h1><p>Message body</p>"

# Send to multiple recipients (comma-separated)
ktx send -from noreply@example.com -to user1@example.com,user2@example.com -subject "Hello" -text "Message"
```

### Email Validation

```bash
# Validate single email
ktx validate email user@example.com

# Batch validate (comma-separated, max 100)
ktx validate batch user1@example.com,user2@example.com,invalid-email
```

### Logs

```bash
# Get logs (domain from config/env)
ktx logs

# With explicit domain and filters
ktx logs -domain example.com -start 2024-01-01 -end 2024-12-31 -sender sender@example.com -recipient recipient@example.com
```

### Suppressions

```bash
# List all suppressions (domain from config/env)
ktx suppressions list

# With explicit domain
ktx suppressions list -domain example.com

# List by type
ktx suppressions list unsubscribe
ktx suppressions list bounce
ktx suppressions list whitelist

# Create whitelist entry
ktx suppressions create-whitelist user@example.com email "Trusted customer"
ktx suppressions create-whitelist -domain example.com user@example.com email "Trusted customer"

# Delete suppressions
ktx suppressions delete unsubscribe 1,2,3
ktx suppressions delete -domain example.com unsubscribe 1,2,3
```

### Webhooks

```bash
# List webhooks (domain from config/env)
ktx webhooks list

# With explicit domain
ktx webhooks list -domain example.com

# Create webhook
ktx webhooks create delivered https://example.com/webhook
ktx webhooks create -domain example.com delivered https://example.com/webhook

# Get webhook details
ktx webhooks get <guid>
ktx webhooks get -domain example.com <guid>

# Delete webhook
ktx webhooks delete <guid>
ktx webhooks delete -domain example.com <guid>

# Test webhook
ktx webhooks test https://example.com/webhook delivered
ktx webhooks test -domain example.com https://example.com/webhook delivered
```

### Quota

```bash
ktx quota
```

## Flags

| Flag       | Default                        | Description                    |
|------------|-------------------------------|--------------------------------|
| `-username`| -                             | KirimEmail username            |
| `-token`  | -                             | KirimEmail API token           |
| `-baseurl`| `https://smtp-app.kirim.email`| API base URL                   |
| `-domain` | -                             | Default domain to use          |

## Priority

Flag > Environment Variable > Config File (`~/.ktxrc`)

When `domain` is set via config file or environment variable, you can omit the `-domain` flag for commands that require it.
