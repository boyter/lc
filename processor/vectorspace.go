// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"math"
)

type Concordance map[string]float64

func (con Concordance) magnitude() float64 {
	total := 0.0

	for _, v := range con {
		total = total + math.Pow(v, 2)
	}

	return math.Sqrt(total)
}

func BuildConcordance(words []string) Concordance {
	con := map[string]float64{}

	for _, key := range words {
		con[key] = con[key] + 1
	}

	return con
}

func Relation(con1 Concordance, con2 Concordance) float64 {
	topValue := 0.0

	for name, count := range con1 {
		_, ok := con2[name]

		if ok {
			topValue = topValue + (count * con2[name])
		}
	}

	mag := con1.magnitude() * con2.magnitude()

	if mag != 0 {
		return topValue / mag
	}

	return 0
}
