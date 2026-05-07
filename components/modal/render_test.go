package modal

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

type staticViewer string

func (v staticViewer) View() string { return string(v) }

func TestGetLinesIgnoresKittyControlSequenceWidth(t *testing.T) {
	_, widest := getLines("\x1b_Ga=T,f=100;abc\x1b\\")
	if widest != 0 {
		t.Fatalf("unexpected widest width for pure kitty control sequence: got %d, want 0", widest)
	}
}

func TestGetLinesCountsVisibleTextAfterKittyControlSequence(t *testing.T) {
	_, widest := getLines("\x1b_Ga=T,f=100;abc\x1b\\HELLO")
	if widest != 5 {
		t.Fatalf("unexpected widest width: got %d, want 5", widest)
	}
}

func TestPlaceModalKeepsRightBorderVisibleWhenWidthCapped(t *testing.T) {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		Padding(0, 1).
		MaxWidth(10)

	bgLine := strings.Repeat(" ", 10)
	background := staticViewer(strings.Join([]string{
		bgLine, bgLine, bgLine, bgLine, bgLine,
	}, "\n"))
	foreground := staticViewer(strings.Repeat("x", 40))

	out := PlaceModal(foreground, background, lipgloss.Left, lipgloss.Top, modalStyle)
	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		t.Fatalf("unexpected output: expected at least 2 lines, got %d", len(lines))
	}

	if !strings.HasSuffix(strings.TrimRight(lines[1], " "), "│") {
		t.Fatalf("expected right border on modal content line, got %q", lines[1])
	}
}
