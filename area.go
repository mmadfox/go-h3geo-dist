package h3geodist

// Supported H3 resolutions.
const (
	Level0 = iota // number of unique indexes 122
	Level1        // number of unique indexes 842
	Level2        // number of unique indexes 5882
	Level3        // number of unique indexes 41162
	Level4        // number of unique indexes 288122
	Level5        // number of unique indexes 2016842
	Level6        // number of unique indexes 14117882
)

// Table of cell areas for H3 resolutions.
var cellAreas = map[int]uint{
	Level0: 122,
	Level1: 842,
	Level2: 5882,
	Level3: 41162,
	Level4: 288122,
	Level5: 2016842,
	Level6: 14117882,
}

// Level0Area returns the area (km2) for level 0.
func Level0Area() uint {
	return cellArea(Level0)
}

// Level1Area returns the area (km2) for level 1.
func Level1Area() uint {
	return cellArea(Level1)
}

// Level2Area returns the area (km2) for level 2.
func Level2Area() uint {
	return cellArea(Level2)
}

// Level3Area returns the area (km2) for level 3.
func Level3Area() uint {
	return cellArea(Level3)
}

// Level4Area returns the area (km2) for level 4.
func Level4Area() uint {
	return cellArea(Level4)
}

// Level5Area returns the area (km2) for level 5.
func Level5Area() uint {
	return cellArea(Level5)
}

// Level6Area returns the area (km2) for level 6.
func Level6Area() uint {
	return cellArea(Level6)
}

// cellArea returns the area (km2) for specified level.
func cellArea(level int) uint {
	area, found := cellAreas[level]
	if !found {
		return 0
	}
	return area
}
