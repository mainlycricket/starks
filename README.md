### Introduction

- **_Starks_** is a PostgreSQL driver written in Go.
- The project is currently in an early stage of development.

### Goals

- To implement a PostgreSQL driver compatible with Go’s `database/sql` package, with minimal external dependencies
- To understand the internal workings of database drivers and the PostgreSQL protocol
- This project is primarily for learning and exploration, and is not intended to be a production-grade driver

### Roadmap

- Authentication ✅
- Simple SQL Commands ⏳
- Database-level Transactions ⏳
- Concurrency support, context cancellation ⏳
- Error handling ⏳

### Theoretical Docs

Detailed documentation is in the `notes` directory:

1. Overview [here](./notes/01-overview.md)
2. Authentication [here](./notes/02-authentication.md)

### Code

- The driver implementation is in the `internal` directory.
- The working of the driver is tested through the application in `cmd/app`.
