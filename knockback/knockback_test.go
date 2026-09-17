package knockback

import (
	"errors"
	"math"
	"testing"
)

func TestCalculateMatchesPMMPMotion(t *testing.T) {
	result, err := Calculate(Input{
		Direction:      Vec3{X: 3, Z: 4},
		CurrentMotion:  Vec3{X: 0.2, Y: 0.6, Z: -0.2},
		Force:          0.4,
		ResistanceRoll: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Applied {
		t.Fatal("expected knockback to be applied")
	}
	assertClose(t, result.Motion.X, 0.34)
	assertClose(t, result.Motion.Y, 0.4)
	assertClose(t, result.Motion.Z, 0.22)
}

func TestCalculateUsesExplicitVerticalLimit(t *testing.T) {
	limit := 0.7
	result, err := Calculate(Input{
		Direction:      Vec3{X: 1},
		CurrentMotion:  Vec3{Y: 1},
		Force:          0.4,
		VerticalLimit:  &limit,
		ResistanceRoll: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result.Motion.Y, 0.7)
}

func TestCalculateUsesExplicitVerticalForce(t *testing.T) {
	verticalForce := 0.39
	result, err := Calculate(Input{
		Direction:      Vec3{X: 1},
		CurrentMotion:  Vec3{Y: -0.2},
		Force:          0.4,
		VerticalForce:  &verticalForce,
		VerticalLimit:  &verticalForce,
		ResistanceRoll: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result.Motion.Y, 0.29)
}

func TestCalculateSkipsResistedAndZeroDirection(t *testing.T) {
	current := Vec3{X: 1, Y: 2, Z: 3}
	for name, input := range map[string]Input{
		"resisted": {
			Direction: Vec3{X: 1}, CurrentMotion: current, Force: 0.4,
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
