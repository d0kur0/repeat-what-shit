package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
)

const appDirName = ".repeat-what-shit"

var (
	logFile       *os.File
	fileLogger    *log.Logger
	consoleLogger *log.Logger
)

func GetAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home dir: %w", err)
	}

	appDir := filepath.Join(homeDir, appDirName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create app dir: %w", err)
	}

	return appDir, nil
}

func Init() error {
	consoleLogger = log.New(os.Stdout, "", log.LstdFlags)
	log.SetOutput(io.Discard)

	appDir, err := GetAppDir()
	if err != nil {
		return fmt.Errorf("failed to get app dir: %w", err)
	}

	logPath := filepath.Join(appDir, "app.log")
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	logFile = f
	fileLogger = log.New(f, "", log.LstdFlags)
	return nil
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

func RecoverWithLog() {
	if r := recover(); r != nil {
		stack := debug.Stack()
		fileLogger.Printf("[FATAL] Panic recovered: %v\n%s", r, stack)
		consoleLogger.Printf("[FATAL] Panic recovered: %v", r)
		panic(r)
	}
}

func Debug(format string, v ...interface{}) {
	consoleLogger.Printf("[DEBUG] "+format, v...)
}

func Info(format string, v ...interface{}) {
	consoleLogger.Printf("[INFO] "+format, v...)
}

func Error(format string, v ...interface{}) {
	consoleLogger.Printf("[ERROR] "+format, v...)
}

func Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fileLogger.Printf("[FATAL] %s", msg)
	consoleLogger.Printf("[FATAL] %s", msg)
	os.Exit(1)
}
