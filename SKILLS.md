# ktx CLI Skills

This document describes how to use the `ktx` CLI for KirimEmail SMTP API management.

## Quick Reference

### Setup
```bash
# Install ktx
curl -fsSL https://raw.githubusercontent.com/kirimemail/ktx/main/install.sh | bash

# Configure credentials
export KIRIM_USERNAME=your-username
export KIRIM_TOKEN=your-api-token
export KIRIM_DOMAIN=your-domain.com
```

### Domains
```bash
# List all domains
ktx domains list

# Create domain
ktx domains create example.com

# Get domain details
ktx domains get example.com

# Delete domain
ktx domains delete example.com

# Verify DNS records
ktx domains verify example.com
```

### Credentials
```bash
# List credentials
ktx credentials list

# Create credential
ktx credentials create myuser

# Get credential
ktx credentials get <guid>

# Delete credential
ktx credentials delete <guid>

# Reset password
ktx credentials reset-password <guid>
```

### Send Email
```bash
# Basic send
ktx send -from noreply@example.com -to user@example.com -subject "Hello" -text "Message"

# HTML email
ktx send -from noreply@example.com -to user@example.com -subject "Hello" -html "<h1>Hello</h1>"

# Multiple recipients
ktx send -from noreply@example.com -to user1@example.com,user2@example.com -subject "Hello" -text "Message"
```

### Email Validation
```bash
# Single email
ktx validate email user@example.com

# Batch (max 100)
ktx validate batch user1@example.com,user2@example.com
```

### Logs
```bash
# Get logs
ktx logs

# With filters
ktx logs -start 2024-01-01 -end 2024-12-31 -sender sender@example.com -recipient recipient@example.com
```

### Suppressions
```bash
# List all
ktx suppressions list

# List by type
ktx suppressions list unsubscribe
ktx suppressions list bounce
ktx suppressions list whitelist

# Create whitelist
ktx suppressions create-whitelist user@example.com email "Trusted"

# Delete
ktx suppressions delete unsubscribe 1,2,3
```

### Webhooks
```bash
# List
ktx webhooks list

# Create
ktx webhooks create delivered https://example.com/webhook

# Get
ktx webhooks get <guid>

# Delete
ktx webhooks delete <guid>

# Test
ktx webhooks test https://example.com/webhook delivered
```

### Quota
```bash
ktx quota
```

## Configuration

### Config File (`~/.ktxrc`)
```json
{
  "username": "your-username",
  "token": "your-api-token",
  "base_url": "https://smtp-app.kirim.email",
  "domain": "example.com"
}
```

### Priority
Flag > Environment Variable > Config File

## Common Workflows

### Send transactional email
```bash
export KIRIM_DOMAIN=example.com
ktx send -from orders@example.com -to customer@email.com -subject "Order Confirmed" -text "Your order #1234 has been confirmed."
```

### Verify domain setup
```bash
ktx domains create example.com
# Add DNS records shown
ktx domains verify example.com
```

### Monitor email logs
```bash
ktx logs -start 2024-01-01 -recipient user@example.com
```

### Manage suppressions
```bash
# Check bounce list
ktx suppressions list bounce

# Whitelist a sender
ktx suppressions create-whitelist trusted@example.com email "Trusted sender"

# Remove from suppression
ktx suppressions delete bounce 1,2,3
```
