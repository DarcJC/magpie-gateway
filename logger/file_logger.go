package logger

import (
    "github.com/gin-gonic/gin"
    "io"
    "log"
    "os"
)

type FileLogger struct {
    filename string
}

func NewFileLogger(filename string) FileLogger {
    return FileLogger{filename: filename}
}

func (f *FileLogger) GetLoggerConfig() gin.LoggerConfig {
    file, err := os.OpenFile(f.filename, os.O_CREATE | os.O_APPEND, 0644)
    if err != nil {
        log.Fatal(err)
    }

    return gin.LoggerConfig{
        Output: io.MultiWriter(file, os.Stdout),
    }
}
