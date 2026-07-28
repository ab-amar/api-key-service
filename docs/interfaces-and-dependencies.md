# Interfaces And Dependencies

Interfaces help separate layers of the application.

Useful dependency flow:

- handler depends on service interface
- service depends on repository interface

Benefits:

- easier unit testing with fakes
- cleaner separation of concerns
- storage implementation can change without rewriting business logic
