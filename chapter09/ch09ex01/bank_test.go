package bank_test

import (
	"fmt"
	"testing"

	bank "github.com/bewuethr/gopl/chapter09/ch09ex01"
)

func TestBank(t *testing.T) {
	done := make(chan struct{})

	// Alice
	go func() {
		bank.Deposit(200)
		fmt.Println("=", bank.Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		bank.Deposit(100)
		done <- struct{}{}
	}()

	// Wait for both transactions.
	<-done
	<-done

	if got, want := bank.Balance(), 300; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}

	if got, want := bank.Withdraw(300), true; got != want {
		t.Errorf("Success = %v, want %v", got, want)
	}

	if got, want := bank.Withdraw(1), false; got != want {
		t.Errorf("Success = %v, want %v", got, want)
	}

	if got, want := bank.Balance(), 0; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}
