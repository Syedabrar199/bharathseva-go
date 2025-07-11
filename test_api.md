# API Testing Guide

This document contains test examples for all API endpoints. Use these to verify your backend is working correctly.

## Prerequisites
- Backend server running on `http://localhost:8080`
- PostgreSQL database configured
- Postman or curl installed

## Test Flow

### 1. Health Check
```bash
curl http://localhost:8080/health
```
**Expected Response:**
```json
{
  "status": "ok",
  "message": "Bharat Seva Space API is running"
}
```

### 2. User Registration
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "phone": "9876543210",
    "password": "password123"
  }'
```
**Expected Response:**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "email": "test@example.com",
    "phone": "9876543210",
    "name": "Test User",
    "role": "user",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. User Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```
**Expected Response:**
```json
{
  "message": "Login successful",
  "user": {
    "id": 1,
    "email": "test@example.com",
    "phone": "9876543210",
    "name": "Test User",
    "role": "user",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 4. Submit Public Query
```bash
curl -X POST http://localhost:8080/api/queries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "9876543211",
    "service": "GST Registration",
    "message": "I need help with GST registration for my business."
  }'
```
**Expected Response:**
```json
{
  "message": "Query submitted successfully",
  "query_id": 1
}
```

### 5. Get User Profile (with token)
```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
**Expected Response:**
```json
{
  "user": {
    "id": 1,
    "email": "test@example.com",
    "phone": "9876543210",
    "name": "Test User",
    "role": "user",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 6. Create Application (with token)
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
**Expected Response:**
```json
{
  "message": "Application created successfully",
  "application_id": 1
}
```

### 7. Get User Applications (with token)
```bash
curl -X GET http://localhost:8080/api/user/applications \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
**Expected Response:**
```json
{
  "applications": [
    {
      "id": 1,
      "user_id": 1,
      "service_type": "GST Registration",
      "status": "pending",
      "progress": "0%",
      "payment_status": "pending",
      "amount": 5000,
      "description": "Need GST registration for my business",
      "assigned_ca": null,
      "notes": "",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "user": {
        "id": 1,
        "email": "test@example.com",
        "phone": "9876543210",
        "name": "Test User",
        "role": "user",
        "is_active": true,
        "created_at": "2024-01-01T00:00:00Z"
      },
      "assigned_ca_user": null,
      "documents": []
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1
  }
}
```

### 8. Get Dashboard Data (with token)
```bash
curl -X GET http://localhost:8080/api/user/dashboard \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
**Expected Response:**
```json
{
  "user": {
    "id": 1,
    "email": "test@example.com",
    "phone": "9876543210",
    "name": "Test User",
    "role": "user",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  },
  "application_stats": {
    "total": 1,
    "pending": 1,
    "in_progress": 0,
    "completed": 0,
    "cancelled": 0
  },
  "recent_applications": [
    {
      "id": 1,
      "user_id": 1,
      "service_type": "GST Registration",
      "status": "pending",
      "progress": "0%",
      "payment_status": "pending",
      "amount": 5000,
      "description": "Need GST registration for my business",
      "assigned_ca": null,
      "notes": "",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "assigned_ca_user": null
    }
  ]
}
```

## Admin Tests

### 9. Admin Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@bharatseva.com",
    "password": "password"
  }'
```

### 10. Get All Queries (admin token)
```bash
curl -X GET http://localhost:8080/api/admin/queries \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

### 11. Get All Users (admin token)
```bash
curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

### 12. Get All Applications (admin token)
```bash
curl -X GET http://localhost:8080/api/admin/applications \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

### 13. Update Application Status (admin token)
```bash
curl -X PUT http://localhost:8080/api/admin/applications/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -d '{
    "status": "in_progress",
    "progress": "25%",
    "assigned_ca": 1,
    "notes": "Application under review"
  }'
```

## Error Testing

### 14. Invalid Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "wrong@example.com",
    "password": "wrongpassword"
  }'
```
**Expected Response:**
```json
{
  "error": "Invalid credentials"
}
```

### 15. Unauthorized Access
```bash
curl -X GET http://localhost:8080/api/user/profile
```
**Expected Response:**
```json
{
  "error": "Authorization header required"
}
```

### 16. Invalid Token
```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer invalid_token"
```
**Expected Response:**
```json
{
  "error": "Invalid or expired token"
}
```

## Testing Checklist

- [ ] Health check endpoint works
- [ ] User registration works
- [ ] User login works
- [ ] JWT token generation works
- [ ] Public query submission works
- [ ] User profile retrieval works
- [ ] Application creation works
- [ ] Application listing works
- [ ] Dashboard data retrieval works
- [ ] Admin login works
- [ ] Admin can view all queries
- [ ] Admin can view all users
- [ ] Admin can view all applications
- [ ] Admin can update application status
- [ ] Error handling works correctly
- [ ] Authentication middleware works
- [ ] Authorization middleware works

## Notes

1. Replace `YOUR_JWT_TOKEN` with the actual token received from login
2. Replace `ADMIN_JWT_TOKEN` with the admin token received from admin login
3. All timestamps will be actual timestamps, not the examples shown
4. IDs will be auto-generated, so they may not match the examples
5. Test with different data to ensure the system works with various inputs 