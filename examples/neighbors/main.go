package main

import (
	"fmt"
	"strconv"
	"strings"

	h3geodist "github.com/mmadfox/go-h3geo-dist"
)

func main() {
	coords := coordsFromString(coordinates)

	h3dist, err := h3geodist.New(h3geodist.Level5)
	if err != nil {
		panic(err)
	}

	var curHost string

	_ = h3dist.Add("127.0.0.1")
	_ = h3dist.Add("127.0.0.2")

	for i := 0; i < len(coords); i++ {
		cord := coords[i]
		target, neighbors, err := h3dist.NeighborsFromLatLon(cord[0], cord[1])
		if err != nil {
			panic(err)
		}
		if curHost != target.Host {
			curHost = target.Host

			neighbor0 := neighbors[0]
			neighbor1 := neighbors[1]

			fmt.Printf("host=%s\ncurrent=%s\nfrom=%s - %.2f%s \n--\n",
				target.Host,
				target.HexID(),
				neighbor0.Cell.HexID(),
				toPercent(neighbor0.DistanceM, neighbor1.DistanceM), "%")
		}
	}
}

func toPercent(v1, v2 float64) float64 {
	if v2 == 0 || v1 == 0 {
		return 0
	}
	return ((v2 - v1) / v2) * 100
}

func coordsFromString(s string) [][2]float64 {
	lines := strings.Split(s, "\n")
	res := make([][2]float64, 0)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		points := strings.Split(line, ",")
		lats := strings.Trim(points[1], " ")
		lons := strings.Trim(points[0], " ")
		if len(lats) == 0 && len(lons) == 0 {
			continue
		}
		lat, err := strconv.ParseFloat(lats, 10)
		if err != nil {
			panic(err)
		}
		lon, err := strconv.ParseFloat(lons, 10)
		if err != nil {
			panic(err)
		}
		res = append(res, [2]float64{lat, lon})
	}
	return res
}

var coordinates = `
-72.2822266, 42.9325219
-72.2821638, 42.9314098
-72.2821208, 42.9304672
-72.2820779, 42.9296816
-72.2820779, 42.9289589
-72.2820779, 42.9282991
-72.2819920, 42.9275135
-72.2819920, 42.9267907
-72.2819491, 42.9261937
-72.2818203, 42.9255338
-72.2818203, 42.9248424
-72.2818203, 42.9240882
-72.2818203, 42.9234912
-72.2818203, 42.9229255
-72.2818203, 42.9221713
-72.2818203, 42.9213856
-72.2816915, 42.9205685
-72.2815627, 42.9197829
-72.2814338, 42.9187772
-72.2813909, 42.9180229
-72.2811762, 42.9173315
-72.2810904, 42.9166715
-72.2810474, 42.9160115
-72.2810003, 42.9155456
-72.2809573, 42.9149484
-72.2809144, 42.9142569
-72.2807427, 42.9131883
-72.2806568, 42.9125283
-72.2806568, 42.9115225
-72.2806568, 42.9105481
-72.2806568, 42.9095737
-72.2806568, 42.9085364
-72.2805709, 42.9075934
-72.2804787, 42.9069159
-72.2803928, 42.9061615
-72.2802640, 42.9053442
-72.2800874, 42.9041813
-72.2797868, 42.9032068
-72.2797868, 42.9022323
-72.2797868, 42.9012892
-72.2796580, 42.9001889
-72.2795292, 42.8990572
-72.2793575, 42.8980511
-72.2791428, 42.8970137
-72.2789710, 42.8960391
-72.2787993, 42.8949701
-72.2785846, 42.8938068
-72.2784558, 42.8925492
-72.2781553, 42.8913544
-72.2778326, 42.8905102
-72.2780472, 42.8893782
-72.2779614, 42.8881834
-72.2776608, 42.8868313
-72.2774461, 42.8857622
-72.2771026, 42.8844415
-72.2768450, 42.8835610
-72.2766303, 42.8824919
-72.2764126, 42.8815181
-72.2761980, 42.8806376
-72.2760262, 42.8795370
-72.2758545, 42.8785621
-72.2756827, 42.8775243
-72.2752534, 42.8763921
-72.2749957, 42.8756059
-72.2747557, 42.8741514
-72.2745840, 42.8724531
-72.2745840, 42.8715410
-72.2744981, 42.8703143
-72.2744981, 42.8692135
-72.2744981, 42.8678295
`
