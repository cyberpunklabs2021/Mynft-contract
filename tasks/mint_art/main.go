package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
)

var (
	senderAddress2 = "0c3881df196c01c9"
	senderPriv2    = "37913a5c7a4632e3f6915b53d1340f68ddd087ac30ccd36cfff9ff5bf659ac4b"
)

const mintArt string = `
//testnet
import FungibleToken from 0x9a0766d93b6608b7
import NonFungibleToken from 0x631e88ae7f1d7c20
import Art,Content,Versus from 0x0c3881df196c01c9




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


`

func main() {
	ctx := context.Background()
	flowClient, err := client.New("access.devnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	referenceBlock, err := flowClient.GetLatestBlock(ctx, false) // 取得區塊最新高度以標誌交易中的塊高 ID
	if err != nil {
		panic(err)
	}
	acctAddress, acctKey, signer := getAccount2(flowClient, senderPriv2)

	tx := flow.NewTransaction().
		SetScript([]byte(mintArt)). // 交易要調用的合約
		SetGasLimit(100). // 測試網具體應該多少不知道, 但填100都是會過得
		SetProposalKey(acctAddress, acctKey.Index, acctKey.SequenceNumber). // 會去用就可以了
		SetReferenceBlockID(referenceBlock.ID). // 標記給交易回朔一個區塊ID
		SetPayer(acctAddress). // 支付這筆交易手續費的人, 大部分是自己支付
		AddAuthorizer(acctAddress) // 驗證的簽名者, 大部分是自己驗證

	if err := tx.AddArgument(cadence.NewAddress(flow.HexToAddress(senderAddress2))); err != nil {
		panic(err)
	}
	if err := tx.AddArgument(cadence.NewString("ExampleArtist")); err != nil {
		panic(err)
	}
	if err := tx.AddArgument(cadence.NewString("Example title")); err != nil {
		panic(err)
	}
	if err := tx.AddArgument(cadence.NewString("image url")); err != nil {
		panic(err)
	}
	if err := tx.AddArgument(cadence.NewString("Description")); err != nil {
		panic(err)
	}

	if err := tx.SignEnvelope(acctAddress, acctKey.Index, signer); err != nil {
		panic(err)
	}

	if err := flowClient.SendTransaction(ctx, *tx); err != nil {
		panic(err)
	}

	// 等待交易完成
	// WaitForSeal(ctx, flowClient, tx.ID())
	fmt.Println("Transaction complet!")
	fmt.Println("tx.ID().String() ---- ", tx.ID().String())
}

func getAccount2(flowClient *client.Client, priveKey string) (flow.Address, *flow.AccountKey, crypto.Signer) {
	servicePrivateKeySigAlgo := crypto.StringToSignatureAlgorithm(crypto.ECDSA_P256.String())
	servicePrivateKeyHex := priveKey
	privateKey, err := crypto.DecodePrivateKeyHex(servicePrivateKeySigAlgo, servicePrivateKeyHex)
	if err != nil {
		panic(err)
	}
	addr := flow.HexToAddress(senderAddress2) // 發送者地址轉換成 flow address 格式
	acc, err := flowClient.GetAccount(context.Background(), addr)
	if err != nil {
		panic(err)
	}
	accountKey := acc.Keys[0]                                           // 大部分地址只會有一個 AccountKey, 雖然 flow 支持一個地址可以很多 AccountKey
	signer := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo) // 傳入私鑰及 AccountKey 加密算法按照方式轉換成簽名者
	return addr, accountKey, signer
}

// 發送交易之後寫一個循環
func WaitForSeal2(ctx context.Context, c *client.Client, id flow.Identifier) *flow.TransactionResult {
	result, err := c.GetTransactionResult(ctx, id)
	if err != nil {
		panic(err)
	}
	for result.Status != flow.TransactionStatusSealed {
		time.Sleep(time.Second)
		result, err = c.GetTransactionResult(ctx, id)
		if err != nil {
			panic(err)
		}
	}

	return result
}
