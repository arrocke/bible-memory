package fsrs

import (
	"math"
)

var params = []float64{0.4, 0.6, 2.4, 5.8, 4.93, 0.94, 0.86, 0.01, 1.49, 0.14, 0.94, 2.18, 0.05, 0.34, 1.26, 0.29, 2.61}
var decay = -0.5
var factor = 19.0 / 81.0

func initialStability(grade int) float64 {
	clampedGrade := max(1, min(4, grade))
	return params[clampedGrade-1]
}

func initialDifficulty(grade int) float64 {
	clampedGrade := max(1, min(4, grade))
	return params[4] - float64(clampedGrade-3)*params[5]
}

func newDifficulty(grade int, difficulty float64) float64 {
	return params[7]*initialDifficulty(3) + (1.0-params[7])*(difficulty-params[6]*float64(grade-3))
}

func retrievabilityAfterDays(days int, stability float64) float64 {
	return math.Pow(1+factor*float64(days)/stability, decay)
}

func intervalForRetrievablity(retrievability float64, stability float64) float64 {
	return stability / factor * (math.Pow(retrievability, 1/decay) - 1.0)
}

func newStability(difficulty float64, stability float64, retrievability float64, grade int) float64 {
	if grade > 1 {
		x := 1.0
		if grade == 2 {
			x = params[15]
		}
		if grade == 4 {
			x = params[16]
		}
		return stability * (math.Exp(params[8])*(11-difficulty)*math.Pow(stability, -params[9])*(math.Exp(params[10]*(1-retrievability))-1)*x + 1)
	} else {
		return params[11] * math.Pow(difficulty, -params[12]) * (math.Pow(stability+1.0, params[13]) - 1) * math.Exp(params[14]*(1-retrievability))
	}
}

type MemoryState struct {
	Stability  float64
	Difficulty float64
}

func InitialMemoryState(grade int) MemoryState {
	return MemoryState{
		Stability:  initialStability(grade),
		Difficulty: initialDifficulty(grade),
	}
}

func (state *MemoryState) Review(grade int, days int) {
	d := newDifficulty(grade, state.Difficulty)
	r := retrievabilityAfterDays(days, state.Stability)
	s := newStability(state.Difficulty, state.Stability, r, grade)
	state.Difficulty = d
	state.Stability = s
}

func (state *MemoryState) Interval(expectedRetrievability float64) float64 {
	return intervalForRetrievablity(expectedRetrievability, state.Stability)
}
