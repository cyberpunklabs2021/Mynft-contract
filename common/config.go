package common

var (
	// Test
	Config = struct {
		Node                    string
		FungibleTokenAddress    string
		NonFungibleTokenAddress string
		FlowTokenAddress        string

		ContractOwnAddress string
		SingerAddress      string
		SingerPriv         string
	}{
		Node:                    "access.devnet.nodes.onflow.org:9000",
		FungibleTokenAddress:    "0x9a0766d93b6608b7",
		NonFungibleTokenAddress: "0x631e88ae7f1d7c20",
		FlowTokenAddress:        "0x7e60df042a9c0868",
		ContractOwnAddress:      "0xdd41871e37c4240a",
		SingerAddress:           "dd41871e37c4240a",
		SingerPriv:              "19f0f7cc0f288e511179b276c0ea9ea4bb1ed10cc16f85d71e0639bbf2ce3db8",
	}
)
