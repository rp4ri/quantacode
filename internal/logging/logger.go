package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type LogEntry struct {
	Timestamp string      `json:"timestamp"`
	Level     string      `json:"level"`
	Module    string      `json:"module"`
	Message   string      `json:"message"`
	Input     interface{} `json:"input,omitempty"`
	Output    interface{} `json:"output,omitempty"`
	Error     string      `json:"error,omitempty"`
	Duration  string      `json:"duration,omitempty"`
}

type Logger struct {
	mu       sync.Mutex
	file     *os.File
	writer   io.Writer
	minLevel LogLevel
	module   string
}

var (
	globalLogger *Logger
	once         sync.Once
)

func Init(logPath string, minLevel LogLevel) error {
	var initErr error
	once.Do(func() {
		dir := filepath.Dir(logPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			initErr = fmt.Errorf("create log dir: %w", err)
			return
		}

		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			initErr = fmt.Errorf("open log file: %w", err)
			return
		}

		globalLogger = &Logger{
			file:     file,
			writer:   file,
			minLevel: minLevel,
			module:   "global",
		}
	})
	return initErr
}

func GetLogger(module string) *Logger {
	if globalLogger == nil {
		return &Logger{
			writer:   io.Discard,
			minLevel: DEBUG,
			module:   module,
		}
	}
	return &Logger{
		file:     globalLogger.file,
		writer:   globalLogger.writer,
		minLevel: globalLogger.minLevel,
		module:   module,
	}
}

func (l *Logger) log(level LogLevel, msg string, input, output interface{}, err error, duration time.Duration) {
	if level < l.minLevel {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Level:     level.String(),
		Module:    l.module,
		Message:   msg,
		Input:     input,
		Output:    output,
	}

	if err != nil {
		entry.Error = err.Error()
	}
	if duration > 0 {
		entry.Duration = duration.String()
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	data, _ := json.Marshal(entry)
	fmt.Fprintln(l.writer, string(data))
}

func (l *Logger) Debug(msg string) {
	l.log(DEBUG, msg, nil, nil, nil, 0)
}

func (l *Logger) Info(msg string) {
	l.log(INFO, msg, nil, nil, nil, 0)
}

func (l *Logger) Warn(msg string) {
	l.log(WARN, msg, nil, nil, nil, 0)
}

func (l *Logger) Error(msg string, err error) {
	l.log(ERROR, msg, nil, nil, err, 0)
}

func (l *Logger) LogIO(msg string, input, output interface{}, duration time.Duration) {
	l.log(INFO, msg, input, output, nil, duration)
}

func (l *Logger) LogOpenRouterCall(input interface{}, output interface{}, err error, duration time.Duration) {
	level := INFO
	if err != nil {
		level = ERROR
	}
	l.log(level, "OpenRouter API call", input, output, err, duration)
}

// LogPriceUpdate logs incoming price updates from Binance
func (l *Logger) LogPriceUpdate(symbol string, price, volume float64) {
	l.log(DEBUG, "Price update", map[string]interface{}{
		"symbol": symbol,
		"price":  price,
		"volume": volume,
	}, nil, nil, 0)
}

// LogIndicatorUpdate logs calculated indicator values
func (l *Logger) LogIndicatorUpdate(rsi, sma, ema float64) {
	l.log(DEBUG, "Indicator update", map[string]interface{}{
		"rsi": rsi,
		"sma": sma,
		"ema": ema,
	}, nil, nil, 0)
}

// LogConnection logs connection events
func (l *Logger) LogConnection(event string, details map[string]interface{}) {
	l.log(INFO, event, details, nil, nil, 0)
}

// LogPairSwitch logs when user switches trading pairs
func (l *Logger) LogPairSwitch(oldPair, newPair string) {
	l.log(INFO, "Pair switch", map[string]interface{}{
		"from": oldPair,
		"to":   newPair,
	}, nil, nil, 0)
}

func Close() {
	if globalLogger != nil && globalLogger.file != nil {
		globalLogger.file.Close()
	}
}
