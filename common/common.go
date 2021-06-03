package common

import (
	"context"
	"fmt"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
)


func WaitForSeal(ctx context.Context, c *client.Client, id flow.Identifier) *flow.TransactionResult {
	result, err := c.GetTransactionResult(ctx, id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Waiting for transaction %s to be sealed...\n", id)

	for result.Status != flow.TransactionStatusSealed {
		time.Sleep(time.Second)
		fmt.Print(".")
		result, err = c.GetTransactionResult(ctx, id)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("%v", result)
	fmt.Printf("Transaction %s sealed\n", id)
	return result
}

func ServiceAccount(flowClient *client.Client, address, privteKey string) (flow.Address, *flow.AccountKey, crypto.Signer) {
	servicePrivateKeySigAlgo := crypto.StringToSignatureAlgorithm(crypto.ECDSA_P256.String())
	servicePrivateKeyHex := privteKey
	privateKey, err := crypto.DecodePrivateKeyHex(servicePrivateKeySigAlgo, servicePrivateKeyHex)
	if err != nil {
		panic(err)
	}

	// addr := flow.ServiceAddress(flow.Emulator)
	addr := flow.HexToAddress(address)
	acc, err := flowClient.GetAccount(context.Background(), addr)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(acc.Keys))
	accountKey := acc.Keys[0]
	signer := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	return addr, accountKey, signer
}

func GetReferenceBlockId(flowClient *client.Client) flow.Identifier {
	block, err := flowClient.GetLatestBlock(context.Background(), false)
	if err != nil {
		panic(err)
	}

	return block.ID
}
