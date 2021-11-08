# game-currency

## Framework

- Web : Echo
- Configuration : Env File
- Database : Postgres

## Architecture
Handler -> Service -> Repository

## Database Diagram
![](game-currency-database-diagram.png?raw=true)


### How To Run This Project

> Make Sure you have run the init.sql in your postgres if not use docker for database
> 
> Fill .env file with your postgres configuration (change POSTGRES_HOST if not using docker as database)

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone git@github.com:mhaqiw/game-currency.git

#move to project
$ cd game-currency

# Build the docker image first
$ make docker

# Run the application
$ make run

# check if the containers are running
$ docker ps


# Stop
$ make stop
```

Request example

```bash
# Add Currency 
curl --location --request POST 'localhost:9090/currency' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "cukong"
}'
```

```bash
# Get All Currency 
curl --location --request GET 'localhost:9090/currency'
```

```bash
# Add Conversion Rate
curl --location --request POST 'localhost:9090/conversion' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from_id": 1,
    "to_id": 2,
    "rate": 29
}'
```

```bash
# Get Rate (Make sure already Add conversion before execute this) 
# path /conversion/rate/:from/:to/:amount
curl --location --request GET 'localhost:9090/conversion/rate/1/2/100'
```

```bash
# Stop Service
$ make stop
```