package renderer

import "github.com/unrolled/render"

// Render is the global renderer for all handlers (the render library is threadsafe)
var Render *render.Render

func init() {
	Render = render.New(render.Options{
		IndentJSON: true,
	})
}
