# HTTP REST API Application

This project is a RESTful API built with Go, designed to handle user registration, authentication, and movie data management. It includes features like rate limiting, structured logging, email notifications, and graceful server shutdown.

## Features

**User Management:**

* User registration with email and password.
* Account activation via email with a token.
* Validation for user input.

**Movie Management:**

* Fetch movie details by ID.
* Handle errors like "record not found" gracefully.

**Rate Limiting:**

* Per-client rate limiting using `rate.Limiter`.

**Logging:**

* Structured JSON logging for errors, info, and fatal events.

**Email Notifications:**  

* Sends activation emails using the `go-mail` package. ####used Goroutines

**Graceful Shutdown:**

* Handles OS signals for clean server shutdown. ####used Goroutines
* Waits for background tasks to complete before exiting.

