package arima

import (
	"testing"

	"github.com/rshade/goarima/timeseries"
)

func TestTwiceDifferencedModelKeepsLinearTrend(t *testing.T) {
	values := make([]float64, 120)
	for i := range values {
		values[i] = 100 + 2*float64(i) + float64((i*37)%11)/10 - 0.5
	}
	model := New(0, 2, 1)
	if err := model.Fit(timeseries.New(values)); err != nil {
		t.Fatalf("fit: %v", err)
	}
	if model.Intercept != 0 {
		t.Errorf("intercept = %v, want 0 when d >= 2", model.Intercept)
	}
	forecast, err := model.Predict(120)
	if err != nil {
		t.Fatalf("predict: %v", err)
	}
	last := values[len(values)-1]
	if !(forecast[119] > forecast[59] && forecast[59] > last) {
		t.Errorf("forecast does not keep rising: last %.2f, h=60 %.2f, h=120 %.2f", last, forecast[59], forecast[119])
	}
}
