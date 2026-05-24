package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"

	bgBlue  = "\033[44m"
	bgCyan  = "\033[46m"
	bgWhite = "\033[47m"
	bgBlack = "\033[40m"

	colorBrightRed    = "\033[91m"
	colorBrightGreen  = "\033[92m"
	colorBrightYellow = "\033[93m"
	colorBrightBlue   = "\033[94m"
	colorBrightCyan   = "\033[96m"
)

func bold(s string) string   { return colorBold + s + colorReset }
func dim(s string) string    { return colorDim + s + colorReset }
func red(s string) string    { return colorRed + s + colorReset }
func green(s string) string  { return colorGreen + s + colorReset }
func yellow(s string) string { return colorYellow + s + colorReset }
func blue(s string) string   { return colorBlue + s + colorReset }
func cyan(s string) string   { return colorCyan + s + colorReset }

func headerBg(s string) string { return bgCyan + colorBold + colorWhite + " " + s + " " + colorReset }

func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		if r < 128 {
			w++
		} else {
			w += 2
		}
	}
	return w
}

func padRight(s string, width int) string {
	sw := displayWidth(s)
	if sw >= width {
		return s
	}
	return s + strings.Repeat(" ", width-sw)
}

func padLeft(s string, width int) string {
	sw := displayWidth(s)
	if sw >= width {
		return s
	}
	return strings.Repeat(" ", width-sw) + s
}

func centerPad(s string, width int) string {
	sw := displayWidth(s)
	if sw >= width {
		return s
	}
	left := (width - sw) / 2
	right := width - sw - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func truncate(s string, maxWidth int) string {
	runes := []rune(s)
	w := 0
	for i, r := range runes {
		rw := 1
		if r >= 128 {
			rw = 2
		}
		if w+rw > maxWidth {
			return string(runes[:i]) + "…"
		}
		w += rw
	}
	return s
}

func statusColor(status string) string {
	s := strings.TrimSpace(status)
	switch s {
	case "成功", "已启用", "正常", "空闲":
		return green(s)
	case "失败", "已失效", "已暂停":
		return red(s)
	case "备份中", "运行中":
		return yellow(s)
	default:
		return s
	}
}

func PrintTable(headers []string, rows [][]string) {
	if len(headers) == 0 || len(rows) == 0 {
		return
	}

	terminalWidth := 120

	colCount := len(headers)
	colWidths := make([]int, colCount)
	for i, h := range headers {
		colWidths[i] = displayWidth(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < colCount {
				w := displayWidth(cell)
				if w > colWidths[i] {
					colWidths[i] = w
				}
			}
		}
	}

	totalWidth := 0
	for _, w := range colWidths {
		totalWidth += w
	}
	totalWidth += (colCount-1)*3 + 2

	if totalWidth > terminalWidth {
		excess := totalWidth - terminalWidth
		for excess > 0 && colCount > 1 {
			longestIdx := 0
			longestW := 0
			for i, w := range colWidths {
				if w > longestW {
					longestW = w
					longestIdx = i
				}
			}
			reduce := longestW / 4
			if reduce < 2 {
				reduce = 2
			}
			if reduce > excess {
				reduce = excess
			}
			colWidths[longestIdx] -= reduce
			excess -= reduce
		}
	}

	totalWidth = 0
	for _, w := range colWidths {
		totalWidth += w
	}
	totalWidth += (colCount-1)*3 + 2

	topBorder := dim("┌") + dim(strings.Repeat("─", totalWidth-2)) + dim("┐")
	sepBorder := dim("├") + dim(strings.Repeat("─", totalWidth-2)) + dim("┤")
	botBorder := dim("└") + dim(strings.Repeat("─", totalWidth-2)) + dim("┘")

	fmt.Println(topBorder)

	headerParts := make([]string, colCount)
	for i, h := range headers {
		headerParts[i] = headerBg(centerPad(h, colWidths[i]))
	}
	fmt.Println(dim("│") + " " + strings.Join(headerParts, " "+dim("│")+" ") + " " + dim("│"))

	fmt.Println(sepBorder)

	for _, row := range rows {
		cells := make([]string, colCount)
		for i, cell := range row {
			if i < colCount {
				text := truncate(cell, colWidths[i])
				padded := padRight(text, colWidths[i])
				cells[i] = padded
			}
		}
		fmt.Println(dim("│") + " " + strings.Join(cells, " "+dim("│")+" ") + " " + dim("│"))
	}

	fmt.Println(botBorder)
}

func PrintDetailTable(fields map[string]string) {
	if len(fields) == 0 {
		return
	}

	maxKeyW := 0
	for k := range fields {
		w := displayWidth(k)
		if w > maxKeyW {
			maxKeyW = w
		}
	}

	terminalWidth := 100
	valueWidth := terminalWidth - maxKeyW - 9
	if valueWidth < 20 {
		valueWidth = 20
	}

	totalWidth := maxKeyW + valueWidth + 7

	topBorder := dim("┌") + dim(strings.Repeat("─", totalWidth-2)) + dim("┐")
	botBorder := dim("└") + dim(strings.Repeat("─", totalWidth-2)) + dim("┘")

	fmt.Println(topBorder)

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}

	for _, k := range keys {
		v := fields[k]
		keyPadded := padRight(bold(k), maxKeyW)

		vRunes := []rune(v)
		if utf8.RuneCountInString(v) > 0 && displayWidth(v) > valueWidth {
			var lines []string
			line := ""
			lineW := 0
			for _, r := range vRunes {
				rw := 1
				if r >= 128 {
					rw = 2
				}
				if lineW+rw > valueWidth {
					lines = append(lines, line)
					line = ""
					lineW = 0
				}
				line += string(r)
				lineW += rw
			}
			if line != "" {
				lines = append(lines, line)
			}
			for i, l := range lines {
				if i == 0 {
					valPadded := padRight(l, valueWidth)
					fmt.Println(dim("│") + " " + keyPadded + " " + dim("│") + " " + valPadded + " " + dim("│"))
				} else {
					emptyKey := padRight("", maxKeyW)
					valPadded := padRight(l, valueWidth)
					fmt.Println(dim("│") + " " + emptyKey + " " + dim("│") + " " + valPadded + " " + dim("│"))
				}
			}
		} else {
			valPadded := padRight(v, valueWidth)
			fmt.Println(dim("│") + " " + keyPadded + " " + dim("│") + " " + valPadded + " " + dim("│"))
		}
	}

	fmt.Println(botBorder)
}

func PrintTitle(title string) {
	width := displayWidth(title) + 4
	if width < 40 {
		width = 40
	}
	fmt.Println()
	fmt.Println(dim("╔" + strings.Repeat("═", width-2) + "╗"))
	fmt.Println(dim("║") + " " + bold(title) + strings.Repeat(" ", width-displayWidth(title)-4) + " " + dim("║"))
	fmt.Println(dim("╚" + strings.Repeat("═", width-2) + "╝"))
}

func PrintSubTitle(title string) {
	fmt.Println()
	fmt.Println(dim("──") + " " + bold(title) + " " + dim(strings.Repeat("─", 60-displayWidth(title))))
}

func PrintInfo(msg string) {
	fmt.Println(cyan("ℹ") + " " + msg)
}

func PrintSuccess(msg string) {
	fmt.Println(green("✅") + " " + green(msg))
}

func PrintError(msg string) {
	fmt.Println(red("❌") + " " + red(msg))
}

func PrintWaiting(msg string) {
	fmt.Println(yellow("⏳") + " " + yellow(msg))
}

func PrintDivider() {
	fmt.Println(dim(strings.Repeat("─", 60)))
}

func PrintBox(content string) {
	lines := strings.Split(content, "\n")
	maxW := 0
	for _, l := range lines {
		w := displayWidth(l)
		if w > maxW {
			maxW = w
		}
	}
	if maxW < 20 {
		maxW = 20
	}
	if maxW > 100 {
		maxW = 100
	}

	fmt.Println(dim("┌" + strings.Repeat("─", maxW+2) + "┐"))
	for _, l := range lines {
		padded := padRight(l, maxW)
		fmt.Println(dim("│") + " " + padded + " " + dim("│"))
	}
	fmt.Println(dim("└" + strings.Repeat("─", maxW+2) + "┘"))
}
