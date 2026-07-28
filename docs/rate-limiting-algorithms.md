# Rate Limiting Algorithms

Common algorithms:

- fixed window
- sliding window
- token bucket

Token bucket is a good choice here because:

- it is practical
- it allows small bursts
- it is stronger than fixed window for interview discussion
- it is still simple enough for a small Go project
