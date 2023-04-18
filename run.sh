# Fetch dependencies if required
# go get github.com/gin-gonic/gin
# go get github.com/streadway/amqp
# go get github.com/go-redis/redis

# export GO111MODULE="off"
export GIN_MODE=release

docker-compose up -d

go build ./...

cd cmd/api && go run .
cd ../queue && go run .
cd ../reporting && go run .
