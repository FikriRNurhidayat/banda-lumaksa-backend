package transaction

type TransactionSpecification interface {
	Call(transaction Transaction) bool
}
