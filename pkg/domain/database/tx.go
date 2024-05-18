package database

type ROTx interface {
	ROTxImpl()
}

type RWTx interface {
	ROTx
	RWTxImpl()
}
