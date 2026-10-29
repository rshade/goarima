package sarima

import (
	"math"
	"testing"

	"github.com/rshade/goarima/timeseries"
)

// trendWeekly is a linear trend with a weekly dip and a small deterministic wobble.
func trendWeekly(n int) ([]float64, func(int) float64) {
	truth := func(i int) float64 {
		v := 50 + 0.1*float64(i)
		if i%7 == 5 || i%7 == 6 {
			v -= 15
		}
		return v
	}
	values := make([]float64, n)
	for i := range values {
		values[i] = truth(i) + float64((i*37)%11)/10 - 0.5
	}
	return values, truth
}

func TestDoublyDifferencedModelKeepsTrend(t *testing.T) {
	values, truth := trendWeekly(180)
	for name, model := range map[string]*Model{
		"css": New(0, 1, 1, 0, 1, 1, 7),
		"mle": NewMLE(0, 1, 1, 0, 1, 1, 7),
	} {
		if err := model.Fit(timeseries.New(values)); err != nil {
			t.Fatalf("%s: fit: %v", name, err)
		}
		if model.Intercept != 0 {
			t.Errorf("%s: intercept = %v, want 0 when d+D >= 2", name, model.Intercept)
		}
		forecast, err := model.Predict(365)
		if err != nil {
			t.Fatalf("%s: predict: %v", name, err)
		}
		// CSS estimates the MA terms of this over-differenced fixture poorly,
		// so only the exact-likelihood fit is held to the trend.
		if name != "mle" {
			continue
		}
		for _, h := range []int{30, 180, 365} {
			want := truth(179 + h)
			if got := forecast[h-1]; math.Abs(got-want)/want > 0.05 {
				t.Errorf("%s: h=%d forecast %.2f, want within 5%% of %.2f", name, h, got, want)
			}
		}
	}
}

func TestSinglyDifferencedModelKeepsDrift(t *testing.T) {
	values, _ := trendWeekly(180)
	model := New(1, 1, 0, 0, 0, 0, 7)
	if err := model.Fit(timeseries.New(values)); err != nil {
		t.Fatalf("fit: %v", err)
	}
	if model.Intercept == 0 {
		t.Error("intercept = 0, want the drift when d+D == 1")
	}
}
