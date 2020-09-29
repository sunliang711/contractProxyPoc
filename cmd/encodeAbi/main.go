package main

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	esdk "github.com/sunliang711/eth/sdk"
)

func main() {
	abiPath := pflag.String("abi", "", "abi file path")
	methodName := pflag.String("method", "", "method name")
	args := pflag.String("args", "", "method args")

	pflag.Parse()
	assertNoEmpty(*abiPath, "abi file path empty")

	abi, err := ioutil.ReadFile(*abiPath)
	if err != nil {
		logrus.Fatalf("read abi file: %v error: %v\n", *abiPath, err)
	}
	data, err := esdk.Pack(string(abi), *methodName, *args)
	if err != nil {
		logrus.Fatalf("Pack abi error: %v", err)
	}
	logrus.Infof("pack result: %x", data)

}

func assertNoEmpty(val string, msg string) {
	if len(val) == 0 {
		logrus.Fatalf("%s\n", msg)
	}
}
