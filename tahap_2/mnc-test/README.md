# MNC Test

Sample program for create transaction

## Installation

Use the golang for running this program

```bash
git clone on github https://github.com/anggaramp/mnc-test
```

## Usage

```bash
# load library data 
go mod tidy

# running docker file
cd docker/
docker-compose up -d  

# update connection DB on main.go 
dsn = "host=0.0.0.0 port=10090 user=mnc_test password=.mnc_test! dbname=mnc_test sslmode=disable"

# update output path on main.go 
"outputPaths": ["./logs/app.log"],
"errorOutputPaths": ["./logs/app.log"],

# running application 
go run main.go

# migration DB 
Hit /migration on postman
```