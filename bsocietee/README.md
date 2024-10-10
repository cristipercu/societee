This is a starter kit for a golang REST API with JWT auth and postgres to store the user info.
I used go version go1.22.0, and I tried to use as few as possible external dependencies. 

You will need a postgres instance running, and if you have docker you can use the bellow command to start an instance:
```
docker run --name my-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
```

Make sure to check cmd/config/env.go for the environment variables that you can use and the default values, maybe you need to change something. Also, those variables are read from the system ENV, so maybe you want to use this feature.

After the postgres is up, you should install the dependencies using:
```
go mod tidy
```

Next step MIGRATIONS. You can find the initial migrations under cmd/migrate/migrations. And you can run the migrations using the make command from the provided Makefile:
```
make migrate-up
```

If all good, you can run the server:
```
go run .
```

The readme is still WIP! I will add more info in the future!
