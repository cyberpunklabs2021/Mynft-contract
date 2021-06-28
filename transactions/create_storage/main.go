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
import Mynft from %s

transaction {
    prepare(signer: AuthAccount) {
        if signer.borrow<&Mynft.Collection>(from: Mynft.CollectionStoragePath) == nil {
            let collection <- Mynft.createEmptyCollection()
            
            signer.save(<-collection, to: Mynft.CollectionStoragePath)

            signer.link<&Mynft.Collection{NonFungibleToken.CollectionPublic, Mynft.MynftCollectionPublic}>(Mynft.CollectionPublicPath, target: Mynft.CollectionStoragePath)
        }
    }
}`, common.Config.NonFungibleTokenAddress, common.Config.ContractOwnAddress)

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

	// acctAddress, acctKey, signer := common.ServiceAccount(flowClient,"b8daf9d5dad74056", "24a3a149b00de3b26911f17603fba9e5e72281425cae91bd88727659fc86621e")
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

	common.WaitForSeal(ctx, flowClient, tx.ID())
}
