package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type CanvasCreateDrawing struct {
	RectangleAt [2]int
	Width       int
	Height      int
	Outline     string
	Fill        string
}

type CanvasCreateRequest []CanvasCreateDrawing

func (canvas CanvasCreateRequest) Validate() error {
	for _, canva := range canvas {
		if canva.Fill == "none" || strings.TrimSpace(canva.Fill) == "" {
			canva.Fill = ` `
		}
		if canva.Outline == "none" || strings.TrimSpace(canva.Outline) == "" {
			canva.Outline = canva.Fill
		}

		if canva.RectangleAt[0] < 0 || canva.RectangleAt[1] < 0 {
			return errors.New("Position of rectangle cannot be lower than 0")
		}

		if canva.RectangleAt[0]+canva.Width > 100 || canva.RectangleAt[1]+canva.Height > 100 {
			return errors.New(fmt.Sprintf("The Canvas size is %vx%v", 100, 100))
		}

		if !isASCII(canva.Fill) || !isASCII(canva.Outline) {
			return errors.New("We only accept ASCII characters :)")
		}

		if canva.Width <= 0 || canva.Height <= 0 {
			return errors.New("Width and Height must be greater than 0")
		}
	}
	return nil
}

func CanvasCreate(req *CanvasCreateRequest) (string,error) {
	err := req.Validate()
	if err != nil {
		return " ", err
	}

	maxHeight, maxWidth := findCanvasDimension(req)
	canvas := func() (elts [][]string) {
		for j := 0; j < maxHeight; j++ {
			elts = append(elts, func(repeated []string, n int) (result []string) {
				for i := 0; i < n; i++ {
					result = append(result, repeated...)
				}
				return result
			}([]string{" "}, maxWidth))
		}
		return
	}()

	for _, rectangle := range *req {
		addRectangle(rectangle, canvas)
	}

	drawing := strings.Join(func() (elts []string) {
		for _, canva := range canvas {
			elts = append(elts, strings.Join(canva, ""))
		}
		return
	}(), "\n")

	return drawing,nil
}

func findCanvasDimension(rectangles *CanvasCreateRequest) (int, int) {
	maxHeight, maxWidth := 0, 0
	for _, rectangle := range *rectangles {
		columnTopLeft := rectangle.RectangleAt[0]
		rowTopLeft := rectangle.RectangleAt[1]
		bottomRight := []int{rowTopLeft + rectangle.Height - 1, columnTopLeft + rectangle.Width - 1}
		maxHeight = func() (m int) {
			for i, e := range []int{bottomRight[0], maxHeight} {
				if i == 0 || e > m {
					m = e
				}
			}
			return
		}()
		maxWidth = func() (m int) {
			for i, e := range []int{bottomRight[1], maxWidth} {
				if i == 0 || e > m {
					m = e
				}
			}
			return
		}()
	}
	return maxHeight + 1, maxWidth + 1
}

func addRectangle(rectangle CanvasCreateDrawing, canvas [][]string) {

	if rectangle.Fill == "none" || strings.TrimSpace(rectangle.Fill) == "" {
		rectangle.Fill = ` `
	}
	if rectangle.Outline == "none" || strings.TrimSpace(rectangle.Outline) == "" {
		rectangle.Outline = rectangle.Fill
	}
	columnTopLeft := rectangle.RectangleAt[0]
	rowTopLeft := rectangle.RectangleAt[1]
	bottomRight := []int{int(rowTopLeft+rectangle.Height) - 1, int(columnTopLeft+rectangle.Width) - 1}
	for row := rowTopLeft; row < bottomRight[0]+1; row++ {
		for col := columnTopLeft; col < bottomRight[1]+1; col++ {
			if func() int {
				for i, v := range []interface{}{rowTopLeft, bottomRight[0]} {
					if v == row {
						return i
					}
				}
				return -1
			}() != -1 || func() int {
				for i, v := range []interface{}{columnTopLeft, bottomRight[1]} {
					if v == col {
						return i
					}
				}
				return -1
			}() != -1 {
				canvas[row][col] = rectangle.Outline
			} else {
				canvas[row][col] = rectangle.Fill
			}
		}
	}
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
