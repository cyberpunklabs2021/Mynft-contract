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

var mintArt string = fmt.Sprintf(`
import FungibleToken from %s
import NonFungibleToken from %s
import Mynft from %s

transaction(recipient: Address,name: String,artist: String,description: String,arLink: String,ipfsLink: String,MD5Hash: String) {
    let minter: &Mynft.NFTMinter

    prepare(signer: AuthAccount) {
        self.minter = signer.borrow<&Mynft.NFTMinter>(from: Mynft.MinterStoragePath)
            ?? panic("Could not borrow a reference to the NFT minter")
    }

    execute {
        let recipient = getAccount(recipient)

        let receiver = recipient
            .getCapability(Mynft.CollectionPublicPath)!
            .borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Could not get receiver reference to the NFT Collection")

        self.minter.mintNFT(recipient: receiver, name: name,artist:artist,description:description,arLink:arLink,ipfsLink:ipfsLink,MD5Hash:MD5Hash)
    }
}`, common.Config.FungibleTokenAddress, common.Config.NonFungibleTokenAddress, common.Config.ContractOwnAddress)

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

	if err := tx.AddArgument(cadence.NewString("Example Name")); err != nil {
		panic(err)
	}

	if err := tx.AddArgument(cadence.NewString("Example artist")); err != nil {
		panic(err)
	}

	if err := tx.AddArgument(cadence.NewString("Example description")); err != nil {
		panic(err)
	}

	if err := tx.AddArgument(cadence.NewString("Example arLink")); err != nil {
		panic(err)
	}

	if err := tx.AddArgument(cadence.NewString("Example ipfsLink")); err != nil {
		panic(err)
	}

	if err := tx.AddArgument(cadence.NewString("Example MD5Hash")); err != nil {
		panic(err)
	}


	if err := tx.AddArgument(cadence.NewString("Example type")); err != nil {
		panic(err)
	}

	if err := tx.SignEnvelope(acctAddress, acctKey.Index, signer); err != nil {
		panic(err)
	}

	if err := flowClient.SendTransaction(ctx, *tx); err != nil {
		panic(err)
	}

	common.WaitForSeal(ctx, flowClient, tx.ID())
	fmt.Println("Transaction complet!")
	fmt.Println("tx.ID().String() ---- ", tx.ID().String())
}
