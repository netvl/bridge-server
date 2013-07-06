/**
 * Date: 06.07.13
 * Time: 22:35
 *
 * @author Vladimir Matveev
 */
package log

import (
    slog "log"
    "fmt"
    "bytes"
)

type Exception interface {
    error
    Cause() error
}

func formatException(message bytes.Buffer, arg interface{}, addColon bool) {
    if addColon {
        message.WriteString(": ")
    }

    if e, ok := arg.(error); ok {
        message.WriteString(e.Error())
        for e != nil {
            message.WriteString("\n  Caused by: ")
            message.WriteString(e.Error())
            if ex, ok := e.(Exception); ok {
                e = ex.Cause()
            } else {
                e = nil
            }
        }
    } else {
        message.WriteString(fmt.Sprint(arg))
    }
}

func formatMessage(message bytes.Buffer, args ...interface{}) {
    if len(args) > 0 {
        last := args[len(args)-1]
        if _, ok := last.(error); ok {
            fmt.Fprint(message, args[:len(args)-1])  // all except the last one
            formatException(message, last, len(args) > 1)
        } else {
            fmt.Fprint(message, args)
        }
    }
}

func formatMessageln(message bytes.Buffer, args ...interface{}) {
    if len(args) > 0 {
        last := args[len(args)-1]
        if _, ok := last.(error); ok {
            fmt.Fprintln(message, args[:len(args)-1])  // all except the last one
            formatException(message, last, len(args) > 1)
        } else {
            fmt.Fprintln(message, args)
        }
    }
}

func formatMessagef(message bytes.Buffer, format string, args ...interface{}) {
    if len(args) > 0 {
        last := args[len(args)-1]
        if _, ok := last.(error); ok {
            fmt.Fprintf(message, format, args[:len(args)-1])
            formatException(message, last, len(args) > 1)
        } else {
            fmt.Fprintf(message, format, args)
        }
    }
}

func prepareLog(level string) (message bytes.Buffer) {
    message.WriteString("[")
    message.WriteString(level)
    message.WriteString("] ")
    return
}

func log(level string, args ...interface{}) {
    message := prepareLog(level)
    formatMessage(message, args)
    slog.Print(message.String())
}

func logln(level string, args ...interface{}) {
    message := prepareLog(level)
    formatMessageln(message, args)
    slog.Print(message.String())
}

func logf(level string, format string, args ...interface{}) {
    message := prepareLog(level)
    formatMessagef(message, format, args)
    slog.Print(message.String())
}

func Info(args ...interface{}) {
    log("INFO", args)
}

func Infoln(args ...interface{}) {
    logln("INFO", args)
}

func Infof(format string, args ...interface{}) {
    logf("INFO", format, args)
}

func Error(args ...interface{}) {
    log("ERROR", args)
}

func Errorln(args ...interface{}) {
    logln("ERROR", args)
}

func Errorf(format string, args ...interface{}) {
    logf("ERROR", format, args)
}

