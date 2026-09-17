package knockback

import (
	"errors"
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func TestCalculateMatchesPMMPMotion(t *testing.T) {
	result, err := Calculate(Input{
		Direction:      mgl64.Vec3{3, 0, 4},
		CurrentMotion:  mgl64.Vec3{0.2, 0.6, -0.2},
		Force:          0.4,
		ResistanceRoll: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Applied {
		t.Fatal("expected knockback to be applied")
	}
	assertClose(t, result.Motion[0], 0.34)
	assertClose(t, result.Motion[1], 0.4)
	assertClose(t, result.Motion[2], 0.22)
}

func TestCalculateUsesExplicitVerticalLimit(t *testing.T) {
	limit := 0.7
	result, err := Calculate(Input{
		Direction:      mgl64.Vec3{1, 0, 0},
		CurrentMotion:  mgl64.Vec3{0, 1, 0},
		Force:          0.4,
		VerticalLimit:  &limit,
		ResistanceRoll: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result.Motion[1], 0.7)
}

func TestCalculateUsesExplicitVerticalForce(t *testing.T) {
	verticalForce := 0.39
	result, err := Calculate(Input{
		Direction:      mgl64.Vec3{1, 0, 0},
		CurrentMotion:  mgl64.Vec3{0, -0.2, 0},
		Force:          0.4,
		VerticalForce:  &verticalForce,
		VerticalLimit:  &verticalForce,
		ResistanceRoll: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result.Motion[1], 0.29)
}

func TestCalculateSkipsResistedAndZeroDirection(t *testing.T) {
	current := mgl64.Vec3{1, 2, 3}
	for name, input := range map[string]Input{
		"resisted": {
			Direction: mgl64.Vec3{1, 0, 0}, CurrentMotion: current, Force: 0.4,
			Resistance: 0.5, ResistanceRoll: 0.5,
		},
		"zero direction": {CurrentMotion: current, Force: 0.4, ResistanceRoll: 1},
	} {
		t.Run(name, func(t *testing.T) {
			result, err := Calculate(input)
			if err != nil {
				t.Fatal(err)
			}
			if result.Applied || result.Motion != current {
				t.Fatalf("unexpected result: %+v", result)
			}
		})
	}
}

func TestCalculateValidatesInputs(t *testing.T) {
	tests := []struct {
		input Input
		want  error
	}{
		{Input{Force: -1}, ErrNegativeForce},
		{Input{Resistance: 1.1}, ErrInvalidResistance},
		{Input{ResistanceRoll: -0.1}, ErrInvalidResistanceRoll},
	}
	for _, test := range tests {
		_, err := Calculate(test.input)
		if !errors.Is(err, test.want) {
			t.Fatalf("expected %v, got %v", test.want, err)
		}
	}
}

func assertClose(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-12 {
		t.Fatalf("got %v, want %v", got, want)
	}
}
