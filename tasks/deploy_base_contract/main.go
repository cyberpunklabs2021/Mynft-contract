package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"

	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
	"google.golang.org/grpc"
)

var (
	// senderAddress = "0c3881df196c01c9"
	// senderPriv    = "37913a5c7a4632e3f6915b53d1340f68ddd087ac30ccd36cfff9ff5bf659ac4b"
	senderAddress = "b8daf9d5dad74056"
	senderPriv    = "24a3a149b00de3b26911f17603fba9e5e72281425cae91bd88727659fc86621e"
)

func main() {
	ctx := context.Background()
	flowClient, err := client.New("access.devnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	referenceBlock, err := flowClient.GetLatestBlock(context.Background(), false)
	if err != nil {
		panic(err)
	}
	serviceAcctAddr, serviceAcctKey, singer := getSenderInfo(flowClient, senderPriv)

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

	// tx.AddAuthorizer(serviceAcctAddr)
	if err := tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, singer); err != nil {
		panic(err)
	}

	if err := flowClient.SendTransaction(ctx, *tx); err != nil {
		panic(err)
	}
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

	accountKey := acc.Keys[0]
	signer := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	return addr, accountKey, signer
}
func NewAccountKey(acckKey *flow.AccountKey) *flow.AccountKey {
	return flow.NewAccountKey().
		SetPublicKey(acckKey.PublicKey).
		SetSigAlgo(acckKey.SigAlgo).
		SetHashAlgo(acckKey.HashAlgo).
		SetWeight(flow.AccountKeyWeightThreshold)
}
