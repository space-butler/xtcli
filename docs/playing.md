# Playing Streams

Requires `vlc_path` to be set in the [config](configuration.md).

## Play a live stream by ID
```
xtcli play <stream-id>
```

## Play a VOD stream by ID
```
xtcli play <stream-id> --type vod
```

By default VOD streams use the `ts` format. To specify a different format:
```
xtcli play <stream-id> --type vod --format mkv
```

## Play a favorite
```
xtcli play --fav <number-or-name>
```

When using `--fav`, the stream type is taken from the favorite automatically.

See [favorites.md](favorites.md) for how to manage favorites.

## Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--type` | `-t` | `live` | Stream type: `live` or `vod` |
| `--format` | `-f` | `ts` | Format/extension (e.g. `ts`, `m3u8`, `mp4`) |
| `--fav` | | | Favorite number or name to play |
