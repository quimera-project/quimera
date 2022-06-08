package structure

type Structure interface {
	Render() (string, error)
	Markdown() (string, error)
	Html() (any, error)
}
