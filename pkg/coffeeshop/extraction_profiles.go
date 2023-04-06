package coffeeshop

import (
	"coffeeshop/pkg/util"
	"math"
)

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
	log                 *util.Logger
	gramsNeededPerOunce float64
	// todo: grind setting? like espresso, drip, etc
	// grindSetting       int
}

type extractionProfiles struct {
	profiles map[ExtractionStrength]extractionProfile
}

// NewExtractionProfiles todo: configurability
func NewExtractionProfiles() IExtractionProfiles {
	log := util.NewLogger("ExtractionProfile")
	return &extractionProfiles{
		profiles: map[ExtractionStrength]extractionProfile{
			Normal: {log: log, gramsNeededPerOunce: normalDripGrams / dripCupOunces},
			Medium: {log: log, gramsNeededPerOunce: mediumDripGrams / dripCupOunces},
			Light:  {log: log, gramsNeededPerOunce: lightDripGrams / dripCupOunces},
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
	vali := int(math.Round(val))
	p.log.Infof("%v grams for %v ounces", vali, ounces)
	return vali
}
