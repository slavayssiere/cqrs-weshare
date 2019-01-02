# Write webservice

## Struct

```language-golang
type User struct {
  Name         string    `json:"username"`
  Email        string    `json:"email"`
  Address      string    `json:"address"`
  Age          int       `json:"age"`
  CreationTime int64     `json:"creation_time"`
  CreateTime   time.Time `json:"create_at"`
  ID           int       `json:"id"`
}
```

## Build

## Test

```language-bash
curl -X POST http://localhost:8080/users -d '{"username":"slavayssiere", "email":"slavayssiere@wescale.fr", "address":"23 rue taitbout 75009", "age":32}'
```

```language-bash
curl -X PUT http://localhost:8080/users/0 -d '{"address":"23 rue taitbout 75009 Paris"}'
```