# Favorites

Favorites let you save streams by name for quick access with `play` and `download`.

## Add a favorite
```
xtcli fav add --name "BBC News" --id <stream-id> --type live
```

For a VOD stream:
```
xtcli fav add --name "Inception" --id <stream-id> --type vod
```

## List favorites
```
xtcli fav list
```

## Remove a favorite by number or name
```
xtcli fav del 1
xtcli fav del "BBC News"
```

Multiple favorites can be removed at once:
```
xtcli fav del 1 2 3
```

## Reorder favorites
```
xtcli fav swap 1 3
```

## Using favorites

See [playing.md](playing.md) and [downloading.md](downloading.md) for full details and available flags.

For example, to play a favorite:
```
xtcli play --fav <number-or-name>
```

To download a favorite:
```
xtcli download --fav <number-or-name>
```

## Flags

### `fav add`

| Flag | Short | Default | Description |
|---|---|---|---|
| `--name` | `-n` | | Shortcut name for the favorite (required) |
| `--id` | `-i` | | Stream ID (required) |
| `--type` | `-t` | `live` | Stream type: `live` or `vod` |

### `fav del`

Takes one or more positional arguments: favorite numbers or names.

### `fav swap`

Takes two positional arguments: the numbers of the two favorites to swap.
