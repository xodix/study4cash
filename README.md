# Study4Cash

This is a website that compares the graduation data, student attending data and average wages data together.

## How to run

create .env file in root folder

```
DB_ROOT_PASSWORD=ROOT_PASSWORD
DB_NAME=study4cash
DB_USER=study4cash_user
DB_PASSWORD=DB_PASSWORD
TOKEN_SECRET=TOKEN_SECRET_THAT_SHOULD_BE_AT_LEAST_32_CHARACTERS_LONG
API_URL=EXTERNAL_ADDRESS_OF_THE_BACKEND_SERVER (for testing purposes, http://localhost:8080, or an IP address or domain name in production)
```

Build the docker images used

```sh
docker compose build
```

Run the docker compose

```sh
docker compose up
```

## Ports

| Service        | Port |
| -------------- | ---- |
| Backend server | 8080 |
| OpenAPI        | 5050 |
| Client         | 3000 |

## Documentation

documentation can be found at http://localhost:5050/ after running the docker container
