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

        err := lidarSensor.Start()
        if err != nil {
                fmt.Println("error starting lidarSensor")
        }

		//EDIT: ARE THE NEXT TWO LOOPS INTENDED TO HELP CALCULATE OUR SPEED?
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
                }
                if count == 50{
						break
                }
        }
        LIDARstart = (LIDARstart/count)
        fmt.Println(LIDARstart)

//			EDIT: IS THIS UST HERE SO THAT WE DON'T HIT THE ABOVE SPEED CALCULATION SECTION AGAIN?
        for{
				
					//Start next to box, back up until edge is passed
                for{
                        lidarReading, err := lidarSensor.Distance()
                        if err != nil{
                                fmt.Println("Error reading LIDAR sensor %+v", err)
                        }

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
				count := 0
                for{
                        count = count +1
                        lidarReading, err := lidarSensor.Distance()
                        if err != nil{
                                fmt.Println("Error reading LIDAR sensor %+v", err)
                        }
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
				loadout := 0
				numCases := 4	//NEED TO CHANGE IF CHANGE NUMBER OF CASES IN SWITCH STATEMENT
				fixTurnAmount := 0

				for loadout < numCases{
					
						//MAY NEED TO SPECIFY INT TYPE
						//MAY NEED TO NAME RETURN VALUE
/*
						leftWheelAmount, rightWheelAmount := func setWheels(loadout int) (returnedLeft int, returnedRight int) {
							switch{
								case loadout = 0:
									return 45, 90
								case loadout = 1:
									return 0, 0 //SOME VALUE for leftWheelAmount, SOME VALUE for rightWheelAmount
								case loadout = 0, 0:
									return //SOME VALUE for leftWheelAmount, SOME VALUE for rightWheelAmount
								case loadout = 3:
									return //SOME VALUE for leftWheelAmount, SOME VALUE for rightWheelAmount
							}
*/
						leftWheelAmount := 0 
						rightWheelAmount := 0
							switch (loadout){
								case 0:
									leftWheelAmount = 45
									rightWheelAmount =  90
								case 1:
									leftWheelAmount = 0
									rightWheelAmount = 0 //SOME VALUE for leftWheelAmount, SOME VALUE for rightWheelAmount
								case 2:
									leftWheelAmount = 0
									rightWheelAmount = 0//SOME VALUE for leftWheelAmount, SOME VALUE for rightWheelAmount
								case 3:
									leftWheelAmount = 0
									rightWheelAmount = 0 //SOME VALUE for leftWheelAmount, SOME VALUE for rightWheelAmount
							}
							
							
							
					
					
					
					
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
							//DO I NEED TO PASS lidarSensor.Distanc(), OR CAN CHECK PARALLEL STILL ACCESS IT?
						correctTurn := 0
							turnCount :=0
										
										
							for turnCount < 50{
								gopigo3.SetMotorDps(g.MOTOR_LEFT,60)
								gopigo3.SetMotorDps(g.MOTOR_RIGHT,60)
								turnCount = count + 1
							}
							gopigo3.SetMotorDps(g.MOTOR_LEFT,0)
							gopigo3.SetMotorDps(g.MOTOR_RIGHT,0)
							
							lidarReading, err := lidarSensor.Distance()
							if err != nil{
									fmt.Println("Error reading LIDAR sensor %+v", err)
							}
							temp := float64(LIDARstart)/float64(lidarReading)
							fmt.Println("Turn check reading  ", lidarReading)
							
							//EDIT:
							//PARALLEL TO BOX
							if temp > .8 || temp < 1.2{
									correctTurn = 1
							}else{
							correctTurn = 0
							}
						
							//RETURN TO POSSITION ROBOT WAS IN BEFORE CHECKING IF PARALLEL
						for turnCount < 50{
								gopigo3.SetMotorDps(g.MOTOR_LEFT,60)
								gopigo3.SetMotorDps(g.MOTOR_RIGHT,60)
								turnCount = count + 1
						}
						gopigo3.SetMotorDps(g.MOTOR_LEFT,0)
						gopigo3.SetMotorDps(g.MOTOR_RIGHT,0)
						
						if correctTurn == 0{
							fixTurnAmount =1
						}else{
							fixTurnAmount = 0
							loadout = loadout + 1
						}
						
				}
					
				

				
					//EDIT: OLD CODE TO TURN AND START OF FIRST ATTEMPT TO CHECK FOR RIGHT ORIENTATION AFTER TURNING
				/*
					//EDIT: 
					//Loop to turn.  May need to know the diameter of the wheel and the length of the robot
					//May need to adjust tuning degrees
				for{
						lidarReading, err := lidarSensor.Distance()
						if err != nil{
								fmt.Println("Error reading LIDAR sensor %+v", err)
						}
						
						temp := float64(LIDARstart)/float64(lidarReading)
						fmt.Println("Turning reading  ", lidarReading)
						
						if temp > .8 || temp < 1.2{
								gopigo3.SetMotorDps(g.MOTOR_LEFT, 0)
								gopigo3.SetMotorDps(g.MOTOR_RIGHT, 0)
								break
						}
						gopigo3.SetMotorDps(g.MOTOR_LEFT,45)				
						gopigo3.SetMotorDps(g.MOTOR_RIGHT, 90)
				}
				
				//EDIT:
						//MAY WANT A LOOP HERE THAT DRIVES THE ROBOT FORWARD FOR A SET DISTANCE (USE A COUNT) AND CHECKS TO MAKE SURE THAT THE ROBOT IS PROPERLY ALIGNED (MAKE SURE THAT THE ROBOT IS NEITHER TOO CLOSE, NOR TOO FAR AWAY AFTER THE DISTANCE IS REACHED.)  RETURN TO POSITION AT START OF COUNT, CHOOSE A DIFFERENT VALUES TO TURN THE WHEELS, AND TRY AGAIN.
						//IF WE DO USE THIS TO CORRECT THE ANGLE OF TURNING, 
							//WE WILL LIKELY NEED TO PUT THIS AND THE TURNING LOOP 
							//INTO ANOTHER FOR LOOP, REPLACE THE HARD-CODED VALUES IN 
							//THE TURNING LOOP WITH VARIABLES THAT HAVE THEIR VALUES 
							//SET TO EITHER POSITIVE OR NEGATIVE, DEPENDING ON WHETHER 
							//THE ROBOT IS GOING FORWARD, OR RETURNING TO ITS 
							//PRE-TURNING STARTING POSITION BEFORE TRYING THE NEXT 
							//SET OF POSSIBLE TURNING DEGREES (WHICH MAY BE REMOVED 
							//ONCE WE FIND THE PROPER SET, AS THEY WOULD BE USED 
							//FOR TESTING PURPOSES, SO THAT WE DON'T NEED TO STOP, 
							//EDIT, RECOMPILE, AND RUN EACH TIME WE WANT TO TEST A NEW SET OF DEGREES.) 
					//MAY NEED TO CHECK SYNTAX FOR "False"
				*/
				/*
				correctTurn := False
				for correctTurn == False{
						
						turnCount :=0
						
						
						for turnCount < 50{
							gopigo3.SetMotorDps(g.MOTOR_LEFT,60)
							gopigo3.SetMotorDps(g.MOTOR_RIGHT,60)
							turnCount = count + 1
						}
						
						lidarReading, err := lidarSensor.Distance()
                        if err != nil{
                                fmt.Println("Error reading LIDAR sensor %+v", err)
                        }
						temp := float64(LIDARstart)/float64(lidarReading)
						fmt.Println("Turn check reading  ", lidarReading)
						
						//EDIT:
						if temp > .8 || temp < 1.2{
								correctTurn = True
								break
						}
						
						
				}
				*/
//					EDIT: END OF OLD STUFF
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


