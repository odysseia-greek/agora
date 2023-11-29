package logging

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	white    = "\033[90;47m"
	purple   = "\033[0;35m"
	cyan     = "\033[0;36m"
	red      = "\033[0;31m"
	green    = "\033[0;32m"
	yellow   = "\033[0;33m"
	blue     = "\033[0;34m"
	reset    = "\033[0m"
	whiteBg  = "\033[47;97m" // White background, bright white text
	purpleBg = "\033[45;97m" // Purple background, bright white text
	cyanBg   = "\033[46;97m" // Cyan background, bright white text
	redBg    = "\033[41;97m" // Red background, bright white text
	greenBg  = "\033[42;97m" // Green background, bright white text
	yellowBg = "\033[43;97m" // Yellow background, bright white text
	blueBg   = "\033[44;97m" // Blue background, bright white text
	resetBg  = "\033[0m"
)

func init() {
	// Set log prefix to empty to avoid timestamp and other information
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

func Debug(message string) {
	log.Print(formatLog("DEBUG", message, blue))
}

func Info(message string) {
	log.Print(formatLog("INFO", message, green))
}

func Api(statusCode int, method, traceId, clientIp, path string, latency time.Duration) {
	log.Print(formatApiLog(statusCode, method, clientIp, path, traceId, latency))
}

func Warn(message string) {
	log.Print(formatLog("WARN", message, yellow))
}

func System(message string) {
	log.Print(formatLog("SYSTEM", message, purple))
}

func Trace(message string) {
	log.Print(formatLog("TRACE", message, cyan))
}

func Error(message string) {
	log.Print(formatLog("ERROR", message, red))
}

func formatLog(level, message, color string) string {
	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s%s | %s | %s%s", color, level, timestamp, message, reset)
}
func formatApiLog(statusCode int, method, clientIp, path, traceId string, latency time.Duration) string {
	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")
	statusColour := getColour(statusCode)
	methodColour := getMethodColor(method)
	return fmt.Sprintf("[API] %v |%s %3d %s| %v | %15s |%s %-8s %s %s %s\n",
		timestamp,
		statusColour, statusCode, reset,
		latency,
		clientIp,
		methodColour, method, reset,
		path,
		traceId,
	)
}
func getColour(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return greenBg
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return whiteBg
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellowBg
	default:
		return redBg
	}
}

func getMethodColor(method string) string {
	switch method {
	case http.MethodGet:
		return blueBg
	case http.MethodPost:
		return cyanBg
	case http.MethodPut:
		return yellowBg
	case http.MethodDelete:
		return redBg
	case http.MethodPatch:
		return greenBg
	case http.MethodHead:
		return purpleBg
	case http.MethodOptions:
		return whiteBg
	default:
		return reset
	}
}
