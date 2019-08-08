package brucecore

import (
	"fmt"
)

// FormatTime - 100 => 1m40s
func FormatTime(time int) string {
	if time > 24*60*60 {
		d := time / (24 * 60 * 60)
		time -= d * (24 * 60 * 60)
		h := time / (60 * 60)
		time -= h * (60 * 60)
		m := time / 60
		time -= m * 60
		s := time

		return fmt.Sprintf("%vd%vh%vm%vs", d, h, m, s)
	} else if time > 60*60 {
		h := time / (60 * 60)
		time -= h * (60 * 60)
		m := time / 60
		time -= m * 60
		s := time

		return fmt.Sprintf("%vh%vm%vs", h, m, s)
	} else if time > 60 {
		m := time / 60
		time -= m * 60
		s := time

		return fmt.Sprintf("%vm%vs", m, s)
	}

	return fmt.Sprintf("%vs", time)
}

// FormatByteSize - 1025 => 1k1b
func FormatByteSize(bytesize int) string {
	if bytesize > 1024*1024*1024 {
		g := bytesize / (1024 * 1024 * 1024)
		bytesize -= g * (bytesize)
		m := bytesize / (1024 * 1024)
		bytesize -= m * (1024 * 1024)
		k := bytesize / 1024
		bytesize -= k * 1024
		b := bytesize

		return fmt.Sprintf("%vg%vm%vk%vb", g, m, k, b)
	} else if bytesize > 1024*1024 {
		m := bytesize / (1024 * 1024)
		bytesize -= m * (1024 * 1024)
		k := bytesize / 1024
		bytesize -= k * 1024
		b := bytesize

		return fmt.Sprintf("%vm%vk%vb", m, k, b)
	} else if bytesize > 1024 {
		k := bytesize / 1024
		bytesize -= k * 1024
		b := bytesize

		return fmt.Sprintf("%vk%vb", k, b)
	}

	return fmt.Sprintf("%vb", bytesize)
}
