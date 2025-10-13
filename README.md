# Go Movies CRUD API
 This is a simple Go-based REST API project for managing a Movies database. It uses
 MySQL as the backend database and Gorilla Mux for routing.
 
✅ Implemented Features
Feature	Status	        Key     Files/Middleware
User Registration	    Done	handlers.RegisterUser
User Login (JWT)	    Done	handlers.LoginUser
JWT Authentication	    Done	middleware.Auth
Protected Routes	    Done	All /movies routes
Token Refresh	        Done	/refresh-token
Logout (Token Revoke)	Done	/logout
Password Hashing	    Done	bcrypt in RegisterUser
Cookie Handling	        Done	http.SetCookie
Context Propagation	    Done	r.WithContext()


 ⚙️ Tech Stack :
✅ Language: Go (Golang)
✅ Database: MySQL
✅ Routing: Gorilla Mux
✅ JSON Parsing: encoding/json


📦 Setup Instructions

1. Install dependencies
go mod tidy

2. Configure MySQL
DB Driver: database/sql with go-sql-driver/mysql
Create a MySQL database:

 ### Prerequisites
- Go installed on your system (1.20 or higher recommended)- MySQL server- MySQL Workbench (optional for DB management)
 ### Database Setup
 Create a new MySQL database and required tables:
 ```
 CREATE DATABASE go_movie_db;
 USE go_movie_db;
 CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    password VARCHAR(255)
 );
 CREATE TABLE movies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    isbn VARCHAR(50),
    title VARCHAR(100),
    director_firstname VARCHAR(100),
    director_lastname VARCHAR(100)
 );
```
 ### Clone the Project
 ```bash
 git clone https://github.com/vivekmodar03/go-movies-crud.git
 cd go-movies-crud
 ```
 ### Project Structure
 ```
go-movies-crud/
├── Handlers/           # Business logic handlers
│   ├── movie_handlers.go
│   └── auth_handlers.go
├── Middleware/         # Authentication middleware
│   └── jwt_auth.go
├── Routes/             # Router setup
│   └── router.go
├── Model/              # Database models
│   └── user.go, movie.go
├── DB/                 # DB connection setup
│   └── db.go
├── go.mod              # Go module
├── main.go             # Application entry point
└── README.md           # Project documentation
 ```--


📬 API Endpoints

🔐 Auth Routes

POST /register – Register a new user
POST /login – Login with email, receive token  to on console

Use that token as Authorization: Bearer <token> header for protected routes

🎬 Movie Routes

All routes below require the OTP token in Authorization header
GET /movies – Get all movies
GET /movies/{id} – Get movie by ID
POST /movies – Create a new movie
PUT /movies/{id} – Update a movie
DELETE /movies/{id} – Delete a movie
DELETE /movies – Delete all movies


 Example Header:
 ```
 Authorization: TOKEN
 ```--


## Run the Server
 ```bash
     go build
     go run main.go
 ```
     The server runs by default on `http://localhost:8080`.--
## Notes- Passwords are hashed using SHA-256 before saving to DB. The token acts as a temporary session token.

🧪 Testing
You can test this API using tools like:
Postman
URL

✍️ Author
Vivek Modar

📜 License
This project is open-source and free to use.