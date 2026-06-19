package reaprime

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
