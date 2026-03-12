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

Play a favorite:
```
xtcli play --fav <number-or-name>
```

Download a favorite:
```
xtcli download --fav <number-or-name>
```
