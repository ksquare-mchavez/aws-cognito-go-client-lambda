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

If your main file is in `cmd/main.go`:

```bash
go run cmd/main.go
```

Or if your main file is in the root:

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
  "access_token": "<access-token>",
  "id_token": "<id-token>",
  "refresh_token": "<refresh-token>",
  "token_type": "Bearer"
}
```



## Notes
- Replace `your.email@domain.com` with your Cognito username.
- Replace `<session-token>` with the value returned from `/initiate-auth`.
- Replace `123456` with the OTP code you received.
- The `/respond-challenge` endpoint returns all Cognito tokens: `access_token`, `id_token`, `refresh_token`, and `token_type`.
- If you move your main file, update the run command accordingly.


## Troubleshooting

- **Invalid session for the user:**
  - Make sure you use the exact session value returned from `/initiate-auth` in `/respond-challenge`.
  - Do not reuse or modify the session value.
  - Ensure the username and challenge name match those used in the initial request.

- **AccessDeniedException in Lambda triggers:**
  - Check the IAM role assigned to your Cognito Lambda triggers (e.g., VerifyAuthChallengeResponse, DefineAuthChallenge, CreateAuthChallenge).
  - Make sure the role has all required permissions for AWS services your Lambda uses (e.g., DynamoDB, SES, SSM).
  - Update the IAM policy to add missing permissions if needed.

- **Lambda location:**
  - Custom Cognito Lambda triggers (e.g., `defineAuthChallengePy`,`createAuthChallengePy`,`verifyAuthChallengeResponsePy`) are located in the `lambda/` directory.
  - Ensure your Lambda code has the correct permissions and logic for your authentication flow.
