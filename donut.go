package main

import (
	"fmt"
	"math"
	"time"
)

const (
    screenWidth  = 80
    screenHeight = 22
    thetaSpacing = 0.07
    phiSpacing   = 0.02
)

// ANSI colors
var colors = []string{
    "\033[31m", // red
    "\033[33m", // yellow
    "\033[32m", // green
    "\033[36m", // cyan
    "\033[34m", // blue
    "\033[35m", // purple
}

func clear() {
    fmt.Print("\033[H\033[2J")
}

func main() {
    fmt.Print("\033[?25l") 
    defer fmt.Print("\033[?25h")

    zBuffer := make([]float64, screenWidth*screenHeight)
    output := make([]byte, screenWidth*screenHeight)
    luminanceChars := []byte(".,-~:;=!*#$@")
    
    colorIndex := 0
    colorTimer := 0.0
    
    A, B := 0.0, 0.0
    
    for {
        for i := range zBuffer {
            zBuffer[i] = 0
            output[i] = ' '
        }
        
        sinA, cosA := math.Sin(A), math.Cos(A)
        sinB, cosB := math.Sin(B), math.Cos(B)

        for theta := 0.0; theta < 2*math.Pi; theta += thetaSpacing {
            for phi := 0.0; phi < 2*math.Pi; phi += phiSpacing {
                sinPhi, cosPhi := math.Sin(phi), math.Cos(phi)
                sinTheta, cosTheta := math.Sin(theta), math.Cos(theta)
                
                h := cosTheta + 2
                D := 1 / (sinPhi*h*sinA + sinTheta*cosA + 5)
                t := sinPhi*h*cosA - sinTheta*sinA
                
                x := int(40 + 30*D*(cosPhi*h*cosB - t*sinB))
                y := int(12 + 15*D*(cosPhi*h*sinB + t*cosB))
                
                N := int(8 * ((sinTheta*sinA - sinPhi*cosTheta*cosA) * cosB - 
                     sinPhi*cosTheta*sinA - sinTheta*cosA - cosPhi*cosTheta*sinB))
                
                pos := x + screenWidth*y
                if y > 0 && y < screenHeight && x > 0 && x < screenWidth {
                    if D > zBuffer[pos] {
                        zBuffer[pos] = D
                        if N > 0 {
                            output[pos] = luminanceChars[N%12]
                        } else {
                            output[pos] = luminanceChars[0]
                        }
                    }
                }
            }
        }
        
        fmt.Print("\033[H")
        
        colorTimer += 0.1
        if colorTimer >= 1.0 {
            colorTimer = 0.0
            colorIndex = (colorIndex + 1) % len(colors)
        }
        
        
        currentColor := colors[colorIndex]
        for i := 0; i < screenHeight; i++ {
            fmt.Print(currentColor) 
            for j := 0; j < screenWidth; j++ {
                fmt.Printf("%c", output[i*screenWidth+j])
            }
            fmt.Println("\033[0m")
        }
        
        A += 0.04
        B += 0.02
        
        time.Sleep(20 * time.Millisecond)
    }
}
