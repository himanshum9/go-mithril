# Authentication Service Documentation

## Overview
The Auth Service is responsible for handling user authentication within the multi-tenant location streaming system. It utilizes AWS Cognito for user management and authentication, ensuring secure access to the system's resources.

## Setup Instructions
1. **AWS Cognito Configuration**: 
   - Create a new Cognito User Pool in the AWS Management Console.
   - Configure the User Pool with appropriate settings (e.g., password policies, multi-factor authentication).
   - Note the User Pool ID and App Client ID, as these will be required for the service configuration.

2. **Environment Variables**:
   - Set the following environment variables in your `.env` file or your deployment environment:
     - `COGNITO_USER_POOL_ID`: Your Cognito User Pool ID.
     - `COGNITO_APP_CLIENT_ID`: Your Cognito App Client ID.
     - `AWS_REGION`: The AWS region where your Cognito User Pool is located.

3. **Run the Service**:
   - Navigate to the `auth-service` directory.
   - Execute the command: `go run main.go`.
   - The service will start on the configured port (default is 8080).

## API Endpoints
- **POST /register**: Register a new user.
  - Request Body: `{ "email": "user@example.com", "password": "yourpassword" }`
  
- **POST /login**: Authenticate a user and return a JWT token.
  - Request Body: `{ "email": "user@example.com", "password": "yourpassword" }`
  
- **GET /me**: Retrieve the authenticated user's information.
  - Headers: `Authorization: Bearer <JWT_TOKEN>`

## Role-Based Access Control (RBAC)
The Auth Service implements a simple RBAC mechanism to differentiate between tenant users and admin users. Ensure that user roles are assigned appropriately during registration.

## Error Handling
The service provides meaningful error messages for various authentication failures, such as invalid credentials or unauthorized access attempts.

## Testing
To test the Auth Service, you can use tools like Postman or cURL to interact with the API endpoints. Ensure that you have valid credentials for testing the login and registration functionalities.

## Dependencies
- AWS SDK for Go
- Gorilla Mux (for routing)

## License
This project is licensed under the MIT License. See the LICENSE file for details.