package anyswap

import "app/utils"

const (
	// LogAnySwapOut (index_topic_1 address token, index_topic_2 address from, index_topic_3 address to, uint256 amount, uint256 fromChainID, uint256 toChainID)
	LogAnySwapOut = "0x97116cf6cd4f6412bb47914d6db18da9e16ab2142f543b86e207c24fbd16b23a"

	// LogAnySwapIn (index_topic_1 bytes32 txhash, index_topic_2 address token, index_topic_3 address to, uint256 amount, uint256 fromChainID, uint256 toChainID)
	LogAnySwapIn = "0xaac9ce45fe3adf5143598c4f18a369591a20a3384aedaf1b525d29127e1fcd55"

	// underlying()
	Underlying = "0x6f307dc3"
)

// LogAnySwapOut (index_topic_1 address token, index_topic_2 address from, index_topic_3 address to, uint256 amount, uint256 fromChainID, uint256 toChainID)

var AnyswapContracts = map[string][]string{
	"eth": {
		"0x6b7a87899490ece95443e979ca9485cbe7e71522",
		"0xba8da9dcf11b50b03fd5284f164ef5cdef910705",
		"0x765277eebeca2e31912c9946eae1021199b39c61",
		"0x7782046601e7b9b05ca55a3899780ce6ee6b8b2b",
		"0xe95fd76cf16008c12ff3b3a937cb16cd9cc20284",
	},
	"bsc": {
		"0xd1c5966f9f5ee6881ff6b261bbeda45972b1b5f3",
		"0xabd380327fe66724ffda91a87c772fb8d00be488",
	},
	"polygon": {
		"0x4f3aff3a747fcade12598081e80c6605a8be192f",
		"0x2eF4A574b72E1f555185AfA8A09c6d1A8AC4025C",
	},
	"fantom": {
		"0x1ccca1ce62c62f7be95d4a67722a8fdbed6eecb4",
	},
	"arbitrum": {
		"0x0cae51e1032e8461f4806e26332c030e34de3adb",
		"0xC931f61B1534EB21D8c11B24f3f5Ab2471d4aB50",
		"0x650Af55D5877F289837c30b94af91538a7504b76",
	},
	"avalanche": {
		"0xB0731d50C681C45856BFc3f7539D5f61d4bE81D8",
		"0x833f307ac507d47309fd8cdd1f835bef8d702a93",
	},
	"optimism": {
		"0x80A16016cC4A2E6a2CACA8a4a498b1699fF0f844",
		"0xDC42728B0eA910349ed3c6e1c9Dc06b5FB591f98",
	},
}

var AnyTokens = map[string]map[string]map[string]string{
	"USDC": {
		"eth": {
			"underlyingToken": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			"anyToken_1":      "0x7EA2be2df7BA6E54B1A9C70676f668455E329d29",
			"anyToken_2":      "0xeA928a8d09E11c66e074fBf2f6804E19821F438D",
			"anyToken_3":      "0x2cb1712fa24aBc7Ce787b8853235C86e38ACca44",
		},
		"bsc": {
			"underlyingToken": "0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d",
			"anyToken":        "0x8965349fb649A33a30cbFDa057D8eC2C48AbE2A2",
			"anyToken_2":      "0xab6290bBd5C2d26881E8A7a10bC98552B9082E7f",
		},
		"avanlanche": {
			"underlyingToken": "0xA7D7079b0FEaD91F3e65f86E8915Cb59c1a4C664",
			"anyToken":        "0xcc9b1F919282c255eB9AD2C0757E8036165e0cAd",
		},
		"polygon": {
			"underlyingToken": "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174",
			"anyToken":        "0xd69b31c3225728CC57ddaf9be532a4ee1620Be51",
		},
		"fantom": {
			"underlyingToken": "0x04068DA6C83AFCFA0e13ba15A6696662335D5B75",
			"anyToken":        "0x95bf7E307BC1ab0BA38ae10fc27084bC36FcD605",
		},
		"arbitrum": {
			"underlyingToken": "0x04068DA6C83AFCFA0e13ba15A6696662335D5B75",
			"anyToken":        "0x3405A1bd46B85c5C029483FbECf2F3E611026e45",
		},
		"optimism": {
			"underlyingToken": "0x7F5c764cBc14f9669B88837ca1490cCa17c31607",
			"anyToken":        "0xf390830DF829cf22c53c8840554B98eafC5dCBc2",
		},
	},
}

func init() {
	for name, chain := range AnyswapContracts {
		AnyswapContracts[name] = utils.StrSliceToLower(chain)
	}
}

type Detail struct {
	SrcTxHash string `json:"src_tx_hash,omitempty"`
}
