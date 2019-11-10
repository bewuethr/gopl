package ch09ex01

import (
	"fmt"
	"testing"
)

func TestBank(t *testing.T) {
	done := make(chan struct{})

	// Alice
	go func() {
		Deposit(200)
		fmt.Println("=", Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		Deposit(100)
		done <- struct{}{}
	}()

	// Wait for both transactions.
	<-done
	<-done

	if got, want := Balance(), 300; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}

	if got, want := Withdraw(300), true; got != want {
		t.Errorf("Success = %v, want %v", got, want)
	}

	if got, want := Withdraw(1), false; got != want {
		t.Errorf("Success = %v, want %v", got, want)
	}

	if got, want := Balance(), 0; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}
