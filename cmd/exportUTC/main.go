package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	esdk "github.com/sunliang711/eth/sdk"
)

func main() {
	utcFile := pflag.StringP("utc", "u", "", "utc file stored in datadir/keystore")
	lvl := pflag.StringP("level", "l", "info", "log level")
	password := pflag.StringP("password", "p", "", "password of utc file")

	pflag.Parse()

	logrus.SetLevel(logLevel(*lvl))
	if len(*utcFile) == 0 {
		logrus.Fatalln("Not specify utc file")
	}

	logrus.Infof("utc file: %s\n", *utcFile)
	account, err := esdk.ExportAccountObject(*utcFile, *password)
	if err != nil {
		logrus.Fatalf("Export utc file error: %v", err)
	}

	fmt.Printf("account: %+v", *account)
}

func logLevel(lvl string) logrus.Level {
	lvl = strings.ToLower(lvl)
	switch lvl {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.ErrorLevel
	}
}
