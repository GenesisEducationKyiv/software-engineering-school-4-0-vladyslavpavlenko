[![Go Report Card](https://goreportcard.com/badge/github.com/vladyslavpavlenko/genesis-api-project)](https://goreportcard.com/report/github.com/vladyslavpavlenko/genesis-api-project)

# Genesis API Project

This service implements an API for subscribing to exchange rate updates via email (the current implementation focuses on the `USD to UAH` exchange rate).
## API Endpoints

### `GET` /rate

This endpoint returns the current `USD to UAH` exchange rate using the Coinbase API.

#### Parameters

`No parameters`

#### Response Codes

```
200: Returns the actual USD to UAH exchange rate.
400: Invalid status value.
```

---

### `POST` /subscribe

This endpoint adds an email address to the database and automatically subscribes it to the USD to UAH exchange rate newsletter.

_The code includes the ability to subscribe to other rates for future development, but this functionality is not currently used to fulfill the requirements._

#### Parameters
``email`` **string** (formData): The email address to be added to the database and the mailing list.

#### Response Codes
```
200: The email address is added to the database and subscribed to the mailing list.
409: The email address already exists.
```

_Not mentioned in the task, but arose during the development process:_
```
400: The provided data (such as email address) is invalid.
500: Internal error status.
```

---

### `POST` /sendEmails

This endpoint sends the current `USD to UAH` exchange rate to subscribed email addresses using goroutines.

#### Parameters
`No parameters`

#### Response Codes

```
200: Emails were sent.
```


## Usage
Clone the repository to your local machine:
```sh
git clone https://github.com/vladyslavpavlenko/genesis-api-project.git
cd genesis-api-project
```
❗️Ensure you have a `.env` file in the root directory with the necessary environment variables. The `.env` should look like this:
```dotenv
GMAIL_EMAIL=<GMAIL_EMAIL>
GMAIL_PASSWORD=<GMAIL_APP_PASSWORD>
DB_HOST=<DB_HOST>
DB_PORT=<DB_PORT>
DB_USER=<DB_USER>
DB_PASSWORD=<DB_PASSWORD>
DB_NAME=<DB_NAME>
```

### Makefile
For Unix-like systems, use the following command to build the application binary:
```sh
make up_build
```
For Windows systems, use the following command:
```cmd
make -f Makefile.windows build_app
```
To start the Docker containers without forcing a build, run:
```sh
make up
```
To stop the Docker containers, run:
```sh
make down
```

### Docker Compose
Alternatively, you can use Docker Compose commands directly:

Start the containers without forcing a build:
```sh
docker-compose up -d
```
Build the application and start the containers:
```sh
docker-compose up --build -d
```
Stop the containers:
```sh
docker-compose down
```

### Accessing the Application
The application will be accessible at [`http://localhost:8080/`](http://localhost:8081/).
