package kmeans

import (
	"image/color"
	"math"
	"math/rand"
	"time"
)

// This implemtation of k-means clustering for detecting dominant colors
// is based on the tutorial from https://zeevgilovitz.com/detecting-dominant-colours-in-python
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// TODO: Rename a bunch of these variables to be more concise
type Cluster struct {
	Pixels   []color.Color
	Centroid color.Color
}

func (c Cluster) AddPoint(clr color.Color) {
	c.Pixels = append(c.Pixels, clr)
}

func (c Cluster) SetNewCentroid() {
	if len(c.Pixels) == 0 {
		return
	}
	var rs, gs, bs []uint32

	for _, clr := range c.Pixels {
		r, g, b, _ := clr.RGBA()
		rs = append(rs, r)
		gs = append(gs, g)
		bs = append(bs, b)
	}
	sum := func(input []uint32) uint32 {
		sum := uint32(0)

		for i := range input {
			sum += input[i]
		}
		return sum
	}

	newR := sum(rs) / uint32(len(rs))
	newG := sum(gs) / uint32(len(gs))
	newB := sum(bs) / uint32(len(bs))
	c.Centroid = color.Color(color.RGBA{R: uint8(newR), G: uint8(newG), B: uint8(newB)})
	c.Pixels = []color.Color{}
}

type Kmeans struct {
	K             int
	MaxIterations int
	MinDistance   int

	Pixels      []color.Color
	Clusters    []Cluster
	OldClusters []color.Color
}

func NewKmeansClustering(k, maxIterations, minDistance int) (km Kmeans) {
	km.K = k
	km.MaxIterations = maxIterations
	km.MinDistance = minDistance
	km.Clusters = make([]Cluster, k)
	return
}
func (km Kmeans) Run(pixels []color.Color) []color.Color {
	// random pixels
	randomPixels := make([]color.Color, km.K)
	randInts := rand.Perm(len(pixels))
	for r := 0; r < km.K; r++ {
		randomPixels[r] = pixels[randInts[r]]
		km.Clusters[r] = Cluster{Centroid: randomPixels[r]}
	}

	iteration := 0
	ok := true
	for ok {
		ok = km.shouldExit(iteration)
		for _, cluster := range km.Clusters {
			km.OldClusters = append(km.OldClusters, cluster.Centroid)
		}
		for _, clr := range km.Pixels {
			km.assignClusters(clr)
		}

		for _, cluster := range km.Clusters {
			cluster.SetNewCentroid()
		}
		iteration++
	}

	centroids := []color.Color{}
	for _, cluster := range km.Clusters {
		centroids = append(centroids, cluster.Centroid)
	}
	return centroids
}

func (km Kmeans) assignClusters(clr color.Color) {
	shortest := math.MaxFloat64
	nearest := km.Clusters[0]
	for _, cluster := range km.Clusters {
		dist := calcDistance(cluster.Centroid, clr)
		if dist < shortest {
			shortest = dist
			nearest = cluster
		}
	}
	nearest.AddPoint(clr)
}
func (km Kmeans) shouldExit(i int) bool {
	if len(km.OldClusters) == 0 {
		return false
	}

	for idx := 0; idx < km.K; idx++ {
		dist := calcDistance(km.Clusters[idx].Centroid, km.OldClusters[idx])
		if dist < float64(km.MinDistance) {
			return true
		}
	}
	if i <= km.MaxIterations {
		return false
	}
	return true
}

func calcDistance(clrA, clrB color.Color) float64 {
	ra, ga, ba, _ := clrA.RGBA()
	rb, gb, bb, _ := clrB.RGBA()
	r := math.Pow(float64(ra-rb), 2)
	g := math.Pow(float64(ga-gb), 2)
	b := math.Pow(float64(ba-bb), 2)

	return math.Sqrt(r + g + b)
}
