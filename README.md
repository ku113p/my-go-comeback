# Go Transition Plan

**Goal:** Successfully transition from Python to Go for professional development, aiming for a Middle Go Developer position.  
**Schedule:** Approx. 2 hours +/- 1 hour on weekdays. Weekends dedicated to Rust pet projects.  
**Status:** Not Started

---

## Phase 1: Refreshing Go Fundamentals (Estimate: 1-3 weeks)

*Goal: Confidently handle basic syntax and key concepts.*

- [ ] **Go Tour / Go by Example:** Complete to quickly recap syntax and standard constructs.
- [ ] **Basic Data Types:** Review (int, float, string, bool).
- [ ] **Data Structures:**
    - [ ] Slices: Creation, `len()`, `cap()`, `append`, pitfalls (modifying the underlying array).
    - [ ] Maps: Creation, adding, getting, deleting, checking key existence.
- [ ] **Control Flow:** `if/else`, `switch`, `for` (various forms).
- [ ] **Functions:** Definition, multiple return values, named results.
- [ ] **Structs:** Definition, embedding, methods for structs.
- [ ] **Pointers:** Syntax (`*`, `&`), when to use them.
- [ ] **Interfaces:**
    - [ ] Defining interfaces.
    - [ ] Implicit implementation.
    - [ ] The empty interface (`interface{}`).
    - [ ] Type Assertion and Type Switch.
- [ ] **Error Handling:**
    - [ ] The `if err != nil` pattern.
    - [ ] `errors.New`, `fmt.Errorf`.
- [ ] **Modules and Dependencies:**
    - [ ] `go mod init`, `go.mod`, `go.sum`.
    - [ ] `go get`, `go mod tidy`.
    - [ ] Create a test project with 1-2 external dependencies.
- [ ] **Concurrency (Basics):**
    - [ ] Starting goroutines (`go func()`).
    - [ ] Channels (creation, sending, receiving, buffering).
    - [ ] `select`.
    - [ ] `sync.WaitGroup`.
    - [ ] `sync.Mutex` / `sync.RWMutex` (basic usage).
- [ ] **Generics (Go 1.18+):**
    - [ ] Syntax (`[T any]`, `[T comparable]`).
    - [ ] Write a simple generic function/type.
- [ ] **Standard Library (Basics):** Review usage examples for `fmt`, `net/http`, `encoding/json`, `os`, `io`, `time`.
- [ ] **Testing (Basics):**
    - [ ] Writing a simple test (`_test.go` files).
    - [ ] `go test`.

**--- CHECKPOINT: Go Fundamentals Refreshed ---**

## Phase 2: Practice Through Projects (Estimate: 3-8 weeks)

*Goal: Apply knowledge by building typical backend applications.*

- **Project 1: CLI Utility (Simple)**
    - [ ] Define functionality (e.g., TODO list, file parser).
    - [ ] Implement core logic.
    - [ ] Handle command-line arguments (`flag` package).
    - [ ] Read/write files (`os`, `io`, `encoding/json` or `encoding/csv`).
    - [ ] Write basic tests.
- **Project 2: Simple REST API (In-memory) (Medium)**
    - [ ] Choose an entity (e.g., books, notes).
    - [ ] Implement CRUD handlers (`net/http` or Gin/Echo).
    - [ ] Handle JSON requests/responses.
    - [ ] Basic project structure (e.g., handlers, models).
    - [ ] Write tests for handlers (using `net/http/httptest`).
- **Project 3: API with Database (Medium+)**
    - [ ] Choose a DB (e.g., PostgreSQL, SQLite).
    - [ ] Set up DB connection (`database/sql`).
    - [ ] Replace in-memory storage with DB operations (CRUD).
    - [ ] Handle DB errors.
    - [ ] (Optional) Explore `sqlx` or GORM.
    - [ ] Write integration tests (require a running DB).
- **Project 4: Concurrent Processing (Medium+)**
    - [ ] Choose a task (e.g., web crawler, task processor).
    - [ ] Design concurrent logic (goroutines, channels).
    - [ ] Implement using `sync.WaitGroup`, channels, `select`.
    - [ ] Ensure safe access to shared data (`Mutex`).
    - [ ] Write tests, possibly checking for race conditions (`go test -race`).

**--- CHECKPOINT: 2-3 Working Projects Created ---**

## Phase 3: Interview Preparation (Go Specific) (Ongoing during Phases 2 & 4)

*Goal: Be ready for Go-specific technical questions.*

- [ ] **Concurrency (Deeper Dive):**
    - [ ] Race conditions (`go test -race`).
    - [ ] Deadlocks (understanding).
    - [ ] Channel patterns.
    - [ ] `context` package (understand usage, apply in API project).
- [ ] **Interfaces:** Deep understanding, use cases.
- [ ] **Error Handling:** `errors.Is`, `errors.As`, wrapping.
- [ ] **GC & Memory Management:** General understanding. Stack vs Heap.
- [ ] **Profiling:** Know about `pprof`.
- [ ] **Project Structure:** Idiomatic approaches.
- [ ] **Gotchas:** Common mistakes (loop variables & goroutines, slice append issues).
- [ ] **Go vs Python:** Articulate your reasons for switching and Go's benefits.
- [ ] **Problem Solving:** Try solving some problems (e.g., LeetCode Easy/Medium) using Go.

**--- CHECKPOINT: Confident Understanding of Key Interview Topics ---**

## Phase 4: Polishing GitHub Portfolio (Final Phase)

*Goal: Prepare 1-2 projects to showcase to potential employers.*

- [ ] **Select 1-2 best projects** from Phase 2.
- [ ] **Refactoring & Code Cleanup:** Improve readability, structure.
- [ ] **Write `README.md`:** Clear description, setup instructions, tech used.
- [ ] **Tests:** Ensure tests exist and pass.
- [ ] **Dependencies:** Check `go.mod`, `go.sum`.
- [ ] **(Optional) Dockerfile:** Add for easy execution.
- [ ] **Commit History:** Clean up if necessary.
- [ ] **Upload/Update on GitHub.**

**--- CHECKPOINT: 1-2 Showcase-Ready Projects on GitHub ---**

## Notes and Reminders

- **Consistency > Intensity:** Try to stick to the schedule, but avoid burnout.  
- **Don't Burn Out:** Take breaks when needed.  
- **Rust:** Remember to allocate time for Rust on weekends.  
- **Flexibility:** This plan is a guideline. Adjust based on your feelings and progress.  
- **Read Code:** Study code from the standard library and popular Go projects.  
- **Community:** Don't hesitate to search Stack Overflow, Go docs, forums for answers.  

---
