package main

import (
	"fmt"
	"time"
)

/*
Assignment

You are provided with an array of five people.
Each person has the following details:
height, weight and bmi
The Person struct also has a function that calculates the bmi.

In the main function, we want to show the average bmi. The problem is,
calculating each person's bmi takes 1.2 seconds.
Calculating in order will take 6 seconds.
Use goroutines to bring this time down to 1.2 seconds

I have added a linear function that takes 6 seconds.
Edit the 'getBMIAvgConcurrent' function that takes 1.2 seconds (using goroutines)
*/

type Person struct {
	Height float32
	Weight float32
	BMI    float32
}

func (p *Person) CalculateBMI() {
	// calculate bmi
	// weight (kg) / [height (m)]2
	p.BMI = p.Weight / (p.Height * p.Height)
	// sleep for 1200 ms = 1.2 seconds
	time.Sleep(time.Millisecond * 1200)
}

func main() {
	var avgBMI float32
	var people = []Person{
		{Height: 1.5, Weight: 50.7},
		{Height: 1.8, Weight: 70.5},
		{Height: 1.2, Weight: 40.0},
		{Height: 1.6, Weight: 100.4},
		{Height: 1.2, Weight: 99.45},
	}

	// linear
	avgBMI = getBMIAvgLinear(people)
	fmt.Printf("avg BMI : %.2f\n\n", avgBMI)

	// concurrent
	avgBMI = getBMIAvgConcurrent(people)
	fmt.Printf("avg BMI : %.2f\n", avgBMI)
}

// linear function
func getBMIAvgLinear(people []Person) (avgBMI float32) {
	start := time.Now()
	defer func() {
		fmt.Printf("linear took : %v\n", time.Since(start))
	}()

	var totalBMI float32
	for _, p := range people {
		p.CalculateBMI()
		totalBMI += p.BMI
	}

	// this returns the average BMI
	return totalBMI / float32(len(people))
}

func getBMIAvgConcurrent(people []Person) (avgBMI float32) {
	start := time.Now()
	defer func() {
		fmt.Printf("concurrent took : %v\n", time.Since(start))
	}()

	size := len(people)
	workers := size

	jobs := make(chan Person, size)
	results := make(chan float32, size)
	var totalBMI float32
	for _, p := range people {
		jobs <- p
	}

	for i := 0; i < workers; i++ {
		go worker(jobs, results)
	}

	for i := 0; i < size; i++ {
		BMI := <-results
		totalBMI += BMI
	}

	// this returns the average BMI
	return totalBMI / float32(len(people))
}

func worker(jobs <-chan Person, results chan<- float32) {
	for j := range jobs {
		j.CalculateBMI()
		results <- j.BMI
	}

}
