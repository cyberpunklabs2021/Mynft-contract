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

var getNftInfo string = fmt.Sprintf(`
import FungibleToken from %s
import NonFungibleToken from  %s
import Mynft from %s

pub struct AddressStatus {

  pub(set) var address:Address
  pub(set) var balance: UFix64
  pub(set) var nft: [Mynft.NftData]
  init (_ address:Address) {
    self.address=address
    self.balance= 0.0
    self.nft= []
  }
}

pub fun main(address:Address) : AddressStatus {
    // get the accounts' public address objects
    let account = getAccount(address)
    let status= AddressStatus(address)
    
    if let vault= account.getCapability(/public/flowTokenBalance).borrow<&{FungibleToken.Balance}>() {
       status.balance=vault.balance
    }

    status.nft= Mynft.getNft(address: address)
    return status
}`, common.Config.FungibleTokenAddress, common.Config.NonFungibleTokenAddress, common.Config.ContractOwnAddress)

var (
	searchAddress = "b8daf9d5dad74056"
)

func main() {
	flowClient, err := client.New(common.Config.Node, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	result, err := flowClient.ExecuteScriptAtLatestBlock(ctx, []byte(getNftInfo), []cadence.Value{cadence.NewAddress(flow.HexToAddress(searchAddress))})
	if err != nil {
		panic(err)
	}

	fmt.Println(common.CadenceValueToJsonString(result))
}
