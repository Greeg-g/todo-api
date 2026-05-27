# To-Do API (Go)

Minimal, practical REST API to manage tasks (To-Do). The project includes JWT authentication, task ownership and sharing, Redis caching hooks and a worker for deadline alerts. It's intended as a learning / small production-ready starter.

---

## Quick start (with Docker)

1. Copy or create a `.env` file in the project root with required variables (examples below).
2. Build and start services:

```bash
docker-compose up -d --build
```

3. Check API logs:

```bash
docker-compose logs -f api
```

To stop:

```bash
docker-compose down
```

Notes:
- The `Dockerfile` builds the Go binary inside the image and runs it, so you don't need to run `go run` locally while using Docker.
- For development you can use a bind mount + hot-reload tool (not included by default).

---

## Routes (current)

Auth
- POST /auth/register — register a new user
- POST /auth/login — login and receive JWT token

Tasks (all protected by JWT middleware — include header `Authorization: Bearer <token>`)
- POST /tasks/create — create a task (owner set from token)
- GET /tasks/ — list tasks owned by current user
- GET /tasks/recent — cached recent tasks for current user
- GET /tasks/shared — list tasks shared with current user
- POST /tasks/complete/:id — mark task completed (owner or shared-with only)
- DELETE /tasks/delete/:id — delete task (owner only)
- POST /tasks/share/:id — share a task with another user (by username or email)
- GET /tasks/category/:category — list current user's tasks in a category
- GET /tasks/owner/:owner — list tasks by owner; requester sees all if owner, otherwise only tasks shared with requester
- GET /tasks/check-deadlines — trigger deadline check (enqueue)

---

## Example JSON bodies (paste into Insomnia / Postman)

1) Register (POST /auth/register)

{
	"username": "alice",
	"email": "alice@example.com",
	"password": "Password123!"
}

2) Login (POST /auth/login)

{
	"email": "alice@example.com",
	"password": "Password123!"
}

Response contains `token` (JWT). Add header to protected requests:

Authorization: Bearer <token>

3) Create Task (POST /tasks/create)

{
	"title": "Buy milk",
	"description": "Two liters of skimmed milk",
	"category": "shopping",
	"deadline": "2025-12-10T15:00:00Z"
}

4) Share Task (POST /tasks/share/:id)

{
	"user": "bob"   // or "bob@example.com"
}

5) Complete Task (POST /tasks/complete/:id)
- no body required; requester must be owner or included in shared_with

---

## Testing locally (without Docker)

Ensure environment variables are set, then run:

```bash
go run ./
```

Then use Insomnia/Postman against `http://localhost:8080`.

---

## Contributing / Next steps

- Add unit/integration tests for auth and task authorization flows.
- Add a frontend.
- Add role-based admin support (admin can manage all tasks).

---

If you want, I can also generate an Insomnia collection JSON you can import directly.