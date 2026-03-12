# Cache

xtcli caches category and stream data locally to avoid repeated server requests. The cache TTL can be configured via `cache_ttl` in the [config](configuration.md) (default: 24 hours).

## Show cache information
```
xtcli config cache info
```

## Refresh the cache
```
xtcli config cache update
```

## Clear the cache
```
xtcli config cache clear
```
