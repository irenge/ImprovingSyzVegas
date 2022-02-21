// Copyright 2015 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package main

import (
	"math/rand"
	"time"

	"github.com/google/syzkaller/pkg/log"
	"github.com/google/syzkaller/pkg/mab"
	"github.com/google/syzkaller/prog"
)

func (proc *Proc) DoGenerate() *mab.ExecResult {
	ts0 := time.Now().UnixNano()
	ct := proc.fuzzer.choiceTable
	p := proc.fuzzer.target.Generate(proc.rnd, prog.RecommendedCalls, ct)
	_, ret := proc.execute(proc.execOpts, p, ProgNormal, StatSmash)
	ret.TimeTotal = float64(time.Now().UnixNano()-ts0) / MABTimeUnit
	return &ret
}

func (proc *Proc) DoMutate() *mab.ExecResult {
	ts0 := time.Now().UnixNano()
	fuzzerSnapshot := proc.fuzzer.snapshot()
	ct := proc.fuzzer.choiceTable
	// MAB seed selection is integrated with chooseProgram
	pidx, _p := fuzzerSnapshot.chooseProgram(proc.rnd)
	p := _p.Clone()
	p.ResetReward()
	p.Mutate(proc.rnd, prog.RecommendedCalls, ct, fuzzerSnapshot.corpus)
	_, ret := proc.execute(proc.execOpts, p, ProgNormal, StatFuzz)
	ret.Pidx = pidx
	ret.TimeTotal = float64(time.Now().UnixNano()-ts0) / MABTimeUnit
	return &ret
}

func (proc *Proc) DoTriage() *mab.TriageResult {
	ts0 := time.Now().UnixNano()
	item := proc.fuzzer.workQueue.dequeue(DequeueOptionTriageOnly)
	switch item := item.(type) {
	case *WorkTriage:
		{
			ret := proc.ProcessItem(item)
			ret.TimeTotal = float64(time.Now().UnixNano()-ts0) / MABTimeUnit
			return ret
		}
	default:
		{
			return nil
		}
	}
}

func (proc *Proc) DoRemoveCall() *mab.ExecResult {
	ts0 := time.Now().UnixNano()
	fuzzerSnapshot := proc.fuzzer.snapshot()
	ct := proc.fuzzer.choiceTable
	// MAB seed selection is integrated with chooseProgram
	pidx, _p := fuzzerSnapshot.chooseProgram(proc.rnd)
	p := _p.Clone()
	p.ResetReward()
	p.MutateRemove(proc.rnd, prog.RecommendedCalls, ct, fuzzerSnapshot.corpus)
	_, ret := proc.execute(proc.execOpts, p, ProgNormal, StatFuzz)
	ret.Pidx = pidx
	ret.TimeTotal = float64(time.Now().UnixNano()-ts0) / MABTimeUnit
	return &ret
}

func (proc *Proc) DoMutateArg() *mab.ExecResult {
	ts0 := time.Now().UnixNano()
	fuzzerSnapshot := proc.fuzzer.snapshot()
	ct := proc.fuzzer.choiceTable
	// MAB seed selection is integrated with chooseProgram
	pidx, _p := fuzzerSnapshot.chooseProgram(proc.rnd)
	p := _p.Clone()
	p.ResetReward()
	p.MutateArg(proc.rnd, prog.RecommendedCalls, ct, fuzzerSnapshot.corpus)
	_, ret := proc.execute(proc.execOpts, p, ProgNormal, StatFuzz)
	ret.Pidx = pidx
	ret.TimeTotal = float64(time.Now().UnixNano()-ts0) / MABTimeUnit
	return &ret
}

func (proc *Proc) DoInsertCall() *mab.ExecResult {
	ts0 := time.Now().UnixNano()
	fuzzerSnapshot := proc.fuzzer.snapshot()
	ct := proc.fuzzer.choiceTable
	// MAB seed selection is integrated with chooseProgram
	pidx, _p := fuzzerSnapshot.chooseProgram(proc.rnd)
	p := _p.Clone()
	p.ResetReward()
	p.InsertCall(proc.rnd, prog.RecommendedCalls, ct, fuzzerSnapshot.corpus)
	_, ret := proc.execute(proc.execOpts, p, ProgNormal, StatFuzz)
	ret.Pidx = pidx
	ret.TimeTotal = float64(time.Now().UnixNano()-ts0) / MABTimeUnit
	return &ret
}

func (proc *Proc) DoSplice() *mab.ExecResult {
	ts0 := time.Now().UnixNano()
	fuzzerSnapshot := proc.fuzzer.snapshot()
	ct := proc.fuzzer.choiceTable
	// MAB seed selection is integrated with chooseProgram
	pidx, _p := fuzzerSnapshot.chooseProgram(proc.rnd)
	p := _p.Clone()
	p.ResetReward()
	p.Splice(proc.rnd, prog.RecommendedCalls, ct, fuzzerSnapshot.corpus)
	_, ret := proc.execute(proc.execOpts, p, ProgNormal, StatFuzz)
	ret.Pidx = pidx
	ret.TimeTotal = float64(time.Now().UnixNano()-ts0) / MABTimeUnit
	return &ret
}
func (proc *Proc) DoSquashAny() *mab.ExecResult {
	ts0 := time.Now().UnixNano()
	fuzzerSnapshot := proc.fuzzer.snapshot()
	ct := proc.fuzzer.choiceTable
	// MAB seed selection is integrated with chooseProgram
	pidx, _p := fuzzerSnapshot.chooseProgram(proc.rnd)
	p := _p.Clone()
	p.ResetReward()
	p.SquashAny(proc.rnd, prog.RecommendedCalls, ct, fuzzerSnapshot.corpus)
	_, ret := proc.execute(proc.execOpts, p, ProgNormal, StatFuzz)
	ret.Pidx = pidx
	ret.TimeTotal = float64(time.Now().UnixNano()-ts0) / MABTimeUnit
	return &ret
}


func (proc *Proc) clearQueue() {
	// Clear the work queue for all non-Triage items
	count := 0
	for {
		item := proc.fuzzer.workQueue.dequeue(DequeueOptionNoTriage)
		if item != nil {
			count++
			log.Logf(0, "Clearing candidate %v from work queue.\n", count)
			proc.ProcessItem(item)
		} else {
			return
		}
	}
}

func (proc *Proc) MABLoop() {
	// Compute weight and proba
	weight := proc.fuzzer.MABStatus.GetTSWeight(true)
	fuzzerSnapshot := proc.fuzzer.snapshot()
	triageCount := 1
	mutateCount := 1
	removeCount := 1
	mutateArgCount := 1
	insertCallCount := 1
	spliceCount := 1
	squashAnyCount := 1
	if len(fuzzerSnapshot.corpus) == 0 { // Check whether mutation is an option
		mutateCount = 0
		removeCount = 0
		mutateArgCount = 0
		insertCallCount = 0
		spliceCount = 0
		squashAnyCount = 0
	}
	proc.fuzzer.workQueue.mu.Lock() // Check whether triage is an option
	triageQueueLen := len(proc.fuzzer.workQueue.triage)
	triageQueueLenCandidate := len(proc.fuzzer.workQueue.triageCandidate)
	proc.fuzzer.workQueue.mu.Unlock()
	if triageQueueLen+triageQueueLenCandidate == 0 {
		triageCount = 0
	}
	W := weight[0] + float64(mutateCount)*weight[1] + float64(triageCount)*weight[2] + float64(removeCount) *weight[3] + float64(mutateArgCount) * weight[4] + float64(insertCallCount) * weight[5] + float64(spliceCount) * weight[6] + float64(squashAnyCount) * weight[7]

	if W == 0.0 {
		log.Fatalf("Error total weight W = 0")
	}
	prGenerate := weight[0] / W
	prMutate := weight[1] / W
	prTriage := weight[2] / W
	prRemove := weight[3] / W
	prMutateArg := weight[4] / W
        prInsert := weight[5] / W
	prSplice := weight[6] / W
	prSquash := weight[7] / W 

	prMutateActual := float64(mutateCount) * prMutate
	prTriageActual := float64(triageCount) * prTriage
	prRemoveActual := float64(removeCount) * prRemove
	prMutateArgActual := float64(mutateAgrCount) * prMutateArg
	prInsertActual := float64(insertCallCount) * prInsert
	prSpliceActual := float64(spliceCount) * prSplice
	prSquashActual := float64(squashAnyCount) * prSquash

	// Use real weight as pr. Consider cases where triage/mutation might be unavailable
	prTasks := []float64{prGenerate, prMutateActual, prTriageActual, prRemoveActual, prMutateArgActual, prInsertActual, prSliceActual, prSquashActual }
	log.Logf(0, "MAB Probability: [%v, %v, %v, %v, %v, %v, %v, %v]\n", prTasks[0], prTasks[1], prTasks[2], prTasks[3], prTasks[4], prTasks[5], prTasks[6], prTasks[7])
	// Choose
	randNum := rand.Float64() * (prGenerate + prMutateActual + prTriageActual + prRemoveActual + prMutateArgActual + prInsertActual + prSliceActual + prSquashActual)

	choice := -1
	if randNum <= prGenerate {
		choice = 0
	} else if randNum <= prGenerate+prMutateActual {
		choice = 1
	} else if randNum <= prGenerate + prMutateActual + prTriageActual {
		choice = 2
	} else if randNum <= prGenerate + prMutateActual + prTriageActual + prRemoveActual {
		choice = 3
	} else if randNum <= prGenerate + prMutateActual + prTriageActual + prRemoveActual + prMutateArgActual {
		choice = 4
	} else if randNum <= prGenerate + prMutateActual + prTriageActual + prRemoveActual + prMutateArgActual + prInsertActual {
		choice = 5
	} else if randNum <= prGenerate + prMutateActual + prTriageActual + prRemoveActual + prMutateArgActual + prInsertActual + prSliceActual {
		choice = 6
	} else {
		choice = 7
	}

	// Handle choices
	var r interface{}
	if choice == 0 {
		r = *proc.DoGenerate()
	} else if choice == 1 {
		r = *proc.DoMutate()
	} else if choice == 2 {
		r = *proc.DoTriage()
	} else if choice == 3 {
		r = *proc.DoRemoveCall()
	} else if choice == 4 {
		r = *proc.DoMutateArg()
	} else if choice == 5 {
		r = *proc.DoInsertCall()
	} else if choice == 6 {
		r = *proc.DoSplice()
	} else if choice == 7 {
		r = *proc.DoSquashAny()
	}

	if r != nil {
		log.Logf(0, "MAB Choice: %v, Result: %+v\n", choice, r)
	}
	// Update Weight
	proc.fuzzer.MABStatus.UpdateWeight(choice, r, prTasks)
}
