# Hashing Vs Encryption For API Keys

API keys should usually be stored as hashes, not encrypted values.

Hashing:

- one-way
- good for verification
- safer if the database is exposed

Encryption:

- two-way
- useful only if the original value must be recovered

For API keys, the service normally needs to verify the presented key, not recover the stored raw key later.
