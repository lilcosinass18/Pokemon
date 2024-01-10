package dbutil

type AnyFetcher interface {
	Getter
	Selecter
}

type DB interface {
	AnyFetcher
	Execer
}
