package prime

import "testing"

func TestPrime(t *testing.T) {
	numbers := []struct {
		number  float64
		isPrime bool
	}{
		{number: 10.0, isPrime: false},
		{number: 11.0, isPrime: true},
		{number: 10.1, isPrime: false},
		{number: 3.0, isPrime: true},
		{number: -1.0, isPrime: false},
		{number: 1.0, isPrime: false},
	}

	for i := 0; i < len(numbers); i += 1 {
		testcase := numbers[i]
		if checkPrime(testcase.number) != testcase.isPrime {
			t.Fatalf("Expected %v to be %v, received: %v", testcase.number, testcase.isPrime, checkPrime(testcase.number))
		}
	}

}
