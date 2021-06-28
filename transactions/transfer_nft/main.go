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

var (
	toAddress = "344d25cddb58ed2b"

	transfCode = fmt.Sprintf(`
	import NonFungibleToken from %s
	import KittyItems from %s

	transaction(recipient: Address, withdrawID: UInt64) {
    prepare(signer: AuthAccount) {
        
        // get the recipients public account object
        let recipient = getAccount(recipient)

        // borrow a reference to the signer's NFT collection
        let collectionRef = signer.borrow<&KittyItems.Collection>(from: KittyItems.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

        // borrow a public reference to the receivers collection
        let depositRef = recipient.getCapability(KittyItems.CollectionPublicPath)!.borrow<&{NonFungibleToken.CollectionPublic}>()!

        // withdraw the NFT from the owner's collection
        let nft <- collectionRef.withdraw(withdrawID: withdrawID)

        // Deposit the NFT in the recipient's collection
        depositRef.deposit(token: <-nft)
    }
}
`,  common.Config.NonFungibleTokenAddress, common.Config.ContractOwnAddress)
)

func main() {
	ctx := context.Background()
	flowClient, err := client.New(common.Config.Node, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	referenceBlock, err := flowClient.GetLatestBlock(ctx, false) // 取得區塊最新高度以標誌交易中的塊高 ID
	if err != nil {
		panic(err)
	}

	acctAddress, acctKey, signer := common.ServiceAccount(flowClient, common.Config.SingerAddress, common.Config.SingerPriv)

	tx := flow.NewTransaction().
		SetScript([]byte(transfCode)). // 交易要調用的合約
		SetGasLimit(100). // 測試網具體應該多少不知道, 但填100都是會過得
		SetProposalKey(acctAddress, acctKey.Index, acctKey.SequenceNumber). // 會去用就可以了
		SetReferenceBlockID(referenceBlock.ID). // 標記給交易回朔一個區塊ID
		SetPayer(acctAddress). // 支付這筆交易手續費的人, 大部分是自己支付
		AddAuthorizer(acctAddress) // 驗證的簽名者, 大部分是自己驗證

	if err := tx.AddArgument(cadence.NewAddress(flow.HexToAddress(toAddress))); err != nil {
		panic(err)
	}

	if err := tx.AddArgument(cadence.NewUInt64(1)); err != nil {
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
