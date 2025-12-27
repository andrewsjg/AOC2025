package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
)

type Point3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Pair struct {
	From     Point3D `json:"from"`
	To       Point3D `json:"to"`
	Distance float64 `json:"distance"`
}

type SceneData struct {
	Points []Point3D `json:"points"`
	Lines  []Pair    `json:"lines"`
}

func readInput(inputFile string) ([]Point3D, error) {

	points := []Point3D{}

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return points, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		var point Point3D

		_, err := fmt.Sscanf(line, "%f,%f,%f", &point.X, &point.Y, &point.Z)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing line: %v\n", err)
			continue
		}
		points = append(points, point)

	}

	return points, nil
}

func makeSceneData(points []Point3D, lines []Pair) SceneData {

	scene := SceneData{
		Points: points,
		Lines:  lines,
	}

	return scene
}

func visualise(scene SceneData) {

	// Work out the circuits

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		data := scene

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlTemplate))
	})

	println("Open http://localhost:8080 to see the visulaisation")
	http.ListenAndServe(":8080", nil)
}

func part1(points []Point3D, numConnections int) (SceneData, int) {
	// Initialize Union-Find structure
	uf := NewUnionFind(len(points))

	type IndexedPair struct {
		i        int
		j        int
		Distance float64
	}

	pairs := []IndexedPair{}

	// Generate all the pairs and their distances
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dx := points[i].X - points[j].X
			dy := points[i].Y - points[j].Y
			dz := points[i].Z - points[j].Z
			dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

			pair := IndexedPair{
				i:        i,
				j:        j,
				Distance: dist,
			}
			pairs = append(pairs, pair)
		}
	}

	// Sort the pairs by distance
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Distance < pairs[j].Distance
	})

	// Connect the shortest pairs
	connectedPairs := []Pair{}

	for i := 0; i < numConnections && i < len(pairs); i++ {
		pair := pairs[i]

		// Try to union the two points
		connected := uf.Union(pair.i, pair.j)

		// Store for visualization if they were connected
		if connected {
			connectedPairs = append(connectedPairs, Pair{
				From:     points[pair.i],
				To:       points[pair.j],
				Distance: pair.Distance,
			})
		}
	}

	// Get circuit sizes
	sizes := uf.GetSizes()
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	//fmt.Printf("After processing %d pairs, circuit sizes: %v\n", numConnections, sizes)

	answer := 1
	for i := 0; i < 3 && i < len(sizes); i++ {
		answer *= sizes[i]
	}

	scene := makeSceneData(points, connectedPairs)
	return scene, answer
}

func part2(points []Point3D) (SceneData, int) {
	// Initialize Union-Find structure
	uf := NewUnionFind(len(points))

	type IndexedPair struct {
		i        int
		j        int
		Distance float64
	}

	pairs := []IndexedPair{}

	// Generate all the pairs and their distances
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dx := points[i].X - points[j].X
			dy := points[i].Y - points[j].Y
			dz := points[i].Z - points[j].Z
			dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

			pair := IndexedPair{
				i:        i,
				j:        j,
				Distance: dist,
			}
			pairs = append(pairs, pair)
		}
	}

	// Sort the pairs by distance
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Distance < pairs[j].Distance
	})

	// Connect pairs until we have only 1 circuit
	connectedPairs := []Pair{}
	var lastPair IndexedPair
	numCircuits := len(points) // Start with each point in its own circuit

	for _, pair := range pairs {
		// Try to union the two points
		if uf.Union(pair.i, pair.j) {
			numCircuits-- // One less circuit after successful union
			lastPair = pair

			// Store for visualization
			connectedPairs = append(connectedPairs, Pair{
				From:     points[pair.i],
				To:       points[pair.j],
				Distance: pair.Distance,
			})

			// Check if we're done (all in one circuit)
			if numCircuits == 1 {
				break
			}
		}
	}

	// Calculate answer: product of X coordinates of the last connection
	x1 := int(points[lastPair.i].X)
	x2 := int(points[lastPair.j].X)
	answer := x1 * x2

	scene := makeSceneData(points, connectedPairs)
	return scene, answer
}

// UnionFind data structure for tracking connected components (circuits)
type UnionFind struct {
	parent map[int]int
	size   map[int]int
}

func NewUnionFind(n int) *UnionFind {
	uf := &UnionFind{
		parent: make(map[int]int),
		size:   make(map[int]int),
	}
	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.size[i] = 1
	}
	return uf
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // Path compression
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	if rootX == rootY {
		return false // Already in same set
	}

	// Union by size
	if uf.size[rootX] < uf.size[rootY] {
		uf.parent[rootX] = rootY
		uf.size[rootY] += uf.size[rootX]
	} else {
		uf.parent[rootY] = rootX
		uf.size[rootX] += uf.size[rootY]
	}
	return true
}

func (uf *UnionFind) GetSizes() []int {
	sizeMap := make(map[int]int)
	for i := range uf.parent {
		root := uf.Find(i)
		sizeMap[root] = uf.size[root]
	}

	sizes := []int{}
	for _, size := range sizeMap {
		sizes = append(sizes, size)
	}
	return sizes
}

func main() {
	// Configuration
	inputFile := "input.txt" // "testinput.txt" or "input.txt"
	numConnections := 1000
	runPart2 := true // Set to true for Part 2, false for Part 1

	input, err := readInput(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	var scene SceneData
	var answer int

	if runPart2 {
		scene, answer = part2(input)
		fmt.Printf("Part 2 - sum of last two connected junctions: %d\n\n", answer)
	} else {
		scene, answer = part1(input, numConnections)
		fmt.Printf("Part 1 - sum of three largest circuits: %d\n\n", answer)
	}

	// Visualise it for fun
	visualise(scene)
}

const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
    <style>
        body { margin: 0; }
        canvas { display: block; }
    </style>
</head>
<body>
<script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r128/three.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/three@0.128.0/examples/js/controls/OrbitControls.js"></script>
<script>
const SCALE = 100000;

const scene = new THREE.Scene();
scene.background = new THREE.Color(0x1a1a2e);

const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 1, SCALE * 10);
camera.position.set(SCALE * 1.5, SCALE * 1.5, SCALE * 1.5);

const renderer = new THREE.WebGLRenderer({ antialias: true });
renderer.setSize(window.innerWidth, window.innerHeight);
document.body.appendChild(renderer.domElement);

const controls = new THREE.OrbitControls(camera, renderer.domElement);
controls.target.set(SCALE / 2, SCALE / 2, SCALE / 2);  // Look at center of space
controls.update();

fetch('/data')
    .then(res => res.json())
    .then(data => {
        // Draw points as spheres (scaled for visibility)
        const pointGeom = new THREE.SphereGeometry(SCALE * 0.005,16, 16);
        const pointMat = new THREE.MeshBasicMaterial({ color: 0x00ff88 });
        
        data.points.forEach(p => {
            const sphere = new THREE.Mesh(pointGeom, pointMat);
            sphere.position.set(p.x, p.y, p.z);
            scene.add(sphere);
        });

        // Draw lines
        const lineMat = new THREE.LineBasicMaterial({ color: 0xff6b6b });
        
        data.lines.forEach(line => {
            const geometry = new THREE.BufferGeometry().setFromPoints([
                new THREE.Vector3(line.from.x, line.from.y, line.from.z),
                new THREE.Vector3(line.to.x, line.to.y, line.to.z)
            ]);
            scene.add(new THREE.Line(geometry, lineMat));
        });
    });

// Grid sized for 0-1000 range
const grid = new THREE.GridHelper(SCALE, 20, 0x444444, 0x222222);
grid.position.set(SCALE / 2, 0, SCALE / 2);  // Center grid in the space
scene.add(grid);

// Add axis lines for reference
const axisLen = SCALE;
const axisLines = [
    { color: 0xff0000, from: [0,0,0], to: [axisLen,0,0] },  // X - red
    { color: 0x00ff00, from: [0,0,0], to: [0,axisLen,0] },  // Y - green
    { color: 0x0000ff, from: [0,0,0], to: [0,0,axisLen] },  // Z - blue
];
axisLines.forEach(axis => {
    const mat = new THREE.LineBasicMaterial({ color: axis.color });
    const geom = new THREE.BufferGeometry().setFromPoints([
        new THREE.Vector3(...axis.from),
        new THREE.Vector3(...axis.to)
    ]);
    scene.add(new THREE.Line(geom, mat));
});

function animate() {
    requestAnimationFrame(animate);
    controls.update();
    renderer.render(scene, camera);
}
animate();

window.addEventListener('resize', () => {
    camera.aspect = window.innerWidth / window.innerHeight;
    camera.updateProjectionMatrix();
    renderer.setSize(window.innerWidth, window.innerHeight);
});
</script>
</body>
</html>`
