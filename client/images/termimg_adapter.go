package images

import (
	"fmt"
	"image"
	"strings"

	termimg "github.com/blacktop/go-termimg"
)

type RenderProtocol string

const (
	ProtocolAuto       RenderProtocol = "auto"
	ProtocolKitty      RenderProtocol = "kitty"
	ProtocolITerm2     RenderProtocol = "iterm2"
	ProtocolSixel      RenderProtocol = "sixel"
	ProtocolHalfblocks RenderProtocol = "halfblocks"
)

type RenderOptions struct {
	Width    int
	Height   int
	Protocol RenderProtocol
}

type Renderer interface {
	Render(img image.Image, opts RenderOptions) (string, error)
}

type TermimgRenderer struct{}

func NewTermimgRenderer() TermimgRenderer {
	return TermimgRenderer{}
}

func (r TermimgRenderer) Render(img image.Image, opts RenderOptions) (string, error) {
	protocol, err := parseProtocol(opts.Protocol)
	if err != nil {
		return "", err
	}

	renderer := termimg.New(img).Protocol(protocol).Scale(termimg.ScaleFit).UseUnicode(true)
	if opts.Width > 0 {
		renderer = renderer.Width(opts.Width)
	}
	if opts.Height > 0 {
		renderer = renderer.Height(opts.Height)
	}

	return renderer.Render()
}

func parseProtocol(protocol RenderProtocol) (termimg.Protocol, error) {
	switch strings.ToLower(string(protocol)) {
	case "", "auto":
		return termimg.Auto, nil
	case "kitty":
		return termimg.Kitty, nil
	case "iterm2":
		return termimg.ITerm2, nil
	case "sixel":
		return termimg.Sixel, nil
	case "halfblocks":
		return termimg.Halfblocks, nil
	default:
		return termimg.Unsupported, fmt.Errorf("unsupported image preview protocol %q", protocol)
	}
}
