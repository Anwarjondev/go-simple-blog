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

### Docker Deployment

This project uses Docker for containerization and GitHub Actions for CI/CD. To set up the deployment:

1. Add the following secrets to your GitHub repository:
   - `DOCKER_HUB_USERNAME`: Your Docker Hub username
   - `DOCKER_HUB_ACCESS_TOKEN`: Your Docker Hub access token
   - `EC2_HOST`: Your EC2 instance's public IP or DNS
   - `EC2_USERNAME`: Your EC2 instance's username
   - `EC2_SSH_KEY`: Your private SSH key for EC2 access

2. On your EC2 instance, install Docker and Docker Compose:
   ```bash
   # Install Docker
   sudo apt update
   sudo apt install docker.io docker-compose
   
   # Add your user to the docker group
   sudo usermod -aG docker $USER
   
   # Start Docker service
   sudo systemctl start docker
   sudo systemctl enable docker
   ```

3. Make sure your EC2 security group allows:
   - SSH (port 22)
   - HTTP (port 80)
   - HTTPS (port 443)
   - Application port (8000)

The CI/CD pipeline will automatically:
- Build a Docker image
- Push it to Docker Hub
- Deploy to your EC2 instance using Docker Compose
- Set up PostgreSQL database
- Restart the containers when new changes are pushed to the main branch

