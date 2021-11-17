[![Go Report Card][goreport-badge]][goreport-url]


## Golang projects - Using Gin
This is one of a series of projects (I hope) for experimenting some go libraries.

For this one, I'll use [Gin](https://github.com/gin-gonic/gin) to build a simple web application and [jackc/pgx](https://github.com/jackc/pgx) to integrate with a PostgreSQL database


### Run the project
```bash
# Clone me!
git clone https://github.com/luanapp/using-gin

# Update dependencies
make install

# Run!
make run
```

### Run database migrations
There are some database migrations to create the database structure.

Jus run:
```bash
make migrate-up
```

To undo the migrations, run:
```bash
make migrate-up
```

Happy coding!

[goreport-badge]:https://goreportcard.com/badge/github.com/luanapp/using-gin
[goreport-url]:https://goreportcard.com/report/github.com/luanapp/using-gin
