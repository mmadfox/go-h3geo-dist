package h3geodist

const (
	Level0 = iota
	Level1
	Level2
	Level3
	Level4
	Level5
	Level6
)

var cellAreas = map[int]uint{
	Level0: 122,
	Level1: 842,
	Level2: 5882,
	Level3: 41162,
	Level4: 288122,
	Level5: 2016842,
	Level6: 14117882,
}

func Level0Area() uint {
	return cellArea(Level0)
}

func Level1Area() uint {
	return cellArea(Level1)
}

func Level2Area() uint {
	return cellArea(Level2)
}

func Level3Area() uint {
	return cellArea(Level3)
}

func Level4Area() uint {
	return cellArea(Level4)
}

func Level5Area() uint {
	return cellArea(Level5)
}

func Level6Area() uint {
	return cellArea(Level6)
}

func cellArea(level int) uint {
	area, found := cellAreas[level]
	if !found {
		return 0
	}
	return area
}
