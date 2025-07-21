package main

import (
	"math"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		l := len(name)
		if l < 0 || l > 42 {
			return
		}

		var n [42]byte
		for i, r := range name {
			n[i] = byte(r)
		}
		if nil == person.name {
			person.name = new([42]byte)
		}
		*person.name = n
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		if mana < 0 || mana > 1000 {
			panic("mana is greater than 1000 or less than 0")
		}
		// cleanup old value
		person.stats &= statsManaMask
		person.stats |= uint32(mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		if health < 0 || health > 1000 {
			panic("health is greater than 1000 or less than 0")
		}
		// cleanup old value
		person.stats &= statsHealthMask
		person.stats |= uint32(health) << statsHealthOffset
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		if strength < 0 || strength > 10 {
			panic("strength is greater than 10 or less than 0")
		}
		// cleanup old value
		person.stats &= statsStrengthMask
		person.stats |= uint32(strength) << statsStrengthOffset
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		if respect < 0 || respect > 10 {
			panic("strength is greater than 10 or less than 0")
		}
		// cleanup old value
		person.stats &= statsRespectMask
		person.stats |= uint32(respect) << statsRespectOffset
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		if experience < 0 || experience > 10 {
			panic("strength is greater than 10 or less than 0")
		}
		// cleanup old value
		person.stats &= statsExperienceMask
		person.stats |= uint32(experience) << statsExperienceOffset
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.level = uint8(level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags = setBit(person.flags, flagsHasHouseBit)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags = setBit(person.flags, flagsHasGunBit)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags = setBit(person.flags, flagsHasFamilyBit)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		if getBit(person.flags, flagsPersonTypeControlBit) {
			// cleanup person type bit
			person.flags &= 0b1111_0000
		}
		person.flags = setBit(person.flags, personType)
		person.flags = setBit(person.flags, flagsPersonTypeControlBit)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	flagsPersonTypeControlBit = iota + 4
	flagsHasFamilyBit
	flagsHasGunBit
	flagsHasHouseBit
)

const (
	statsManaMask         = 0b1111_1111_1111_1111111111_0000000000
	statsManaMaskInverted = 0b0000_0000_0000_0000000000_1111111111

	statsHealthMask         = 0b1111_1111_1111_0000000000_1111111111
	statsHealthMaskInverted = 0b0000_0000_0000_1111111111_0000000000
	statsHealthOffset       = 10

	statsStrengthMask         = 0b1111_1111_0000_1111111111_1111111111
	statsStrengthMaskInverted = 0b0000_0000_1111_0000000000_0000000000
	statsStrengthOffset       = 20

	statsRespectMask         = 0b1111_0000_1111_1111111111_1111111111
	statsRespectMaskInverted = 0b0000_1111_0000_0000000000_0000000000
	statsRespectOffset       = 24

	statsExperienceMask         = 0b0000_1111_1111_1111111111_1111111111
	statsExperienceMaskInverted = 0b1111_0000_0000_0000000000_0000000000
	statsExperienceOffset       = 28
)

type GamePerson struct {
	x, y, z int32  // 12 bytes
	gold    uint32 // 4 bytes
	// Exp	    Respect   Strength   Health (10-bit)   Mana (10-bit)
	// [0000]   [0000]    [0000]     [0000000000]      [0000000000]
	stats uint32    // 4 bytes
	name  *[42]byte // 42 bytes
	level uint8     // 1 byte
	// [0, hasHouse, hasGun, hasFamily] + [personTypeControlBit, Warrior, Blacksmith, Builder]
	flags uint8 // 1 byte
}

func getBit(num uint8, idx int) bool {
	return ((num & (1 << idx)) != 0)
}

func setBit(num uint8, idx int) uint8 {
	return num | 1<<idx
}

func NewGamePerson(options ...Option) GamePerson {
	p := GamePerson{}
	for _, option := range options {
		option(&p)
	}
	return p
}

func (p *GamePerson) Name() string {
	if nil == p.name {
		return ""
	}
	b := strings.Builder{}
	b.Write(p.name[:])
	return b.String()
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.stats & statsManaMaskInverted)
}

func (p *GamePerson) Health() int {
	return int((p.stats & statsHealthMaskInverted) >> statsHealthOffset)
}

func (p *GamePerson) Strength() int {
	return int((p.stats & statsStrengthMaskInverted) >> statsStrengthOffset)
}

func (p *GamePerson) Respect() int {
	return int((p.stats & statsRespectMaskInverted) >> statsRespectOffset)
}

func (p *GamePerson) Experience() int {
	return int((p.stats & statsExperienceMaskInverted) >> statsExperienceOffset)
}

func (p *GamePerson) Level() int {
	return int(p.level)
}

func (p *GamePerson) HasHouse() bool {
	return getBit(p.flags, flagsHasHouseBit)
}

func (p *GamePerson) HasGun() bool {
	return getBit(p.flags, flagsHasGunBit)
}

func (p *GamePerson) HasFamily() bool {
	return getBit(p.flags, flagsHasFamilyBit)
}

func (p *GamePerson) Type() int {
	switch int(p.flags & 0b00000111) {
	case 1 << BuilderGamePersonType:
		return BuilderGamePersonType
	case 1 << BlacksmithGamePersonType:
		return BlacksmithGamePersonType
	case 1 << WarriorGamePersonType:
		return WarriorGamePersonType
	}
	return -1
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = WarriorGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
