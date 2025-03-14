package utils

import "unicode"

func MaskOnlyNumbers(in string, let int) string {
	runes := []rune(in)

	digitCount := 0
	for _, r := range runes {
		if unicode.IsDigit(r) {
			digitCount++
		}
	}

	if digitCount <= let {
		return in
	}

	masked := make([]rune, len(runes))
	digitsMasked := 0

	for i, r := range runes {
		if unicode.IsDigit(r) {
			if digitsMasked < digitCount-let {
				masked[i] = '*'
				digitsMasked++
			} else {
				masked[i] = r
			}
		} else {
			masked[i] = r
		}
	}

	return string(masked)
}

func MaskAll(in string, let int) string {
	runes := []rune(in)

	visibleCount := 0
	for _, r := range runes {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			visibleCount++
		}
	}

	if visibleCount <= let {
		return in
	}

	masked := make([]rune, len(runes))
	visibleRemaining := visibleCount - let

	for i, r := range runes {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if visibleRemaining == 0 {
				masked[i] = '*'
			} else {
				masked[i] = r
				visibleRemaining--
			}
		} else {
			masked[i] = r
		}
	}

	return string(masked)
}
