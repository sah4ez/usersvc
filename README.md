# About

Example using go-kit

# usersvc

Run the example with the optional port address for the service: 

```bash
$ make build && ./bin/usersvc -http.addr=:8080
rm -rf ./bin/
CC=gcc vgo build -o ./bin/usersvc ./cmd/usersvc 
ts=2018-05-09T18:54:43.826157153Z caller=main.go:47 transport=HTTP addr=:8080
```
Create user:
```bash
$ curl -X POST localhost:8080/users/ -d '{"name":"1231","email":"2@mail.com","password":"fff"}'
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjU5Nzg0OTEsIm5hbWUiOiIxMjMxIn0.ipobn3w2PlODb1rAUkj-uCn_gQPVREsRdh-8GSYLqnE","id":"827c4aa9-d210-45ad-9596-0822d1672b6d"}
```

Get all users:
```bash
$ curl -X GET localhost:8080/users/
{"users":[{"id":"827c4aa9-d210-45ad-9596-0822d1672b6d","name":"1231","email":"2@mail.com","password":"fff"}]}
```

Update user:
``` bash
$ curl -X PATCH -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjU5Nzg0OTEsIm5hbWUiOiIxMjMxIn0.ipobn3w2PlODb1rAUkj-uCn_gQPVREsRdh-8GSYLqnE" localhost:8080/users/827c4aa9-d210-45ad-9596-0822d1672b6d -d '{"name":"1231321","email":"1231231","password":"12312312"}'
{"id":"827c4aa9-d210-45ad-9596-0822d1672b6d"}
```

Get user by id:
```bash
$ curl -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjU5Nzg0OTEsIm5hbWUiOiIxMjMxIn0.ipobn3w2PlODb1rAUkj-uCn_gQPVREsRdh-8GSYLqnE" localhost:8080/users/827c4aa9-d210-45ad-9596-0822d1672b6d 
{"user":{"id":"","name":"1231321","email":"1231231","password":"12312312"}}
```

