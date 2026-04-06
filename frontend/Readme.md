# CodeForge — Frontend

React/Next.js frontend for the CodeForge competitive programming platform.

## Tech Stack

- **Framework**: Next.js 15 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS v4
- **Editor**: Monaco Editor (`@monaco-editor/react`)
- **Icons**: Lucide React

## Prerequisites

- Node.js 18+
- The backend judge service running on `http://localhost:7000` (see `../judge/`)

## Getting Started

Install dependencies:
```bash
npm install
```

Run the development server:
```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Project Structure

```
src/
├── app/
│   ├── page.tsx                  # Problems dashboard (home)
│   ├── layout.tsx                # Root layout with Navbar
│   ├── globals.css
│   └── problems/
│       └── [id]/
│           └── page.tsx          # Problem detail + IDE
├── components/
│   ├── Navbar.tsx
│   └── CodeEditor.tsx            # Monaco editor with verdict polling
└── lib/
    └── api.ts                    # Judge API client
```

## API Integration

All requests go to the judge master node. The base URL is configured in `src/lib/api.ts`:

```ts
export const API_BASE = "http://localhost:7000/api";
```

Update this to the production URL when deploying.

## Supported Runtimes

| Label     | Runtime ID |
|-----------|------------|
| Python 3  | `python3`  |
| C++ 17    | `c++17`    |
| C++ 20    | `c++20`    |
| Node.js   | `node-25`  |
