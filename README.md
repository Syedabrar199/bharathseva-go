# Bharat Seva Space - Backend API

A comprehensive backend API for Bharat Seva Space, built with Go, Gin, and PostgreSQL. This API handles user authentication, query management, application tracking, and admin operations.

## Features

### 🔐 Authentication & Authorization
- JWT-based authentication
- Role-based access control (User/Admin)
- Password hashing with bcrypt
- Secure token management

### 📝 Query Management
- Public query submission from website
- Admin query management and assignment
- Query status tracking (New, Contacted, Converted, Closed)
- Query statistics and analytics

### 📋 Application Management
- User application creation and tracking
- Admin application management
- Progress tracking and status updates
- CA assignment to applications
- Payment status tracking

### 👥 User Management
- User registration and login
- Profile management
- Admin user management
- User statistics and analytics

### 📊 Dashboard & Analytics
- User dashboard with application overview
- Admin dashboard with comprehensive stats
- Real-time statistics and reports

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin (Web framework)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT tokens
- **Password Hashing**: bcrypt

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## Installation & Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd bharath-Go
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Database Setup
1. Create a PostgreSQL database:
```sql
CREATE DATABASE bharat_seva_space;
```

2. Update the database configuration in `config.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=bharat_seva_space
```

### 4. Environment Configuration
Update `config.env` with your settings:
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=bharat_seva_space

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=24h

# Server Configuration
PORT=8080
ENV=development

# File Upload Configuration
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=10485760
```

### 5. Run the Application
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Public Endpoints (No Authentication Required)

#### Query Submission
- `POST /api/queries` - Submit a public query

#### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login

### User Endpoints (Authentication Required)

#### Profile Management
- `GET /api/user/profile` - Get user profile
- `PUT /api/user/profile` - Update user profile

#### Dashboard
- `GET /api/user/dashboard` - Get dashboard data

#### Applications
- `GET /api/user/applications` - Get user's applications
- `POST /api/user/applications` - Create new application
- `GET /api/user/applications/:id` - Get specific application
- `GET /api/user/applications/stats` - Get application statistics

### Admin Endpoints (Admin Authentication Required)

#### Query Management
- `GET /api/admin/queries` - Get all queries
- `GET /api/admin/queries/:id` - Get specific query
- `PUT /api/admin/queries/:id` - Update query
- `GET /api/admin/queries/stats` - Get query statistics

#### User Management
- `GET /api/admin/users` - Get all users
- `GET /api/admin/users/:id` - Get specific user
- `PUT /api/admin/users/:id` - Update user
- `GET /api/admin/users/stats` - Get user statistics

#### Application Management
- `GET /api/admin/applications` - Get all applications
- `PUT /api/admin/applications/:id` - Update application
- `GET /api/admin/applications/stats` - Get application statistics

## API Examples

### User Registration
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "9876543210",
    "password": "password123"
  }'
```

### User Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Submit Query
```bash
curl -X POST http://localhost:8080/api/queries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "email": "jane@example.com",
    "phone": "9876543211",
    "service": "GST Registration",
    "message": "I need help with GST registration for my business."
  }'
```

### Create Application (with token)
```bash
curl -X POST http://localhost:8080/api/user/applications \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "service_type": "GST Registration",
    "description": "Need GST registration for my business",
    "amount": 5000
  }'
```

## Database Schema

### Users Table
- `id` - Primary key
- `email` - Unique email address
- `phone` - Unique phone number
- `name` - User's full name
- `password` - Hashed password
- `role` - User role (user/admin)
- `is_active` - Account status
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Queries Table
- `id` - Primary key
- `name` - Query submitter's name
- `email` - Query submitter's email
- `phone` - Query submitter's phone
- `service` - Requested service
- `message` - Query message
- `status` - Query status (new/contacted/converted/closed)
- `assigned_to` - Admin user ID assigned to query
- `notes` - Admin notes
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Applications Table
- `id` - Primary key
- `user_id` - Foreign key to users table
- `service_type` - Type of service
- `status` - Application status (pending/in_progress/completed/cancelled)
- `progress` - Progress percentage
- `payment_status` - Payment status (pending/paid/refunded)
- `amount` - Service amount
- `description` - Application description
- `assigned_ca` - Admin user ID assigned as CA
- `notes` - Admin notes
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Documents Table
- `id` - Primary key
- `application_id` - Foreign key to applications table
- `file_name` - Original file name
- `file_path` - File storage path
- `file_size` - File size in bytes
- `file_type` - File MIME type
- `description` - Document description
- `uploaded_at` - Upload timestamp
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

## Default Admin Account

When the application starts for the first time, a default admin account is created:

- **Email**: admin@bharatseva.com
- **Password**: password

**Important**: Change the default password after first login!

## Development

### Project Structure
```
bharath-Go/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── config.env              # Environment configuration
├── README.md              # This file
├── config/
│   └── config.go          # Configuration management
├── database/
│   └── database.go        # Database connection and migrations
├── models/
│   ├── user.go            # User model
│   ├── query.go           # Query model
│   ├── application.go     # Application model
│   └── document.go        # Document model
├── handlers/
│   ├── auth.go            # Authentication handlers
│   ├── query.go           # Query handlers
│   ├── application.go     # Application handlers
│   └── user.go            # User management handlers
├── middleware/
│   └── auth.go            # Authentication middleware
├── routes/
│   └── routes.go          # Route definitions
└── utils/
    └── auth.go            # Authentication utilities
```

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o bharat-seva-space main.go
```

## Security Considerations

1. **JWT Secret**: Change the JWT secret in production
2. **Database**: Use strong passwords and secure connections
3. **CORS**: Configure CORS properly for production
4. **Rate Limiting**: Implement rate limiting for API endpoints
5. **Input Validation**: All inputs are validated
6. **SQL Injection**: Protected by GORM

## Deployment

### Docker Deployment
```dockerfile
FROM golang:1.21-alpine

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

EXPOSE 8080
CMD ["./main"]
```

### Environment Variables for Production
- Set `ENV=production`
- Use strong `JWT_SECRET`
- Configure production database
- Set appropriate `PORT`

## Support

For support and questions:
- Create an issue in the repository
- Contact the development team

## License

This project is proprietary software for Bharat Seva Space. 