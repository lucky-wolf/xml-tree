package etc

// helper to panic if an error is not nil
// usage: etc.Require(err)
func Require(err error) {
	if err != nil {
		panic(err)
	}
}

// helper to panic if a function returns an error
// can be used with defer to ensure cleanup functions succeed
// usage: defer etc.DeferredRequire(file.Close)()
func DeferredRequire(f func() error) func() {
	return func() {
		Require(f())
	}
}

// helper to set an error if a function returns an error
// can be used with defer to ensure cleanup functions succeed
// usage: defer etc.DeferredError(&err, file.Close)()
func DeferredError(err *error, f func() error) func() {
	if err == nil {
		panic("nil error pointer passed to DeferredError")
	}
	return func() {
		inner := f()
		if *err == nil && inner != nil {
			*err = inner
		}
	}
}
