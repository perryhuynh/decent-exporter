package decent

import (
	"testing"
	"time"
)

func TestStoreReadyRequiresFreshMachineStream(t *testing.T) {
	now := time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC)
	store := NewStore(func() time.Time { return now })

	if store.Ready(30 * time.Second) {
		t.Fatal("store should not be ready before machine data")
	}

	store.SetStreamConnected("machine", true)
	store.SetMachine(MachineSnapshot{State: "idle"})
	if !store.Ready(30 * time.Second) {
		t.Fatal("store should be ready with fresh machine data")
	}
}

func TestStoreCountsStateTransitions(t *testing.T) {
	now := time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC)
	store := NewStore(func() time.Time { return now })
	store.SetMachine(MachineSnapshot{State: "idle"})
	store.SetMachine(MachineSnapshot{State: "espresso"})
	store.SetMachine(MachineSnapshot{State: "steam"})

	snap := store.Snapshot()
	if snap.Transitions["espresso"] != 1 {
		t.Fatalf("espresso transitions = %d", snap.Transitions["espresso"])
	}
	if snap.Transitions["steam"] != 1 {
		t.Fatalf("steam transitions = %d", snap.Transitions["steam"])
	}
}

func TestStoreSummarizesShot(t *testing.T) {
	now := time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC)
	store := NewStore(func() time.Time { return now })

	store.SetMachine(MachineSnapshot{State: "idle"})

	// Espresso shot: pressure peaks at 9, flow peaks at 3, mean flow = 2.
	store.SetMachine(MachineSnapshot{State: "espresso", Pressure: 2, Flow: 1})
	now = now.Add(15 * time.Second)
	store.SetMachine(MachineSnapshot{State: "espresso", Pressure: 9, Flow: 3})
	now = now.Add(15 * time.Second)
	store.SetMachine(MachineSnapshot{State: "espresso", Pressure: 6, Flow: 2})

	// Leave espresso -> finalize.
	store.SetMachine(MachineSnapshot{State: "idle"})

	snap := store.Snapshot()
	if snap.ShotsTotal != 1 {
		t.Fatalf("shots total = %d", snap.ShotsTotal)
	}
	if snap.LastShot == nil {
		t.Fatal("expected a last shot summary")
	}
	if snap.LastShot.Duration != 30*time.Second {
		t.Fatalf("duration = %v", snap.LastShot.Duration)
	}
	if snap.LastShot.PeakPressure != 9 {
		t.Fatalf("peak pressure = %v", snap.LastShot.PeakPressure)
	}
	if snap.LastShot.PeakFlow != 3 {
		t.Fatalf("peak flow = %v", snap.LastShot.PeakFlow)
	}
	if snap.LastShot.AverageFlow != 2 {
		t.Fatalf("average flow = %v", snap.LastShot.AverageFlow)
	}
}

func TestStoreDiscardsShortShot(t *testing.T) {
	now := time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC)
	store := NewStore(func() time.Time { return now })

	store.SetMachine(MachineSnapshot{State: "idle"})
	store.SetMachine(MachineSnapshot{State: "espresso", Pressure: 5, Flow: 2})
	now = now.Add(time.Second) // shorter than minShotDuration
	store.SetMachine(MachineSnapshot{State: "idle"})

	snap := store.Snapshot()
	if snap.ShotsTotal != 0 {
		t.Fatalf("shots total = %d, want 0", snap.ShotsTotal)
	}
	if snap.LastShot != nil {
		t.Fatalf("expected no shot summary, got %#v", snap.LastShot)
	}
}
