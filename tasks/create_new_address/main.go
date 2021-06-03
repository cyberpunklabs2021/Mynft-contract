package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"

	"github.com/onflow/flow-go-sdk/templates"
	"google.golang.org/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"

	"Mynft-contractf/common"
)

func main() {
	CreateAccount()
}

func CreateAccount() {
	ctx := context.Background()
	// server
	flowClient, err := client.New(common.Config.Node, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	// random seed
	seed := make([]byte, crypto.MinSeedLength)
	_, err = rand.Read(seed)
	if err != nil {
		panic(err)
	}

	privateKey1, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
	if err != nil {
		panic(err)
	}

	log.Println("privateKey====", hex.EncodeToString(privateKey1.Encode()))

	myAcctKey := flow.NewAccountKey().
		SetPublicKey(privateKey1.PublicKey()).
		SetSigAlgo(privateKey1.Algorithm()).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)
	serviceAcctAddr, serviceAcctKey, serviceSigner := common.ServiceAccount(flowClient, common.Config.SingerAddress, common.Config.SingerPriv)
	referenceBlockID := common.GetReferenceBlockId(flowClient)
	createAccountTx := templates.CreateAccount([]*flow.AccountKey{myAcctKey}, nil, serviceAcctAddr)
	createAccountTx.SetProposalKey(
		serviceAcctAddr,
		serviceAcctKey.Index,
		serviceAcctKey.SequenceNumber,
	)
	createAccountTx.SetReferenceBlockID(referenceBlockID)
	createAccountTx.SetPayer(serviceAcctAddr)

	// All new accounts must be created by an existing account
	err = createAccountTx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	if err != nil {
		panic(err)
	}

	// Send the transaction to the network
	err = flowClient.SendTransaction(ctx, *createAccountTx)
	if err != nil {
		panic(err)
	}

	accountCreationTxRes := common.WaitForSeal(ctx, flowClient, createAccountTx.ID())
	var myAddress flow.Address
	for _, event := range accountCreationTxRes.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent := flow.AccountCreatedEvent(event)
			myAddress = accountCreatedEvent.Address()
		}
	}

	fmt.Println("Account created with address:", myAddress.Hex())
}
