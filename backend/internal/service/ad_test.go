package service

import (
	"backend/internal/model"
	"testing"
)

const (
	wantA     = 1
	wantB     = -1
	wantEqual = 0
)

func TestCompareAdCandidates(t *testing.T) {
	type testCase struct {
		name string
		a    model.AdCandidate
		b    model.AdCandidate
		want int
	}
	tests := []testCase{
		{
			"equal models has equal scores",
			model.AdCandidate{},
			model.AdCandidate{},
			wantEqual,
		},
		{
			"equal models has equal scores",
			model.AdCandidate{CostPerClick: 1.4, CostPerImpression: 1.7, MlScore: 1000},
			model.AdCandidate{CostPerClick: 1.4, CostPerImpression: 1.7, MlScore: 1000},
			wantEqual,
		},
		{
			"not viewed is greater that viewed",
			model.AdCandidate{CostPerImpression: 42, Viewed: true},
			model.AdCandidate{CostPerImpression: 42},
			wantB,
		},
		{
			"viewed is less that not viewed",
			model.AdCandidate{CostPerImpression: 42, MlScore: 100},
			model.AdCandidate{CostPerImpression: 42, Viewed: true, MlScore: 100},
			wantA,
		},
		{
			"not clicked is greater that clicked",
			model.AdCandidate{CostPerImpression: 1.5, Viewed: true, CostPerClick: 0.5, Clicked: true},
			model.AdCandidate{CostPerImpression: 1.5, Viewed: true, CostPerClick: 0.5},
			wantB,
		},
		{
			"clicked is greater that not clicked",
			model.AdCandidate{CostPerImpression: 1.5, Viewed: true, CostPerClick: 0.5},
			model.AdCandidate{CostPerImpression: 1.5, Viewed: true, CostPerClick: 0.5, Clicked: true},
			wantA,
		},
		{
			"cost increasing",
			model.AdCandidate{CostPerImpression: 5},
			model.AdCandidate{CostPerImpression: 10},
			wantB,
		},
		{
			"click cost decreasing",
			model.AdCandidate{CostPerClick: 50},
			model.AdCandidate{CostPerClick: 10},
			wantA,
		},
	}
	for _, tt := range tests {
		if got := compareAdCandidates(tt.a, tt.b); got != tt.want {
			t.Errorf("%s: compareAdCandidates(%v, %v) = %v, want %v", tt.name, tt.a, tt.b, got, tt.want)
		}
	}
}
