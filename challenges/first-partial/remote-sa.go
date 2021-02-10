package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Point struct {
	X, Y float64
}

/* used to get orientation */
type Vector struct {
	A, B float64
}

func main() {
	http.HandleFunc("/", handler) //calling our web server
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//generatePoints array
func generatePoints(s string) ([]Point, error) {

	points := []Point{}

	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	vals := strings.Split(s, ",")
	if len(vals) < 2 {
		return []Point{}, fmt.Errorf("Point [%v] was not well defined", s)
	}

	var x, y float64

	for idx, val := range vals {

		if idx%2 == 0 {
			x, _ = strconv.ParseFloat(val, 64)
		} else {
			y, _ = strconv.ParseFloat(val, 64)
			points = append(points, Point{x, y})
		}
	}
	return points, nil
}

// getArea gets the area inside from a given shape
func getArea(points []Point) float64 {
	// Your code goes here
	return 0.0
}

// getPerimeter gets the perimeter from a given array of connected points
func getPerimeter(points []Point) float64 {
	// Your code goes here
	return 0.0
}

/* construct vector */
func vector_make(po, q Point) Vector {
	var v Vector
	v.A = q.X - po.X
	v.B = q.Y - po.Y
	return v
}

/* cross product in 2D */
func vector_cross(v, u Vector) float64 {
	return (v.A * u.B) - (v.B * u.A)
}

/* counterclockwise */
func ccw(p, q, r Point) bool {
	var pq, pr Vector
	pq = vector_make(p, q)
	pr = vector_make(p, r)
	return vector_cross(pq, pr) > 0
}

/*
 *  Check for collisions (p1,q1) and (p2,q2)
 *  – (p1, q1, p2) and (p1, q1, q2) have different orientations and
 *  – (p2, q2, p1) and (p2, q2, q1) have different orientations.
 */
func areIntersecting(v, u, i, k Point) bool {
	if (ccw(v, u, i) != ccw(v, u, k)) ||
		(ccw(i, k, v) != ccw(i, k, u)) {
		return true
	}
	return false
}

// handler handles the web request and reponds it
func handler(w http.ResponseWriter, r *http.Request) {

	var vertices []Point
	for k, v := range r.URL.Query() {
		if k == "vertices" {
			points, err := generatePoints(v[0])
			if err != nil {
				fmt.Fprintf(w, fmt.Sprintf("error: %v", err))
				return
			}
			vertices = points
			break
		}
	}

	// Results gathering
	area := getArea(vertices)
	perimeter := getPerimeter(vertices)

	// Logging in the server side /*this prints in the server side*/
	log.Printf("Received vertices array: %v", vertices)

	// Response construction
	response := fmt.Sprintf("Welcome to the Remote Shapes Analyzer\n")
	response += fmt.Sprintf(" - Your figure has : [%v] vertices\n", len(vertices))
	response += fmt.Sprintf(" - Vertices        : %v\n", vertices)
	response += fmt.Sprintf(" - Perimeter       : %v\n", perimeter)
	response += fmt.Sprintf(" - Area            : %v\n", area)

	// Send response to client
	fmt.Fprintf(w, response)
}
