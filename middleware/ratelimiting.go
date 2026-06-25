package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

func GetClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]
	if !exists {
		limiter := rate.NewLimiter(5, 10) // 5 request/s, brust 10 vi du 1,5 => cu sau 1 giay cap nhat 1 token
		newClient := &Client{limiter, time.Now()}
		clients[ip] = newClient
		log.Printf("A clients[%s] - {limiter: %+v, lastSeen: %s}", ip, newClient.limiter, newClient.lastSeen)

		return limiter
	}

	log.Printf("A clients[%s] - {limiter: %+v, lastSeen: %s}", ip, client.limiter, client.lastSeen)
	client.lastSeen = time.Now()
	return client.limiter
}

func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// go install github.com/rakyll/hey@latest
// hey -n 20 -c 1 -H "X-API-Key:d98dd177-ab53-4ecb-9c40-eccfa6d9bcb2" http://localhost:8080/api/v1/categories/golang
func RateLimitingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := GetClientIP(ctx)

		limiter := getRateLimiter(ip)

		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many request",
				"message": "Ban da gui qua nhieu request. Hay gui lai sau",
			})
		}
	}
}
