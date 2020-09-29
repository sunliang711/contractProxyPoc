package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	esdk "github.com/sunliang711/eth/sdk"
)

type ByteCode struct {
	Object    string `json:"object"`
	Opcodes   string `json:"opcodes"`
	SourceMap string `json:"sourceMap"`
}

func main() {
	rpcURL := pflag.String("rpc", "", "rpc url")
	sk := pflag.String("sk", "", "private key")
	bytecodeFile := pflag.StringP("bytecodefile", "b", "", "bytecode file path")
	gasPrice := pflag.Uint64("gasprice", 0, "gas price")
	gasLimit := pflag.Uint64("gaslimit", 3e8, "gas limit")
	output := pflag.StringP("output", "o", "", "output file of contract details")

	pflag.Parse()

	if len(*rpcURL) == 0 {
		logrus.Fatalf("no abi file")
	}

	if len(*bytecodeFile) == 0 {
		logrus.Fatalf("no byte code file path")
	}

	data, err := ioutil.ReadFile(*bytecodeFile)
	if err != nil {
		logrus.Fatalf("read bytecode file:%v error: %v", *bytecodeFile, err)
	}

	var bc ByteCode
	err = json.Unmarshal(data, &bc)
	if err != nil {
		logrus.Fatalf("Unmarshal bytecode file error: %v", err)
	}

	bytecode, err := hex.DecodeString(bc.Object)
	if err != nil {
		logrus.Fatalf("decode bytecode error: %v", err)
	}
	// esdk.CreateContract(*rpcURL, sk, data, *gasPrice, *gaslimit)

	txMan, err := esdk.New(*rpcURL, *gasPrice, *gasLimit, 0, 0)
	if err != nil {
		logrus.Fatalf("New txManager error: %s", err)
	}
	defer txMan.Close()

	address, hash, gasUsed, err := txMan.CreateContractSync(*sk, bytecode, *gasPrice, 0, 0)
	if err != nil {
		logrus.Fatalf("create contract error: %s", err)
	}
	result := fmt.Sprintf("contract address: %s\nhash: %s\ngasUsed: %d\n", address, hash, gasUsed)
	if len(*output) > 0 {
		f, err := os.OpenFile(*output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			logrus.Errorf("open file: %v error: %v\n", *output, err)
			fmt.Print(result)
		} else {
			defer f.Close()
			f.WriteString(result)
		}
	} else {
		fmt.Printf(result)
	}
}
