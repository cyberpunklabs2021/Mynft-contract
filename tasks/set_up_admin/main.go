package main

import (
	"context"
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"google.golang.org/grpc"
)

var (
	senderAddress = "0c3881df196c01c9"
	senderPriv    = "37913a5c7a4632e3f6915b53d1340f68ddd087ac30ccd36cfff9ff5bf659ac4b"
)

const versus1 string = `
import FungibleToken from 0x9a0766d93b6608b7
import NonFungibleToken from 0x631e88ae7f1d7c20
import Content, Art, Auction, Versus from 0x0c3881df196c01c9

transaction() {

    prepare(account: AuthAccount) {
        //create versus admin client
        account.save(<- Versus.createAdminClient(), to:Versus.VersusAdminStoragePath)
        account.link<&{Versus.AdminPublic}>(Versus.VersusAdminPublicPath, target: Versus.VersusAdminStoragePath)
    }
}


`

const versus2 string = `
import FungibleToken from 0x9a0766d93b6608b7
import NonFungibleToken from 0x631e88ae7f1d7c20
import Content, Art, Auction, Versus from 0x0c3881df196c01c9

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
 

`

func main() {
	set2()
}

func set1() {
	ctx := context.Background()
	flowClient, err := client.New("access.devnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	referenceBlock, err := flowClient.GetLatestBlock(context.Background(), false)
	if err != nil {
		panic(err)
	}

	fmt.Println("referenceBlock.Height --- ", referenceBlock.Height)

	acctAddress, acctKey, signer := getSenderInfo(flowClient, senderPriv)
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

	fmt.Println("Transaction complete!")
	fmt.Println("tx ID is ---- ", tx.ID().String())
}

func set2(){
	ctx := context.Background()
	flowClient, err := client.New("access.devnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	referenceBlock, err := flowClient.GetLatestBlock(context.Background(), false)
	if err != nil {
		panic(err)
	}

	fmt.Println("referenceBlock.Height --- ", referenceBlock.Height)

	acctAddress, acctKey, signer := getSenderInfo(flowClient, senderPriv)
	tx := flow.NewTransaction().
		SetScript([]byte(versus2)).
		SetGasLimit(100).
		SetProposalKey(acctAddress, acctKey.Index, acctKey.SequenceNumber).
		SetReferenceBlockID(referenceBlock.ID).
		SetPayer(acctAddress).
		AddAuthorizer(acctAddress)
	tx.AddArgument(cadence.NewAddress(flow.HexToAddress(senderAddress)))
	if err := tx.SignEnvelope(acctAddress, acctKey.Index, signer); err != nil {
		panic(err)
	}

	if err := flowClient.SendTransaction(ctx, *tx); err != nil {
		panic(err)
	}

	fmt.Println("Transaction complete!")
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

	fmt.Println(len(acc.Keys))
	accountKey := acc.Keys[0]
	signer := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	return addr, accountKey, signer
}
