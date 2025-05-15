package accounting_test

import (
	"fmt"
	"time"

	"github.com/HashemJaafar7/accounting"
)

// Set up in-memory storage
var inventoryStore = make(map[accounting.AccountAddress]accounting.Inventory)
var lastEntry accounting.AccountingEntry
var journal []accounting.AccountingEntry

// Set up required helper functions
func getInventory(addr accounting.AccountAddress) (accounting.Inventory, error) {
	return inventoryStore[addr], nil
}

func setInventory(addr accounting.AccountAddress, inv accounting.Inventory) error {
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

func Example_purchaseInventory() {
	const (
		capital   accounting.AccountAddress = -1001
		USD       accounting.AccountAddress = 2001
		inventory accounting.AccountAddress = 1001
		COGS      accounting.AccountAddress = 3001
		revenue   accounting.AccountAddress = -4001
	)
	{
		// Create an entry to start capital
		entry := accounting.AccountingEntry{
			EntryNumber: 1,
			TimeUnix:    time.Now().Unix(),
			DoubleEntry: []accounting.SingleEntry{
				{
					CostFlowType:   accounting.INFLOW, //
					AccountAddress: capital,           // capital
					Quantity:       0,                 //
					Amount:         1000,              // $1000 total
				},
				{
					CostFlowType:   accounting.INFLOW, // Cash going in
					AccountAddress: USD,               // Cash account
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
					AccountAddress: inventory, // Inventory account
					Quantity:       50,        // 50 units
					Amount:         500,       // $500 total
				},
				{
					CostFlowType:   accounting.WAC,
					AccountAddress: USD, // Cash account
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
					AccountAddress: inventory,       // Inventory account
					Quantity:       5,               // Sell 5 units
					Amount:         50,              // Cost of goods sold ($10/unit)
				},
				{
					CostFlowType:   accounting.INFLOW, //
					AccountAddress: COGS,              // COGS account
					Quantity:       5,                 //
					Amount:         50,                //
				},
				{
					CostFlowType:   accounting.INFLOW, //
					AccountAddress: USD,               // cash account
					Quantity:       80,                // Sell 60 units
					Amount:         80,                // Cost of goods sold ($10/unit)
				},
				{
					CostFlowType:   accounting.INFLOW, //
					AccountAddress: revenue,           // revenue account
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

	i := accounting.EntryNumber(1)
	iterOnJournalFunction := func() (accounting.AccountingEntry, bool, error) {
		a := journal[i]
		i++
		return a, len(journal) == int(i)+1, nil
	}

	err := accounting.CheckAllTheJournal(setInventory, iterOnJournalFunction)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	//output:
	// address: -1001  inventory:[{1 0 1000}]
	// address: 2001   inventory:[{1 1000 1000}]
	// ____________________________________________
	// address: -1001  inventory:[{1 0 1000}]
	// address: 2001   inventory:[{2 500 500}]
	// address: 1001   inventory:[{2 50 500}]
	// ____________________________________________
	// address: 2001   inventory:[{2 500 500} {3 80 80}]
	// address: 1001   inventory:[{2 45 450}]
	// address: 3001   inventory:[{3 5 50}]
	// address: -4001  inventory:[{3 5 80}]
	// address: -1001  inventory:[{1 0 1000}]
	// ____________________________________________
}

func printInventory() {
	for k, v := range inventoryStore {
		fmt.Printf("address: %v\tinventory:%v\n", k, v)
	}
	fmt.Println("____________________________________________")
}
