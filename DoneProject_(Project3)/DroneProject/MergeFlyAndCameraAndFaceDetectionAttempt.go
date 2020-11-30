/*
Seems to start the camera after taking off
set up mplayer, start mplayer and takeoff on ConnectedEvent
This also seems to work.
Will we need to use ParseFlightData ?
*/

package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"os/exec"
	"sync"
	"time"
	//
)

func main() {
	var mutex = &sync.Mutex{}
	var batteryLevel int8

	var inFlight = 0
	drone := tello.NewDriver("8890")
	work := func() {
		mplayer := exec.Command("mplayer", "-fps", "25", "-")
		mplayerIn, _ := mplayer.StdinPipe()
		if err := mplayer.Start(); err != nil {
			fmt.Println(err)
			return
		}

		drone.On(tello.FlightDataEvent, func(data interface{}) {
			// TODO: protect flight data from race condition
			mutex.Lock()
			flightData := data.(*tello.FlightData)
			mutex.Unlock()
			//checkBattery(drone, d)
			checkBattery := func() int8 {
				return flightData.BatteryPercentage
			}
			batteryLevel = checkBattery()
		})

		drone.On(tello.ConnectedEvent, func(data interface{}) {

			fmt.Println("Connected")
			var connectedWaitGroup sync.WaitGroup
			connectedWaitGroup.Add(2)
			go func() {
				drone.StartVideo()
				drone.SetVideoEncoderRate(4)
				connectedWaitGroup.Done()
			}()

			go func() {
				drone.TakeOff()
				time.Sleep(time.Second * 3)
				inFlight += 1
				connectedWaitGroup.Done()
			}()

			//turn in place
			go func() {
				var count = 0
				for count < 10 { //may need to eventually make this a for-ever loop and remove count
					drone.Clockwise(30)
					time.Sleep((time.Second * 2))
					fmt.Println("batteryLevel: ", batteryLevel)
					if inFlight == 1 && batteryLevel < 20 {
						//mutex.Lock()
						fmt.Println("Battery low, please charge.")
						inFlight = 0
						drone.Land()
						time.Sleep(time.Second * 2)
						break

					}
					count += 1
					fmt.Println("count is ", count)
				}
				inFlight = 0
				drone.Land()
				//drone.Halt()  //this throws, "Error:  read udp 0.0.0.0:11111: use of closed network connection"
			}()

			gobot.Every(100*time.Millisecond, func() {
				drone.StartVideo()
			})

			gobot.After(1*time.Second, func() { //I don't seem to be using this.
				/*
					drone.Left(10)
					time.Sleep(time.Second*3)
					drone.Right(10)
					time.Sleep(time.Second*3)
				*/

				//drone.Land()
				//inFlight = 0
				time.Sleep(time.Second * 3)

			})
			connectedWaitGroup.Wait()
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := mplayerIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})

	}
	robot := gobot.NewRobot("tello",

		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)
	robot.Start()
}
