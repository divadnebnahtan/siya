package state_test

import (
	"math"
	"sync"
	"testing"
	"time"

	"github.com/divadnebnahtan/siya/backend/state"
)

func TestStoreUpsert(t *testing.T) {
	var store state.Store
	now := time.Now().UTC()
	motion := false
	position := state.Position{ObservedAt: now, ReceivedAt: now, Motion: &motion}
	location := state.Location{ID: 2, Name: "Tracker", Position: &position}

	_, changed := store.Snapshot()
	invalid := position
	invalid.Latitude = math.NaN()

	err := store.Upsert(location, state.Location{ID: 1, Position: &invalid})
	if err == nil {
		t.Fatal("accepted invalid latitude")
	}

	_, ok := store.Get(2)
	if ok {
		t.Fatal("partially applied invalid batch")
	}

	err = store.Upsert(location, state.Location{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-changed:
	default:
		t.Fatal("change notification was missed")
	}

	motion = true
	position.Latitude = 42
	snapshot, changed := store.Snapshot()
	if len(snapshot) != 2 || snapshot[0].ID != 1 || snapshot[0].Position != nil || snapshot[1].ID != 2 {
		t.Fatalf("unexpected snapshot: %+v", snapshot)
	}

	p := snapshot[1].Position
	if p.Latitude != 0 || *p.Motion {
		t.Fatalf("input mutation reached store: %+v", p)
	}
	*snapshot[1].Position.Motion = true

	position.ObservedAt = now.Add(-time.Second)
	err = store.Upsert(location)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-changed:
		t.Fatal("stale report triggered a change")
	default:
	}

	err = store.Upsert(state.Location{ID: 2, Name: "Renamed"})
	if err != nil {
		t.Fatal(err)
	}
	got, ok := store.Get(2)
	if !ok || got.Name != "Renamed" || got.Position == nil || *got.Position.Motion || got.Position.Latitude != 0 {
		t.Fatalf("unexpected location after metadata update: %+v", got)
	}
}

func TestStoreConcurrentBatches(t *testing.T) {
	var store state.Store
	now := time.Now().UTC()
	var workers sync.WaitGroup
	workers.Go(func() {
		for i := range 100 {
			stamp := now.Add(time.Duration(i) * time.Second)
			p := state.Position{ObservedAt: stamp, ReceivedAt: stamp}

			err := store.Upsert(state.Location{ID: 1, Position: &p}, state.Location{ID: 2, Position: &p})
			if err != nil {
				t.Error(err)
				return
			}
		}
	})
	workers.Go(func() {
		for range 100 {
			snapshot, _ := store.Snapshot()
			if len(snapshot) == 0 {
				continue
			}
			if len(snapshot) != 2 || !snapshot[0].Position.ObservedAt.Equal(snapshot[1].Position.ObservedAt) {
				t.Error("snapshot included a partial batch")
				return
			}
			store.Get(1)
		}
	})
	workers.Wait()
}
