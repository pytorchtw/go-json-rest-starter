# go-json-rest-starter

#### This is a simple demo json rest server in golang with code adapted from https://github.com/ant0ine/go-json-rest 


Installation

- Git clone to local directory

- Build, test and start up the service with docker compose
```
docker-compose up
```

- Open a new console and test with the following:

```
curl -i -H 'Content-Type: application/json' \
    -d '{"Code":"FR","Name":"France"}' http://127.0.0.1:8080/countries
curl -i -H 'Content-Type: application/json' \
    -d '{"Code":"US","Name":"United States"}' http://127.0.0.1:8080/countries
curl -i http://127.0.0.1:8080/countries/FR
curl -i http://127.0.0.1:8080/countries/US
curl -i http://127.0.0.1:8080/countries
curl -i -X DELETE http://127.0.0.1:8080/countries/FR
curl -i http://127.0.0.1:8080/countries
curl -i -X DELETE http://127.0.0.1:8080/countries/US
curl -i http://127.0.0.1:8080/countries
```

You should be able the see the incoming requests in the docker console

