package main

import (
	"fmt"
)

func main() {
	PrintMessage()
	var x int = 5
	var y int = 10
	var sumStr = fmt.Sprintf("%v", x+y)
	fmt.Println("Sum of x and y as string:", string(sumStr))

	var a int = 3
	var b *int = &a
	fmt.Println("Value of b:", *b)
	fmt.Println("Value of a:", a)
	*b = 5
	fmt.Println("Value of b:", *b)
	fmt.Println("Value of a:", a)
}

// PrintMessage prints a message to the console.
func PrintMessage() {
	fmt.Println("Hello, this is a test message from the testing package!")
}
