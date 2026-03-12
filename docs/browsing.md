# Browsing

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
