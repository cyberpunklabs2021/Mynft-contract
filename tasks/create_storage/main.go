package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"

	"Mynft-contractf/common"
)

var createStorage string = fmt.Sprintf(`
import NonFungibleToken from %s
import Art from %s
transaction() {
  prepare(signer: AuthAccount) {
    if signer.borrow<&Art.Collection>(from: Art.CollectionStoragePath) == nil {
      signer.save(<-Art.createEmptyCollection(), to: Art.CollectionStoragePath)
      signer.link<&Art.Collection{Art.CollectionPublic}>(
        Art.CollectionPublicPath,
        target: Art.CollectionStoragePath
      )
    }
  }
}
`, common.Config.NonFungibleTokenAddress, common.Config.ContractOwnAddress)

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
		SetScript([]byte(createStorage)).
		SetGasLimit(100).
		SetProposalKey(acctAddress, acctKey.Index, acctKey.SequenceNumber).
		SetReferenceBlockID(referenceBlock.ID).
		SetPayer(acctAddress).
		AddAuthorizer(acctAddress)

	if err := tx.SignEnvelope(acctAddress, acctKey.Index, signer); err != nil {
		panic(err)
	}

	if err := flowClient.SendTransaction(ctx, *tx); err != nil {
		panic(err)
	}

	// 等待交易完成
	common.WaitForSeal(ctx, flowClient, tx.ID())
}
