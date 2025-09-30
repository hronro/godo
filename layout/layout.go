package layout

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func Render(ctx context.Context, w io.Writer, title string, content templ.Component) error {
	return Page(title, content).Render(ctx, w)
}
