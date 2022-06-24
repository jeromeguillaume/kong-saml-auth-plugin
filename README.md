# custom-plugins
Init repo for setting up custom plugin in Go

## Prepare Kong environment : Docker network and PostgreSQL (start & prepare PostgreSQL)
1) Create a custom Docker network to allow the containers to discover and communicate with each other:
``` 
docker network create kong-net
``` 
2) Start a PostgreSQL container
```
docker run -d --name kong-database-saml \
  --network=kong-net \
  -p 5432:5432 \
  -e "POSTGRES_USER=kong" \
  -e "POSTGRES_DB=kong" \
  -e "POSTGRES_PASSWORD=kongpass" \
  postgres:9.6
```
3) Prepare the Kong database:
```
docker run --rm --network=kong-net \
  -e "KONG_DATABASE=postgres" \
  -e "KONG_PG_HOST=kong-database" \
  -e "KONG_PG_PASSWORD=kongpass" \
  -e "KONG_PASSWORD=test" \
 kong/kong-gateway:2.8.1.1-alpine kong migrations bootstrap
```
## Build Kong gateway image with Go custom plugin
Build Kong Gateway from project root with name as richemont/kong-gateway-saml (note the trailing .)
```
docker build -t richemont/kong-gateway-saml .
```
## Run Kong gateway container with Go custom plugin
Start Kong Gateway
```
./start-kong.sh
```

## Environment
Once the environment is up and running, both Kong Manager and the Admin API can be accessed on localhost for testing
Manager or the Admin API:
![manager](https://user-images.githubusercontent.com/11033758/138120791-4caf86cf-9204-40d1-a234-fa2962f106cb.png)
![adminapi](https://user-images.githubusercontent.com/11033758/138120631-c8294e8f-6498-4536-a734-02f69811691b.png)
