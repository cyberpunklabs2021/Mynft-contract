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
		ContractOwnAddress:      "0x0c3881df196c01c9",
		SingerAddress:           "0c3881df196c01c9",
		SingerPriv:              "37913a5c7a4632e3f6915b53d1340f68ddd087ac30ccd36cfff9ff5bf659ac4b",
	}
)
