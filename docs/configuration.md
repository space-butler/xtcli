# Configuration

The config file is stored at `~/.xtcli/config.json`. Generate a default config with:

```
xtcli config create
```

Then edit the file to fill in your provider credentials and VLC path. The structure is:

```json
{
    "default_provider": "myprovider",
    "providers": [
        {
            "name": "myprovider",
            "username": "your_username",
            "password": "your_password",
            "host": "https://your.xtream.server"
        }
    ],
    "vlc_path": "/Applications/VLC.app/Contents/MacOS/VLC",
    "cache_ttl": 24
}
```

| Field | Description |
|---|---|
| `default_provider` | The provider used when `--provider` is not specified |
| `providers` | List of configured IPTV providers |
| `vlc_path` | Path to the VLC binary; required for `play` and `download` |
| `cache_ttl` | How long cached data is considered fresh, in hours (default: 24) |

## Providers

See [providers.md](providers.md) for managing providers via the CLI.

## Cache

See [cache.md](cache.md) for managing the local cache.
