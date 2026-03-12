# Downloading Streams

Requires `vlc_path` to be set in the [config](configuration.md).

## Download a VOD stream
```
xtcli download <stream-id>
```

By default the output file is named `stream_<id>.mkv`. To specify the output path:
```
xtcli download <stream-id> --output /path/to/movie.mkv
```

## Download as a different format
```
xtcli download <stream-id> --format mp4
```

## Download a live stream
```
xtcli download <stream-id> --type live --format ts
```

## Download a favorite
```
xtcli download --fav <number-or-name>
```

See [favorites.md](favorites.md) for how to manage favorites.

## Quiet mode (no progress output)
```
xtcli download <stream-id> --quiet
```

## Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--output` | `-o` | `stream_<id>.<ext>` | Output file path |
| `--type` | `-t` | `vod` | Stream type: `live` or `vod` |
| `--format` | `-f` | `mkv` | Container format, e.g. `mkv`, `mp4`, `ts` |
| `--quiet` | `-q` | false | Suppress progress output |
| `--fav` | | | Favorite number or name to download |
