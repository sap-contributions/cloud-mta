package main

//go:generate go run ./tools/embed.go -source=./validations/schema/schema_v3.2.yaml -target=./validations/mta_schema.go -name=schemaDef -package=validate