# 🧰 Toolkit v2
A lightweight utility library for Go projects, designed to simplify logging, common tasks across JSON handling, file operations, validation, encryption, temporal types — and now production‑grade error handling.
Built for clarity, composability, and long‑term maintainability.


## ✨ Features

### ⚙️ Core Utilities
- [X] Logging - with structured logging for prod and easy to read logging for development
- [X] Faults — Structured Errors With Stack Traces
- [X] Read & write JSON
- [X] Error JSON responses with optional status codes
- [X] Unified error handling helper (HandleError)
- [X] File upload & static file download
- [X] Random string generation
- [X] HTTP JSON POST helper
- [X] XML writer
- [X] Directory creation utilities
- [X] URL‑safe slug generation
- [X] Validation helpers
- [X] Encryption & decryption utilities

### 📦 Installation

`go get -u github.com/fouched/toolkit/v2`

#### ⚠️ Error handling
Toolkit v2 includes a lightweight but powerful error system designed for real‑world services.
Key capabilities
- [X] Automatic stack capture at the point of failure
- [X] Context‑rich error wrapping (faults.Wrap)
- [X] Annotation without stack pollution (faults.Annotate)
- [X] Root‑cause extraction (faults.Root)
- [X] Stack inspection (faults.Stack)
- [X] Pretty stack formatting with %+v
- [X] Drop‑in compatibility with errors.Is and errors.As

#### Examples
```go
if err != nil {
    return faults.Wrap(err, "repo: failed to insert user")
}
```

To attach a stack to a foreign error:
```go
return faults.WithStack(err)
```

To add context without changing the origin:
```go
return faults.Annotate(err, "service: user creation failed")
```

🖨️ Pretty Logging Integration

Toolkit v2 includes a production and development‑friendly slog handler that automatically detect faults, error values and prints:
- [X] the full error chain
- [X] the captured stack trace
- [X] file + line + function for each frame

Sample usage of dev handler:
```go
logger := slog.New(logging.NewPrettyDevHandler())
```

Example output:

```markdown
ERROR 2026-04-02T14:14:29+02:00 failed to accept relationship request
err: repo: failed to insert relationship: ERROR: duplicate key...
stack:
    /internal/repo/relationship_repo.go:56  (*RelationshipRepo).Insert
    /internal/services/relationship_service.go:37  (*RelationshipService).Add
...
```


Sample usage of structure logging with a prod handler:
```go
// 1. Parse configured log level
lvl, err := logging.ParseLevel(app.Logging.Level)
if err != nil {
// Fallback to INFO but log the issue
fmt.Printf("Invalid log level %q, defaulting to INFO\n", app.Logging.Level)
    lvl = slog.LevelInfo
}

// 2. Set the global production level
logging.ProdLevel.Set(lvl)

// 3. Build the JSON handler with dynamic level
jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: logging.ProdLevel, 
})

// 4. Wrap it with ProdHandler
prodHandler := &logging.ProdHandler{
    Handler: jsonHandler,
}

logger = slog.New(prodHandler)
```

🕒 Temporal Types

Includes two production‑ready temporal primitives:
DateOnly
TimeOnly
Both provide:
- [X] JSON marshalling/unmarshalling
- [X] SQL scanning & value support
- [X] Nullable semantics
- [X] Formatting & comparison helpers

Designed to avoid zero‑value ambiguity while remaining ergonomic.

#### 📅 Temporal Types Philosophy
The DateOnly and TimeOnly types are designed to:
- Represent nullable date/time values
- Integrate seamlessly with JSON, SQL, and domain logic
- Avoid zero-value ambiguity by using *time.Time internally
- Provide comparison, conversion, and formatting utilities


