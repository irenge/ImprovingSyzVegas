// Copyright 2015 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package main

import (
	"math"
	"sync"

	"github.com/google/syzkaller/pkg/log"
	"github.com/google/syzkaller/pkg/mab"
)

// Time unit is important since MAB ended up converting coverage
// gained into time, before feeding it to the exponential growth
// algorithm. Converting time unit from ns to s can prevent overflow
// of float64 when computing exponentials
const (
	MABTimeUnit = 1000000000.0
	MABLogLevel = 4
)

type MABStatus struct {
	fuzzer *Fuzzer

	MABMu     sync.RWMutex
	SSEnabled bool
	TSEnabled bool

	SSGamma float64
	SSEta   float64
	TSGamma float64
	TSEta   float64

	Round          int         // How many MAB choices have been made
	Exp31Round     int         // Round # for Exp3.1.
	Exp31Threshold float64     // Threshold based on Round.
	CorpusUpdate   map[int]int // Track seed priority update
	Reward         mab.TotalReward
}

func (status *MABStatus) GetTSWeight(lock bool) []float64 {
	if lock {
		status.MABMu.Lock()
		defer status.MABMu.Unlock()
	}
	x := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	y := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}

	weight := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	eta := status.TSEta
	const (
		MABWeightThresholdMax = 1.0e+300
		MABWeightThresholdMin = 1.0e-300
	)
	x[0] = eta * status.Reward.EstimatedRewardGenerate
	x[1] = eta * status.Reward.EstimatedRewardMutate
	x[2] = eta * status.Reward.EstimatedRewardTriage
	x[3] = eta * status.Reward.EstimatedRewardRemoveCall
	x[4] = eta * status.Reward.EstimatedRewardMutateArg
	x[5] = eta * status.Reward.EstimatedRewardInsertCall 
	x[6] = eta * status.Reward.EstimatedRewardSplice
	x[7] = eta * status.Reward.EstimatedRewardSquashAny  

	for a:= 0; a < 8; a++ {
		y[a] = x[a]
	}

	log.Logf(MABLogLevel, "MABWeight %v\n", x)
	sort.Float64s(y) 
	// Compute median to prevent overflow
	median := (y[3] + y[4]) / 2
	
	for i := 0; i < 8; i++ {
		weight[i] = math.Exp(x[i] - median)
		if weight[i] > MABWeightThresholdMax {
			weight[i] = MABWeightThresholdMax
		}
		if weight[i] < MABWeightThresholdMin {
			weight[i] = MABWeightThresholdMin
		}
	}
	return weight
}
