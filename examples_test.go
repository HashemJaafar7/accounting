package accounting_test

import (
	"fmt"
	"time"

	"github.com/HashemJaafar7/accounting"
)

// Set up required helper functions
type myDB struct {
	myInv         accounting.AccountIDAndInventory
	myEntries     []accounting.AccountingEntry
	lastEntryTime accounting.TimeUnix
	i             int
}

func (s *myDB) GetInventory(key accounting.AccountID) (accounting.Inventory, error) {
	inv, ok := s.myInv[key]
	if !ok {
		// return nil, fmt.Errorf("ID %v not found", key)
	}
	return inv, nil
}
func (s *myDB) SetInventory(key accounting.AccountID, value accounting.Inventory) error {
	s.myInv[key] = value
	return nil
}
func (s *myDB) GetLastEntryTime() (accounting.TimeUnix, error) {
	return s.lastEntryTime, nil
}
func (s *myDB) SetEntry(value accounting.AccountingEntry) error {
	s.lastEntryTime = value.TimeUnix
	s.myEntries = append(s.myEntries, value)
	return nil
}
func (s *myDB) IterOnJournal() (accounting.AccountingEntry, bool, error) {
	if len(s.myEntries) == s.i {
		return accounting.AccountingEntry{}, false, nil
	}

	a := s.myEntries[s.i]
	s.i++
	return a, true, nil
}

var kk myDB

func Example_purchaseInventory() {
	// Set up in-memory storage
	kk.myInv = make(map[accounting.AccountID]accounting.Inventory)

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
			TimeUnix: time.Now().UnixMicro(),
			DoubleEntry: []accounting.SingleEntry{
				{
					CostFlowType: accounting.INFLOW, //
					AccountID:    capital,           // capital
					Quantity:     0,                 //
					Amount:       1000,              // $1000 total
				},
				{
					CostFlowType: accounting.INFLOW, // Cash going in
					AccountID:    USD,               // Cash account
					Quantity:     1000,              //
					Amount:       1000,              // $1000 total
				},
			},
		}

		// Add the entry to the journal
		err := accounting.AddToJournal(entry, &kk)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Print the resulting inventory balance
		printInventory()
	}

	{
		entry := accounting.AccountingEntry{
			TimeUnix: time.Now().UnixMicro(),
			DoubleEntry: []accounting.SingleEntry{
				{
					CostFlowType: accounting.INFLOW,
					AccountID:    inventory, // Inventory account
					Quantity:     50,        // 50 units
					Amount:       500,       // $500 total
				},
				{
					CostFlowType: accounting.WAC,
					AccountID:    USD, // Cash account
					Quantity:     500,
					Amount:       500,
				},
			},
		}

		err := accounting.AddToJournal(entry, &kk)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Print the resulting inventory balance
		printInventory()
	}

	{
		// Now sell 60 units at $15 each using FIFO costing
		entry := accounting.AccountingEntry{
			TimeUnix: time.Now().UnixMicro(),
			DoubleEntry: []accounting.SingleEntry{
				{
					CostFlowType: accounting.FIFO, // Use FIFO costing
					AccountID:    inventory,       // Inventory account
					Quantity:     5,               // Sell 5 units
					Amount:       50,              // Cost of goods sold ($10/unit)
				},
				{
					CostFlowType: accounting.INFLOW, //
					AccountID:    COGS,              // COGS account
					Quantity:     5,                 //
					Amount:       50,                //
				},
				{
					CostFlowType: accounting.INFLOW, //
					AccountID:    USD,               // cash account
					Quantity:     80,                // Sell 60 units
					Amount:       80,                // Cost of goods sold ($10/unit)
				},
				{
					CostFlowType: accounting.INFLOW, //
					AccountID:    revenue,           // revenue account
					Quantity:     5,                 // Sell 5 units price 16
					Amount:       80,                //
				},
			},
		}

		err := accounting.AddToJournal(entry, &kk)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Print the resulting inventory balance
		printInventory()
	}

	err := accounting.CheckAllTheJournal(&kk)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
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
	for k, v := range kk.myInv {
		fmt.Printf("ID: %v\tinventory:%v\n", k, v)
	}
	fmt.Println("____________________________________________")
}
