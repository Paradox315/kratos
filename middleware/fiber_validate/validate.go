package fiber_validate

type validator interface {
	Validate() error
}
