# AWS Cognito Go Client Lambda (Fiber)

This project provides two HTTP endpoints using Fiber for AWS Cognito custom authentication:
- `/initiate-auth`: Initiates custom authentication
- `/respond-challenge`: Responds to the Cognito challenge (e.g., OTP)

## Prerequisites
- Go 1.18+
- Set environment variables:
  - `COGNITO_CLIENT_ID`
  - `COGNITO_CLIENT_SECRET`

## Running the Server

```bash
go run main.go
```

## Example CURL Commands

### 1. Initiate Custom Auth

```bash
curl -X POST http://localhost:8080/initiate-auth \
  -H "Content-Type: application/json" \
  -d '{"username": "your.email@domain.com"}'
```

**Response:**
```json
{
  "challenge_name": "CUSTOM_CHALLENGE",
  "session": "<session-token>"
}
```

### 2. Respond to Challenge (OTP)

```bash
curl -X POST http://localhost:8080/respond-challenge \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your.email@domain.com",
    "otp": "123456",
    "challenge_name": "CUSTOM_CHALLENGE",
    "session": "<session-token>"
  }'
```

**Response:**
```json
{
  "access_token": "<access-token>"
}
```

## Notes
- Replace `your.email@domain.com` with your Cognito username.
- Replace `<session-token>` with the value returned from `/initiate-auth`.
- Replace `123456` with the OTP code you received.
- Replace `<access-token>` with the value returned from `/respond-challenge`.
