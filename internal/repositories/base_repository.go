package repositories

// RepositoryInterface defines common methods for all repositories
type RepositoryInterface interface {
	WithTx(tx interface{}) RepositoryInterface
	GetDB() interface{}
}
