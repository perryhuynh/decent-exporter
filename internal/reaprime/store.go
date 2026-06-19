package reaprime

import (
	"sync"
	"time"
)

type Store struct {
	mu  sync.RWMutex
	now func() time.Time

	streams map[string]StreamState

	machine *MachineSnapshot
	scale   *ScaleSnapshot
	shot    *ShotSettings
	water   *WaterLevels
	display *DisplayState
	devices *DevicesSnapshot

	lastMachineState string
	transitions      map[string]uint64
}

func NewStore(now func() time.Time) *Store {
	return &Store{
		now:         now,
		streams:     map[string]StreamState{},
		transitions: map[string]uint64{},
	}
}

func (s *Store) SetStreamConnected(name string, connected bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := s.streams[name]
	if connected && !st.Connected {
		st.Reconnects++
	}
	st.Connected = connected
	s.streams[name] = st
}

func (s *Store) StreamError(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := s.streams[name]
	st.Errors++
	st.Connected = false
	s.streams[name] = st
}

func (s *Store) streamMessageLocked(name string) {
	st := s.streams[name]
	st.Messages++
	st.LastMessageTime = s.now()
	s.streams[name] = st
}

func (s *Store) SetMachine(v MachineSnapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streamMessageLocked("machine")
	if s.lastMachineState != "" && s.lastMachineState != v.State {
		s.transitions[v.State]++
	}
	s.lastMachineState = v.State
	s.machine = &v
}

func (s *Store) SetScale(v ScaleSnapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streamMessageLocked("scale")
	s.scale = &v
}

func (s *Store) SetShot(v ShotSettings) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streamMessageLocked("shot_settings")
	s.shot = &v
}

func (s *Store) SetWater(v WaterLevels) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streamMessageLocked("water_levels")
	s.water = &v
}

func (s *Store) SetDisplay(v DisplayState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streamMessageLocked("display")
	s.display = &v
}

func (s *Store) SetDevices(v DevicesSnapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streamMessageLocked("devices")
	s.devices = &v
}

func (s *Store) Snapshot() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	streams := make(map[string]StreamState, len(s.streams))
	for k, v := range s.streams {
		streams[k] = v
	}
	transitions := make(map[string]uint64, len(s.transitions))
	for k, v := range s.transitions {
		transitions[k] = v
	}

	return Snapshot{
		Now:         s.now(),
		Streams:     streams,
		Machine:     clone(s.machine),
		Scale:       clone(s.scale),
		Shot:        clone(s.shot),
		Water:       clone(s.water),
		Display:     clone(s.display),
		Devices:     clone(s.devices),
		Transitions: transitions,
	}
}

func (s *Store) Ready(maxAge time.Duration) bool {
	snap := s.Snapshot()
	machine, ok := snap.Streams["machine"]
	return ok && machine.Connected && !machine.LastMessageTime.IsZero() && snap.Now.Sub(machine.LastMessageTime) <= maxAge
}

func clone[T any](in *T) *T {
	if in == nil {
		return nil
	}
	out := *in
	return &out
}
