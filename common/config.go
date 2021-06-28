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
		ContractOwnAddress:      "0xbe1c37911f64d57e",
		SingerAddress:           "be1c37911f64d57e",
		SingerPriv:              "c16385c9b7543c98d2c18c2af839a3e4b858a323b600ef718ea8c26ec9e0b091",
	}
)
