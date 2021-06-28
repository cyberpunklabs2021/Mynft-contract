package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

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
	flowClient, err := client.New(common.Config.Node, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	seed := "qreqwrewqrewqrweqrwqerewqrqwrewqrqwewqerewqrewqrqwerewqrqwrewqrbg"

	privateKey1, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, []byte(seed))
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

	if err := createAccountTx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner);err != nil {
		panic(err)
	}

	if err := flowClient.SendTransaction(ctx, *createAccountTx);err != nil {
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
