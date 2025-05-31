# Golang project for a simple blog 
Test how create a CRUD using Go

### Installation
1. Install PostgreSQL if you haven't already
2. Create the database using the bd.sql file:
```bash
psql -U postgres -f bd.sql
```

After that run the project using the next command [on the root folder]:
```bash
go mod tidy
go run ./main.go
```

### Front 
The styles are using the BootstrapV5 framework 

### Folders
/views
    /components
    /home
    /posts

### Packages
For PostgreSQL connection: 
- github.com/lib/pq
- database/sql

For Logs
- log

For the http server
- net/http

For the Front views
- text/template

For get the format for the dates
- time

