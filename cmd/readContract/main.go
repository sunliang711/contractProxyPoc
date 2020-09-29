package main

import (
	"encoding/hex"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	esdk "github.com/sunliang711/eth/sdk"
)

func main() {
	rpcURL := pflag.String("rpc", "", "rpc url")
	fromAddr := pflag.String("fromaddr", "", "from address")
	contractAddress := pflag.String("addr", "", "contract address")
	abiPath := pflag.String("abi", "", "abi file path")
	methodName := pflag.String("method", "", "method name")
	args := pflag.String("args", "", "method args")
	gasPrice := pflag.Uint64("gasprice", 0, "gas price, 0 for suggest gas price")
	gasLimit := pflag.Uint64("gaslimit", 3e8, "gas limit")

	pflag.Parse()

	assertNoEmpty(*rpcURL, "no rpc url")
	assertNoEmpty(*fromAddr, "from address empty")
	assertNoEmpty(*contractAddress, "contract address empty")
	assertNoEmpty(*abiPath, "abi file path empty")
	assertNoEmpty(*methodName, "method name empty")

	abi, err := ioutil.ReadFile(*abiPath)
	if err != nil {
		logrus.Fatalf("read abi file:%v error: %v", abiPath, abi)
	}

	txMan, err := esdk.New(*rpcURL, *gasPrice, *gasLimit, 0, 0)
	if err != nil {
		logrus.Fatalf("new tx manager error: %v", err)
	}
	defer txMan.Close()

	result, err := txMan.ReadContract(*fromAddr, *contractAddress, string(abi), *methodName, *args, *gasPrice, *gasLimit)
	if err != nil {
		logrus.Fatalf("read contract error: %v", err)
	}
	logrus.Infof("read contract result: %v", hex.EncodeToString(result))
}

func assertNoEmpty(val string, msg string) {
	if len(val) == 0 {
		logrus.Fatalf("%s\n", msg)
	}
}
