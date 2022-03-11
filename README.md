# Postal Code service

This service is used in Grippenet to map from postal codes to official municipality codes. 

## Build database

We use dataset provided by [LaPoste](https://datanova.legroupe.laposte.fr/explore/dataset/laposte_hexasmal) from [data.gouv.fr](https://www.data.gouv.fr/fr/datasets/base-officielle-des-codes-postaux) dataset, released under  [Licence Ouverte / Open License 2.0](https://www.etalab.gouv.fr/licence-ouverte-open-licence/)

To update the database, download the hexasmal csv file , recompile the json using build command.

The expected format is a semi column separated file, with at least 3 columns :
 - 1 : Municipality code 
 - 2 : Municipality label
 - 3 : Postal code associated with the municipality

Only the first municipality will be used for the municipality label.

The database is compiled into a optmized json file with postal codes already mapped to a list of municipalities entries (using a numerical index)

Run 

```go
    go run cmd/build/build.go laposte_hexasmal.csv
```

## Run the service

To run the service
```go
    go run cmd/server/server.go 
```

## Configuration

Environment variable "PORT" can be used to change the default port (8080)
