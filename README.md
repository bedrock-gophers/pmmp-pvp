# pmmp-pvp

`pmmp-pvp` provides small, deterministic Go implementations of PocketMine-MP's damage modifier and knockback motion calculations. It has no Dragonfly dependency, so adapters can translate framework types without tying the formulas to a particular server release.

The behavior is pinned to [PocketMine-MP 5.43.1](https://github.com/pmmp/PocketMine-MP/tree/5.43.1) (`763354d`). Random values are supplied by the caller, making combat behavior reproducible in tests.

## Install

```sh
go get github.com/bedrock-gophers/pmmp-pvp/damage
go get github.com/bedrock-gophers/pmmp-pvp/knockback
```

## Damage

```go
result, err := damage.Calculate(damage.Input{
	Cause:                       damage.CauseEntityAttack,
	BaseDamage:                  8,
	ArmorPoints:                 10,
	ResistanceLevel:             1,
	EnchantmentProtectionFactor: 7,
	ProtectionRoll:              73, // PMMP: mt_rand(50, 100)
	Absorption:                  2,
})
if err != nil {
	return err
}
health.Reduce(result.FinalDamage)
```

Damage modifiers are applied in PMMP order: previous-hit cooldown, armor, resistance, armor enchantments, absorption, then falling-block helmet reduction.

## Knockback

```go
result, err := knockback.Calculate(knockback.Input{
	Direction:      knockback.Vec3{X: targetX - attackerX, Z: targetZ - attackerZ},
	CurrentMotion:  knockback.Vec3{X: vx, Y: vy, Z: vz},
	Force:          0.4,
	Resistance:     0.0,
	ResistanceRoll: randomFloat, // PMMP: mt_rand()/mt_getrandmax()
})
if err != nil {
	return err
}
if result.Applied {
	setMotion(result.Motion)
}
```

The core intentionally excludes arena boundaries, combo rules, item parsing, and framework-specific event handling.

`VerticalForce` may be supplied separately for servers with independently configurable horizontal and vertical knockback. Leaving it nil preserves stock PMMP behaviour.

## License

MIT
