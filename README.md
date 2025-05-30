# Accounting Library for Go

A robust Go library for handling accounting transactions with support for various inventory costing methods and double-entry bookkeeping.

## Features

- **Double-Entry Accounting**: Enforces double-entry bookkeeping principles ensuring balanced debits and credits
- **Inventory Costing Methods**:
  - FIFO (First In, First Out)
  - LIFO (Last In, First Out)
  - WAC (Weighted Average Cost)
  - HIFO (Highest In, First Out)
  - LOFO (Lowest In, First Out)
- **Inventory Management**:
  - Track quantities and amounts separately
  - Support for inventory write-downs
  - Proper handling of zero-quantity and zero-amount cases
- **Data Validation**:
  - Entry number sequence validation
  - Time sequence validation
  - Balance validation
  - Duplicate account prevention
  - Amount and quantity validation

## Installation

```bash
go get github.com/HashemJaafar7/accounting
```

## Quick Start

Here's a simple example of recording a purchase and sale using FIFO costing:

```go
package main

import (
	"fmt"
	"time"

	"github.com/HashemJaafar7/accounting"
)

// Set up in-memory storage
var inventoryStore = make(map[accounting.AccountID]accounting.Inventory)
var lastEntry accounting.AccountingEntry
var journal []accounting.AccountingEntry

// Set up required helper functions
func getInventory(addr accounting.AccountID) (accounting.Inventory, error) {
	return inventoryStore[addr], nil
}

func setInventory(addr accounting.AccountID, inv accounting.Inventory) error {
	inventoryStore[addr] = inv
	return nil
}

func getLastEntry() (accounting.AccountingEntry, error) {
	return lastEntry, nil
}

func setEntry(entry accounting.AccountingEntry) error {
	lastEntry = entry
	journal = append(journal, entry)
	return nil
}

func main() {
	const (
		capital   accounting.AccountID = -1001
		USD       accounting.AccountID = 2001
		inventory accounting.AccountID = 1001
		COGS      accounting.AccountID = 3001
		revenue   accounting.AccountID = -4001
	)
	{
		// Create an entry to start capital
		entry := accounting.AccountingEntry{
			EntryNumber: 1,
			TimeUnix:    time.Now().Unix(),
			DoubleEntry: []accounting.SingleEntry{
				{
					CostFlowType:   accounting.INFLOW, //
					AccountID: capital,           // capital
					Quantity:       0,                 //
					Amount:         1000,              // $1000 total
				},
				{
					CostFlowType:   accounting.INFLOW, // Cash going in
					AccountID: USD,               // Cash account
					Quantity:       1000,              //
					Amount:         1000,              // $1000 total
				},
			},
		}

		// Add the entry to the journal
		err := accounting.AddToJournal(entry, getInventory, setInventory, getLastEntry, setEntry)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Print the resulting inventory balance
		printInventory()
	}

	{
		purchaseEntry := accounting.AccountingEntry{
			EntryNumber: 2,
			TimeUnix:    time.Now().Unix(),
			DoubleEntry: []accounting.SingleEntry{
				{
					CostFlowType:   accounting.INFLOW,
					AccountID: inventory, // Inventory account
					Quantity:       50,        // 50 units
					Amount:         500,       // $500 total
				},
				{
					CostFlowType:   accounting.WAC,
					AccountID: USD, // Cash account
					Quantity:       500,
					Amount:         500,
				},
			},
		}

		err := accounting.AddToJournal(purchaseEntry, getInventory, setInventory, getLastEntry, setEntry)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Print the resulting inventory balance
		printInventory()
	}

	{
		// Now sell 60 units at $15 each using FIFO costing
		saleEntry := accounting.AccountingEntry{
			EntryNumber: 3,
			TimeUnix:    time.Now().Unix(),
			DoubleEntry: []accounting.SingleEntry{
				{
					CostFlowType:   accounting.FIFO, // Use FIFO costing
					AccountID: inventory,       // Inventory account
					Quantity:       5,               // Sell 5 units
					Amount:         50,              // Cost of goods sold ($10/unit)
				},
				{
					CostFlowType:   accounting.INFLOW, //
					AccountID: COGS,              // COGS account
					Quantity:       5,                 //
					Amount:         50,                //
				},
				{
					CostFlowType:   accounting.INFLOW, //
					AccountID: USD,               // cash account
					Quantity:       80,                // Sell 60 units
					Amount:         80,                // Cost of goods sold ($10/unit)
				},
				{
					CostFlowType:   accounting.INFLOW, //
					AccountID: revenue,           // revenue account
					Quantity:       5,                 // Sell 5 units price 16
					Amount:         80,                //
				},
			},
		}

		err := accounting.AddToJournal(saleEntry, getInventory, setInventory, getLastEntry, setEntry)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Print the resulting inventory balance
		printInventory()
	}

	//output:
	// ID: -1001  inventory:[{1 0 1000}]
	// ID: 2001   inventory:[{1 1000 1000}]
	// ____________________________________________
	// ID: -1001  inventory:[{1 0 1000}]
	// ID: 2001   inventory:[{2 500 500}]
	// ID: 1001   inventory:[{2 50 500}]
	// ____________________________________________
	// ID: 2001   inventory:[{2 500 500} {3 80 80}]
	// ID: 1001   inventory:[{2 45 450}]
	// ID: 3001   inventory:[{3 5 50}]
	// ID: -4001  inventory:[{3 5 80}]
	// ID: -1001  inventory:[{1 0 1000}]
	// ____________________________________________
}

func printInventory() {
	for k, v := range inventoryStore {
		fmt.Printf("ID: %v\tinventory:%v\n", k, v)
	}
	fmt.Println("____________________________________________")
}
```

## Usage

### Account IDes

- Positive account IDes are debit-nature accounts
- Negative account IDes are credit-nature accounts

### Cost Flow Types

The library supports multiple cost flow types:

- `INFLOW`: For purchases and other additions to inventory
- `WAC`: Weighted Average Cost method
- `FIFO`: First In, First Out method
- `LIFO`: Last In, First Out method
- `HIFO`: Highest In, First Out method
- `LOFO`: Lowest In, First Out method
- `NONE`: For non-inventory transactions

### Recording Transactions

Every transaction must:

1. Have balanced debits and credits
2. Have a sequential entry number
3. Have a timestamp greater than the previous entry

### Error Handling

The library provides detailed error messages for common issues:

- Invalid entry numbers
- Time sequence violations
- Unbalanced entries
- Duplicate accounts in a single entry
- Insufficient inventory
- Amount mismatches
- Invalid cost flow types

## Examples

Check the [examples_test.go](examples_test.go) file for comprehensive examples of:

- Recording inventory purchases
- Selling inventory using different costing methods
- Handling inventory write-downs
- Using weighted average cost method

## Best Practices

1. Always validate the return error from `AddToJournal`
2. Use appropriate cost flow types:
   - `INFLOW` for purchases
   - `FIFO/LIFO/WAC` for sales
   - `NONE` if you want to do by hand (Specific Identification Method)
3. Keep track of your account IDes consistently
4. Implement proper storage for inventory and journal entries

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](LICENSE)
