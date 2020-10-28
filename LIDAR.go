//NEW: Added 10/26/20
//	calculates and displays output of P.D. feedback control calculation (OUTPUT IS NOT YET USED)
//	Do we need to reset lastTime back to 0.0 when we start measuring a new edge?  (May be a stupid question.)
//CTRL F EDIT: for edited parts and parts we may want to double-check
//SOME VARIABLES STILL NEED VALUES ASSIGNED (WILL FALIL TO COMPILE)
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

        count := 0
        LIDARstart := 0
//NEW: 
		lastError := 0.0
		lastTime := 0.0
		const kP := 0.0	//Proportionality Constant
		const kD := 1.1	//Derivative Proportionality Constant

        err := lidarSensor.Start()
        if err != nil {
                fmt.Println("error starting lidarSensor")
        }

		//EDIT: ARE THE NEXT TWO LOOPS INTENDED TO HELP CALCULATE OUR SPEED?
		//NEW: Don't need this in a loop (Do we need this at all, or is it just to give a short pause at the begining?)
        for { //loop forever
                lidarReading, err := lidarSensor.Distance()
                if err != nil {
						fmt.Println("Error reading lidar sensor %+v", err)
                }
                message := fmt.Sprintf("Lidar Reading: %d", lidarReading)


                fmt.Println(message)
                time.Sleep(time.Second * 3)
                break
        }
        for{     
				lidarReading, err := lidarSensor.Distance()
                if err != nil{
						fmt.Println("Error reading LIDAR sensor %+v", err)
                }
                if count != 50{
                        LIDARstart = LIDARstart + lidarReading
                        count = count + 1
                }else{
				//NEW: Changed to use else instead of a new if statement
                //if count == 50{
						break
                }
        }
        LIDARstart = (LIDARstart/count)
        fmt.Println(LIDARstart)

//			EDIT: IS THIS JUST HERE SO THAT WE DON'T HIT THE ABOVE LIDARstart CALCULATION SECTION AGAIN?
        for{
				
					//Start next to box, back up until edge is passed
                for{
                        lidarReading, err := lidarSensor.Distance()
                        if err != nil{
                                fmt.Println("Error reading LIDAR sensor %+v", err)
                        }
						
						//NEW: Get error for feedback control
						currentError := LIDARstart - lidarReading
						currentTime := time.Now()
						output := kP*currentError + kD*((lastError - currentError) / (lastTime - currentTime))
						lastError = currentError
						lastTime = currentTime
						fmt.Println("Backup feedback output: ", output)

                        gopigo3.SetMotorDps(g.MOTOR_LEFT, -60)
                        gopigo3.SetMotorDps(g.MOTOR_RIGHT, -60)
                        temp := float64(LIDARstart)/float64(lidarReading)
                        fmt.Println("Lidar reading ", lidarReading)


                        if temp < .8 || temp > 1.2{
                                gopigo3.SetMotorDps(g.MOTOR_LEFT, 0)
                                gopigo3.SetMotorDps(g.MOTOR_RIGHT, 0)
                                break
                        }
                }
				
					//EDIT:
					//Move forward until edge of box is reached.  May want to reduce the amount we move the wheels forward.
                for{
                        lidarReading, err := lidarSensor.Distance()
                        if err != nil{
                                fmt.Println("Error reading LIDAR sensor %+v", err)
                        }
						
						//NEW: Get error for feedback control
						currentError := LIDARstart - lidarReading
						currentTime := time.Now()
						output := kP*currentError + kD*((lastError - currentError) / (lastTime - currentTime))
						lastError = currentError
						lastTime = currentTime
						fmt.Println("Advance to front edge feedback output: ", output)

                        gopigo3.SetMotorDps(g.MOTOR_LEFT, 60)
                        gopigo3.SetMotorDps(g.MOTOR_RIGHT, 60)
                        temp := float64(LIDARstart)/float64(lidarReading)
                        fmt.Println("Moving to star, LIDAR reading is  ", lidarReading)
                        if temp > .8 || temp < 1.2{
                                break
                        }
                }
				
					//EDIT:
					//Start following the edge of the box forward.  This is the loop to put the measurement code in.
					//May need to remove the count at some point (It may put us too far forward.
				NEW: I changed "count := 0" to "count = 0" to reset the variable, without recreating it.
				count = 0
                for{
                        count = count + 1
                        lidarReading, err := lidarSensor.Distance()
                        if err != nil{
                                fmt.Println("Error reading LIDAR sensor %+v", err)
                        }
						
						//NEW: Get error for feedback control
						currentError := LIDARstart - lidarReading
						currentTime := time.Now()
						output := kP*currentError + kD*((lastError - currentError) / (lastTime - currentTime))
						lastError = currentError
						lastTime = currentTime
						fmt.Println("Forward feedback output: ", output)
						
                        gopigo3.SetMotorDps(g.MOTOR_LEFT,60)
                        gopigo3.SetMotorDps(g.MOTOR_RIGHT,60)
                        temp := float64(LIDARstart)/float64(lidarReading)
                        fmt.Println("Forward reading  ", lidarReading)
                        if count > 20{
                                if temp < .8 || temp > 1.2{
                                        gopigo3.SetMotorDps(g.MOTOR_LEFT, 0)
                                        gopigo3.SetMotorDps(g.MOTOR_RIGHT, 0)
                                        break
                                }
                        }														
                }
				
					//EDIT: CODE FOR TURNING AND CHECKING IF ANGLE WAS CORRECT
				fixTurnAmount := 0

				for{
						leftWheelAmount := 45 
						rightWheelAmount := 90
						
						if fixTurnAmount == 1{
							leftWheelAmount = leftWheelAmount * -1
							rightWheelAmount = rightWheelAmount * -1
							fixTurnAmount = 0
						}

							//EDIT: 
							//Loop to turn.  May need to know the diameter of the wheel and the length of the robot
							//May need to adjust tuning degrees
						for{
								lidarReading, err := lidarSensor.Distance()
								if err != nil{
										fmt.Println("Error reading LIDAR sensor %+v", err)
								}
								
								//NEW: Get error for feedback control
								currentError := LIDARstart - lidarReading
								currentTime := time.Now()
								output := kP*currentError + kD*((lastError - currentError) / (lastTime - currentTime))
								lastError = currentError
								lastTime = currentTime
								fmt.Println("Turning feedback output: ", output)
								
								temp := float64(LIDARstart)/float64(lidarReading)
								fmt.Println("Turning reading  ", lidarReading)
								
								if temp > .8 || temp < 1.2{
										gopigo3.SetMotorDps(g.MOTOR_LEFT, 0)
										gopigo3.SetMotorDps(g.MOTOR_RIGHT, 0)
										break
								}
								gopigo3.SetMotorDps(g.MOTOR_LEFT,leftWheelAmount)				
								gopigo3.SetMotorDps(g.MOTOR_RIGHT, rightWheelAmount)
						}
						
						
						
							//DRIVE FORWARD AND CHECK TO SEE IF ROBOT IS (MOSTLY) PARALLEL TO THE BOX
							//DO I NEED TO PASS lidarSensor.Distance(), OR CAN CHECK PARALLEL STILL ACCESS IT?
						correctTurn := 0
						turnCount :=0
										
										
						for turnCount < 50{
							gopigo3.SetMotorDps(g.MOTOR_LEFT,60)
							gopigo3.SetMotorDps(g.MOTOR_RIGHT,60)
							turnCount = turnCount + 1
						}
						gopigo3.SetMotorDps(g.MOTOR_LEFT,0)
						gopigo3.SetMotorDps(g.MOTOR_RIGHT,0)
						
						lidarReading, err := lidarSensor.Distance()
						if err != nil{
								fmt.Println("Error reading LIDAR sensor %+v", err)
						}
						
						//NEW: Get error for feedback control
						currentError := LIDARstart - lidarReading
						currentTime := time.Now()
						output := kP*currentError + kD*((lastError - currentError) / (lastTime - currentTime))
						lastError = currentError
						lastTime = currentTime
						fmt.Println("Parallel check feedback output: ", output)
						
						temp := float64(LIDARstart)/float64(lidarReading)
						fmt.Println("Turn check reading  ", lidarReading)
						
						//EDIT:
						//PARALLEL TO BOX
						if temp > .8 || temp < 1.2{
								correctTurn = 1
						}else{
								correctTurn = 0
						}
						
							//RETURN TO POSITION ROBOT WAS IN BEFORE CHECKING IF PARALLEL
						for turnCount > 0{
								gopigo3.SetMotorDps(g.MOTOR_LEFT, -60)
								gopigo3.SetMotorDps(g.MOTOR_RIGHT, -60)
								turnCount = turnCount - 1
						}
						gopigo3.SetMotorDps(g.MOTOR_LEFT,0)
						gopigo3.SetMotorDps(g.MOTOR_RIGHT,0)
						
						if correctTurn == 0{
							fixTurnAmount = 1
						}else{
							break
						}
						
				}
					
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
