
# Notes API

This is a simple RESTful API for managing notes and authors. It provides basic CRUD (Create, Read, Update, Delete) operations for these resources.

## Overview

The API is built using the Go programming language and leverages the standard `net/http` package for handling HTTP requests. It uses `github.com/google/uuid` for generating unique identifiers and `encoding/json` for handling JSON encoding and decoding. Logging is done using the `log/slog` package.

The API defines two main resources:

-   **Notes**: Represent individual notes with attributes like name, description, creation timestamp, and the ID of the author.
-   **Authors**: Represent users who can create notes, with attributes like username, first name, and second name.

## API Endpoints

The following endpoints are available:

### Ping

-   `GET /ping`: A simple health check endpoint that returns a 200 OK status.

### Notes

-   `GET /notes`: Lists all notes.
-   `POST /notes`: Creates a new note. The request body should be a JSON object representing the note.
-   `GET /notes/{id}`: Retrieves a specific note by its ID (UUID).
-   `PUT /notes/{id}`: Updates an existing note. The request body should be a JSON object representing the updated note.
-   `DELETE /notes/{id}`: Deletes a specific note by its ID (UUID).

### Authors

-   `GET /authors`: Lists all authors.
-   `POST /authors`: Creates a new author. The request body should be a JSON object representing the author.
-   `GET /authors/{id}`: Retrieves a specific author by their ID (UUID).
-   `PUT /authors/{id}`: Updates an existing author. The request body should be a JSON object representing the updated author.
-   `DELETE /authors/{id}`: Deletes a specific author by their ID (UUID).

**Path Parameter:**

-   `{id}`: Represents the UUID of the specific resource (note or author).

**Request and Response Format:**

-   All request and response bodies are in JSON format.

## Running the Application

1.  **Prerequisites:** Make sure you have Go installed on your system.
2.  **Clone the repository** (if you haven't already).
3.  **Navigate to the project directory** in your terminal.
4.  **Run the application:**

    ```bash
    go run main.go
    ```

    The server will start and listen on `http://127.0.0.1:8090`.

## Example Usage (using `curl`)

### Get all notes

```bash
curl [http://127.0.0.1:8090/notes](http://127.0.0.1:8090/notes)
```

### Create a new note

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name": "My First Note", "description": "This is a test note.", "author_id": "your-author-uuid"}' [http://127.0.0.1:8090/notes](http://127.0.0.1:8090/notes)
```

Replace `"your-author-uuid"` with an actual author ID.

### Get a specific author

```bash
curl [http://127.0.0.1:8090/authors/your-author-uuid](http://127.0.0.1:8090/authors/your-author-uuid)
```

Replace `"your-author-uuid"` with the ID of the author you want to retrieve.

## Error Handling

The API returns standard HTTP status codes to indicate the outcome of requests. Common error responses include:

-   `400 Bad Request`: Indicates that the request was malformed (e.g., invalid UUID).
-   `404 Not Found`: Indicates that the requested resource could not be found.
-   `405 Method Not Allowed`: Indicates that the HTTP method used is not supported for the given endpoint.
-   `500 Internal Server Error`: Indicates an unexpected error on the server. The error details are logged to the standard error output.

Error responses are typically returned as JSON objects with an `error` field and optionally a `message` field providing more context.

```json
{
  "error": "bad_request",
  "message": "invalid uuid"
}
```

## Data Storage

The API currently uses an in-memory storage mechanism (`models.Table`). This means that all data will be lost when the server stops. For a production environment, a persistent storage solution (like a database) would be necessary.

## Further Development

This is a basic implementation and can be extended with features such as:

-   Persistent data storage (e.g., using a database like PostgreSQL, MySQL, or SQLite).
-   Input validation to ensure data integrity.
-   More sophisticated error handling and logging.
-   Authentication and authorization to secure the API.
-   Pagination for listing large numbers of resources.
