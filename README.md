### Server which receives and lists payments
I made this project as an interview assignment. I used https://github.com/kyleconroy/sqlc
to generate /postgres folder from sqlc.yaml.

The one thing I would do differently if this project would be deployed in production and used by customers, is not using float64 (in Go) and float(in postgres) for money type.
I should have used numeric in postgres. However, since Golang does not have an equivalent for big numbers (such as numeric)
it stores these values as strings. Therefore I decided I do not want to implement basic integer operations (such as +, -, *, /) on strings.
As this project will not be deployed and used anywhere, I decided to use float64 and float for money type.


####To run the app in production mode

```
make up
make createdb
make bringtofg
```

Go to http://localhost:5000/payments/ or use curl to make requests.

##### To stop the container
```
make down
```

###To run the app in development mode
```
make updb
make createdb
go run main.go -dbport=5431 -dbhost=localhost
```

#####To log to db in container
``` 
psql -h localhost -p 5432 -U postgres -d dev
```

###Client

##### Create payment
```
curl -X POST http://localhost:5000/payments/ -d '{"amount": 10, "sender":"saska","description":"1st payment","datetime": "2020-03-10T12:06:56.720Z" }'
```

##### Get payment with id
```
curl -X GET http://localhost:5000/payments/{id}
```

##### Get all payments
```
curl -X GET http://localhost:5000/payments/
```

### Linter
```
golangci-lint run --fix
```