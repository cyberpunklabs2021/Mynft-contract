package main

import (
	"context"
	"fmt"

	"github.com/onflow/cadence"

	"google.golang.org/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"

	"Mynft-contractf/common"
)

var (
	removeCode = `
transaction(name: String) {
    prepare(signer: AuthAccount) {
        signer.contracts.remove(name: name)
    }
}
`
)

func main() {
	ctx := context.Background()
	flowClient, err := client.New(common.Config.Node, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	referenceBlock, err := flowClient.GetLatestBlock(ctx, false)
	if err != nil {
		panic(err)
	}

	acctAddress, acctKey, signer := common.ServiceAccount(flowClient, common.Config.SingerAddress, common.Config.SingerPriv)

	tx := flow.NewTransaction().
		SetScript([]byte(removeCode)).
		SetGasLimit(100).
		SetProposalKey(acctAddress, acctKey.Index, acctKey.SequenceNumber).
		SetReferenceBlockID(referenceBlock.ID).
		SetPayer(acctAddress).
		AddAuthorizer(acctAddress)

	if err := tx.AddArgument(cadence.NewString("Auction")); err != nil {
		panic(err)
	}

	if err := tx.SignEnvelope(acctAddress, acctKey.Index, signer); err != nil {
		panic(err)
	}

	if err := flowClient.SendTransaction(ctx, *tx); err != nil {
		panic(err)
	}

	// 等待交易完成
	common.WaitForSeal(ctx, flowClient, tx.ID())

	fmt.Println("Transaction complet!")
	fmt.Println("tx.ID().String() ---- ", tx.ID().String())
}
