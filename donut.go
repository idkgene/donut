/*
3D Donut implements a mathematical torus model to create an animated ASCII art.

The visualization occurs through projecting 3D coordinates onto a 2D place with lighting calculations.

Visual Representation of the Torus:
    
    Top View:           Side View:
      ******             ---
    **      **         /    \
   *          *       |      |
   *          *       |      |
    **      **         \    /
      ******             ---

Mathematical Model:
- Torus is generated using R1 (major radius) and R2 (minor radius)
- Each point P(θ,φ) = (R1 + R2·cos(φ))(cos(θ), sin(θ), R2·sin(φ))

Quick Implementation Overview:
- Uses parametric equations to generate points on a torus surface
- Applies rotation matrices for animation
- Projects 3D points onto 2D screen space
- Implements basic illumination calculations
- Depth management for correct rendering
- Mapping depth and lighting values to ASCII characters
*/

package main

import (
	"fmt"
	"math"
	"time"
)

// Screen and sampling constant for our visualization
const (
    screenWidth  = 40 // Width of the output canvas in characters
    screenHeight = 20 // Height of the output canvas in characters
    thetaSpacing = 0.07 // Angular step for the main torus circle (affects detail level)
    phiSpacing   = 0.02 // Angular step for the torus tube (affects smoothness)
)

// ANSI color sequences for kind of rainbow effect
// Each color is represented by its escape sequence
var colors = []string{
    "\033[31m", // Red
    "\033[33m", // Yellow
    "\033[32m", // Green
    "\033[36m", // Cyan
    "\033[34m", // Blue
    "\033[35m", // Magenta
}


func main() {
    // Initialize rendering buffers
    zBuffer := make([]float64, screenWidth*screenHeight) // Depth buffer for 3D projection
    output := make([]byte, screenWidth*screenHeight)     // Character buffer for ASCII output
    
    // Luminance mapping characters from darkest to brightest
    // Provides visual depth through ASCII character density
    luminanceChars := []byte(".,-~:;=!*#$@")
    
    // Rotation angles for 3D transformation
    // A: rotation around X-axis
    // B: rotation around Z-axis
    A, B := 0.0, 0.0

    colorIndex := 0
    
    for {
        // Clear buffers for new frame
        // Prevents ghosting and ensures cleaner rendering
        for i := range output {
            output[i] = ' '
            zBuffer[i] = 0
        }
        
        // Pre-calculate trigonometric values for optimization
        // Reduces redundant calculations in the rendering loop
        sinA, cosA := math.Sin(A), math.Cos(A)
        sinB, cosB := math.Sin(B), math.Cos(B)

        // Iterates through all points on the torus surface
        for theta := 0.0; theta < 2*math.Pi; theta += thetaSpacing {
            for phi := 0.0; phi < 2*math.Pi; phi += phiSpacing {
                // Calculate 3D coordinates on torus surface
                // Using parametric equations for torus generation
                sinPhi, cosPhi := math.Sin(phi), math.Cos(phi)
                sinTheta, cosTheta := math.Sin(theta), math.Cos(theta)
                
                // Calculate donut surface point
                h := cosTheta + 2 // Distance from center to torus tube
                
                // Calculate depth (D) and transformation (t) values
                // D: Used for z-buffering and perspective
                                // t: Used for rotation transformation
                                D := 1 / (sinPhi*h*sinA + sinTheta*cosA + 5) // 5 is the viewing distance
                                t := sinPhi*h*cosA - sinTheta*sinA
                
                                // Project 3D coordinates to 2D screen space
                                // Applies perspective division and screen space transformation
                                x := int(screenWidth/2 + 15*D*(cosPhi*h*cosB - t*sinB))
                                y := int(screenHeight/2 + 7*D*(cosPhi*h*sinB + t*cosB))
                
                                // Calculate surface normal for lighting
                                // N determines the luminance value for each point
                                N := int(8 * ((sinTheta*sinA - sinPhi*cosTheta*cosA) * cosB - 
                                     sinPhi*cosTheta*sinA - sinTheta*cosA - cosPhi*cosTheta*sinB))
                
                                // Screen space boundary check and z-buffer comparison
                                pos := x + screenWidth*y
                                if y >= 0 && y < screenHeight && x >= 0 && x < screenWidth {
                                    if D > zBuffer[pos] { // Z-buffer check for depth ordering
                                        zBuffer[pos] = D
                                        // Map normal to ASCII character based on luminance
                                        if N > 0 {
                                            output[pos] = luminanceChars[N%12]
                                        } else {
                                            output[pos] = luminanceChars[0]
                                        }
                                    }
                                }
                            }
                        }
                
                        // Render frame to terminal
                        fmt.Print("\033[H") // Reset cursor to home position
                        
                        // Output rendered frame with color
                        currentColor := colors[colorIndex]
                        for i := 0; i < screenHeight; i++ {
                            fmt.Print(currentColor)
                            for j := 0; j < screenWidth; j++ {
                                fmt.Printf("%c", output[i*screenWidth+j])
                            }
                            fmt.Println("\033[0m") // Reset color at end of line
                        }
                
                        // Update rotation angles for next frame
                        // Controls the rotation speed and direction
                        A += 0.07 // X-axis rotation increment
                        B += 0.03 // Z-axis rotation increment
                        
                        // 50ms delay provides animation at around ~20 FPS
                        time.Sleep(50 * time.Millisecond)
                    }
                }
