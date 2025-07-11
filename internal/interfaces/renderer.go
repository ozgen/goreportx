package interfaces

type RendererInterface interface {
	Render(filename string) ([]byte, error)
	WithTimestamp(enable bool) RendererInterface
	SetTimestampFormat(format string) RendererInterface
}
