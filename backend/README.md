# Create your first migration

To create your first migration, you can use the following command:

```bash
go run cmd/migrate/main.go create <name-of-migration>
```

Output

````bash
migrations/
    20260802143001_create_users.up.sql
    20260802143001_create_users.down.sql
    ```
````

## Run migrations
```bash 
go run cmd/migrate/main.go up
```

other imortant commands

```bash
go run cmd/migrate/main.go up

go run cmd/migrate/main.go down

go run cmd/migrate/main.go version

go run cmd/migrate/main.go force 5
```