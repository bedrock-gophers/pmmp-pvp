package damage

import (
	"errors"
	"math"
	"testing"
)

func TestCalculateAppliesModifiersInPMMPOrder(t *testing.T) {
	result, err := Calculate(Input{
		Cause:                       CauseEntityAttack,
		BaseDamage:                  10,
		ArmorPoints:                 10,
		ResistanceLevel:             2,
		EnchantmentProtectionFactor: 12,
		ProtectionRoll:              50,
		Absorption:                  1,
	})
	if err != nil {
		t.Fatal(err)
	}

	// 10 -> armour 6 -> resistance 3.6 -> EPF 2.736 -> absorption 1.736.
	assertClose(t, result.Modifiers.Armor, -4)
	assertClose(t, result.Modifiers.Resistance, -2.4)
	assertClose(t, result.Modifiers.ArmorEnchantments, -0.864)
	assertClose(t, result.Modifiers.Absorption, -1)
	assertClose(t, result.FinalDamage, 1.736)
}

func TestCalculatePreviousDamageCooldown(t *testing.T) {
	previous := 4.0
	result, err := Calculate(Input{
		Cause:              CauseEntityAttack,
		BaseDamage:         4,
		PreviousBaseDamage: &previous,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Cancelled {
		t.Fatal("expected equal repeated damage to be cancelled")
	}
	assertClose(t, result.FinalDamage, 0)
}

func TestCalculateIncludesExistingEventModifiers(t *testing.T) {
	result, err := Calculate(Input{
		Cause:                CauseEntityAttack,
		BaseDamage:           5,
		InitialModifierTotal: 3,
		ArmorPoints:          5,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result.Modifiers.Armor, -1.6)
	assertClose(t, result.FinalDamage, 6.4)
}

func TestCalculateCauseExceptionsAndHelmet(t *testing.T) {
	result, err := Calculate(Input{
		Cause:           CauseFallingBlock,
		BaseDamage:      8,
		ArmorPoints:     5,
		WearingHelmet:   true,
		ResistanceLevel: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	// 8 -> armour 6.4 -> resistance 5.12 -> helmet 3.84.
	assertClose(t, result.FinalDamage, 3.84)

	voidResult, err := Calculate(Input{
		Cause:           CauseVoid,
		BaseDamage:      8,
		ArmorPoints:     20,
		ResistanceLevel: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, voidResult.FinalDamage, 8)
}

func TestCalculateValidatesInputs(t *testing.T) {
	_, err := Calculate(Input{BaseDamage: -1})
	if !errors.Is(err, ErrNegativeInput) {
		t.Fatalf("expected ErrNegativeInput, got %v", err)
	}
	_, err = Calculate(Input{BaseDamage: 1, EnchantmentProtectionFactor: 1, ProtectionRoll: 49})
	if !errors.Is(err, ErrInvalidProtectionRoll) {
		t.Fatalf("expected ErrInvalidProtectionRoll, got %v", err)
	}
}

func assertClose(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-12 {
		t.Fatalf("got %v, want %v", got, want)
	}
}
