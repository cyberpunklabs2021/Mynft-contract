package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/onflow/flow-go-sdk/client"

	"github.com/onflow/flow-go-sdk/templates"
	"google.golang.org/grpc"

	"Mynft-contractf/common"
)

func main() {
	ctx := context.Background()
	flowClient, err := client.New(common.Config.Node, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	referenceBlock, err := flowClient.GetLatestBlock(context.Background(), false)
	if err != nil {
		panic(err)
	}
	serviceAcctAddr, serviceAcctKey, singer := common.ServiceAccount(flowClient, common.Config.SingerAddress, common.Config.SingerPriv)

	name := "Content"
	// name := "Art"
	// name := "Auction"
	// name := "Versus"

	contractPath := fmt.Sprintf("../contracts/%s.cdc", name)
	code, err := ioutil.ReadFile(contractPath)
	if err != nil {
		panic(err)
	}
	tx := templates.AddAccountContract(serviceAcctAddr, templates.Contract{
		Name:   name,
		Source: string(code),
	})
	tx.SetProposalKey(
		serviceAcctAddr,
		serviceAcctKey.Index,
		serviceAcctKey.SequenceNumber,
	)
	// we can set the same reference block id. We shouldn't be to far away from it
	tx.SetReferenceBlockID(referenceBlock.ID)
	tx.SetPayer(serviceAcctAddr)
	tx.SetGasLimit(9999)

	if err := tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, singer); err != nil {
		panic(err)
	}

	if err := flowClient.SendTransaction(ctx, *tx); err != nil {
		panic(err)
	}

	common.WaitForSeal(ctx, flowClient, tx.ID())
}
