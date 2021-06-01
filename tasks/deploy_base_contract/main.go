package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"

	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
	"google.golang.org/grpc"
)

var (
	senderAddress = "b8daf9d5dad74056"
	senderPriv    = "24a3a149b00de3b26911f17603fba9e5e72281425cae91bd88727659fc86621e"
)

func main() {
	ctx := context.Background()
	flowClient, err := client.New("access.devnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	referenceBlock, err := flowClient.GetLatestBlock(context.Background(), false)
	if err != nil {
		panic(err)
	}
	serviceAcctAddr, serviceAcctKey, singer := getSenderInfo(flowClient, senderPriv)
	var contracts []templates.Contract
	// contractPath2 := "/Users/jay/Desktop/selfcode/Mynft-contract/contracts/Content.cdc"
	// code2, err := ioutil.ReadFile(contractPath2)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// contracts = append(contracts, templates.Contract{
	// 	Name:   "Content",
	// 	Source: string(code2),
	// })

	contractPath := "/Users/jay/Desktop/selfcode/Mynft-contract/contracts/Art.cdc"
	code, err := ioutil.ReadFile(contractPath)
	if err != nil {
		panic(err)
	}

	contracts = append(contracts, templates.Contract{
		Name:   "Art",
		Source: string(code),
	})


	tx := templates.CreateAccount(nil, contracts, serviceAcctAddr)
	tx.SetProposalKey(
		serviceAcctAddr,
		serviceAcctKey.Index,
		serviceAcctKey.SequenceNumber,
	)
	// we can set the same reference block id. We shouldn't be to far away from it
	tx.SetReferenceBlockID(referenceBlock.ID)
	tx.SetPayer(serviceAcctAddr)

	if err := tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, singer); err != nil {
		panic(err)
	}

	if err := flowClient.SendTransaction(ctx, *tx); err != nil {
		panic(err)
	}
	fmt.Println("tx ID is ---- ", tx.ID().String())
}

func getSenderInfo(flowClient *client.Client, privKeyStr string) (flow.Address, *flow.AccountKey, crypto.Signer) {
	privateKeySigAlgo := crypto.StringToSignatureAlgorithm(crypto.ECDSA_P256.String())
	privateKey, err := crypto.DecodePrivateKeyHex(privateKeySigAlgo, privKeyStr)
	if err != nil {
		panic(err)
	}

	addr := flow.HexToAddress(senderAddress)
	acc, err := flowClient.GetAccount(context.Background(), addr)
	if err != nil {
		panic(err)
	}

	accountKey := acc.Keys[0]
	signer := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	return addr, accountKey, signer
}
