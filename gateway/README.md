# Setup

* Install [docker](https://docs.docker.com/engine/install/)
* Define necessary env variables
```
POSTGRES_PASSWORD=MAKE_UP_SOMETHING_RANDOM
POSTGRES_USER=labdao
POSTGRES_DB=labdao
POSTGRES_HOST=localhost
```

# Start the database

```
docker compose up -d
```

Note: New docker installation include docker compose, older installations required you install docker-compose separately and run `docker-compose up -d`

# Frontend start command

from ./gateway/frontend
```
npm run dev
```

# Backend start command

from ./gateway
```
go run app.go
```
