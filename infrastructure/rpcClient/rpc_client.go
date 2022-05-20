package rpcClient

type Client interface {
	SendRawTransaction(data string) (string, error)
	GetAccountNonce(account, status string) (uint64, error)
	EstimateEnergyPrice() (int64, error)
}
