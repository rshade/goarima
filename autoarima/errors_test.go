package autoarima

import (
	"errors"
	"testing"

	"github.com/rshade/goarima/timeseries"
)

func TestAutoARIMAShortSeriesReturnsError(t *testing.T) {
	result, err := AutoARIMA(timeseries.New([]float64{100, 104, 109, 113, 118, 122}), nil)
	if !errors.Is(err, ErrSeriesTooShort) {
		t.Fatalf("expected ErrSeriesTooShort, got %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestAutoARIMANilSeriesReturnsError(t *testing.T) {
	if _, err := AutoARIMA(nil, nil); !errors.Is(err, ErrSeriesTooShort) {
		t.Fatalf("expected ErrSeriesTooShort, got %v", err)
	}
}

func TestAutoARIMAConstantSeriesReturnsError(t *testing.T) {
	values := make([]float64, 24)
	for i := range values {
		values[i] = 250
	}
	result, err := AutoARIMA(timeseries.New(values), nil)
	if !errors.Is(err, ErrNoModel) {
		t.Fatalf("expected ErrNoModel, got %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestPredictWithIntervalWithoutModelReturnsError(t *testing.T) {
	var nilResult *Result
	if _, _, _, err := nilResult.PredictWithInterval(3, 0.95); !errors.Is(err, ErrNilResult) {
		t.Fatalf("nil result: expected ErrNilResult, got %v", err)
	}
	if _, _, _, err := (&Result{}).PredictWithInterval(3, 0.95); !errors.Is(err, ErrNilResult) {
		t.Fatalf("empty result: expected ErrNilResult, got %v", err)
	}
	if _, err := (&Result{}).Predict(3); !errors.Is(err, ErrNilResult) {
		t.Fatalf("Predict: expected ErrNilResult, got %v", err)
	}
}
