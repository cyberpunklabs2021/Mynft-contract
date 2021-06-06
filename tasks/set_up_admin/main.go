package main

import (
	"context"
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"

	"Mynft-contractf/common"
)

var versus1 string = fmt.Sprintf(
	`
import FungibleToken from %s
import NonFungibleToken from %s
import Content, Art, Auction, %s

transaction() {

    prepare(account: AuthAccount) {
        //create versus admin client
        account.save(<- Versus.createAdminClient(), to:Versus.VersusAdminStoragePath)
        account.link<&{Versus.AdminPublic}>(Versus.VersusAdminPublicPath, target: Versus.VersusAdminStoragePath)
    }
}
`, common.Config.FungibleTokenAddress, common.Config.NonFungibleTokenAddress, common.Config.ContractOwnAddress)

var versus2 string = fmt.Sprintf(
	`
import FungibleToken from %s
import NonFungibleToken from %s
import Content, Art, Auction, Versus from %s

transaction(ownerAddress: Address) {

    //versus account
    prepare(account: AuthAccount) {

        let owner= getAccount(ownerAddress)
        let client= owner.getCapability<&{Versus.AdminPublic}>(Versus.VersusAdminPublicPath)
                .borrow() ?? panic("Could not borrow admin client")

        let versusAdminCap=account.getCapability<&Versus.DropCollection>(Versus.CollectionPrivatePath)
        client.addCapability(versusAdminCap)

    }
}
`, common.Config.FungibleTokenAddress, common.Config.NonFungibleTokenAddress, common.Config.ContractOwnAddress)

func main() {
	set1()
	set2()
}

func set1() {
	ctx := context.Background()
	flowClient, err := client.New(common.Config.Node, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	referenceBlock, err := flowClient.GetLatestBlock(context.Background(), false)
	if err != nil {
		panic(err)
	}

	fmt.Println("referenceBlock.Height --- ", referenceBlock.Height)

	acctAddress, acctKey, signer := common.ServiceAccount(flowClient, common.Config.SingerAddress, common.Config.SingerPriv)
	tx := flow.NewTransaction().
		SetScript([]byte(versus1)).
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

	fmt.Println("Transaction complete!")
	fmt.Println("tx ID is ---- ", tx.ID().String())
}

func set2() {
	ctx := context.Background()
	flowClient, err := client.New(common.Config.Node, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	referenceBlock, err := flowClient.GetLatestBlock(context.Background(), false)
	if err != nil {
		panic(err)
	}

	fmt.Println("referenceBlock.Height --- ", referenceBlock.Height)

	acctAddress, acctKey, signer := common.ServiceAccount(flowClient, common.Config.SingerAddress, common.Config.SingerPriv)
	tx := flow.NewTransaction().
		SetScript([]byte(versus2)).
		SetGasLimit(100).
		SetProposalKey(acctAddress, acctKey.Index, acctKey.SequenceNumber).
		SetReferenceBlockID(referenceBlock.ID).
		SetPayer(acctAddress).
		AddAuthorizer(acctAddress)
	if err := tx.AddArgument(cadence.NewAddress(flow.HexToAddress(common.Config.SingerAddress))); err != nil {
		panic(err)
	}

	if err := tx.SignEnvelope(acctAddress, acctKey.Index, signer); err != nil {
		panic(err)
	}

	if err := flowClient.SendTransaction(ctx, *tx); err != nil {
		panic(err)
	}

	common.WaitForSeal(ctx, flowClient, tx.ID())

	fmt.Println("Transaction complete!")
	fmt.Println("tx ID is ---- ", tx.ID().String())
}
