package main

import "os"

func handleInput(quitChan chan interface{}, inputChan chan string){

	
	go func() {
		// defer wg.Done()
		for {
			var buf [1]byte
			_, err := os.Stdin.Read(buf[:]) // Read a single byte
			if err != nil {
				break // Exit on error
			}
			switch buf[0] {
			case 'a': // Move left
				inputChan <- "left"
			case 'd': // Move right
				inputChan <- "right"
			case 's': // Move down
				inputChan <- "down"
			case 'e': // rotate
				inputChan <- "rotc"
			case 'q': // rotate
				inputChan <- "rota"
			case 27: // Quit
				close(quitChan)
				return
			}
		}
	}()

}