### Overview

- The client/frontend (application) and server/backend (PostgreSQL database) communicate at the TCP level.
- The messages are passed in the form of bytes.
- Generally, all the messages, sent by either the client or the server, are in the following format:
  - single byte to indicate the type of the message
  - int32 (4 bytes) to indicate the length the message (n), including self
  - the actual message, of length n - 4 bytes
- The start-up message, sent by the client, doesn't have the first single byte for historical reasons, as per the [PostgreSQL docs](https://www.postgresql.org/docs/18/protocol-overview.html#PROTOCOL-MESSAGE-CONCEPTS).

### Chapters

1. [Authentication](./02-authentication.md)

