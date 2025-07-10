package interfaces

type RendererInterface interface {
	Render(filename string) ([]byte, error)
}
