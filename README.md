# 🎟️ BookingApp (Go + Gin)

Short, practical demo of a thread-safe ticket booking API with a simple frontend.

## What it does

- In-memory store with RWMutex for concurrency safety.
- 15s seat holds (reservations) with live countdown and cancel/confirm.
- Fair FIFO wait queue per conference (Join Queue → Claim Now when first).
- Each user can have only one active reservation per conference.
- Users are unique by email (case-insensitive).
- Conferences are returned sorted by ID; UI shows on-hold and queue badges.
-

## Run locally (Windows cmd)

```bat
go build ./...
go test  ./...

REM If 8080 is busy, choose another port
set PORT=8081
go run .

REM Open the app
REM Frontend: http://127.0.0.1:%PORT%
REM API:      http://127.0.0.1:%PORT%/api/v1
```

Tip: You can override the API base from the browser with `?api=http://127.0.0.1:8081` if you ran on a non-default port.

## API (quick glance)

- GET /api/v1/health
- GET /api/v1/conferences // includes stats: reserved and queue size
- POST /api/v1/users // {name, email}
- POST /api/v1/reservations // {user_id, conference_id, ticket_count}
- GET /api/v1/reservations/:id
- POST /api/v1/reservations/:id/confirm
- DELETE /api/v1/reservations/:id
- GET /api/v1/bookings // testing/demo list
- GET /api/v1/users/:userID/bookings
- GET /api/v1/users/:userID/reservations
- POST /api/v1/queue/enqueue // {user_id, conference_id, ticket_count}
- GET /api/v1/queue/:conferenceID/position?user_id=...
- POST /api/v1/queue/claim // {user_id, conference_id}

## Docker (optional)

```bat
docker build -t booking-app .
docker run -p 8080:8080 booking-app
```

## Project structure

- main.go – routes/server
- models/models.go – User, Conference, Booking, SeatReservation
- database/database.go – in-memory data + business rules + wait queue
- handlers/handlers.go – HTTP handlers
- index.html – test UI (join, book, queue, timers)
- Dockerfile, docker-compose.yml

## Notes

- All data is in-memory for demo purposes; restarting clears state.
- The UI polls every 2s for queue position and every 1s for timers.
- CORS is enabled for quick local testing.
  cd c:\Users\ojade\Downloads\GOProject\BookingApp

# 2. Build the Docker image

docker build -t booking-system .

# 3. Run the container

docker run -p 8080:8080 booking-system

# 4. Test the API

curl http://localhost:8080/api/v1/health

````

### **Docker Commands Reference**

```bash
# Build image
docker build -t booking-system .

# Run container
docker run -p 8080:8080 booking-system

# Run container (background)
docker run -d -p 8080:8080 --name booking-api booking-system

# View running containers
docker ps

# View logs
docker logs booking-api

# Stop container
docker stop booking-api

# Remove container
docker rm booking-api

# Remove image
docker rmi booking-system
````

### **Docker Troubleshooting**

**Issue: "Docker daemon not running"**

```bash
# Solution: Start Docker Desktop application
# Check if Docker is running:
docker ps
```

**Issue: "Port already in use"**

```bash
# Solution: Use different port
docker run -p 8081:8080 booking-system

# Or stop existing process
taskkill /f /im main.exe
```

**Issue: "Build context too large"**

```bash
# Solution: Add .dockerignore file with:
echo "tests/" > .dockerignore
echo "*.md" >> .dockerignore
echo ".git/" >> .dockerignore
```

### **Docker Compose (Optional)**

Create `docker-compose.yml`:

```yaml
version: "3.8"
services:
  booking-api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
    restart: unless-stopped
```

Run with:

```bash
docker-compose up -d
docker-compose down
```

---



---

"# GoBookingApp" 
