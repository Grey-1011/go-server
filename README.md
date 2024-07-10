# Chirpy Web Server

Chirpy is a web server application built using Go that provides an API for creating, retrieving, and managing chirps.

## Getting Started

### Prerequisites

- Go (1.13 or higher)
- Git

### Installation

1. Clone the repository:
git clone https://github.com/Grey-1011/go-server-chirpy.git

2. Navigate to the project directory:
cd go-server-chirpy

3. Create a `.env` file in the project directory and set the required environment variables:

4. Build and run the application:
go build -o out && ./out --debug

## Usage

The Chirpy web server provides endpoints to manage chirps and users.

### API Endpoints

- **POST /api/users**: Create a new user.
- **PUT /api/users**: Update a user's information.

- **POST /api/login**: Authenticate user login and generate JWT.
- **POST /api/revoke**: Revoke a JWT.
- **POST /api/refresh**: Refresh an expired JWT.

- **POST /api/chirps**: Create a new chirp.
- **GET /api/chirps**: Retrieve chirps.
- **GET /api/chirps/{chirpID}**: Retrieve a specific chirp by ID.
- **DELETE /api/chirps/{chirpID}**: Delete a chirp.


- **POST /api/polka/webhooks**: Handle webhook for Polka verification.

- **GET /admin/metrics**: Retrieve server metrics.

## Configuration

The application can be configured using environment variables in the `.env` file.

- `JWT_SECRET`: Secret key for JWT generation and validation.
- `POLKA_API`: URL for the Polka API.

## Contributing

Feel free to contribute to this project by submitting pull requests or issues.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.



---


# API for Chirps

## User resource

```json
{
  "id": 1,
  "email": "walt@breakingbad.com",
  "is_chirpy_red": false
}
```

### POST /api/users
Request Body:
```json
{
  "email": "walt@breakingbad.com",
  "password": "123456"
}
```

Status: 201
Response Body:
```json
{
  "id": 1,
  "email": "walt@breakingbad.com",
  "is_chirpy_red": false
}
```

### PUT /api/users
Headers:
```json
{
  "Authorization": "Bearer ${jwtToken1}"
}
```
Request Body:
```json
{
  "email": "mike@bettercall.com",
  "password": "654321"
}
```
Status: 200


### POST /api/login
Request Body:
```json
{
  "email": "walt@breakingbad.com",
  "password": "123456"
}
```

Returns:  Status: 200
```json
{
  "id": 1,
  "email": "walt@breakingbad.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiIxIiwiZXhwIjoxNzIwNjA5MDQ0LCJpYXQiOjE3MjA2MDU0NDR9.MMH-2bpPhfxdWytK1Uw78ia1TUtlLdMGWwWiHP4AzPo",
  "refresh_token": "db6a8cde5d2030b96419c2d608a371a1c4c1b4c7d9d5a35e5b9fac296a028483"
}
```

### POST /api/chirps
Headers:
```json
{
  "Authorization": "Bearer ${jwtToken1}"
}
```

Request Body:
```json
{
  "body": "I'm the one who knocks!"
}
```

Status: 201
Returns: Chirps
```json
{
  "id": 1,
  "body": "I'm the one who knocks!",
  "author_id": 1
}
```


### GET /api/chirps
Status: 200
Returns an array of chirps

### GET /api/chirps/{chirpID}
Returns chirps by id


### GET /api/chirps?author_id=1
Status: 200
Returns an array of chirps by author_id



### GET /api/chirps?sort=desc
### GET /api/chirps?sort=asc
Status: 200
Returns an array of chirps
by id in *ascending* OR *descending* order


### POST /api/polka/webhooks
Request Body:
```json
{
  "data": {
    "user_id": 1
  },
  "event": "user.upgraded"
}
```
Status: 204
Logged user `is_chirpy_red` to be equal to `true`



### DELETE /api/chirps/{chirpID}
Headers:
```json
{
  "Authorization": "Bearer ${jwtToken1}"
}
```
Status: 204


###  POST /api/refresh
Headers:
```json
Authorization: Bearer <token>
```
Status: 200
Returns:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

###  POST /api/revoke
Headers:
```json
Authorization: Bearer <token>
```
Respond with a 204 status code. A 204 status means the request was successful but no body is returned.

###  GET /app/*

###  GET /api/healthz
Returns a 200 OK status code indicating that it has started up successfully and is listening for traffic

###  GET /admin/metrics
Returns the number of visits(`fileserverHits`) to the website: `/app/`

###  GET /api/reset
Reset your `fileserverHits` back to ``0`.

