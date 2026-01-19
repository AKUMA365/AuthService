AuthService

AuthService is a microservice for user authentication and registration written in Go. The project utilizes REST API, a PostgreSQL database, and structured logging.
🛠 Tech Stack

    Language: Go (Golang) 1.22+

    Router:

    Database: PostgreSQL

    DB Driver:

    Configuration: (YAML + ENV)

    Logging: log/slog (Structured Logging)

    JWT:

    Validation:

    Testing: testing, testify

📂 Project Structure

The project follows the standard Go project layout:
🚀 Installation & Run
1. Prerequisites

    Go installed (version 1.21 or higher)

    PostgreSQL (running locally or in Docker)

2. Database Setup

Create a database and the users table. Example SQL:
3. Configuration Setup

Ensure the config/local.yaml file is configured correctly:
4. Running the Application

Note: If you have updated main.go to accept flags, the command might look like: go run cmd/authservic/main.go --config=./config/local.yaml
🔌 API Endpoints
Register User

URL: /register (or /auth/register depending on router setup) Method: POST

Request Body (JSON):

    email: Required, must be a valid email address.

    password: Required, minimum 8 characters.

Success Response (200 OK):

Error Responses:

    400 Bad Request: Invalid data format or validation failed.

    409 Conflict: User with this email already exists.

    500 Internal Server Error: Server-side error.

Login (Authorization)

(Logic is implemented, waiting to be connected to router)

URL: /login Method: POST

Request Body (JSON):

Success Response (200 OK):
🧪 Testing

The project includes unit tests for handlers. To run them:
📝 License

This project is distributed under the MIT License. See the LICENSE file for details.
