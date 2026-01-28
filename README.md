AuthService

AuthService is a high-performance microservice for user registration and authentication, written in Go (Golang). It follows Clean Architecture principles and utilizes a REST API, PostgreSQL database, JWT authentication, and Docker containerization.
🚀 Features

    User Registration: Secure creation of users with password hashing (bcrypt).

    Authentication: Login mechanism returning JWT (JSON Web Tokens).

    Protected Routes: Middleware to protect endpoints using Bearer tokens.

    Clean Architecture: Separation of concerns (Handlers -> Services -> Storage).

    Structured Logging: Uses Go's log/slog for JSON and text logging.

    Configuration: Flexible config via YAML and Environment Variables.

    Dockerized: Ready-to-deploy with Docker and Docker Compose.

    Migrations: Automatic database initialization.

🛠 Tech Stack

    Language: Go 1.22+

    Router: chi

    Database: PostgreSQL

    Driver:

    Config:

    JWT:

    Validation: 

📂 Project Structure
⚙️ Configuration

The application is configured using config/local.yaml. Values can be overridden by environment variables (useful for Docker).
🚀 Getting Started
Option 1: Run with Docker (Recommended)

This is the easiest way to run the service and the database together.

    Clone the repository.

    Run Docker Compose:

    The service will start on port 8080. The database will be initialized automatically using the migration files.

Option 2: Run Locally

Prerequisites:

    Go 1.22+ installed.

    PostgreSQL running locally.

    Create the database: Create a database named PulseUserDB (or update config/local.yaml to match your DB).

    Apply Migrations: Run the SQL commands found in migrations/1_init.up.sql in your database manager.

    Install Dependencies:

    Run the App:

📡 API Endpoints
1. Register User

Create a new user account.

    URL: /register

    Method: POST

    Body:

    Response (200 OK):

2. Login

Authenticate and receive a JWT access token.

    URL: /login

    Method: POST

    Body:

    Response (200 OK):

3. Get User Info (Protected)

An example of a protected route requiring a token.

    URL: /me

    Method: GET

    Headers:

    Response (200 OK):

🧪 Testing

To run the unit tests (e.g., for the Register handler):
📄 License

This project is licensed under the MIT License. See the file for details.
