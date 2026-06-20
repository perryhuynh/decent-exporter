package decent

import (
	"sync"
	"time"
)

type Store struct {
	mu  sync.RWMutex
	now func() time.Time

	streams map[string]StreamState

	machine *MachineSnapshot
	shot    *ShotSettings
	water   *WaterLevels
	display *DisplayState
	devices *DevicesSnapshot

	lastMachineState string
	transitions      map[string]uint64

	accum      shotAccum
	lastShot   *ShotSummary
	shotsTotal uint64
}

// shotState is the machine state that marks an in-progress espresso shot.
const shotState = "espresso"

// minShotDuration discards espresso episodes shorter than this (flushes, aborts).
const minShotDuration = 3 * time.Second

type shotAccum struct {
	active       bool
	startTime    time.Time
	peakPressure float64
	peakFlow     float64
	flowSum      float64
	flowCount    float64
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
	s.updateShotLocked(v)
	s.lastMachineState = v.State
	s.machine = &v
}

// updateShotLocked accumulates per-shot aggregates as machine samples arrive and
// finalizes a shot summary when the machine leaves the espresso state. Caller holds s.mu.
func (s *Store) updateShotLocked(v MachineSnapshot) {
	switch {
	case v.State == shotState && !s.accum.active:
		s.accum = shotAccum{active: true, startTime: s.now()}
		fallthrough
	case v.State == shotState:
		if v.Pressure > s.accum.peakPressure {
			s.accum.peakPressure = v.Pressure
		}
		if v.Flow > s.accum.peakFlow {
			s.accum.peakFlow = v.Flow
		}
		s.accum.flowSum += v.Flow
		s.accum.flowCount++
	case s.accum.active:
		s.finalizeShotLocked()
	}
}

func (s *Store) finalizeShotLocked() {
	duration := s.now().Sub(s.accum.startTime)
	if duration >= minShotDuration {
		var avg float64
		if s.accum.flowCount > 0 {
			avg = s.accum.flowSum / s.accum.flowCount
		}
		s.lastShot = &ShotSummary{
			Duration:     duration,
			PeakPressure: s.accum.peakPressure,
			PeakFlow:     s.accum.peakFlow,
			AverageFlow:  avg,
		}
		s.shotsTotal++
	}
	s.accum = shotAccum{}
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
		Shot:        clone(s.shot),
		Water:       clone(s.water),
		Display:     clone(s.display),
		Devices:     clone(s.devices),
		Transitions: transitions,
		LastShot:    clone(s.lastShot),
		ShotsTotal:  s.shotsTotal,
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
