# Create your first migration

To create your first migration, you can use the following command at the root of your project:

```bash
make migrate-create name=create_users
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
make migrate-up
```

## Run down migrations
```bash
make migrate-down
```
