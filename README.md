# Gin Gonic CRUD API with JWT Authentication

A RESTful API built with Go and Gin framework for user management with complete CRUD operations and JWT authentication.

## Features

- ✅ Complete CRUD operations (Create, Read, Update, Delete)
- ✅ JWT Authentication & Authorization
- ✅ Password hashing with bcrypt
- ✅ Input validation with Gin binding
- ✅ Consistent API response format
- ✅ Error handling
- ✅ Auto migration with GORM
- ✅ Soft delete support

## Database Schema

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    born_date TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## Authentication

### JWT Token Usage

Include the JWT token in the Authorization header for protected endpoints:

```
Authorization: Bearer <your_jwt_token>
```

## API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

---

## 🔐 AUTHENTICATION ENDPOINTS

### **1. POST /auth/login** - User Login
- **Description**: Login dengan email dan password, mengembalikan JWT token
- **Headers**: `Content-Type: application/json`
- **Request Body**:
```json
{
  "email": "john@example.com",
  "password": "mypassword123"
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "address": "123 Main St",
      "email": "john@example.com",
      "born_date": "1990-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

---

## 👥 USER ENDPOINTS

### **2. GET /users** - Get All Users
- **Method**: `GET`
- **Auth**: Optional (JWT token)
- **Description**: Mengambil semua data user
- **Query Parameters** (optional):
  - `page`: Page number (default: 1)
  - `limit`: Items per page (default: 10, max: 100)
  - `sort_by`: Sort field (id, name, email, etc.)
  - `order`: Sort order (ASC, DESC)

### **3. GET /users/:id** - Get User by ID
- **Method**: `GET`
- **Auth**: Optional (JWT token)
- **Parameters**: `id` (path parameter)

### **4. POST /users** - Create User
- **Method**: `POST`
- **Auth**: None
- **Headers**: `Content-Type: application/json`
- **Request Body**:
```json
{
  "name": "John Doe",
  "address": "123 Main Street",
  "email": "john@example.com",
  "password": "mypassword123",
  "born_date": "1990-01-01"
}
```
- **Validation Rules**:
  - `name`: required, min 2 chars, max 100 chars
  - `email`: required, valid email format, must be unique
  - `password`: required, min 6 characters
  - `born_date`: required, format YYYY-MM-DD

### **5. PUT /users/:id** - Update User
- **Method**: `PUT`
- **Auth**: Optional (JWT token)
- **Parameters**: `id` (path parameter)
- **Request Body**: All fields optional

### **6. DELETE /users/:id** - Delete User
- **Method**: `DELETE`
- **Auth**: Optional (JWT token)
- **Parameters**: `id` (path parameter)

---

## 🛡️ PROTECTED ENDPOINTS (Require JWT Token)

### **7. GET /users/profile** - Get User Profile
- **Method**: `GET`
- **Auth**: **Required** (JWT token)
- **Description**: Mengambil profile user yang sedang login
- **Headers**: `Authorization: Bearer <token>`

---

## 📚 BOOK ENDPOINTS (Placeholder)

### **8. GET /books** - Get All Books
- **Method**: `GET`
- **Auth**: Optional

---

## ❌ ERROR RESPONSES

### **Validation Error**
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "Key: 'CreateUserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

### **Authentication Error**
```json
{
  "success": false,
  "message": "Authorization header required"
}
```

### **Invalid Credentials**
```json
{
  "success": false,
  "message": "Invalid email or password"
}
```

### **Token Expired**
```json
{
  "success": false,
  "message": "Invalid or expired token",
  "error": "token is expired"
}
```

---

## ✅ DATABASE MIGRATION COMPLETED

**Migration Status**: ✅ **COMPLETED**
- All existing users have been migrated with default password
- Auto migration runs smoothly on application startup
- No manual intervention required

**Migration Details:**
- **12 existing users** successfully migrated
- **Default password**: `defaultpassword123`
- **Security**: All passwords are properly hashed with bcrypt

**⚠️ Important**: Users should change their default password after first login for security.

---

## 🚀 QUICK START

### 1. Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "born_date": "1990-01-01"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 3. Use JWT Token
```bash
# Save token from login response
TOKEN="your_jwt_token_here"

# Access protected endpoint
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer $TOKEN"
```

### 4. CRUD Operations
```bash
# Get all users
curl http://localhost:8080/api/v1/users

# Get user by ID
curl http://localhost:8080/api/v1/users/1

# Update user
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe"}'

# Delete user
curl -X DELETE http://localhost:8080/api/v1/users/1
```

---

## 🔧 SECURITY FEATURES

- **Password Hashing**: Menggunakan bcrypt dengan cost default
- **JWT Tokens**: Expires in 24 hours
- **Input Validation**: Comprehensive validation untuk semua inputs
- **SQL Injection Prevention**: GORM parameterized queries
- **Unique Constraints**: Email uniqueness di database level

## 📝 NOTES

- JWT secret key should be changed in production
- Password minimum length is 6 characters
- Token expires in 24 hours
- All passwords are hashed before storage
- Email must be unique across all users

## 🔄 LEGACY ENDPOINTS

For backward compatibility:
- `GET /user` → `GET /api/v1/users`
- `GET /book` → `GET /api/v1/books`