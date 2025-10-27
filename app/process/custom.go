package process

import (
	"image"
	"math"
	"strings"
	"bufio"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"

	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
	"github.com/Zebbeni/ansizalizer/controls/settings/size"
)

func (m Renderer) processCustom(input image.Image) string {
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

	_, _, useFgBg, chars := m.Settings.Characters.Selected()
	if len(chars) == 0 {
		return "Enter at least one custom character"
	}

	content := ""
	rows := make([]string, height)
	row := make([]string, width)

	for y := 0; y < height*2; y += 2 {
		for x := 0; x < width*2; x += 2 {
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

			if useFgBg == characters.TwoColor {
				fg, bg, brightness := m.fgBgBrightness(r1, r2, r3, r4)

				lipFg := lipgloss.Color(fg.Hex())
				lipBg := lipgloss.Color(bg.Hex())
				style := lipgloss.NewStyle().Foreground(lipFg).Background(lipBg).Bold(true)

				index := min(int(brightness*float64(len(chars))), len(chars)-1)
				char := chars[index]
				charString := string(char)

				row[x/2] = style.Render(charString)
			} else {
				fg := m.avgColTrue(r1, r2, r3, r4)
				brightness := math.Min(1.0, math.Abs(fg.DistanceLuv(black)))
				if isPaletted {
					fg, _ = colorful.MakeColor(palette.Colors().Convert(fg))
				}
				lipFg := lipgloss.Color(fg.Hex())
				style := lipgloss.NewStyle().Foreground(lipFg).Bold(true)
				index := min(int(brightness*float64(len(chars))), len(chars)-1)
				char := chars[index]
				charString := string(char)
				row[x/2] = style.Render(charString)
			}
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
