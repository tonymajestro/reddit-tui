package colors

import "github.com/charmbracelet/lipgloss"

// https://github.com/catppuccin/catppuccin

type Color int

const (
	Red = iota
	Maroon
	Pink
	Orange
	Yellow
	Green
	Blue
	Purple
	Indigo
	Lavender
	Text
	Subtext
	Sand
)

type Palette struct {
	Red      string
	Maroon   string
	Pink     string
	Orange   string
	Yellow   string
	Green    string
	Blue     string
	Purple   string
	Indigo   string
	Lavender string
	Text     string
	Subtext  string
	Sand     string
}

func (p Palette) ToHex(color Color) string {
	switch color {
	case Red:
		return p.Red
	case Maroon:
		return p.Maroon
	case Pink:
		return p.Pink
	case Orange:
		return p.Orange
	case Yellow:
		return p.Yellow
	case Green:
		return p.Green
	case Blue:
		return p.Blue
	case Purple:
		return p.Purple
	case Indigo:
		return p.Indigo
	case Lavender:
		return p.Lavender
	case Text:
		return p.Text
	case Subtext:
		return p.Subtext
	case Sand:
		return p.Sand
	default:
		return p.Text
	}
}

// catppuccin-macchiato
var Dark = Palette{
	Red:      "#ed8796",
	Maroon:   "#ee99a0",
	Pink:     "#f5bde6",
	Orange:   "#f5a97f",
	Yellow:   "#eed49f",
	Green:    "#a6da95",
	Blue:     "#8aadf4",
	Purple:   "#c6a0f6",
	Indigo:   "#5f5fd7 ",
	Lavender: "#b7bdf8",
	Text:     "#cad3f5",
	Subtext:  "#b8c0e0",
	Sand:     "#dddddd",
}

// catppuccin-latte
var Light = Palette{
	Red:      "#d20f39",
	Maroon:   "#e64553",
	Pink:     "#ea76cb",
	Orange:   "#fe640b",
	Yellow:   "#df8e1d",
	Green:    "#40a02b",
	Blue:     "#1e66f5",
	Purple:   "#8839ef",
	Lavender: "#7287fd",
	Text:     "#4c4f69",
	Subtext:  "#5c5f77",
	Sand:     "#dddddd",
}

func AdaptiveColors(light, dark Color) lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: Light.ToHex(light),
		Dark:  Dark.ToHex(dark),
	}
}

func AdaptiveColor(color Color) lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: Light.ToHex(color),
		Dark:  Dark.ToHex(color),
	}
}
