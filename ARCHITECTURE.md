# LogAggregator — Architecture & Design

## System Overview

```mermaid
graph TB
    subgraph CLI["Command Line Interface"]
        MAIN["cmd/logmerge/main.go<br/>Entry Point & Flag Parsing"]
    end

    subgraph CORE["Core Pipeline (internal/)"]
        SRC["source/<br/>File Reader"]
        FANIN["fanin/<br/>Fan-In Merge"]
        PARSER["parser/<br/>Line Parser"]
        FILTER["filter/<br/>Level & Keyword"]
        ORDER["ordering/<br/>Time Sort"]
        OUT["output/<br/>Formatter"]
        MODELS["models/<br/>LogEvent, Storage"]
    end

    subgraph WEB["Web Server (web/)"]
        SRV["server.go<br/>HTTP Server + Routing"]
        AUTH_MW["Auth Middleware<br/>requireAuth()"]
        API["REST API<br/>/api/logs, /api/stats"]
        STATIC["Static Files<br/>index.html, login.html"]
    end

    subgraph SEC["Security (internal/auth/)"]
        BCRYPT["UserStore<br/>bcrypt Hashing"]
        SESSION["SessionManager<br/>Token Sessions"]
    end

    subgraph CFG["Configuration"]
        CONFIG["config/<br/>JSON Config + Export"]
    end

    MAIN -->|"--web flag"| SRV
    MAIN -->|"CLI mode"| SRC
    SRV --> SRC

    SRC --> FANIN
    FANIN --> PARSER
    PARSER --> MODELS
    PARSER --> FILTER
    FILTER --> ORDER
    ORDER --> OUT

    SRV --> AUTH_MW
    AUTH_MW --> BCRYPT
    AUTH_MW --> SESSION
    AUTH_MW --> API
    API --> FILTER
    API --> MODELS
    SRV --> STATIC

    MAIN --> CONFIG

    style CLI fill:#1e293b,stroke:#3b82f6,color:#e2e8f0
    style CORE fill:#1e293b,stroke:#10b981,color:#e2e8f0
    style WEB fill:#1e293b,stroke:#8b5cf6,color:#e2e8f0
    style SEC fill:#1e293b,stroke:#ef4444,color:#e2e8f0
    style CFG fill:#1e293b,stroke:#f59e0b,color:#e2e8f0
```

---

## Concurrent Data Flow Pipeline

```mermaid
flowchart LR
    subgraph FILES["Log Files"]
        F1["app.log"]
        F2["server.log"]
        F3["error.log"]
    end

    subgraph GOROUTINES["Goroutines (1 per file)"]
        G1["goroutine 1"]
        G2["goroutine 2"]
        G3["goroutine 3"]
    end

    subgraph FANIN["Fan-In (WaitGroup)"]
        MERGED["merged channel"]
    end

    subgraph PROCESS["Processing"]
        PARSE["Parser<br/>raw → LogEvent"]
        FILT["Filter Chain<br/>Level + Keyword"]
        SORT["Time Ordering<br/>Buffer + Sort"]
    end

    subgraph OUTPUT["Output"]
        TERM["Terminal<br/>Color-coded"]
        FILE["File Export<br/>.log"]
        JSON["JSON Export<br/>.json"]
        WEBAPI["Web API<br/>/api/logs"]
    end

    F1 --> G1
    F2 --> G2
    F3 --> G3
    G1 -->|"chan string"| MERGED
    G2 -->|"chan string"| MERGED
    G3 -->|"chan string"| MERGED
    MERGED --> PARSE
    PARSE --> FILT
    FILT --> SORT
    SORT --> TERM
    SORT --> FILE
    SORT --> JSON
    SORT --> WEBAPI

    style FILES fill:#0f172a,stroke:#64748b,color:#e2e8f0
    style GOROUTINES fill:#1e293b,stroke:#3b82f6,color:#e2e8f0
    style FANIN fill:#1e293b,stroke:#f59e0b,color:#e2e8f0
    style PROCESS fill:#1e293b,stroke:#10b981,color:#e2e8f0
    style OUTPUT fill:#1e293b,stroke:#8b5cf6,color:#e2e8f0
```

---

## Web Authentication Flow

```mermaid
sequenceDiagram
    participant U as Browser/UI
    participant S as HTTP Server
    participant MW as Auth Middleware
    participant US as UserStore (bcrypt)
    participant SM as SessionManager

    U->>S: GET /login
    S-->>U: login.html

    U->>S: POST /login (username, password)
    S->>US: Authenticate(user, pass)
    US->>US: bcrypt.CompareHashAndPassword()
    US-->>S: true/false

    alt Valid Credentials
        S->>SM: CreateSession(username)
        SM->>SM: crypto/rand → 32-byte token
        SM-->>S: session token
        S-->>U: Set-Cookie: session_token (HttpOnly)
        S-->>U: 302 Redirect → /
    else Invalid Credentials
        S-->>U: 302 Redirect → /login?error=1
    end

    U->>S: GET /api/logs (Cookie: session_token)
    S->>MW: requireAuth()
    MW->>SM: ValidateSession(token)
    SM-->>MW: username, valid
    MW-->>S: proceed
    S-->>U: JSON response

    U->>S: GET /logout
    S->>SM: DestroySession(token)
    S-->>U: Clear cookie + redirect
```

---

## Project Directory Structure

```
LogAggregator/
├── cmd/logmerge/
│   └── main.go              ← Entry point, CLI flags, pipeline orchestration
├── internal/
│   ├── auth/                ← 🔐 bcrypt hashing + session management
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── config/              ← JSON config loading + export
│   ├── fanin/               ← Fan-in channel merge pattern
│   ├── filter/              ← Interface-based filter chain (Level, Keyword)
│   ├── models/              ← LogEvent, LogLevel, LogStorage, Buffer
│   ├── ordering/            ← Time-based sorting
│   ├── output/              ← Terminal formatting, file writing
│   ├── parser/              ← Raw line → LogEvent parser
│   ├── pointers/            ← Pointer demonstrations
│   └── source/              ← ⚡ Concurrent file reader (goroutines + channels)
├── web/
│   ├── server.go            ← HTTP server, routing, middleware
│   └── static/
│       ├── index.html       ← Dashboard UI (stats, filters, log table)
│       └── login.html       ← Authentication page
└── logs/                    ← Sample log files
```

---

## Key Design Patterns

| Pattern | Where | Purpose |
|---------|-------|---------|
| **Fan-In** | `source/` | Merge N file channels → 1 output channel |
| **Pipeline** | `main.go` | Source → Parse → Filter → Order → Output |
| **Middleware** | `server.go` | `requireAuth()` wraps handlers |
| **Interface Polymorphism** | `filter/`, `parser/`, `source/` | Swappable implementations |
| **Observer (channel)** | `source/` | Goroutines emit lines, main loop consumes |
| **Strategy** | `filter/` | `ChainFilters()` composes filter strategies |
