# API Key Generation Notes

API keys should be:

- hard to guess
- long enough to resist brute force
- generated from secure randomness

For this project, the right default is:

- generate the raw key with `crypto/rand`
- encode it into a transport-safe string

Avoid:

- predictable IDs
- timestamps
- math/rand for secret material
