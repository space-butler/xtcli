# Browsing

Browsing information from your provider is primarily done through the `list` command. Such as categories, streams, VOD titles, EPG data, and stream URLs.

## List live TV categories
```
xtcli list categories
```

Aliases: `list c`, `list cat`

## List VOD categories
```
xtcli list categories --type vod
```

## List streams in a category

Use the category ID from `list categories`:
```
xtcli list streams <category-id>
```

Alias: `list s <category-id>`

## List VOD streams in a category
```
xtcli list streams <category-id> --type vod
```

## Show details for a single stream
```
xtcli list stream <stream-id>
```

For a VOD stream:
```
xtcli list stream <stream-id> --type vod
```

## Show EPG data for a stream
```
xtcli list epg <stream-id>
```

Alias: `list e <stream-id>`

Returns the next 4 EPG entries by default. To retrieve more:
```
xtcli list epg <stream-id> --limit 10
```

## Get the stream URL
```
xtcli list url <stream-id>
```

Alias: `list u <stream-id>`

To get a specific format:
```
xtcli list url <stream-id> --format m3u8
```

## Flags

### `list categories`

| Flag | Short | Default | Description |
|---|---|---|---|
| `--type` | `-t` | `live` | Category type: `live` or `vod` |

### `list streams`

| Flag | Short | Default | Description |
|---|---|---|---|
| `--type` | `-t` | `live` | Stream type: `live` or `vod` |

### `list stream`

| Flag | Short | Default | Description |
|---|---|---|---|
| `--type` | `-t` | `live` | Stream type: `live` or `vod` |

### `list epg`

| Flag | Short | Default | Description |
|---|---|---|---|
| `--limit` | `-l` | `4` | Number of EPG entries to retrieve |

### `list url`

| Flag | Short | Default | Description |
|---|---|---|---|
| `--format` | `-f` | `ts` | Stream format (e.g. `ts`, `m3u8`) |
