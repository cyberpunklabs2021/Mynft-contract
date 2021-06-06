package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"

	"Mynft-contractf/common"
)

var mintArt string = fmt.Sprintf(
	`
//testnet
import FungibleToken from %s
import NonFungibleToken from %s
import Art,Content,Versus from %s

transaction(
    artist: Address,
    artistName: String, 
    artName: String, 
    content: String, 
    description: String) {

    let artistCollection: Capability<&{Art.CollectionPublic}>
    let client: &Versus.Admin

    prepare(account: AuthAccount) {

        self.client = account.borrow<&Versus.Admin>(from: Versus.VersusAdminStoragePath) ?? panic("could not load versus admin")
        self.artistCollection= getAccount(artist).getCapability<&{Art.CollectionPublic}>(Art.CollectionPublicPath)
    }

    execute {
        let art <-  self.client.mintArt(artist: artist, artistName: artistName, artName: artName, content:content, description: description)
        self.artistCollection.borrow()!.deposit(token: <- art)
    }
}
`, common.Config.FungibleTokenAddress, common.Config.NonFungibleTokenAddress, common.Config.ContractOwnAddress)

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
		SetScript([]byte(mintArt)).
		SetGasLimit(100).
		SetProposalKey(acctAddress, acctKey.Index, acctKey.SequenceNumber).
		SetReferenceBlockID(referenceBlock.ID).
		SetPayer(acctAddress).
		AddAuthorizer(acctAddress)

	if err := tx.AddArgument(cadence.NewAddress(flow.HexToAddress(common.Config.SingerAddress))); err != nil {
		panic(err)
	}
	if err := tx.AddArgument(cadence.NewString("ExampleArtist222")); err != nil {
		panic(err)
	}
	if err := tx.AddArgument(cadence.NewString("Example title22")); err != nil {
		panic(err)
	}
	if err := tx.AddArgument(cadence.NewString("image url22")); err != nil {
		panic(err)
	}
	if err := tx.AddArgument(cadence.NewString("Description22")); err != nil {
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
