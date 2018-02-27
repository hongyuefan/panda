package chain

type ChainOp interface {
	DoTransaction(string, string, string) (string, error)
	QueryTransaction(string) (int64, error)
	GenKeyPair() (string, string, error)
}
