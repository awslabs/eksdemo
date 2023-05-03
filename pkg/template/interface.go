package template

type Template interface {
	Render(data interface{}) (string, error)
}
