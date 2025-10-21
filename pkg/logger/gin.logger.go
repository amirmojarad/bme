package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

var (
	internalLogger *zap.Logger
	once           sync.Once
	SkipPaths      []string = []string{"/metrics", "/health", "/healthcheck"}
	SystemField             = zap.String("system", "grpc")
	ServerField             = zap.String("span.kind", "server")
)

type CustomResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

type GinLogEntry struct {
	Namespace     string `json:"namespace"`
	Logger        string `json:"logger"`
	Timestamp     string `json:"timestamp"`
	StatusCode    int    `json:"status_code"`
	Latency       string `json:"latency,omitempty"`
	ClientIP      string `json:"client_ip,omitempty"`
	Method        string `json:"method"`
	Path          string `json:"path"`
	Args          string `json:"args,omitempty"`
	UserAgent     string `json:"user_agent,omitempty"`
	BodySize      int    `json:"body_size,omitempty"`
	RealIP        string `json:"real_ip,omitempty"`
	XForwardedFor string `json:"x_forwarded_for,omitempty"`
	OriginUri     string `json:"origin_uri,omitempty"`
	UserId        string `json:"user_id,omitempty"`
	Error         string `json:"error,omitempty"`
}

func Logger() *zap.Logger {
	once.Do(initLogger)

	if internalLogger == nil {
		log.Fatal("Logger custom not initialized")
	}

	return internalLogger
}

func initLogger() {
	var (
		cfg zap.Config
		err error
	)

	cfg = zap.NewProductionConfig()
	cfg.Encoding = "json"

	if level, exists := os.LookupEnv("LOG_LEVEL"); exists {
		parsedLevel := zap.NewAtomicLevel()
		if err := parsedLevel.UnmarshalText([]byte(level)); err == nil {
			cfg.Level = parsedLevel
		} else {
			fmt.Printf("Invalid log level %s, using default level\n", level)
		}
	}

	internalLogger, err = cfg.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
}

func extractUserFromJWT(req *http.Request) string {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return ""
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		for key, value := range claims {
			if key == "I" {
				if user, ok := value.(float64); ok {
					return strconv.FormatFloat(user, 'f', -1, 64)
				}
			}
		}
	}

	return ""
}

func GinLogFormatter(param gin.LogFormatterParams) string {
	var (
		logFunc func(*zap.Logger, string, ...zap.Field)
		msg     string
	)

	switch {
	case param.StatusCode >= http.StatusInternalServerError:
		msg = "Internal server error"
		logFunc = func(logger *zap.Logger, msg string, fields ...zap.Field) {
			logger.Error(msg, fields...)
		}

	case param.StatusCode >= http.StatusBadRequest:
		msg = "Client error"
		logFunc = func(logger *zap.Logger, msg string, fields ...zap.Field) {
			logger.Warn(msg, fields...)
		}
	default:
		msg = "Request processed successfully"
		logFunc = func(logger *zap.Logger, msg string, fields ...zap.Field) {
			logger.Info(msg, fields...)
		}
	}

	if os.Getenv("DEV_MODE") == "true" {
		return GinDevLogFormatter(param, logFunc, msg)
	} else {
		return GinJSONLogFormatter(param, logFunc, msg)
	}
}

func GinJSONLogFormatter(param gin.LogFormatterParams, logFunc func(*zap.Logger, string, ...zap.Field), msg string) string {
	now := time.Now()
	fields := []zap.Field{
		zap.String("logger", "GIN"),
		zap.String("timestamp", now.Format(time.RFC3339)),
		zap.Int("status_code", param.StatusCode),
		zap.String("latency", param.Latency.String()),
		zap.String("client_ip", param.ClientIP),
		zap.String("method", param.Method),
		zap.String("path", param.Request.URL.Path),
		zap.String("args", param.Request.URL.RawQuery),
		zap.String("user_agent", param.Request.UserAgent()),
		zap.Int("body_size", param.BodySize),
		zap.String("real_ip", param.Request.Header.Get("X-Real-IP")),
		zap.String("x_forwarded_for", param.Request.Header.Get("X-Forwarded-For")),
		zap.String("origin_uri", param.Request.Header.Get("X-Original-Uri")),
		zap.String("user_id", extractUserFromJWT(param.Request)),
		zap.String("error", param.ErrorMessage),
	}

	logFunc(Logger(), msg, fields...)

	return ""
}

func GinDevLogFormatter(param gin.LogFormatterParams, logFunc func(*zap.Logger, string, ...zap.Field), msg string) string {
	now := time.Now()

	logMsg := fmt.Sprintf(
		"[GIN] %s | %3d | %13v | %-7s %#v | params=%s | user_id=%s | %s\n",
		now.Format(time.RFC3339),
		param.StatusCode,
		param.Latency,
		param.Method,
		param.Request.URL.Path,
		param.Request.URL.RawQuery,
		extractUserFromJWT(param.Request),
		param.ErrorMessage,
	)

	return logMsg
}
