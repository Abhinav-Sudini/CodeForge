# CodeForge

CodeForge is an online judge platform for competitive programming practice. It has a Next.js frontend, a Go master service, Go worker services, PostgreSQL storage, and Dockerized runtime containers for compiling and running submitted code.

The main project goal is the complete judge flow:

1. Browse problems from the frontend.
2. Open a problem and write code in the Monaco editor.
3. Submit code with a selected runtime.
4. Store the submission in PostgreSQL.
5. Queue the job on the master service.
6. Send the job to a free worker for that runtime.
7. Compile/run the code against test cases.
8. Store and display the final verdict.

## Tech Stack

- **Frontend:** Next.js, React, TypeScript, Tailwind CSS, Monaco Editor, KaTeX
- **Backend:** Go
- **Database:** PostgreSQL
- **Containerization:** Docker and Docker Compose
- **Judge runtimes:** C++17, C++20, C++23, C17, Python 3, Node/JavaScript

## Repository Layout

```text
CodeForge/
|-- frontend/              # Next.js user interface
|-- judge/
|   |-- master/            # API server, scheduler, database handlers
|   |-- worker/            # Worker server and compile-run-judge logic
|   |-- Dockerfiles/       # Runtime-specific worker images
|   |-- CSES_Scraper/      # Tools for collecting CSES problem/test data
|   `-- Makefile           # Judge build helpers
|-- db/                    # PostgreSQL environment file
|-- compose.yaml           # Local compose setup
|-- prod.compose.yaml      # Production-like compose setup
|-- Makefile               # Root build/run helpers
`-- CodeForge_Project_Report.tex
```

## Architecture

CodeForge is split into three working layers.

### Frontend

The frontend lives in `frontend/`. It provides:

- problem dashboard with search, filters, sorting, and pagination
- problem detail page
- split problem/editor workspace
- Monaco-based code editor
- runtime selection
- submission polling and verdict display
- previous submission display for a problem

The frontend API base URL is currently defined in:

```text
frontend/src/lib/api.ts
```

For local development, update `API_BASE` to point to the local master service, usually:

```ts
export const API_BASE = "http://localhost:7000/api";
```

### Master Service

The master service lives in `judge/master/`. It is responsible for:

- serving user-facing judge APIs
- validating submitted question IDs and runtimes
- creating submission records
- keeping per-runtime job queues
- tracking available workers
- dispatching jobs to workers
- receiving verdicts from workers
- updating verdict and runtime statistics in PostgreSQL

The master listens on port `7000` inside the container.

### Worker Service

The worker service lives in `judge/worker/`. Each worker runs for one runtime. A worker:

- registers itself with the master
- accepts a job only when it is free
- writes submitted code into its working directory
- compiles the code if the runtime needs compilation
- runs the code against the problem test cases
- compares actual output with expected output
- reports verdict, time, memory, stderr, and wrong-answer details back to the master

Workers expose:

```text
POST /judge/
GET  /healthcheck/
```

## Supported Runtimes

The master and frontend currently use these runtime IDs:

| Runtime ID | Language |
| --- | --- |
| `c++17` | C++ 17 |
| `c++20` | C++ 20 |
| `c++23` | C++ 23 |
| `gcc-c17` | C 17 |
| `python3` | Python 3 |
| `node-25` | JavaScript / Node |

## Database

The main PostgreSQL schema is in:

```text
judge/master/db/schema.sql
```

Main tables:

- `users`
- `questions`
- `submissions`
- `verdict_stats`

The master service stores a submission first with verdict `queued`. Once the worker finishes, the master updates the submission verdict and inserts the execution stats into `verdict_stats`.

## API Overview

User-facing APIs are served under `/api`.

| Method | Route | Purpose |
| --- | --- | --- |
| `GET` | `/api/question/` | List all problems |
| `GET` | `/api/question/{q_id}/` | Get one problem |
| `POST` | `/api/judge/` | Submit code |
| `GET` | `/api/submissions/{submission_id}/` | Fetch verdict for one submission |
| `GET` | `/api/question/submissions/{q_id}/` | Fetch submissions for one question |

Worker/internal APIs:

| Method | Route | Purpose |
| --- | --- | --- |
| `POST` | `/api/worker/register/` | Register a free worker with the scheduler |
| `POST` | `/api/worker/verdict/` | Receive final verdict from a worker |

Example submission body:

```json
{
  "QuestionId": 1068,
  "runtime": "c++23",
  "code": "#include <iostream>\nusing namespace std;\n\nint main() {\n    cout << \"Hello\";\n    return 0;\n}"
}
```

## Setup

### Prerequisites

- Docker
- Docker Compose
- Go `1.25.3` or compatible
- Node.js and npm
- GNU Make, if using the provided Makefiles

### Environment Files

Runtime `.env` files are ignored by git, so a fresh clone needs them to be created before using Docker Compose.

The database env file already exists at:

```text
db/env.dev
```

It contains:

```env
POSTGRES_DB=dev_db
POSTGRES_PASSWORD=dev_password
```

Create the master env file:

```text
judge/Dockerfiles/master/master.env
```

Example:

```env
PG_DATABASE_URL=postgres://postgres:dev_password@database:5432/dev_db
QUESTIONS_DIR=/questionsDir
INTERNAL_API_KEY=dev_internal_key
```

Create one env file per worker runtime.

Example for C++23:

```text
judge/Dockerfiles/cpp23/cpp.env
```

```env
WORKER_RUNTIME=c++23
QUESTIONS_DIR=/questionsDir
CODE_FILE_DIR=/codeDir
MASTER_IP_ADDR=http://judge.master
MASTER_PORT_CF=7000
WORKER_MAX_COMPILATION_TIME=1000
INTERNAL_API_KEY=dev_internal_key
```

Use the same shape for the other worker env files and only change `WORKER_RUNTIME`.

Expected env file paths from `compose.yaml`:

```text
judge/Dockerfiles/cpp/cpp.env
judge/Dockerfiles/cpp20/cpp.env
judge/Dockerfiles/cpp23/cpp.env
judge/Dockerfiles/c/c.env
judge/Dockerfiles/python/python.env
judge/Dockerfiles/node/node.env
```

Runtime values:

```text
c++17
c++20
c++23
gcc-c17
python3
node-25
```

### Test Case Data

The judge expects problem test cases under:

```text
judge/scraped_questions/
```

For each problem, the worker expects a directory named by question ID:

```text
judge/scraped_questions/1068/
|-- 1.in
|-- 1.out
|-- 2.in
|-- 2.out
`-- ...
```

This directory is gitignored because test data can be large. The scraper tools under `judge/CSES_Scraper/` are included for collecting CSES problem data and test files.

## Running the Project

### 1. Build and start judge services

From the repository root:

```bash
make judge
```

This runs the judge build step and then starts Docker Compose.

Equivalent manual flow:

```bash
cd judge
make build
cd ..
docker compose up --build
```

The master API should be available at:

```text
http://localhost:7000/api
```

### 2. Run database migration

The migration command is inside `judge/master/`.

```bash
cd judge/master
PG_DATABASE_URL=postgres://postgres:dev_password@localhost:5432/dev_db go run cmd/migrate_db/migrate.go
```

On Windows PowerShell:

```powershell
cd judge/master
$env:PG_DATABASE_URL="postgres://postgres:dev_password@localhost:5432/dev_db"
go run cmd/migrate_db/migrate.go
```

### 3. Seed sample data

Optional dev seed:

```bash
cd judge/master
PG_DATABASE_URL=postgres://postgres:dev_password@localhost:5432/dev_db go run cmd/seed_dev_db/seed.go
```

PowerShell:

```powershell
cd judge/master
$env:PG_DATABASE_URL="postgres://postgres:dev_password@localhost:5432/dev_db"
go run cmd/seed_dev_db/seed.go
```

### 4. Run the frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend runs at:

```text
http://localhost:3000
```

Before testing locally, update `frontend/src/lib/api.ts` so `API_BASE` points to the local master API:

```ts
export const API_BASE = "http://localhost:7000/api";
```

## Development Commands

Root commands:

```bash
make judge          # Build judge binaries and start local compose
make restart_master # Rebuild and restart only the master container
make judge_prod     # Run prod.compose.yaml
```

Judge commands:

```bash
cd judge
make build          # Build master and worker binaries
make it_cpp_dev     # Build and run a C++ worker dev container
make master_it_dev  # Build and run master dev container
```

Worker commands:

```bash
cd judge/worker
make prod_bin       # Build worker server and execution binary
make worker         # Run worker server directly
```

Master commands:

```bash
cd judge/master
make prod_bin       # Build master binary
make run            # Run master directly
```

Frontend commands:

```bash
cd frontend
npm run dev
npm run build
npm run lint
```

## Verdicts

The judging code maps execution results into verdict messages such as:

- `Accepted`
- `Wrong Answer`
- `Compilation Error`
- `Time Limit Exceeded`
- `Memory Limit Exceeded`
- `Internal Server Error`

For wrong answers, the response can include the failed test case number, expected output, actual output, and stderr.

## Notes for Evaluators

- The main backend for the project is under `judge/`; the top-level `backend/` directory currently only contains a README placeholder.
- The frontend may point to a hosted API until `frontend/src/lib/api.ts` is changed for local development.
- The Docker Compose setup depends on env files and testcase folders that are intentionally not committed.
- The important implemented path is the online judge path: problem listing, submission, scheduling, worker execution, verdict persistence, and verdict polling.

## Team Members

- Abhinav Sudini (2301CS03)
- Ankesh Kumar (2301CS06)
- Aarsh Sanghavi (2301CS01)
