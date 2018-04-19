
## Build & Run
```
cd dummy-mongo-stress-test
export GOPATH=$(pwd)
cd src/github.com
go get ./...
go build -o main ./app && ./main mongodb://localhost Flaconi12 20
```

## Run
```
./main <database address> <campaign id> <number of fetches>
```

## (Cross) build for ubuntu
```
GOOS=linux GOARCH=amd64 go build -o main ./app
```
The executable file `main` now can be used to run on ubuntu servers