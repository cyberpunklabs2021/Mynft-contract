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
)

const getInfo string = `
import FungibleToken from 0x9a0766d93b6608b7
import NonFungibleToken from  0x631e88ae7f1d7c20
import Art,Auction,Content from 0x0c3881df196c01c9

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

`

var (
	senderAddress = "0c3881df196c01c9"
)

func main() {
	flowClient, err := client.New("access.devnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	result, err := flowClient.ExecuteScriptAtLatestBlock(ctx, []byte(getInfo), []cadence.Value{cadence.NewAddress(flow.HexToAddress(senderAddress))})
	if err != nil {
		panic(err)
	}

	fmt.Println(CadenceValueToJsonString(result))
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

	// dictionaryValue, isDictionary := field.(cadence.Dictionary)
	// structValue, isStruct := field.(cadence.Struct)
	// arrayValue, isArray := field.(cadence.Array)
	// if isStruct {
	// 	subStructNames := structValue.StructType.Fields
	// 	result := map[string]interface{}{}
	// 	for j, subField := range structValue.Fields {
	// 		result[subStructNames[j].Identifier] = CadenceValueToInterface(subField)
	// 	}
	// 	return result
	// } else if isDictionary {
	// 	result := map[string]interface{}{}
	// 	for _, item := range dictionaryValue.Pairs {
	// 		result[item.Key.String()] = CadenceValueToInterface(item.Value)
	// 	}
	// 	return result
	// } else if isArray {
	// 	result := []interface{}{}
	// 	for _, item := range arrayValue.Values {
	// 		result = append(result, CadenceValueToInterface(item))
	// 	}
	// 	return result
	// }
	// result, err := strconv.Unquote(field.String())
	// if err != nil {
	// 	return field.String()
	// }
	// return result
}
