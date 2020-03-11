####To run the app in production mode

```
make up
make createdb
```

Go to http://localhost:5000/payments/

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

#####To log in to db in container
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