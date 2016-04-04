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
/id
/squid 
```
## Returns a new UUID and saves it under the provided namespace
```
/{namespace}/id
/{namespace}/squid
```

## Returns true if an id exists under that namespace
```
/{namespace}/exists/{uuid}
```
