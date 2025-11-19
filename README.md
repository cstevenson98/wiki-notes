# Wiki Notes

A full-stack wiki-style note-taking application built with SvelteKit, Go, and PostgreSQL.

## Features

- **Markdown Support**: Write your notes in markdown format
- **LaTeX Equations**: Support for inline ($...$) and block ($$...$$) mathematical equations
- **Wiki Links**: Create links between pages using `[[Page Name]]` syntax
- **Search**: Search through all pages by name
- **Minimalistic Design**: Clean black and white interface with focused content layout

## Tech Stack

### Frontend
- SvelteKit 5
- TipTap rich text editor
- Tailwind CSS
- KaTeX for LaTeX rendering

### Backend
- Go with Gin framework
- PostgreSQL database

### Deployment
- Docker & Docker Compose

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Node.js 20+ (for local development)
- Go 1.21+ (for local development)

### Running with Docker (Recommended)

1. Clone the repository:
```bash
git clone <repo-url>
cd wiki-notes
```

2. Start all services:
```bash
docker-compose up --build
```

3. Access the application:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432

### Local Development

#### Backend

```bash
cd backend
go mod download
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/wikinotes?sslmode=disable"
go run main.go
```

#### Frontend

```bash
cd frontend
npm install
npm run dev
```

#### Database

You'll need a running PostgreSQL instance. You can use the Docker one:

```bash
docker-compose up db
```

Then initialize the database with the schema:

```bash
psql -h localhost -U postgres -d wikinotes -f backend/init.sql
```

## API Endpoints

- `GET /api/pages` - List all pages
- `GET /api/page/:id` - Get a single page by ID
- `GET /api/page/by-name/:name` - Get a page by name
- `POST /api/page` - Create a new page
- `PATCH /api/page/:id` - Update a page
- `DELETE /api/page/:id` - Delete a page

## Database Schema

```sql
CREATE TABLE pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## Usage

1. The application starts with a default "Home" page
2. Click "Edit" to modify any page
3. Use `[[Page Name]]` to create links to other pages
4. Clicking a blue link navigates to an existing page
5. Clicking a red link prompts you to create that page
6. Use the search bar at the top to find pages by name

## License

MIT

