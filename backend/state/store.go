package state

import (
	"cmp"
	"fmt"
	"slices"
	"sync"
)

// latest location per device id
// MUST NOT BE COPIED AFTER FIRST USE
type Store struct {
	mu        sync.RWMutex
	locations map[int64]Location
	watchOnce sync.Once
	changed   chan struct{}
}

// applies batch atomically
// nil position updates metadata, retains existing position
// positions ordered by observation time then server time
func (s *Store) Upsert(locations ...Location) error {
	batch := make([]Location, 0, len(locations))
	for _, location := range locations {
		if location.ID <= 0 {
			return fmt.Errorf("device ID must be positive: %d", location.ID)
		}
		if location.Position != nil {
			err := location.Position.validate()
			if err != nil {
				return fmt.Errorf("device %d: %w", location.ID, err)
			}
		}
		location = location.clone()
		p := location.Position
		if p != nil {
			p.ObservedAt = p.ObservedAt.UTC()
			p.ServerTime = p.ServerTime.UTC()
			p.ReceivedAt = p.ReceivedAt.UTC()
		}
		batch = append(batch, location)
	}
	if len(batch) == 0 {
		return nil
	}

	s.initWatch()
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.locations == nil {
		s.locations = make(map[int64]Location, len(batch))
	}
	changed := false
	for _, location := range batch {
		current := s.locations[location.ID]
		if !newerPosition(location.Position, current.Position) {
			location.Position = current.Position
		}
		if location != current {
			s.locations[location.ID] = location
			changed = true
		}
	}
	if changed {
		close(s.changed)
		s.changed = make(chan struct{})
	}
	return nil
}

// returns detached copy
func (s *Store) Get(deviceID int64) (Location, bool) {
	s.mu.RLock()
	location, ok := s.locations[deviceID]
	s.mu.RUnlock()
	return location.clone(), ok
}

// returns detached list sorted by device ID
// + channel (closes when store changes after snapshot)
// call again after channel closes
// multiple updates can coalesce for slow reader
func (s *Store) Snapshot() ([]Location, <-chan struct{}) {
	s.initWatch()
	s.mu.RLock()
	locations := make([]Location, 0, len(s.locations))
	for _, location := range s.locations {
		locations = append(locations, location)
	}
	changed := s.changed
	s.mu.RUnlock()

	// stored positions are immutable
	// can copy outside lock
	for i, location := range locations {
		locations[i] = location.clone()
	}
	slices.SortFunc(locations, func(a, b Location) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return locations, changed
}

func (s *Store) initWatch() {
	s.watchOnce.Do(func() {
		s.changed = make(chan struct{})
	})
}

func newerPosition(incoming, current *Position) bool {
	if incoming == nil {
		return false
	}
	if current == nil {
		return true
	}
	if incoming.ObservedAt.Equal(current.ObservedAt) {
		return incoming.ServerTime.After(current.ServerTime)
	}
	return incoming.ObservedAt.After(current.ObservedAt)
}
