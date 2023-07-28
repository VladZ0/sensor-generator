package mocks

import (
	"sensors-generator/internal/group"
	"sensors-generator/internal/sensor"
	"sensors-generator/internal/spiece"
)

const (
	alphaTime   = 30
	betaTime    = 45
	gammaTime   = 40
	deltaTime   = 30
	epsilonTime = 35
)

var (
	CreateSensorGroups = []group.CreateSensorGroupDTO{
		{
			Name: "alpha",
		},
		{
			Name: "beta",
		},
		{
			Name: "gamma",
		},
		{
			Name: "delta",
		},
		{
			Name: "epsilon",
		},
	}

	CreateSensors = []sensor.CreateSensorDTO{
		{
			CodeName: sensor.Codename{
				GroupName: "alpha",
				Index:     1,
			},
			Coords: sensor.Coordinates{
				X: 34.33,
				Y: 27.24,
				Z: 4,
			},
			DataOutputRate: alphaTime,
		},

		{
			CodeName: sensor.Codename{
				GroupName: "alpha",
				Index:     2,
			},
			Coords: sensor.Coordinates{
				X: 30,
				Y: 21.88,
				Z: 5,
			},
			DataOutputRate: alphaTime,
		},

		{
			CodeName: sensor.Codename{
				GroupName: "alpha",
				Index:     3,
			},
			Coords: sensor.Coordinates{
				X: 55.33,
				Y: 32.24,
				Z: 5,
			},
			DataOutputRate: alphaTime,
		},

		{
			CodeName: sensor.Codename{
				GroupName: "beta",
				Index:     1,
			},
			Coords: sensor.Coordinates{
				X: 16,
				Y: 27.88,
				Z: 8,
			},
			DataOutputRate: betaTime,
		},

		{
			CodeName: sensor.Codename{
				GroupName: "beta",
				Index:     2,
			},
			Coords: sensor.Coordinates{
				X: 33.33,
				Y: 33.24,
				Z: 7,
			},
			DataOutputRate: betaTime,
		},

		{
			CodeName: sensor.Codename{
				GroupName: "gamma",
				Index:     1,
			},
			Coords: sensor.Coordinates{
				X: 123.33,
				Y: 46.24,
				Z: 12,
			},
			DataOutputRate: gammaTime,
		},

		{
			CodeName: sensor.Codename{
				GroupName: "gamma",
				Index:     2,
			},
			Coords: sensor.Coordinates{
				X: 102,
				Y: 38,
				Z: 13,
			},
			DataOutputRate: gammaTime,
		},

		{
			CodeName: sensor.Codename{
				GroupName: "gamma",
				Index:     3,
			},
			Coords: sensor.Coordinates{
				X: 144.66,
				Y: 34.11,
				Z: 13,
			},
			DataOutputRate: gammaTime,
		},

		{
			CodeName: sensor.Codename{
				GroupName: "delta",
				Index:     1,
			},
			Coords: sensor.Coordinates{
				X: 87.33,
				Y: 68.24,
				Z: 2,
			},
			DataOutputRate: deltaTime,
		},

		{
			CodeName: sensor.Codename{
				GroupName: "delta",
				Index:     2,
			},
			Coords: sensor.Coordinates{
				X: 101.99,
				Y: 57.76,
				Z: 2,
			},
			DataOutputRate: deltaTime,
		},

		{
			CodeName: sensor.Codename{
				GroupName: "epsilon",
				Index:     1,
			},
			Coords: sensor.Coordinates{
				X: 213.45,
				Y: 66.24,
				Z: 7,
			},
			DataOutputRate: epsilonTime,
		},

		{
			CodeName: sensor.Codename{
				GroupName: "epsilon",
				Index:     2,
			},
			Coords: sensor.Coordinates{
				X: 189.87,
				Y: 56.64,
				Z: 7,
			},
			DataOutputRate: epsilonTime,
		},

		{
			CodeName: sensor.Codename{
				GroupName: "epsilon",
				Index:     3,
			},
			Coords: sensor.Coordinates{
				X: 192.48,
				Y: 44.59,
				Z: 8,
			},
			DataOutputRate: epsilonTime,
		},
	}

	CreateSpieces = []spiece.CreateSpieceDTO{
		{
			Name: "Atlantic Bluefin Tuna",
		},

		{
			Name: "Atlantic Cod",
		},

		{
			Name: "Atlantic Goliath Grouper",
		},

		{
			Name: "Banded Butterflyfish",
		},

		{
			Name: "Beluga Sturgeon",
		},

		{
			Name: "Blue Marlin",
		},

		{
			Name: "Blue Tang",
		},

		{
			Name: "Bluebanded Goby",
		},

		{
			Name: "Bluehead Wrasse",
		},

		{
			Name: "California Grunion",
		},

		{
			Name: "Clown Triggerfish",
		},

		{
			Name: "Coelacanth",
		},

		{
			Name: "Flashlight Fish",
		},

		{
			Name: "French Angelfish",
		},

		{
			Name: "John Dory",
		},

		{
			Name: "Nassau Grouper",
		},

		{
			Name: "Ocean Sunfish",
		},

		{
			Name: "Pacific Herring",
		},

		{
			Name: "Patagonian Toothfish",
		},

		{
			Name: "Sailfish",
		},
	}
)
