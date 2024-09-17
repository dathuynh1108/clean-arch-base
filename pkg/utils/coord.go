package utils

func RectIOU(rect1, rect2 []float64) float64 {
	// Calculate the area of the int ersection rectangle
	intersectionArea := IntersecArea(rect1, rect2)

	// Calculate the area of both rectangles
	rect1Area := (rect1[2] - rect1[0]) * (rect1[3] - rect1[1])
	rect2Area := (rect2[2] - rect2[0]) * (rect2[3] - rect2[1])

	// Calculate the union area
	unionArea := rect1Area + rect2Area - intersectionArea
	if unionArea == 0 {
		return 0
	}

	// Calculate the overlap percentage
	return (intersectionArea / unionArea) * 100
}

func IntersecArea(rect1, rect2 []float64) float64 {
	// Calculate the (x, y) coordinates of the intersection rectangle
	xOverlap := max(0, min(rect1[2], rect2[2])-max(rect1[0], rect2[0]))
	yOverlap := max(0, min(rect1[3], rect2[3])-max(rect1[1], rect2[1]))

	return xOverlap * yOverlap
}

func ReactSquare(rect []float64) float64 {
	return (rect[2] - rect[0]) * (rect[3] - rect[1])
}
