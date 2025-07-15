# CineHub

Cinehub is a full stack movie explorer application built using Vanilla JavaScript, HTML, and CSS for the frontend, and Go for the backend. It uses Postgres as the database and implements JWT authentication for user sessions.

## Features

- Landing page displaying top and random movies
- Movie details page showing movie details and cast
- User JWT-based Authentication (Register and Login)
- Favorite and Watchlist Management

## Get started locally

### Prerequisites

- Go installed on your system
- Postgres database either locally or on a cloud provider
- Database populated with movies data ( use database-dump.sql to populate the database )
- Air installed on your system ( https://github.com/air-verse/air )

1. Clone the repository

```
git clone https://github.com/junaidshaikh-js/CineHubServer.git
```

2. Navigate to the project directory

```
cd CineHubServer
```

3. Set up environment variables

```
cp .env.example .env
```

4. Run the server

```
air
```

Open your browser and navigate to `http://localhost:5555` to access the application.
