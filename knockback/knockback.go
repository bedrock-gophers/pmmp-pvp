// Package knockback implements PocketMine-MP's living-entity knockback motion.
package knockback

import (
	"errors"
	"math"
)

// Vec3 is a three-dimensional velocity or direction vector.
type Vec3 struct {
	X float64
	Y float64
	Z float64
}

// Input contains the state used by PMMP to calculate knockback. Direction.Y is
// ignored. ResistanceRoll corresponds to mt_rand()/mt_getrandmax() and must be
// in [0, 1]. A nil VerticalForce uses Force, preserving stock PMMP behaviour.
// A nil VerticalLimit also uses Force, as PMMP does by default.
type Input struct {
	Direction      Vec3
	CurrentMotion  Vec3
	Force          float64
	VerticalForce  *float64
	VerticalLimit  *float64
	Resistance     float64
	ResistanceRoll float64
}

// Result is the resulting motion and whether knockback was applied.
type Result struct {
	Motion  Vec3
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
	length := math.Hypot(in.Direction.X, in.Direction.Z)
	if length <= 0 || in.ResistanceRoll <= in.Resistance {
		return result, nil
	}

	verticalForce := in.Force
	if in.VerticalForce != nil {
		verticalForce = *in.VerticalForce
	}
	result.Motion = Vec3{
		X: in.CurrentMotion.X/2 + in.Direction.X/length*in.Force,
		Y: in.CurrentMotion.Y/2 + verticalForce,
		Z: in.CurrentMotion.Z/2 + in.Direction.Z/length*in.Force,
	}
	limit := in.Force
	if in.VerticalLimit != nil {
		limit = *in.VerticalLimit
	}
	if result.Motion.Y > limit {
		result.Motion.Y = limit
	}
	result.Applied = true
	return result, nil
}
