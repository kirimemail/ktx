# ktx - KirimEmail SMTP CLI

CLI tool for managing KirimEmail SMTP API — domains, credentials, sending, validation, logs, suppressions, webhooks, quota.

---

## Setup & Configuration

### Authentication priority (flag > env > config file)

```
1. Flag:      -username <user> -token <token>
2. Env var:   KIRIM_USERNAME / KIRIM_TOKEN
3. Config:    ~/.ktxrc (JSON)
```

### Config file (~/.ktxrc)

```json
{
  "username": "your-username",
  "token": "your-api-token",
  "base_url": "https://smtp-app.kirim.email",
  "domain": "example.com"
}
```

### Domain resolution (for commands that need it)

Commands that operate on a domain (`credentials`, `logs`, `suppressions`, `webhooks`) accept the domain in two ways:

1. **As a positional arg** after the subcommand (most explicit)
2. **From config/env** via `-domain` flag, `KIRIM_DOMAIN` env, or `domain` field in `~/.ktxrc`

**Send** is different — it takes `-domain` as a flag, not a positional arg.

---

## Domains

### Commands

| Subcommand | Positional args | Description |
| --- | --- | --- |
| `list` | *(none)* | List all domains |
| `create` | `<domain>` | Create domain (DKIM 2048, auto) |
| `get` | `<domain>` | Get domain details |
| `delete` | `<domain>` | Delete domain |
| `verify` | `<domain>` | Verify mandatory DNS records (DKIM, SPF, MX) |
| `setup-auth-domain` | `<domain> [auth_domain]` | Set up auth sending domain (defaults to same domain) |
| `verify-auth-domain` | `<domain>` | Verify auth domain DNS records |
| `delete-auth-domain` | `<domain>` | Delete auth domain |
| `setup-tracklink` | `<domain> [tracking_domain]` | Set up tracking link domain (defaults to same domain) |
| `verify-tracklink` | `<domain>` | Verify tracking link CNAME |
| `delete-tracklink` | `<domain>` | Delete tracking link |

### Examples

```bash
ktx domains list
ktx domains create mydomain.com
ktx domains get mydomain.com
ktx domains verify mydomain.com
ktx domains setup-auth-domain mydomain.com auth.mydomain.com
ktx domains delete mydomain.com
```

### Usage patterns

```
ktx domains <subcommand> [args...]
```

### Output format (list)

Tab-separated columns with headers:

```
DOMAIN         STATUS  VERIFIED
mydomain.com   active  true
```

### Output format (get)

```
Domain: mydomain.com
Status: active
Verified: true
```

### Output format (verify)

```
DKIM: true
SPF: true
MX: true
```

### Error messages

| Error | Cause |
| --- | --- |
| `domains command requires subcommand: ...` | No subcommand given |
| `usage: ktx domains create <domain>` | Missing domain arg |
| `usage: ktx domains verify <domain>` | Missing domain arg |
| API error string | API call failed (auth, not found, already exists) |

---

## Credentials

### Commands

| Subcommand | Positional args | Description |
| --- | --- | --- |
| `list` | `<domain>` | List credentials for a domain |
| `create` | `<domain> <username>` | Create SMTP credential (password printed once) |
| `get` | `<domain> <guid>` | Get credential details by GUID |
| `delete` | `<domain> <guid>` | Delete credential |
| `reset-password` | `<domain> <guid>` | Reset password (new password printed once) |

> **IMPORTANT for agent**: `create` and `reset-password` print the password to stdout **only once** and it's not retrievable later. The agent **must** capture and relay this immediately.

### Examples

```bash
ktx credentials list mydomain.com
ktx credentials create mydomain.com myuser
ktx credentials get mydomain.com <guid>
ktx credentials reset-password mydomain.com <guid>
ktx credentials delete mydomain.com <guid>
```

### Usage patterns

```
ktx credentials <subcommand> <domain> [args...]
```

### Output format (list)

```
USERNAME  DOMAIN
myuser    mydomain.com
```

### Output format (create)

```
Credential created:
  Username: myuser
  Password: smtp-password-here
```

### Output format (get)

```
Username: myuser
```

### Output format (reset-password)

```
New password: new-password-here
```

### Error messages

| Error | Cause |
| --- | --- |
| `credentials command requires subcommand: ...` | No subcommand given |
| `usage: ktx credentials list <domain>` | Missing domain |
| `usage: ktx credentials create <domain> <username>` | Missing domain or username |
| `usage: ktx credentials get <domain> <guid>` | Missing domain or GUID |
| `usage: ktx credentials delete <domain> <guid>` | Missing domain or GUID |
| `usage: ktx credentials reset-password <domain> <guid>` | Missing domain or GUID |

---

## Send Email

### Flags

| Flag | Required | Description |
| --- | --- | --- |
| `-domain` | depends | Sending domain (required if not in config/env) |
| `-from` | **yes** | Sender email address |
| `-to` | **yes** | Recipient(s), comma-separated for multiple |
| `-subject` | **yes** | Email subject |
| `-text` | one of | Plain text body |
| `-html` | one of | HTML body |

Either `-text`, `-html`, or both must be provided.

### Examples

```bash
ktx send -domain mydomain.com -from noreply@mydomain.com -to user@example.com -subject "Hello" -text "Message body"
ktx send -domain mydomain.com -from noreply@mydomain.com -to user@example.com -subject "Hello" -html "<h1>Hello</h1><p>Body</p>"
ktx send -domain mydomain.com -from noreply@mydomain.com -to user1@example.com,user2@example.com -subject "Hi" -text "Message"
```

### Usage patterns

```
ktx send -domain <domain> -from <from> -to <to> -subject <subject> [-text <text>] [-html <html>]
```

### Output format

```
Message sent: <message-id>
```

### Error messages

| Error | Cause |
| --- | --- |
| `usage: ktx send -domain <domain> -from <from> ...` | Missing required flag (-from, -to, or -subject) |
| `domain is required (use -domain flag or KIRIM_DOMAIN env)` | No domain resolved |
| API error | Sending failed (invalid address, quota, etc.) |

---

## Email Validation

### Commands

| Subcommand | Args | Description |
| --- | --- | --- |
| `email` | `<email>` | Validate a single email address |
| `batch` | `<email1,email2,...>` | Batch validate (max 100) |

### Examples

```bash
ktx validate email user@example.com
ktx validate batch user1@example.com,user2@example.com,invalid-email
```

### Usage patterns

```
ktx validate <subcommand> <args...>
```

### Output format (email)

```
Email: user@example.com
Valid: true
Is spamtrap: false
Spamtrap score: 0.00
```

### Output format (batch)

```
Total: 3
Valid: 2
Invalid: 1
```

### Error messages

| Error | Cause |
| --- | --- |
| `validate command requires subcommand: email, batch` | No subcommand given |
| `usage: ktx validate email <email>` | Missing email |
| `usage: ktx validate batch <email1,email2,...>` | Missing email list |

### Best practices for agent

- Use `grep -i` to check `Valid:` value from output
- For batch validation, keep lists under 100 addresses
- Combine with file reads to validate subscriber lists from CSV files

---

## Logs

### Flags

| Flag | Description |
| --- | --- |
| `-start` | Start date (YYYY-MM-DD) |
| `-end` | End date (YYYY-MM-DD) |
| `-sender` | Filter by sender email |
| `-recipient` | Filter by recipient email |
| `-subject` | Filter by subject |
| `-eventType` | Filter by event type (e.g. `delivered`, `bounce`, `open`, `click`) |
| `-tags` | Filter by tags |
| `-csv` | Output as CSV instead of table |

### Examples

```bash
ktx logs -domain mydomain.com
ktx logs -domain mydomain.com -start 2024-01-01 -end 2024-12-31
ktx logs -domain mydomain.com -recipient user@example.com
ktx logs -domain mydomain.com -eventType bounce
ktx logs -domain mydomain.com -csv
```

> **Note:** `-domain` is required unless set in config/env. Unlike other commands, logs takes `-domain` as a flag, not a positional arg.

### Usage patterns

```
ktx logs -domain <domain> [flags...]
```

### Output format (table)

```
TIMESTAMP                  EVENT      MESSAGE GUID
2024-01-15T10:30:00Z       delivered  msg-abc-123
```

### Output format (CSV)

```
TIMESTAMP,EVENT,MESSAGE GUID
2024-01-15T10:30:00Z,delivered,msg-abc-123
```

### Error messages

| Error | Cause |
|---|---|
| `domain is required (use -domain flag or KIRIM_DOMAIN env)` | No domain resolved |

---

## Suppressions

### Commands

| Subcommand | Positional args | Description |
| --- | --- | --- |
| `list` | `<domain> [type]` | List suppressions (optionally by type: `unsubscribe`, `bounce`, `whitelist`) |
| `create-whitelist` | `<domain> <recipient> <type> [description]` | Add a whitelist entry |
| `delete` | `<domain> <type> <id1,id2,...>` | Delete suppression entries by comma-separated IDs |

### Examples

```bash
ktx suppressions list mydomain.com
ktx suppressions list mydomain.com bounce
ktx suppressions list mydomain.com whitelist
ktx suppressions create-whitelist mydomain.com user@example.com email "Trusted customer"
ktx suppressions delete mydomain.com bounce 1,2,3
```

### Usage patterns

```
ktx suppressions <subcommand> <domain> [args...]
```

### Output format (list)

```
RECIPIENT              TYPE        DESCRIPTION
user@example.com       bounce      Some error
```

### Output format (create-whitelist)

```
Whitelist entry created: user@example.com
```

### Output format (delete)

```
Deleted suppressions
```

### Error messages

| Error | Cause |
| --- | --- |
| `suppressions command requires subcommand: ...` | No subcommand given |
| `usage: ktx suppressions list <domain> [type]` | Missing domain |
| `usage: ktx suppressions create-whitelist <domain> <recipient> <type> [description]` | Missing required args |
| `usage: ktx suppressions delete <domain> <type> <id1,id2,...>` | Missing required args |
| `unknown suppression type: <type>` | Invalid type for delete (must be `unsubscribe`, `bounce`, or `whitelist`) |

### Suppression types

| Type | Description |
| --- | --- |
| `unsubscribe` | Unsubscribed recipients |
| `bounce` | Bounced emails |
| `whitelist` | Whitelisted recipients |

---

## Webhooks

### Commands

| Subcommand | Positional args | Description |
| --- | --- | --- |
| `list` | `<domain>` | List webhooks |
| `create` | `<domain> <type> <url>` | Create webhook |
| `get` | `<domain> <guid>` | Get webhook details |
| `delete` | `<domain> <guid>` | Delete webhook |
| `test` | `<domain> <url> <event_type>` | Send a test webhook to a URL |
| `update` | `<domain> <guid> <field> <value>` | Update webhook field (`type` or `url`) |

### Webhook event types

| Type | Description |
| --- | --- |
| `delivered` | Email delivered |
| `bounce` | Email bounced |
| `spam` | Marked as spam |
| `open` | Email opened |
| `click` | Link clicked |
| `unsubscribe` | Recipient unsubscribed |

### Examples

```bash
ktx webhooks list mydomain.com
ktx webhooks create mydomain.com delivered https://example.com/webhook
ktx webhooks get mydomain.com <guid>
ktx webhooks update mydomain.com <guid> url https://new-url.com/webhook
ktx webhooks delete mydomain.com <guid>
ktx webhooks test mydomain.com https://example.com/webhook delivered
```

### Usage patterns

```
ktx webhooks <subcommand> <domain> [args...]
```

### Output format (list)

```
GUID                                  TYPE        URL
abc-123-def                          delivered   https://example.com/webhook
```

### Output format (create)

```
Webhook created: abc-123-def
```

### Output format (get)

```
GUID: abc-123-def
Type: delivered
URL: https://example.com/webhook
```

### Output format (test)

```
Test successful! Response time: 123ms
```

```
Test failed with status: 404
```

### Output format (delete)

```
Webhook deleted: abc-123-def
```

### Output format (update)

```
Webhook updated: abc-123-def
```

### Error messages

| Error | Cause |
| --- | --- |
| `webhooks command requires subcommand: ...` | No subcommand given |
| `usage: ktx webhooks list <domain>` | Missing domain |
| `usage: ktx webhooks create <domain> <type> <url>` | Missing args |
| `usage: ktx webhooks get <domain> <guid>` | Missing args |
| `usage: ktx webhooks delete <domain> <guid>` | Missing args |
| `usage: ktx webhooks test <domain> <url> <event_type>` | Missing args |
| `usage: ktx webhooks update <domain> <guid> <field> <value>` | Missing args |
| `unknown field: <field> (use type or url)` | Invalid update field |

---

## Quota

### Usage

```bash
ktx quota
```

### Output format

```
Current quota: 5000
Max quota: 10000
Usage: 50.0%
```

### Error messages

| Error | Cause |
|---|---|
| API error | API call failed (auth, network) |

---

## Common Workflows

### 1. Onboard a new sending domain

```bash
# 1. Create domain
ktx domains create example.com

# 2. Verify DNS (DKIM, SPF, MX)
ktx domains verify example.com

# 3. Set up auth domain (if needed)
ktx domains setup-auth-domain example.com auth.example.com
ktx domains verify-auth-domain example.com

# 4. Set up tracking links (if needed)
ktx domains setup-tracklink example.com track.example.com
ktx domains verify-tracklink example.com

# 5. Create SMTP credentials
ktx credentials create example.com welcome-sender
```

### 2. Send a transactional email

```bash
ktx send -domain example.com -from orders@example.com -to customer@email.com -subject "Order Confirmed" -text "Your order #1234 has been confirmed."
```

### 3. Investigate delivery issues

```bash
# Check recent logs for a recipient
ktx logs -domain example.com -recipient user@example.com -eventType bounce

# Check if recipient is suppressed
ktx suppressions list example.com bounce
ktx suppressions list example.com unsubscribe

# Whitelist if needed
ktx suppressions create-whitelist example.com user@example.com email "Trusted"

# Check quota
ktx quota
```

### 4. Rotate credentials

```bash
# 1. Create new credential
ktx credentials create example.com new-user
# → Save the password immediately

# 2. Update application with new credentials
# (agent writes the new credentials to the app config)

# 3. Delete old credential
ktx credentials delete example.com <old-guid>
```

### 5. Set up delivery tracking

```bash
# Create webhooks for delivery events
ktx webhooks create example.com delivered https://app.example.com/webhooks/delivered
ktx webhooks create example.com bounce https://app.example.com/webhooks/bounce
ktx webhooks create example.com open https://app.example.com/webhooks/open

# Test each webhook
ktx webhooks test example.com https://app.example.com/webhooks/delivered delivered
ktx webhooks test example.com https://app.example.com/webhooks/bounce bounce

# Monitor logs
ktx logs -domain example.com -eventType bounce
```

### 6. Clean email list before campaign

```bash
# Check bounce list
ktx suppressions list example.com bounce

# Validate list from file (agent reads file, constructs batch)
ktx validate batch addr1@...,addr2@...,...
```

### 7. Post-deployment smoke test

```bash
# Verify domain DNS is still correct
ktx domains verify example.com

# Send a test email
ktx send -domain example.com -from test@example.com -to dev-team@company.com -subject "Deployment test $(date)" -text "Smoke test passed"

# Check it arrived
ktx logs -domain example.com -recipient dev-team@company.com -eventType delivered
```

---

## Key Differences in Subcommand Argument Patterns

Not all commands handle the domain the same way. Watch out for:

| Command | Domain placement |
| --- | --- |
| `domains` | Some subcommands take `<domain>` as a positional arg (`get`, `delete`, `verify`, etc.) |
| `credentials` | **Always** takes `<domain>` as first positional arg after subcommand |
| `send` | `-domain` is a **flag**, not positional |
| `validate` | No domain needed |
| `logs` | `-domain` is a **flag**, not positional |
| `suppressions` | **Always** takes `<domain>` as first positional arg after subcommand |
| `webhooks` | **Always** takes `<domain>` as first positional arg after subcommand |
| `quota` | No domain needed |

## Output Parsing for Agents

Since ktx uses plain tabwriter tables (not JSON), parse output like this:

```bash
# Grep a specific value from a known line
ktx domains get example.com | grep -i "verified:" | awk '{print $2}'
ktx quota | grep "Usage:" | awk '{print $2}'

# Get the first result's GUID from a list
ktx credentials list example.com | tail -n +2 | head -1 | awk '{print $1}'

# Check if a validation passed
ktx validate email test@example.com | grep -q "Valid: true" && echo "valid"
```
