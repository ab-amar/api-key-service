# Schema And Index Notes

Likely fields for an API key record:

- `id`
- `name`
- `key_hash`
- `created_at`
- `updated_at`
- `revoked_at`

Important schema decisions:

- primary key on `id`
- unique constraint on `key_hash`
- index for key lookup path

These choices support correctness and lookup performance.
