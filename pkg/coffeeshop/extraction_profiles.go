package coffeeshop

import "math"

type ExtractionStrength int

const (
	Normal ExtractionStrength = iota
	Medium
	Light
)

const dripCupOunces = 10.5822 // 300g water for standard pour over
const lightDripGrams = 15
const mediumDripGrams = 18
const normalDripGrams = 21 // roughly typical grind amount for a 300g pour over

var _ IExtractionProfile = &extractionProfile{}

type IExtractionProfile interface {
	GramsFromOunces(ounces int) int
}

type IExtractionProfiles interface {
	GetProfile(kind ExtractionStrength) IExtractionProfile
}

type extractionProfile struct {
	gramsNeededPerOunce float64
	// todo: grind setting? like espresso, drip, etc
	//grindSetting       int
}

type extractionProfiles struct {
	profiles map[ExtractionStrength]extractionProfile
}

// NewExtractionProfiles todo: configurability
func NewExtractionProfiles() IExtractionProfiles {
	return &extractionProfiles{
		profiles: map[ExtractionStrength]extractionProfile{
			Normal: {gramsNeededPerOunce: normalDripGrams / dripCupOunces},
			Medium: {gramsNeededPerOunce: mediumDripGrams / dripCupOunces},
			Light:  {gramsNeededPerOunce: lightDripGrams / dripCupOunces},
		},
	}
}

func (p *extractionProfiles) GetProfile(kind ExtractionStrength) IExtractionProfile {
	profile, ok := p.profiles[kind]
	if !ok {
		profile = p.profiles[Normal]
	}

	return &profile
}

// GramsFromOunces computes grams using float and returns a rounded int gram amount
func (p *extractionProfile) GramsFromOunces(ounces int) int {
	val := p.gramsNeededPerOunce * float64(ounces)
	return int(math.Round(val))
}
