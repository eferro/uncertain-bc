// Package domain holds the pure core of ubc: the language and its evaluation.
package domain

import "strconv"

func Evaluate(program string) (float64, error) {
	return strconv.ParseFloat(program, 64)
}
