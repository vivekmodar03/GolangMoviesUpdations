# FilmStack Go

This project is a secure RESTful API built with Go for managing a personal movie collection. It uses a MySQL database for data storage and leverages Firebase Authentication for robust user management and API security.

Each user can register, log in, and then perform Create, Read, Update, and Delete (CRUD) operations on their own list of movies. The API ensures that users can only access and modify the data they own.

---

## ✨ Features

-   **User Authentication**: Secure user registration and login powered by Firebase.
-   **JWT Protection**: All movie-related endpoints are protected using JSON Web Tokens (JWT) issued by Firebase.
-   **CRUD Operations for Movies**: Full support for creating, reading, updating, and deleting movies.
-   **Data Scoping**: Users can only view and manage movies that they have personally created.
-   **MySQL Database**: Utilizes a MySQL database for persistent data storage.
-   **Structured Codebase**: The project follows a clear and modular structure, separating concerns like handlers, models, and routes.

---

## 🛠️ Technologies Used

-   **Go**: The core programming language.
-   **MySQL**: The relational database for storing movie data.
-   **Firebase Authentication**: For user management and token-based authentication.
-   **gorilla/mux**: A powerful URL router and dispatcher for Go.
-   **go-sql-driver/mysql**: The MySQL driver for Go.

---

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

-   [Go](https://golang.org/doc/install) (version 1.18 or higher)
-   [MySQL Server](https://dev.mysql.com/downloads/mysql/)
-   A [Firebase Project](https://console.firebase.google.com/)
-   [Postman](https://www.postman.com/downloads/) or a similar API testing tool.

---

## 🚀 Setup and Installation

Follow these steps to get the project running on your local machine.

### 1. Clone the Repository

```bash
git clone [https://github.com/your-username/go-movies-crud.git](https://github.com/your-username/go-movies-crud.git)
cd go-movies-crud
````

### 2\. Set up the MySQL Database

1.  Connect to your MySQL server.

2.  Create a new database for the project.

    ```sql
    CREATE DATABASE go_movie_db;
    ```

3.  Create the `movies` table within the `go_movie_db` database. This table includes a `user_id` column to associate movies with Firebase users.

    ```sql
    USE go_movie_db;

    CREATE TABLE movies (
        id INT AUTO_INCREMENT PRIMARY KEY,
        isbn VARCHAR(255),
        title VARCHAR(255),
        director_firstname VARCHAR(255),
        director_lastname VARCHAR(255),
        user_id VARCHAR(255) NOT NULL
    );
    ```

### 3\. Set up Firebase

1.  Go to the [Firebase Console](https://console.firebase.google.com/) and create a new project.
2.  In your Firebase project, go to **Authentication** \> **Sign-in method** and enable the **Email/Password** provider.
3.  Go to **Project settings** \> **Service accounts**.
4.  Click **"Generate new private key"** and download the resulting JSON file.
5.  **Rename** this file to `firebase-credentials.json` and place it in the **root directory** of the project.

### 4\. Configure Environment Variables

The application uses an environment variable to securely handle your Firebase Web API Key.

1.  In your Firebase project settings (under the **General** tab), find your **Web API Key**.

2.  Set it as an environment variable named `FIREBASE_API_KEY`.

      * **On macOS/Linux:**
        ```bash
        export FIREBASE_API_KEY="YOUR_API_KEY_HERE"
        ```
      * **On Windows (Command Prompt):**
        ```bash
        set FIREBASE_API_KEY="YOUR_API_KEY_HERE"
        ```

### 5\. Update Database Connection String

Open the `internal/app/db/mysql.go` file and ensure the connection string matches your MySQL setup (username, password, port, etc.).

```go
// internal/app/db/mysql.go
DB, err = sql.Open("mysql", "root:system@0987@tcp(127.0.0.1:3305)/go_movie_db")
```

### 6\. Install Dependencies and Run

```bash
# Tidy up the dependencies
go mod tidy

# Run the application
go run ./cmd
```

The server should now be running on `http://localhost:8080`.

-----

## 🔩 API Endpoints

Here is a detailed guide to all the available API endpoints.

### User Authentication

#### `POST /register`

Registers a new user in Firebase.

  - **Body:**

    ```json
    {
      "email": "testuser@example.com",
      "password": "password123"
    }
    ```

  - **Success Response (201 Created):**

    ```json
    {
      "message": "User created successfully",
      "uid": "FIREBASE_USER_ID"
    }
    ```

#### `POST /login`

Logs in a registered user and returns a Firebase ID token.

  - **Authorization:** Use **Basic Auth** with the user's registered email and password.

  - **Success Response (200 OK):**

    ```json
    {
      "idToken": "FIREBASE_ID_TOKEN",
      "email": "testuser@example.com",
      "refreshToken": "...",
      "expiresIn": "3600",
      "localId": "..."
    }
    ```

### Movie Management

**Note:** All movie endpoints require an `Authorization` header with a Bearer Token.

  - **Header:** `Authorization: Bearer FIREBASE_ID_TOKEN`

#### `POST /movies`

Creates a new movie for the authenticated user.

  - **Body:**

    ```json
    {
      "isbn": "12345",
      "title": "Inception",
      "director": {
        "firstname": "Christopher",
        "lastname": "Nolan"
      }
    }
    ```

  - **Success Response (201 Created):** The newly created movie object.

#### `GET /movies`

Retrieves all movies created by the authenticated user.

  - **Success Response (200 OK):** An array of movie objects.

#### `GET /movies/{id}`

Retrieves a single movie by its ID.

  - **Success Response (200 OK):** The requested movie object.

#### `PUT /movies/{id}`

Updates an existing movie by its ID.

  - **Body:** The updated movie data.
  - **Success Response (200 OK):** The updated movie object.

#### `DELETE /movies/{id}`

Deletes a single movie by its ID.

  - **Success Response (200 OK):** A success message.

#### `DELETE /movies`

Deletes all movies created by the authenticated user.

  - **Success Response (200 OK):** A success message.

<!-- end list -->

```
```
## Notes- Passwords are hashed using SHA-256 before saving to DB. The token acts as a temporary session token.

🧪 Testing
You can test this API using tools like:
Postman
URL

✍️ Author
Vivek Modar

📜 License
This project is open-source and free to use.
