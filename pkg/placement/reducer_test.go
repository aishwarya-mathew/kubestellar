/*
Copyright 2023 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package placement

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	edgeapi "github.com/kcp-dev/edge-mc/pkg/apis/edge/v1alpha1"
)

func exerciseSinglePlacementSliceSetReducer(rng *rand.Rand, initialWhere ResolvedWhere, iterations int, changesPerIteration int, extraPerIteration func(), reducer SinglePlacementSliceSetReducer, uider UIDer, consumer SinglePlacementSet) func(*testing.T) {
	return func(t *testing.T) {
		input := initialWhere
		for iteration := 1; iteration <= iterations; iteration++ {
			prevInput := input
			for change := 1; change <= changesPerIteration; change++ {
				input = reviseSinglePlacementSliceSlice(rng, input)
			}
			reducer.Set(input)
			extraPerIteration()
			checker := NewSinglePlacementSet()
			reference := NewSimplePlacementSliceSetReducer(uider, checker)
			reference.Set(input)
			if consumer.Equals(checker) {
				continue
			}
			t.Errorf("Unexpected result: excess=%v, missing=%v, iteration=%d, prevInput=%v, input=%v", consumer.Sub(checker), checker.Sub(consumer), iteration, prevInput, input)
		}
	}
}

func reviseSinglePlacementSliceSlice(rng *rand.Rand, slices ResolvedWhere) ResolvedWhere {
	ans := make(ResolvedWhere, 0, len(slices))
	copy(ans, slices)
	if len(ans) == 0 || rng.Intn(20) == 1 {
		// Add a new slice
		sliceLen := (rng.Intn(12) + 2) / 3
		slice := edgeapi.SinglePlacementSlice{Destinations: []edgeapi.SinglePlacement{}}
		for dest := 1; dest <= sliceLen; dest++ {
			slice.Destinations = append(slice.Destinations, genSinglePlacement(rng))
		}
		sliceIdx := rng.Intn(len(ans) + 1)
		ans = append(ans[:sliceIdx], append([]*edgeapi.SinglePlacementSlice{&slice}, ans[sliceIdx:]...)...)
	} else if rng.Intn(20) != 1 {
		// modify an existing slice
		sliceIdx := rng.Intn(len(ans))
		slice := *ans[sliceIdx]
		newDestinations := make([]edgeapi.SinglePlacement, 0, len(slice.Destinations))
		copy(newDestinations, slice.Destinations)
		if len(slice.Destinations) == 0 || rng.Intn(20) == 1 {
			// Add a new entry
			destIdx := rng.Intn(len(newDestinations) + 1)
			dest := genSinglePlacement(rng)
			newDestinations = append(newDestinations[:destIdx], append([]edgeapi.SinglePlacement{dest}, newDestinations[destIdx:]...)...)
		} else if rng.Intn(20) != 1 {
			// modify an existing SinglePlacement
			destIdx := rng.Intn(len(newDestinations))
			newDestinations[destIdx] = reviseSinglePlacement(rng, newDestinations[destIdx])
		} else {
			// delete an existing entry
			destIdx := rng.Intn(len(newDestinations))
			newDestinations = append(newDestinations[:destIdx], newDestinations[destIdx+1:]...)
		}
		slice.Destinations = newDestinations
	} else {
		// Delete an existing slice
		sliceIdx := rng.Intn(len(ans))
		ans = append(ans[:sliceIdx], ans[sliceIdx+1:]...)
	}
	return ans
}

func reviseSinglePlacement(rng *rand.Rand, sp edgeapi.SinglePlacement) edgeapi.SinglePlacement {
	switch rng.Intn(3) {
	case 0:
		sp.Location.Workspace = fmt.Sprintf("ws%d", rng.Intn(1000))
	case 1:
		sp.Location.Name = fmt.Sprintf("loc%d", rng.Intn(1000))
	default:
		sp.SyncTargetName = fmt.Sprintf("st%d", rng.Intn(1000))
	}
	return sp
}

func genSinglePlacement(rng *rand.Rand) edgeapi.SinglePlacement {
	return edgeapi.SinglePlacement{
		Location: edgeapi.ExternalName{
			Workspace: fmt.Sprintf("ws%d", rng.Intn(1000)),
			Name:      fmt.Sprintf("loc%d", rng.Intn(1000)),
		},
		SyncTargetName: fmt.Sprintf("st%d", rng.Intn(1000)),
	}
}

func TestSimplePlacementSliceSetReducer(t *testing.T) {
	rs := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(rs)
	var wg sync.WaitGroup
	testUIDer := NewTestUIDer(rng, &wg)
	testConsumer := NewSinglePlacementSet()
	testReducer := NewSimplePlacementSliceSetReducer(testUIDer, testConsumer)
	locRef1 := edgeapi.ExternalName{Workspace: "ws-a", Name: "loc-a"}
	asp1 := edgeapi.SinglePlacement{Location: locRef1, SyncTargetName: "st-a"}
	sp1 := SinglePlacement{asp1, "u-a"}
	testUIDer.Set(sp1.SyncTargetRef(), sp1.SyncTargetUID)
	rw1 := ResolvedWhere{&edgeapi.SinglePlacementSlice{
		Destinations: []edgeapi.SinglePlacement{asp1},
	}}
	testReducer.Set(rw1)
	if actual, expected := len(testConsumer), 1; actual != expected {
		t.Errorf("Wrong size after first Set: actual=%d, expected=%d", actual, expected)
	}
	if actual, expected := testConsumer[sp1.SyncTargetRef()], sp1.Details(); actual != expected {
		t.Errorf("Wrong details after first Set: actual=%#v, expected=%#v", actual, expected)
	}
	sp1a := SinglePlacement{asp1, "u-aa"}
	testUIDer.Set(sp1.SyncTargetRef(), sp1a.SyncTargetUID)
	wg.Wait()
	if actual, expected := len(testConsumer), 1; actual != expected {
		t.Errorf("Wrong size after first tweak: actual=%d, expected=%d", actual, expected)
	}
	if actual, expected := testConsumer[sp1.SyncTargetRef()], sp1a.Details(); actual != expected {
		t.Errorf("Wrong details after first tweak: actual=%#v, expected=%#v", actual, expected)
	}
	exerciseSinglePlacementSliceSetReducer(rng, rw1, 20, 10, testUIDer.TweakNAndWait(rng, 4), testReducer, testUIDer, testConsumer)(t)
}
