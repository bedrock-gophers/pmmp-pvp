// Package knockback implements PocketMine-MP's living-entity knockback motion.
package knockback

import (
	"errors"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// Input contains the state used by PMMP to calculate knockback. Direction.Y is
// ignored. ResistanceRoll corresponds to mt_rand()/mt_getrandmax() and must be
// in [0, 1]. A nil VerticalForce uses Force, preserving stock PMMP behaviour.
// A nil VerticalLimit also uses Force, as PMMP does by default.
type Input struct {
	Direction      mgl64.Vec3
	CurrentMotion  mgl64.Vec3
	Force          float64
	VerticalForce  *float64
	VerticalLimit  *float64
	Resistance     float64
	ResistanceRoll float64
}

// Result is the resulting motion and whether knockback was applied.
type Result struct {
	Motion  mgl64.Vec3
	Applied bool
}

var (
	ErrNegativeForce         = errors.New("force must not be negative")
	ErrInvalidResistance     = errors.New("resistance must be between 0 and 1 inclusive")
	ErrInvalidResistanceRoll = errors.New("resistance roll must be between 0 and 1 inclusive")
)

// Calculate applies the Living.KnockBack algorithm from PocketMine-MP 5.43.1.
func Calculate(in Input) (Result, error) {
	if in.Force < 0 {
		return Result{}, ErrNegativeForce
	}
	if in.VerticalForce != nil && *in.VerticalForce < 0 {
		return Result{}, ErrNegativeForce
	}
	if in.Resistance < 0 || in.Resistance > 1 {
		return Result{}, ErrInvalidResistance
	}
	if in.ResistanceRoll < 0 || in.ResistanceRoll > 1 {
		return Result{}, ErrInvalidResistanceRoll
	}

	result := Result{Motion: in.CurrentMotion}
	length := math.Hypot(in.Direction[0], in.Direction[2])
	if length <= 0 || in.ResistanceRoll <= in.Resistance {
		return result, nil
	}

	verticalForce := in.Force
	if in.VerticalForce != nil {
		verticalForce = *in.VerticalForce
	}
	result.Motion = mgl64.Vec3{
		in.CurrentMotion[0]/2 + in.Direction[0]/length*in.Force,
		in.CurrentMotion[1]/2 + verticalForce,
		in.CurrentMotion[2]/2 + in.Direction[2]/length*in.Force,
	}
	limit := in.Force
	if in.VerticalLimit != nil {
		limit = *in.VerticalLimit
	}
	if result.Motion[1] > limit {
		result.Motion[1] = limit
	}
	result.Applied = true
	return result, nil
}
