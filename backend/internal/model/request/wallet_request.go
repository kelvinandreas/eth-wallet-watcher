package request

type WalletRequest struct {
	Address string `json:"address"`
	Label   string `json:"label"`
}
