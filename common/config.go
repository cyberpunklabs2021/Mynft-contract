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
		Node:                    "",
		FungibleTokenAddress:    "",
		NonFungibleTokenAddress: "",
		FlowTokenAddress:        "",
		ContractOwnAddress:      "",
		SingerAddress:           "",
		SingerPriv:              "",
	}
)
