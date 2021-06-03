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

var (
	toAddress = "344d25cddb58ed2b"
)

var transferScript string = fmt.Sprintf(`
import FungibleToken from %s
import FlowToken from %s

transaction(amount: UFix64, recipient: Address) {
  let sentVault: @FungibleToken.Vault
  prepare(signer: AuthAccount) {
    let vaultRef = signer.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
      ?? panic("failed to borrow reference to sender vault")

    self.sentVault <- vaultRef.withdraw(amount: amount)
  }

  execute {
    let receiverRef =  getAccount(recipient)
      .getCapability(/public/flowTokenReceiver)
      .borrow<&{FungibleToken.Receiver}>()
        ?? panic("failed to borrow reference to recipient vault")

    receiverRef.deposit(from: <-self.sentVault)
  }
}
`, common.Config.FungibleTokenAddress, common.Config.FlowTokenAddress)

func main() {
	transferFlow()
}

func transferFlow() {
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
		SetScript([]byte(transferScript)).
		SetGasLimit(100).
		SetProposalKey(acctAddress, acctKey.Index, acctKey.SequenceNumber).
		SetReferenceBlockID(referenceBlock.ID).
		SetPayer(acctAddress).
		AddAuthorizer(acctAddress)

	toAmount, err := cadence.NewUFix64("10.0")
	if err != nil {
		panic(err)
	}

	toAddress := cadence.NewAddress(flow.HexToAddress(toAddress))

	if err := tx.AddArgument(toAmount); err != nil {
		panic(err)
	}

	if err := tx.AddArgument(toAddress); err != nil {
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
