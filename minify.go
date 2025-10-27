package viewmodel

import (
	"fmt"
	"io"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

func minifyHTML(w io.Writer, r io.Reader) error {
	if err := html.Minify(minify.New(), w, r, nil); err != nil {
		return fmt.Errorf("failed to minify html: %w", err)
	}
	return nil
}
