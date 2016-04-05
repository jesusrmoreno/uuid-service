# uuid-service
Very simple uuid generator as a service

Defaults to port 3000
# Install
## Using go get
```bash
go get github.com/jesusrmoreno/uuid-service
uuid-service --port="8080" --dbPath="databasePath"
```
## From release
Download the version for your OS from the releases page.
[Releases](https://github.com/jesusrmoreno/uuid-service/releases)

#API
## Get a new ID
```
/new/uuid_v4
/new/squid
/new/simplesquid
```
## Returns a new UUID and saves it under the provided namespace
```
/{namespace}/new/uuid_v4
/{namespace}/new/squid
/{namespace}/new/simplesquid
```

## Returns all ids in the given namespace
```
/{namespace}/
```
## Returns a new UUID and saves it under the provided namespace
```
/{namespace}/id
/{namespace}/squid
```

## Returns true if an id exists under that namespace
```
/{namespace}/contains/{uuid}
```
