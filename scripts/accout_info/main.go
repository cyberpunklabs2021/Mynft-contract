package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"

	"Mynft-contractf/common"
)

var getArt string = fmt.Sprintf(`
import FungibleToken from %s
import NonFungibleToken from  %s
import Art,Auction,Content from %s

pub struct AddressStatus {

  pub(set) var address:Address
  pub(set) var balance: UFix64
  pub(set) var art: [Art.ArtData]
  init (_ address:Address) {
    self.address=address
    self.balance= 0.0
    self.art= []
  }
}

/*
  This script will check an address and print out its FT, NFT and Versus resources
 */
pub fun main(address:Address) : AddressStatus {
    // get the accounts' public address objects
    let account = getAccount(address)
    let status= AddressStatus(address)
    
    if let vault= account.getCapability(/public/flowTokenBalance).borrow<&{FungibleToken.Balance}>() {
       status.balance=vault.balance
    }

    status.art= Art.getArt(address: address)
    return status
}
`, common.Config.FungibleTokenAddress, common.Config.NonFungibleTokenAddress, common.Config.ContractOwnAddress)
var getContent string = fmt.Sprintf(`
import FungibleToken from %s
import NonFungibleToken from  %s
import Art,Auction,Content from %s

/*
  This script will check an address and print out its FT, NFT and Versus resources
 */
pub fun main(address:Address, artId:UInt64) : String? {
	 let account=getAccount(address)
        if let artCollection= account.getCapability(Art.CollectionPublicPath).borrow<&{Art.CollectionPublic}>()  {
            return artCollection.borrowArt(id: artId)!.content()
        }
     return nil
}
`, common.Config.FungibleTokenAddress, common.Config.NonFungibleTokenAddress, common.Config.ContractOwnAddress)

var (
	searchAddress = "344d25cddb58ed2b"
)

func main() {
	flowClient, err := client.New(common.Config.Node, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	result, err := flowClient.ExecuteScriptAtLatestBlock(ctx, []byte(getArt), []cadence.Value{cadence.NewAddress(flow.HexToAddress(searchAddress))})
	if err != nil {
		panic(err)
	}

	fmt.Println(CadenceValueToJsonString(result))
	result2, err := flowClient.ExecuteScriptAtLatestBlock(ctx, []byte(getContent), []cadence.Value{cadence.NewAddress(flow.HexToAddress(searchAddress)), cadence.NewUInt64(3)})
	if err != nil {
		panic(err)
	}
	fmt.Println(CadenceValueToJsonString(result2))

}

func CadenceValueToJsonString(value cadence.Value) string {
	result := CadenceValueToInterface(value)
	json1, _ := json.MarshalIndent(result, "", "    ")
	return string(json1)
}

func CadenceValueToInterface(field cadence.Value) interface{} {
	switch field.(type) {
	case cadence.Dictionary:
		result := map[string]interface{}{}
		for _, item := range field.(cadence.Dictionary).Pairs {
			result[item.Key.String()] = CadenceValueToInterface(item.Value)
		}
		return result
	case cadence.Struct:
		result := map[string]interface{}{}
		subStructNames := field.(cadence.Struct).StructType.Fields
		for j, subField := range field.(cadence.Struct).Fields {
			result[subStructNames[j].Identifier] = CadenceValueToInterface(subField)
		}
		return result
	case cadence.Array:
		result := []interface{}{}
		for _, item := range field.(cadence.Array).Values {
			result = append(result, CadenceValueToInterface(item))
		}
		return result
	default:
		result, err := strconv.Unquote(field.String())
		if err != nil {
			return field.String()
		}
		return result
	}
}
