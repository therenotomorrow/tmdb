package fp

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func Silent[T any](_ T, _ error) {}
