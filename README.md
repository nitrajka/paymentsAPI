####To run the app in production mode

Download: https://github.com/vishnubob/wait-for-it

```
make up
make createdb
```

Go to localhost:5000/payments/

###To run the app in development mode
```
make updb
make createdb
go run main.go
```

#####To log in to db in container
``` 
psql -h localhost -p 5432 -U postgres -d dev
```
