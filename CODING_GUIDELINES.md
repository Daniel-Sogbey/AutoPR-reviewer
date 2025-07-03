# Go Coding Guidelines

1. Avoid magic numbers. Use named constants.
2. Don’t use `fmt.Println` in production. Use a logger.
3. Don’t ignore returned errors.
4. Add context to errors using `fmt.Errorf`.
5. Extract duplicate logic into reusable functions.
6. Avoid using `panic` except in `main()`.
7. Validate user input and handle edge cases.
8. Avoid hard-coded credentials or tokens.
9. Avoid global mutable state.
