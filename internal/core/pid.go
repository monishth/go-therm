package core

import (
	"time"

	"github.com/monishth/go-therm/pkg/utils"
)

// PIDController struct
type PIDController struct {
	Kp, Ki, Kd float64 // PID coefficients

	setpoint  float64   // Desired value
	integral  float64   // Integral accumulator
	prevError float64   // Previous error
	prevTime  time.Time // Time of the previous measurement
}

// NewPIDController creates a new PID controller
func NewPIDController(Kp, Ki, Kd, setpoint float64) *PIDController {
	return &PIDController{
		Kp:       Kp,
		Ki:       Ki,
		Kd:       Kd,
		setpoint: setpoint,
		prevTime: time.Now(),
	}
}

// SetSetpoint sets the desired value for the PID controller
func (pid *PIDController) SetSetpoint(setpoint float64) {
	if setpoint != pid.setpoint {
		pid.setpoint = setpoint
		pid.integral = 0
		pid.prevError = 0
		pid.prevTime = time.Now()
	}
}

// Calculate calculates the control variable
func (pid *PIDController) Calculate(measurement float64) float64 {
	// Calculate time difference
	currentTime := time.Now()
	deltaTime := currentTime.Sub(pid.prevTime).Seconds()

	// Calculate error
	error := pid.setpoint - measurement

	// Proportional term
	proportional := pid.Kp * error

	// Integral term
	pid.integral += error * deltaTime * 0.1
	integral := utils.Min(50, utils.Max(pid.Ki*pid.integral, -50))

	// Derivative term
	derivative := 0.0
	if deltaTime > 0 {
		derivative = pid.Kd * (error - pid.prevError) / deltaTime
	}

	// Compute output
	output := proportional + integral + derivative

	// Store current error and time for next calculation
	pid.prevError = error
	pid.prevTime = currentTime

	return utils.Max(-100, utils.Min(output, 100))
}
