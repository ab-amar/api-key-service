# Persistence Notes

Storage decisions matter because API keys are security-sensitive records.

What persistence needs to support:

- create key metadata
- find keys by hash
- revoke keys
- track timestamps

An in-memory repository is useful for:

- early learning
- fast tests

PostgreSQL is useful for:

- durability
- indexing
- real integration testing
