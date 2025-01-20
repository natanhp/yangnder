# Yangder Technical Documentation
Yang (/jeɪŋ/) means couple in Javanese slang. Yangder is a couple finder application made to match you with your potential mate.

## Github Repository URL
https://github.com/natanhp/yangnder

## Requirements
### Functional
1. User can Sign-up to the app.
2. User can Login to the app.
3. User can do the action of view and swipe up to 10 profiles a day.
4. The same profiles must not appear more than once in the same day.
5. User can purchase premium packages to add a verified label to their profile.

### Non-functional
1. Password with Argon2id with minimum configuration of 19 MiB of memory, an iteration count of 2, and 1 degree of parallelism as recommended by OWASP.

## Tech Stacks
1. Go as the programming language because it is a simple yet fast language with huge community support. Also because the test requirements prefers Go and I want to learn it too.
2. Gin as the web framework because it is simple but loaded with complete features.
3. SQLite as the database because it requires no server and suited for a simple application.