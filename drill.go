package main

import "fmt"

// TODO 1: Define an interface named 'PaymentProcessor' with one method: Pay(amount int) string

type Stripe struct{}

// TODO 2: Write a method on 'Stripe' named 'Pay(amount int) stringss' that returns "Paid via Stripe"

func Checkout(p PaymentProcessor, total int) {
	// TODO 3: Call p.Pay(total) and print the result with fmt.Println
}

func main() {
	s := Stripe{}QW
	Checkout(s, 50)
}
