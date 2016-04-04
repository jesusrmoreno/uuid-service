# uuid-service
Very simple uuid generator as a service

Defaults to port 3000
```bash
go get github.com/jesusrmoreno/uuid-service
uuid-service --port="8080" --dbPath="databasePath"
```

#API

## Get a new ID
```
/id
```
## Returns a new UUID and saves it under the provided namespace
```
/{namespace}/id
```

## Returns true if an id exists under that namespace
```
/{namespace}/id/exists/{uuid}
```
