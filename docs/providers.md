# Providers

Providers represent the IPTV servers you have access to. They can be managed without editing the config file manually.

## Add or update a provider
```
xtcli provider add --name myprovider --username user --password pass --host https://your.server
```

## List providers
```
xtcli provider list
```

## Set the default provider
```
xtcli provider default myprovider
```

## Remove a provider
```
xtcli provider del myprovider
```

## Use a non-default provider for a single command

Any command accepts a global `--provider` flag:
```
xtcli --provider myprovider list categories
```
