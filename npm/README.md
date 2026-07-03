# @kirimemail/ktx

KirimEmail SMTP CLI client — manage domains, credentials, send emails, validate addresses, view logs and webhooks, check quota.

## Install

```bash
npm install -g @kirimemail/ktx
```

Or run without installing:

```bash
npx @kirimemail/ktx quota
```

## Usage

```bash
# List domains
ktx domains list -username user -token token

# Send email
ktx send -domain example.com -from noreply@example.com -to user@example.com -subject "Hello" -text "Body"

# Check quota
ktx quota
```

Configuration via `~/.ktxrc`, environment variables, or flags. See [full documentation](https://github.com/kirimemail/ktx).

## Authentication Priority

Flag > Environment Variable > Config File (`~/.ktxrc`)

Flags: `-username user -token token`
Env: `KIRIM_USERNAME` / `KIRIM_TOKEN`
Config: `~/.ktxrc` (JSON with `username`, `token`, `base_url`, `domain`)

## Supported Platforms

- Linux x64 & arm64
- macOS x64 & arm64 (Apple Silicon)
- Windows x64 & arm64

## License

MIT
