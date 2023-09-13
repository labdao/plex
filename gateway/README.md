# Setup

* Install [docker](https://docs.docker.com/engine/install/)
* Define necessary env variables
```
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
POSTGRES_PASSWORD=MAKE_UP_SOMETHING_RANDOM
POSTGRES_USER=labdao
POSTGRES_DB=labdao
POSTGRES_HOST=localhost
```
* Recommended: Install [direnv](https://direnv.net/). With it installed you can create `.env` file with the above environment variables and have them automagically set when you descend into the folder. 

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
