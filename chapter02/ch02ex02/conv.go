package main

import "fmt"

type Celsius float64
type Fahrenheit float64
type Feet float64
type Meters float64
type Pounds float64
type Kilograms float64

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC     Celsius = 0
	BoilingC      Celsius = 100
)

func (c Celsius) String() string    { return fmt.Sprintf("%.3g°C", c) }
func (f Fahrenheit) String() string { return fmt.Sprintf("%.3g°F", f) }
func (ft Feet) String() string      { return fmt.Sprintf("%.3g ft", ft) }
func (m Meters) String() string     { return fmt.Sprintf("%.3g m", m) }
func (lbs Pounds) String() string   { return fmt.Sprintf("%.3g lbs", lbs) }
func (kg Kilograms) String() string { return fmt.Sprintf("%.3g kg", kg) }

// CToF converts a Celsius temperature to Fahrenheit.
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// FToC converts a Fahrenheit temperature to Celsius.
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

// FtToM converts a distance in feet into meters.
func FtToM(ft Feet) Meters { return Meters(ft / 3.281) }

// MToFt converts a distance in meters into feet.
func MToFt(m Meters) Feet { return Feet(m * 3.281) }

// LbsToKg converts a mass in pounds into kilograms.
func LbsToKg(lbs Pounds) Kilograms { return Kilograms(lbs / 2.2) }

// KgToLbs converts a mass in kilograms into pounds.
func KgToLbs(kg Kilograms) Pounds { return Pounds(kg * 2.2) }
