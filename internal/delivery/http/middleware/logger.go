package middleware

import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start  := time.Now()
        path   := c.Request.URL.Path
        method := c.Request.Method

        c.Next()

        latency := time.Since(start)
        status  := c.Writer.Status()

        statusColour := "\033[32m" // green
        if status >= 400 { statusColour = "\033[33m" } // yellow
        if status >= 500 { statusColour = "\033[31m" } // red
        reset := "\033[0m"

        fmt.Printf("%s[%d]%s %-7s %-30s %v\n",
            statusColour, status, reset,
            method, path, latency)
    }
}

func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                fmt.Printf("\033[31m[PANIC] %v\033[0m\n", err)
                c.AbortWithStatusJSON(500, gin.H{
                    "success": false,
                    "error":   "internal server error",
                })
            }
        }()
        c.Next()
    }
}