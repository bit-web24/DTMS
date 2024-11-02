# DTMS - `Distributed Task Management System`

### compile .proto files
```bash
make
```

### services/user/.env
```bash
DB_HOST=localhost
DB_USER=bittu
DB_PASSWORD=bittu
DB_NAME=users
DB_PORT=5432
DB_SSLMODE=disable
DB_TIME_ZONE=UTC
```

### services/task/.env
```bash
DB_HOST=localhost
DB_USER=bittu
DB_PASSWORD=bittu
DB_NAME=tasks
DB_PORT=5432
DB_SSLMODE=disable
DB_TIME_ZONE=UTC
```

### To start the server
```bash
sudo docker-compose up --build
```

### To stop the server
```bash
sudo dokcer-compose down
```

## Design
### Main Interface (For Admin)
![Screenshot from 2024-08-04 14-48-18](https://github.com/user-attachments/assets/8d7dd527-0cae-479a-875e-66a578d3dea5)

### Add task window
![Screenshot from 2024-08-04 14-50-25](https://github.com/user-attachments/assets/bd8f19b9-903f-42e2-b87c-6ca4e4edeadd)

### Main Interface (Workers)
![Screenshot from 2024-08-04 14-52-41](https://github.com/user-attachments/assets/97cdb813-dd75-404c-9cd9-c7dccad505dd)
