package main

import (
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis"
)

var redisClient *redis.Client
const (
    REDIS_ADDRESS = "localhost:6379"
    PORT = ":8090"
)

func main() {
    router := gin.Default()

    // Redis connection
    redisClient = redis.NewClient(&redis.Options{
        Addr: REDIS_ADDRESS,
    })

    // Endpoint definition
    router.GET("/message/list", listMessages)

    router.Run(PORT)
}

func listMessages(c *gin.Context) {
    sender := c.Query("sender")
    receiver := c.Query("receiver")

    // Redis search
    messages, err := redisClient.Sort("messages:"+sender+"->"+receiver, &redis.Sort{
        Order: "DESC",
        Alpha: false,
    }).Result()

    if err != nil {
        c.JSON(500, gin.H{
            "error": "Redis error: " + err.Error(),
        })
        return
    }

    c.JSON(200, gin.H{
        "messages": messages,
    })
}
