// Package bank adds a Withdraw function to the bank package described in
// Chapter 9, Section 1.
package bank

type withdrawalMsg struct {
	amount      int
	txSucceeded chan bool
}

var (
	deposits    = make(chan int)           // send amount to deposit
	balances    = make(chan int)           // receive balance
	withdrawals = make(chan withdrawalMsg) // indicate success of withdrawals
)

// Deposit deposits amount to the account.
func Deposit(amount int) { deposits <- amount }

// Balance returns the current balance of the account.
func Balance() int { return <-balances }

// Withdraw attempts to withdraw amount and returns whether this was successful
// or not.
func Withdraw(amount int) bool {
	txSucceeded := make(chan bool)
	withdrawals <- withdrawalMsg{
		amount:      amount,
		txSucceeded: txSucceeded,
	}
	return <-txSucceeded
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case wdMsg := <-withdrawals:
			if wdMsg.amount > balance {
				wdMsg.txSucceeded <- false
				continue
			}
			balance -= wdMsg.amount
			wdMsg.txSucceeded <- true
		case balances <- balance:
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}
