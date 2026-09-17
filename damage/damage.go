// Package damage implements PocketMine-MP's living-entity damage modifiers.
package damage

import (
	"errors"
	"math"
)

// Cause identifies a PocketMine-MP EntityDamageEvent cause.
type Cause uint8

const (
	CauseContact Cause = iota
	CauseEntityAttack
	CauseProjectile
	CauseSuffocation
	CauseFall
	CauseFire
	CauseFireTick
	CauseLava
	CauseDrowning
	CauseBlockExplosion
	CauseEntityExplosion
	CauseVoid
	CauseSuicide
	CauseMagic
	CauseCustom
	CauseStarvation
	CauseFallingBlock
)

// ReducedByArmor reports whether PMMP applies armour-point reduction for c.
func (c Cause) ReducedByArmor() bool {
	switch c {
	case CauseFireTick, CauseSuffocation, CauseDrowning, CauseStarvation,
		CauseFall, CauseVoid, CauseMagic, CauseSuicide:
		return false
	default:
		return true
	}
}

// Input contains the state used by PMMP to calculate damage modifiers.
// ResistanceLevel is one-based, matching PMMP's EffectInstance.GetEffectLevel.
// PreviousBaseDamage is set while the target's attack cooldown is active.
type Input struct {
	Cause      Cause
	BaseDamage float64
	// InitialModifierTotal is the sum of modifiers already present on the
	// EntityDamageEvent, such as critical or weapon-enchantment damage.
	InitialModifierTotal        float64
	PreviousBaseDamage          *float64
	ArmorPoints                 int
	ResistanceLevel             int
	EnchantmentProtectionFactor int
	ProtectionRoll              int
	Absorption                  float64
	WearingHelmet               bool
}

// Modifiers contains each modifier in PMMP application order. Values reducing
// damage are negative.
type Modifiers struct {
	PreviousDamageCooldown float64
	Armor                  float64
	Resistance             float64
	ArmorEnchantments      float64
	Absorption             float64
	ArmorHelmet            float64
}

// Result is the outcome of a damage calculation.
type Result struct {
	FinalDamage float64
	Cancelled   bool
	Modifiers   Modifiers
}

var (
	ErrNegativeInput         = errors.New("damage inputs must not be negative")
	ErrInvalidProtectionRoll = errors.New("protection roll must be between 50 and 100 inclusive")
)

// Calculate applies the Living.ApplyDamageModifiers algorithm from
// PocketMine-MP 5.43.1. ProtectionRoll must match mt_rand(50, 100) whenever
// EnchantmentProtectionFactor is greater than zero.
func Calculate(in Input) (Result, error) {
	if in.BaseDamage < 0 || in.ArmorPoints < 0 || in.ResistanceLevel < 0 ||
		in.EnchantmentProtectionFactor < 0 || in.Absorption < 0 ||
		(in.PreviousBaseDamage != nil && *in.PreviousBaseDamage < 0) {
		return Result{}, ErrNegativeInput
	}
	if in.EnchantmentProtectionFactor > 0 && (in.ProtectionRoll < 50 || in.ProtectionRoll > 100) {
		return Result{}, ErrInvalidProtectionRoll
	}

	result := Result{}
	modifierTotal := in.InitialModifierTotal
	finalDamage := func() float64 {
		return math.Max(0, in.BaseDamage+modifierTotal)
	}
	apply := func(modifier float64) {
		modifierTotal += modifier
	}

	if in.PreviousBaseDamage != nil {
		if *in.PreviousBaseDamage >= in.BaseDamage {
			result.Cancelled = true
		}
		result.Modifiers.PreviousDamageCooldown = -*in.PreviousBaseDamage
		apply(result.Modifiers.PreviousDamageCooldown)
	}
	if in.Cause.ReducedByArmor() {
		result.Modifiers.Armor = -finalDamage() * float64(in.ArmorPoints) * 0.04
		apply(result.Modifiers.Armor)
	}
	if in.ResistanceLevel > 0 && in.Cause != CauseVoid && in.Cause != CauseSuicide {
		result.Modifiers.Resistance = -finalDamage() * math.Min(1, 0.2*float64(in.ResistanceLevel))
		apply(result.Modifiers.Resistance)
	}
	if in.EnchantmentProtectionFactor > 0 {
		points := math.Ceil(math.Min(float64(in.EnchantmentProtectionFactor), 25) * float64(in.ProtectionRoll) / 100)
		result.Modifiers.ArmorEnchantments = -finalDamage() * math.Min(points, 20) * 0.04
		apply(result.Modifiers.ArmorEnchantments)
	}
	result.Modifiers.Absorption = -math.Min(in.Absorption, finalDamage())
	apply(result.Modifiers.Absorption)
	if in.Cause == CauseFallingBlock && in.WearingHelmet {
		result.Modifiers.ArmorHelmet = -finalDamage() / 4
		apply(result.Modifiers.ArmorHelmet)
	}

	result.FinalDamage = finalDamage()
	return result, nil
}
