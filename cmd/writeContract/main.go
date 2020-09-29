package main

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	esdk "github.com/sunliang711/eth/sdk"
)

func main() {
	rpcURL := pflag.String("rpc", "", "rpc url")
	sk := pflag.String("sk", "", "private key")
	contractAddress := pflag.String("addr", "", "contract address")
	abiPath := pflag.String("abi", "", "abi file path")
	methodName := pflag.String("method", "", "method name")
	args := pflag.String("args", "", "method args")
	gasPrice := pflag.Uint64("gasprice", 0, "gas price, 0 for suggest gas price")
	nonce := pflag.Uint64("nonce", 0, "nonce, 0 for auto nonce")
	gasLimit := pflag.Uint64("gaslimit", 3e8, "gas limit")

	pflag.Parse()

	assertNoEmpty(*rpcURL, "no rpc url")
	assertNoEmpty(*sk, "private key empty")
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

	hash, gasUsed, err := txMan.WriteContractSync(*sk, *contractAddress, string(abi), *methodName, *args, *gasPrice, *nonce, *gasLimit)
	if err != nil {
		logrus.Fatalf("write contract error: %v", err)
	}
	fmt.Printf("hash: %v\ngasUsed: %d\n", hash, gasUsed)
}

func assertNoEmpty(val string, msg string) {
	if len(val) == 0 {
		logrus.Fatalf("%s\n", msg)
	}
}
