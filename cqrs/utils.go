package cqrs

func hashKey(x Handler) string {
	return x.Version() + "/" + x.Resource() + "/" + x.Action()
}

type (
	message interface {
		Version() string
		Resource() string
		Action() string
		Key() *string
	}
)

func partitionKey(x message) []byte {
	if x.Key() == nil {
		return nil
	}

	k := x.Resource() + "/" + *x.Key()

	return []byte(k)
}
