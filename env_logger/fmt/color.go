package fmt

import (
	"strconv"
)

var (
	Red    = Style{fg: 31, bg: 41}
	Green  = Style{fg: 32, bg: 42}
	Yellow = Style{fg: 33, bg: 43}
	Blue   = Style{fg: 34, bg: 44}
	Purple = Style{fg: 35, bg: 45}
)

const RESET = "\x1B[0m"

type Style struct {
	fg        int
	bg        int
	stylecode int
}

func (s Style) IsPlain() bool {
	plain := true
	if s.fg > 0 || s.bg > 0 || s.stylecode > 0 {
		plain = false
	}
	return plain
}

func (s Style) Paint(m string) string {
	if s.IsPlain() {
		return m
	}

	buf := "\x1B["
	written_anything := false

	if s.stylecode > 0 && s.stylecode < 10 {
		if written_anything {
			buf = buf + ";"
		}
		buf = buf + strconv.Itoa(s.stylecode)
		written_anything = true
	}

	if s.fg > 0 {
		if written_anything {
			buf = buf + ";"
		}
		buf = buf + strconv.Itoa(s.fg)
		written_anything = true
	}

	return buf + "m" + m + RESET
}

func (s Style) Dimmed() Style {
	s.stylecode = 2
	return s
}
