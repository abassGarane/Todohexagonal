## Hexagonal Architechture base Todo Web application.
### Requirements
- Environmental Variables
  - _`URL_DB`_ - The selected database ie "redis" or "postgres" or "mongodb" or "mysql"
  - _`DATABASE_URL`_ - Incase you choose any db apart from redis
  - _`[REDIS]/[POSTGRES]/[MONGODB]/[MYSQL]_URL`_- dsn for the database eg `REDIS_URL` for redis
  - _`PORT`_ - Port to run application defaults to :8000
### ***Example***
- Linux bash variables
```bash
export DB_URL=redis
export REDIS_URL=localhost:6329
```

## Running application
```bash
go run cmd/web/main.go
```

## Functions
1. Adding a todo
```bash
http :8000 content="Create a website" status="started"
```
2. Getting a todo
```bash
http :8000/{todo_id}
```
3. Delete a todo
```bash
http DELETE :8000/{todo_id}
```
4. Update a todo
```bash
http PATCH :8080/{todo_id} <parts to be updated eg content="stuff">
```
