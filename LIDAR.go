package main

import (
    "fmt"
    "gobot.io/x/gobot"
    "gobot.io/x/gobot/drivers/aio"
    "gobot.io/x/gobot/drivers/i2c"
    g "gobot.io/x/gobot/platforms/dexter/gopigo3"
    "gobot.io/x/gobot/platforms/raspi"
    "time"
)

func robotMainLoop(piProcessor *raspi.Adaptor, gopigo3 *g.Driver, lidarSensor *i2c.LIDARLiteDriver) {

        count := 0		// this will be used to find the average
        LIDARstart := 0
		
/*		
		lastError := 0
		lastTime := time.Now()
		const kP := 0.1	//Proportionality Constant
		const kD := 1.1	//Derivative Proportionality Constant
*/	

        err := lidarSensor.Start()
        if err != nil {
                fmt.Println("error starting lidarSensor")
        }
        for{							// take in 50 readings to get an average distance
                lidarReading, err := lidarSensor.Distance()	// this will then be used to let us know when we are
                if err != nil{					// out of range of the box
                        fmt.Println("Error reading LIDAR sensor %+v", err)
                }
                if count != 50{
                        LIDARstart = LIDARstart + lidarReading
                        count = count + 1
                }
                if count == 50{
                        break
                }
        }
        LIDARstart = (LIDARstart/count)
        fmt.Println(LIDARstart)
        internalCount := 0					// an internal counter which is equal to the amount of sides
        for{							// main run loop
                count = 0
                if internalCount == 2{				// if its 2 we have checked both sides, therefore stop, and loop forever
                        gopigo3.SetMotorDps(g.MOTOR_LEFT, 0)
                        gopigo3.SetMotorDps(g.MOTOR_RIGHT, 0)
                        for{
                        }
                }
                fmt.Println("finding edge")						// finds the edge of the box
                for{ //go backwards to find edge					// because we start at some point on the box we need a 
                        count = count + 1						// start point. To do so we have the robot move backwards
                        lidarReading, err := lidarSensor.Distance()			// untill it takes in a reading that is greater than 50% wrong
                        if err != nil{							// 50% < x < 150%. This will allow for some change in the sensor
                                fmt.Println("Error reading LIDAR sensor %+v", err)	// as it was not perfect
                        }
						
/*						
						//Get error for feedback control (will copy into other places where needed when this functions)
						currentError = LIDARstart - lidarReading
						currentTime = time.Now()
						output := float64(kP*currentError) + float64(kD*((lastError - currentError)) / float64(lastTime - currentTime))
						lastError = currentError
						lastTime = currentTime
						fmt.Println("Backup feedback output: ", output)
*/

                        gopigo3.SetMotorDps(g.MOTOR_LEFT, 60)
                        gopigo3.SetMotorDps(g.MOTOR_RIGHT, 60)
                        temp := float64(LIDARstart)/float64(lidarReading)

                        if temp < .5 || temp > 1.5{					// stopping condition is found by taking out average distance
                                gopigo3.SetMotorDps(g.MOTOR_LEFT, 0)			// divided by our current distance which creates a %
                                gopigo3.SetMotorDps(g.MOTOR_RIGHT, 0)			// if out of bounds, stop
                                break
                        }
                }
                fmt.Println("edge found")						// the edge of the box has been found
                fmt.Println("starting timer")
                fmt.Println("measuring distance")
                //start timer
                start := time.Now()							// the timer has started
                count = 0
                for{ //move forward
                        count = count + 1						// this counter is in the code so that the robot moves forward
                        lidarReading, err := lidarSensor.Distance()			// some min. distance before checking if out of bounds.
                        if err != nil{							// the code will process so fast that the robot cannot move forward 
                                fmt.Println("Error reading LIDAR sensor %+v", err)	// again before the next reading is taken in. Because the robot will
                        }								// still be at the edge of the box, the LIDAR will read no box found
                        gopigo3.SetMotorDps(g.MOTOR_LEFT,-60)				// which will make it think it has hit the opposite edge and turn.
                        gopigo3.SetMotorDps(g.MOTOR_RIGHT,-60)				// therefore we made the robot move for a min. of 50 iterations
                        temp := float64(LIDARstart)/float64(lidarReading)
                        if count > 50 {							// if count == 50, then check for stopping conditions
                                if temp < .4 || temp > 1.5{
                                        gopigo3.SetMotorDps(g.MOTOR_LEFT, 0)
                                        gopigo3.SetMotorDps(g.MOTOR_RIGHT, 0)
                                        break
                                }
                        }
                }// stop move forward loop
                duration := float64(time.Since(start))					// get the total time
                total := float64(3.145 * duration)					// ((DPS/360)*circumpherence)*time = distance in cm
                total = total *1000
                fmt.Println("face one is ", total)
                fmt.Println("Calculating distance")
                //stop timer
                // maths and print length
                fmt.Println("Finding new side")						// finding a new side
                count = 0
                for{ // turn loop
                        gopigo3.SetMotorDps(g.MOTOR_LEFT, -30)				// turn for some amount of iterations
                        gopigo3.SetMotorDps(g.MOTOR_RIGHT, -60)				// 40,000 times is about 90 degrees
                        count = count+1							// if i knew how to calculate an arch i would
                        if count == 40000{
                                break
                        }
                }
                gopigo3.SetMotorDps(g.MOTOR_LEFT,0)					// after turning stop
                gopigo3.SetMotorDps(g.MOTOR_RIGHT,0)
                count = 0
                LIDARstart = 0
                for{									// find the new average, we are no longer looking at a box
                        lidarReading, err:= lidarSensor.Distance()			// to find a box we need to get an average of whatever the LIDAR
                        if err != nil{							// is viewing. we will then check for when something gets too close.
                                fmt.Println("Error reading lidar sensor %+v", err)	// that object will be the edge of the box
                        }
                        if count != 50{
                                LIDARstart = LIDARstart + lidarReading
                                count = count + 1
                        }
                        if count == 50{
                        break
                        }
                }
                LIDARstart = (LIDARstart/count)
                count = 0
                for{									// move forward untill edge is found
                        lidarReading, err := lidarSensor.Distance()			// once edge is found, keep moving until count is > any number
                        if err != nil{							// this way we are absolutly certain we are looking at the box
                                fmt.Println("Error reading LIDAR sensor %+v", err)
                        }
                        gopigo3.SetMotorDps(g.MOTOR_LEFT, -60)
                        gopigo3.SetMotorDps(g.MOTOR_RIGHT, -60)
                        temp := float64(LIDARstart)/float64(lidarReading)
                        if temp < .5 || temp > 1.5{
                                count = count + 1
                                if count >80 {
                                        gopigo3.SetMotorDps(g.MOTOR_LEFT, 0)
                                        gopigo3.SetMotorDps(g.MOTOR_RIGHT, 0)
                                        break
                                }
                        }
                }
                count = 0
                LIDARstart = 0
                for {										// take in a new average for the new side
                        lidarReading, err := lidarSensor.Distance()
                        if err != nil{
                                fmt.Println("Error reading LIDAR sensor %+v", err)
                        }
                        if count != 50{
                                LIDARstart = LIDARstart + lidarReading
                                count = count +1
                        }
                        if count == 50{
                                break
                        }
                }
                LIDARstart = (LIDARstart/count)
                // find edge again
                // move forward a little
                // find new average
                // go back to top
                gopigo3.SetMotorDps(g.MOTOR_LEFT,0)
                gopigo3.SetMotorDps(g.MOTOR_RIGHT,0)
                internalCount = internalCount +1						// account for 1 side being finished, and repeat
        }
}

func main() {
    raspberryPi := raspi.NewAdaptor()
    gopigo3 := g.NewDriver(raspberryPi)
    lidarSensor := i2c.NewLIDARLiteDriver(raspberryPi)
    lightSensor := aio.NewGroveLightSensorDriver(gopigo3, "AD_2_1")
    workerThread := func() {
        robotMainLoop(raspberryPi, gopigo3, lidarSensor)
    }
    robot := gobot.NewRobot("Gopigo Pi4 Bot",
        []gobot.Connection{raspberryPi},
        []gobot.Device{gopigo3, lidarSensor, lightSensor},
        workerThread)



    robot.Start()

}
