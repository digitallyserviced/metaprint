package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func GetRoundedFloat(value float64, rounding int) string {
	return fmt.Sprintf("%."+strconv.Itoa(rounding)+"f", value)
}

type UnitBits int64

const (
	Bb UnitBits = iota
	Kb
	Mb
	Gb
	Tb
	Pb
)

type HumanizedSpeed struct {
	speed      float64
	unit       UnitBits
	unitSpeed  float64
	unitString string
}

func UnitBitsSuffixes() []string {
	return []string{"b", "k", "m", "g", "t", "p", "e"}
}
func NewHSpeed(value float64) *HumanizedSpeed {
	hs := HumanizedSpeed{
		speed: value,
	}
	hs.BestGuess(1000)
	return &hs
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func (hs *HumanizedSpeed) BestGuess(base float64) {
	if hs.speed < 10 {
		hs.unitSpeed = hs.speed
		hs.unit = 0
		hs.unitString = "b"
		return
	}
	e := math.Floor(logn(float64(hs.speed), base))
	hs.unit = UnitBits(e)
	hs.unitString = UnitBitsSuffixes()[int(e)]
	val := math.Floor(float64(hs.speed)/math.Pow(base, e)*10+0.5) / 10
	hs.unitSpeed = val
}

func (hs *HumanizedSpeed) Parts() (float64, string) {
	speed := hs.unitSpeed
	unitString := hs.unitString
    return speed, unitString
}

func (hs *HumanizedSpeed) String() string {
	speed := hs.unitSpeed
	unitString := hs.unitString
	return fmt.Sprintf("%0.1f%s", speed, unitString)
}

func (hs *HumanizedSpeed) ToUnit(ub UnitBits) {
	hs.unit = ub
	hs.unitSpeed = hs.speed / ub.Multiplier()
	hs.unitString = ub.String()
}

func (ub UnitBits) EnumIndex() UnitBits {
	return UnitBits(ub)
}

func (ub UnitBits) Units() map[UnitBits]string {
	units := UnitBitsSuffixes()
	out := make(map[UnitBits]string)
	for i, v := range units {
		tub := UnitBits(i)
		out[tub] = v
	}
	return out
}

func (ub UnitBits) String() string {
	return ub.Units()[ub]
}

func StringToUnitBit(unit string) UnitBits {
    mainChar := unit[0:1]
    for i, v := range UnitBitsSuffixes() {
        if strings.ToLower(mainChar) == v {
            return UnitBits(i)
        }
    }
    return UnitBits(0)
}

func (ub UnitBits) GetFloat() float64 {
	return float64(ub)
}

func (ub UnitBits) GetInt() int64 {
	return int64(ub)
}

func (ub UnitBits) Multiplier() float64 {
	return float64((2 ^ ub.GetInt())*10)
}

func ValueToBitUnit(val float64, unit UnitBits) HumanizedSpeed {
	hs := NewHSpeed(val)
	hs.unit = unit
    hs.unitString = unit.String()
    hs.unitSpeed = hs.speed / unit.Multiplier()
	return *hs
}

func GetInUnit(val float64, unit string) float64 {
	switch strings.ToLower(unit) {
	case "kb":
		return val / Kb.Multiplier()
	case "mb":
		return val / Mb.Multiplier()
	case "gb":
		return val / Gb.Multiplier()
	}

	return val
}

func ReplaceVariables(str string, variables map[string]interface{}) string {
	for k, v := range variables {
		str = strings.ReplaceAll(str, "%"+k+"%", fmt.Sprintf("%v", v))
	}

	return str
}
