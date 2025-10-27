package process

import (
	"strings"
	"bufio"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"

	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
	"github.com/Zebbeni/ansizalizer/controls/settings/size"
)

var unicodeShadeChars = []rune{' ', '░', '▒', '▓'}

func (m Renderer) processUnicode(input image.Image) string {
	imgW, imgH := float32(input.Bounds().Dx()), float32(input.Bounds().Dy())

	dimensionType, width, height, charRatio := m.Settings.Size.Info()
	if dimensionType == size.Fit {
		fitHeight := float32(width) * (imgH / imgW) * float32(charRatio)
		fitWidth := (float32(height) * (imgW / imgH)) / float32(charRatio)
		if fitHeight > float32(height) {
			width = int(fitWidth)
		} else {
			height = int(fitHeight)
		}
	}

	resizeFunc := m.Settings.Advanced.SamplingFunction()
	refImg := resize.Resize(uint(width)*2, uint(height)*2, input, resizeFunc)

	isTrueColor, _, palette := m.Settings.Colors.GetSelected()
	isPaletted := !isTrueColor

	doDither, doSerpentine, matrix := m.Settings.Advanced.Dithering()
	if doDither && isPaletted {
		ditherer := dither.NewDitherer(palette.Colors())
		ditherer.Matrix = matrix
		if doSerpentine {
			ditherer.Serpentine = true
		}
		refImg = ditherer.Dither(refImg)
	}

	content := ""
	rows := make([]string, height)
	row := make([]string, width)
	for y := 0; y < height*2; y += 2 {
		for x := 0; x < width*2; x += 2 {
			// r1 r2
			// r3 r4
			r1, r1Alpha := colorful.MakeColor(refImg.At(x, y))
			r2, r2Alpha := colorful.MakeColor(refImg.At(x+1, y))
			r3, r3Alpha := colorful.MakeColor(refImg.At(x, y+1))
			r4, r4Alpha := colorful.MakeColor(refImg.At(x+1, y+1))

			// if the fourth argument (a [for alpha]) returned from MakeColor is false, this is a transparent pixel
			if m.Settings.Alpha.ShouldOutputAlpha() && (!r1Alpha || !r2Alpha || !r3Alpha || !r4Alpha) {
				// use a placeholder to designate the transparent "pixel"
				row[x/2] = AlphaPlaceholder
				continue
			}

			// pick the block, fg and bg color with the lowest total difference
			// convert the colors to ansi, render the block and add it at row[x]
			r, fg, bg := m.getBlock(r1, r2, r3, r4)

			pFg, _ := colorful.MakeColor(fg)
			pBg, _ := colorful.MakeColor(bg)

			lipFg := lipgloss.Color(pFg.Hex())
			lipBg := lipgloss.Color(pBg.Hex())

			style := lipgloss.NewStyle().Foreground(lipFg)
			if _, _, mode, _ := m.Settings.Characters.Selected(); mode == characters.TwoColor {
				style = style.Copy().Background(lipBg)
			}

			row[x/2] = style.Render(string(r))
		}
		rows[y/2] = lipgloss.JoinHorizontal(lipgloss.Top, row...)
	}
	if m.Settings.Alpha.ShouldOutputAlpha() {
		// replace ALPHA placeholder with a blank square (space)
		contentAlpha := strings.ReplaceAll(lipgloss.JoinVertical(lipgloss.Left, rows...), AlphaPlaceholder, MagicTransparentPixel)
		// iterate through the return of JoinVertical, separating by lines, trimming whitespace, and then recombining
		reader := strings.NewReader(contentAlpha)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			content += strings.TrimSpace(scanner.Text()) + "\n"
		}
	} else {
		content += lipgloss.JoinVertical(lipgloss.Left, rows...)
	}
	return content
}

// find the best block character and foreground and background colors to match
// a set of 4 pixels. return
func (m Renderer) getBlock(r1, r2, r3, r4 colorful.Color) (r rune, fg, bg colorful.Color) {
	var blockFuncs map[rune]blockFunc
	switch _, charSet, _, _ := m.Settings.Characters.Selected(); charSet {
	case characters.UnicodeFull:
		blockFuncs = m.fullBlockFuncs
	case characters.UnicodeHalf:
		blockFuncs = m.halfBlockFuncs
	case characters.UnicodeQuart:
		blockFuncs = m.quarterBlockFuncs
	case characters.UnicodeShadeLight:
		blockFuncs = m.shadeLightBlockFuncs
	case characters.UnicodeShadeMed:
		blockFuncs = m.shadeMedBlockFuncs
	case characters.UnicodeShadeHeavy:
		blockFuncs = m.shadeHeavyBlockFuncs
	}

	minDist := 100.0
	for bRune, bFunc := range blockFuncs {
		f, b, dist := bFunc(r1, r2, r3, r4)
		if dist < minDist {
			minDist = dist
			r, fg, bg = bRune, f, b
		}
	}
	return
}

func (m Renderer) avgCol(colors ...colorful.Color) (colorful.Color, float64) {
	rSum, gSum, bSum := 0.0, 0.0, 0.0
	for _, col := range colors {
		rSum += col.R
		gSum += col.G
		bSum += col.B
	}
	count := float64(len(colors))
	avg := colorful.Color{R: rSum / count, G: gSum / count, B: bSum / count}

	if m.Settings.Colors.IsLimited() {
		_, _, palette := m.Settings.Colors.GetSelected()

		paletteAvg := palette.Colors().Convert(avg)
		avg, _ = colorful.MakeColor(paletteAvg)
	}

	// compute sum of squares
	totalDist := 0.0
	for _, col := range colors {
		totalDist += math.Pow(col.DistanceCIEDE2000(avg), 2)
	}
	return avg, totalDist
}
