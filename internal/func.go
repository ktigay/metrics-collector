package internal

func Quite(f func() error) {
	_ = f()
}
