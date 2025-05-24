package accounting

import (
	"fmt"
	"testing"

	"github.com/HashemJaafar7/goerrors"
	"github.com/HashemJaafar7/testutils"
)

func fTest[t any](actual t, expected t) {
	testutils.Test(true, false, true, 10, "v", actual, expected)
}

func Test_GetTotalInventory(t *testing.T) {
	type input struct {
		Inventory Inventory
	}
	type output struct {
		Quantity Quantity
		Amount   Amount
	}
	tests := []struct {
		line   string
		input  input
		output output
	}{
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: []InventoryRecord{
					{TimeUnix: 0, Quantity: 10, Amount: 100},
					{TimeUnix: 0, Quantity: 20, Amount: 200},
					{TimeUnix: 0, Quantity: 30, Amount: 300},
				},
			},
			output: output{
				Quantity: 60,
				Amount:   600,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: []InventoryRecord{
					{TimeUnix: 0, Quantity: 5, Amount: 50},
					{TimeUnix: 0, Quantity: 15, Amount: 150},
				},
			},
			output: output{
				Quantity: 20,
				Amount:   200,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: []InventoryRecord{},
			},
			output: output{
				Quantity: 0,
				Amount:   0,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: []InventoryRecord{
					{TimeUnix: 0, Quantity: 0, Amount: 0},
					{TimeUnix: 0, Quantity: 0, Amount: 0},
				},
			},
			output: output{
				Quantity: 0,
				Amount:   0,
			},
		},
	}
	for _, tt := range tests {
		var output output
		output.Quantity, output.Amount = GetTotalInventory(tt.input.Inventory)

		testutils.TestCase("+v", tt.line, tt.input, tt.output, output)
	}
}

func Test_GetStatus(t *testing.T) {
	type input struct {
		CostFlowType CostFlowType
		AccountID    AccountID
	}
	type output struct {
		IsDebit IsDebit
	}
	tests := []struct {
		line   string
		input  input
		output output
	}{
		{
			line: testutils.GetLine(),
			input: input{
				CostFlowType: INFLOW,
				AccountID:    AccountID(1),
			},
			output: output{
				IsDebit: true,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				CostFlowType: INFLOW,
				AccountID:    AccountID(-1),
			},
			output: output{
				IsDebit: false,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				CostFlowType: WAC,
				AccountID:    AccountID(1),
			},
			output: output{
				IsDebit: false,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				CostFlowType: WAC,
				AccountID:    AccountID(-1),
			},
			output: output{
				IsDebit: true,
			},
		},
	}
	for _, tt := range tests {
		var output output
		output.IsDebit = GetStatus(tt.input.CostFlowType, tt.input.AccountID)

		testutils.TestCase("+v", tt.line, tt.input, tt.output, output)
	}
}

func Test_SortInventoryByTime(t *testing.T) {
	type input struct {
		Inventory Inventory
	}
	type output struct {
		Inventory Inventory
	}
	tests := []struct {
		line   string
		input  input
		output output
	}{
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: Inventory{
					{TimeUnix: 3, Quantity: 10, Amount: 100},
					{TimeUnix: 1, Quantity: 20, Amount: 200},
					{TimeUnix: 2, Quantity: 30, Amount: 300},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 20, Amount: 200},
					{TimeUnix: 2, Quantity: 30, Amount: 300},
					{TimeUnix: 3, Quantity: 10, Amount: 100},
				},
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 100},
					{TimeUnix: 1, Quantity: 20, Amount: 200},
				},
			},
			output: output{
				Inventory: Inventory{{TimeUnix: 1, Quantity: 10, Amount: 100},
					{TimeUnix: 1, Quantity: 20, Amount: 200},
				},
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: Inventory{},
			},
			output: output{
				Inventory: Inventory{},
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 100},
				},
			},
			output: output{
				Inventory: Inventory{{TimeUnix: 1, Quantity: 10, Amount: 100}},
			},
		},
	}
	for _, tt := range tests {
		var output output
		SortInventoryByTime(tt.input.Inventory)
		output.Inventory = tt.input.Inventory

		testutils.TestCase("+v", tt.line, tt.input, tt.output, output)
	}
}

func Test_addQuantityAndAmountOnInventory(t *testing.T) {
	type input struct {
		TimeUnix  TimeUnix
		Quantity  Quantity
		Amount    Amount
		Inventory Inventory
	}
	type output struct {
		Inventory Inventory
		err       error
	}
	tests := []struct {
		line   string
		input  input
		output output
	}{
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix:  1,
				Quantity:  100,
				Amount:    0,
				Inventory: nil,
			},
			output: output{
				Inventory: Inventory{{1, 100, 0}},
				err:       nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 3,
				Quantity: 50,
				Amount:   0,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 100, Amount: 0},
					{TimeUnix: 2, Quantity: 200, Amount: 0},
				},
			},
			output: output{
				Inventory: Inventory{{TimeUnix: 3, Quantity: 350, Amount: 0}},
				err:       nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 3,
				Quantity: 60,
				Amount:   0,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 100},
					{TimeUnix: 2, Quantity: 20, Amount: 200},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 3, Quantity: 90, Amount: 300},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 0,
				Quantity: -150,
				Amount:   0,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 100, Amount: 100},
				},
			},
			output: output{
				Inventory: nil,
				err:       fmt.Errorf("ErrInsufficientQuantityInInventory : You want to withdraw quantity = 150 but you do not have enough quantity because your total quantity = 100"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 5,
				Quantity: 30,
				Amount:   0,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 50, Amount: 500},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 5, Quantity: 80, Amount: 500},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 0,
				Quantity: 0,
				Amount:   0,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 100},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 100},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix:  1,
				Quantity:  -5,
				Amount:    -50,
				Inventory: Inventory{{TimeUnix: 1, Quantity: 10, Amount: 100}},
			},
			output: output{
				Inventory: Inventory{{TimeUnix: 1, Quantity: 5, Amount: 50}},
				err:       nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix:  0,
				Quantity:  -15,
				Amount:    -50,
				Inventory: Inventory{{TimeUnix: 1, Quantity: 10, Amount: 100}},
			},
			output: output{
				Inventory: nil,
				err:       fmt.Errorf("ErrInsufficientQuantityInInventory : You want to withdraw quantity = 15 but you do not have enough quantity because your total quantity = 10"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix:  0,
				Quantity:  -5,
				Amount:    -150,
				Inventory: Inventory{{TimeUnix: 1, Quantity: 10, Amount: 100}},
			},
			output: output{
				Inventory: nil,
				err:       fmt.Errorf("ErrInsufficientAmountInInventory : You want to withdraw amount = 150 but you do not have enough amount because your total amount = 100"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix:  0,
				Quantity:  -10,
				Amount:    -100,
				Inventory: Inventory{{TimeUnix: 1, Quantity: 10, Amount: 100}},
			},
			output: output{
				Inventory: Inventory{{TimeUnix: 0, Quantity: 0, Amount: 0}},
				err:       nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix:  0,
				Quantity:  -0,
				Amount:    -0,
				Inventory: Inventory{{TimeUnix: 1, Quantity: 10, Amount: 100}},
			},
			output: output{
				Inventory: Inventory{{TimeUnix: 1, Quantity: 10, Amount: 100}},
				err:       nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix:  1,
				Quantity:  0,
				Amount:    100,
				Inventory: Inventory{},
			},
			output: output{
				Inventory: Inventory{{1, 0, 100}},
				err:       nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 2,
				Quantity: 0,
				Amount:   50,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 0, Amount: 100},
					{TimeUnix: 2, Quantity: 0, Amount: 200},
				},
			},
			output: output{
				Inventory: Inventory{{2, 0, 350}},
				err:       nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 3,
				Quantity: 0,
				Amount:   60,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 100},
					{TimeUnix: 2, Quantity: 20, Amount: 200},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 3, Quantity: 30, Amount: 360},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 4,
				Quantity: 0,
				Amount:   30,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 5, Amount: 50},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 4, Quantity: 5, Amount: 80},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 5,
				Quantity: 0,
				Amount:   90,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 0, Amount: 0},
					{TimeUnix: 2, Quantity: 10, Amount: 100},
					{TimeUnix: 3, Quantity: 20, Amount: 200},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 5, Quantity: 30, Amount: 390},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 6,
				Quantity: 0,
				Amount:   90,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 0, Amount: 0},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 6, Quantity: 0, Amount: 90},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 7,
				Quantity: 0,
				Amount:   -90,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 0, Amount: 0},
					{TimeUnix: 2, Quantity: 10, Amount: 100},
					{TimeUnix: 3, Quantity: 20, Amount: 200},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 7, Quantity: 30, Amount: 210},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 8,
				Quantity: 0,
				Amount:   -50,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 100},
					{TimeUnix: 2, Quantity: 20, Amount: 200},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 8, Quantity: 30, Amount: 250},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 9,
				Quantity: 0,
				Amount:   -50,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 0, Amount: 100},
					{TimeUnix: 2, Quantity: 0, Amount: 200},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 9, Quantity: 0, Amount: 250},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 10,
				Quantity: 0,
				Amount:   -100,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 0, Amount: 50},
					{TimeUnix: 2, Quantity: 10, Amount: 100},
					{TimeUnix: 3, Quantity: 20, Amount: 200},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 10, Quantity: 30, Amount: 250},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 11,
				Quantity: 0,
				Amount:   -30,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 5, Amount: 50},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 11, Quantity: 5, Amount: 20},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 0,
				Quantity: 0,
				Amount:   -50,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 0},
					{TimeUnix: 2, Quantity: 20, Amount: 0},
				},
			},
			output: output{
				Inventory: nil,
				err:       fmt.Errorf("ErrInsufficientAmountInInventory : You want to withdraw amount = 50 but you do not have enough amount because your total amount = 0"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix:  0,
				Quantity:  0,
				Amount:    -50,
				Inventory: Inventory{},
			},
			output: output{
				Inventory: nil,
				err:       fmt.Errorf("ErrInsufficientAmountInInventory : You want to withdraw amount = 50 but you do not have enough amount because your total amount = 0"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				TimeUnix: 10,
				Quantity: 0,
				Amount:   -100,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 0, Amount: 50},
					{TimeUnix: 2, Quantity: 1000, Amount: 100},
					{TimeUnix: 3, Quantity: 25, Amount: 200},
					{TimeUnix: 4, Quantity: 25, Amount: 1},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 10, Quantity: 1050, Amount: 251},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		var output output
		output.Inventory, output.err = addQuantityAndAmountOnInventory(tt.input.TimeUnix, tt.input.Quantity, tt.input.Amount, tt.input.Inventory)
		output.err = goerrors.NormalizeTheError(output.err)
		testutils.TestCase("+v", tt.line, tt.input, tt.output, output)
	}
}

func Test_decreaseInventory(t *testing.T) {
	type input struct {
		Quantity  Quantity
		Amount    Amount
		Inventory Inventory
	}
	type output struct {
		Inventory Inventory
		err       error
	}
	tests := []struct {
		line   string
		input  input
		output output
	}{
		{
			line: testutils.GetLine(),
			input: input{
				Quantity: 10,
				Amount:   100,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 5, Amount: 50},
					{TimeUnix: 2, Quantity: 5, Amount: 50},
				},
			},
			output: output{
				Inventory: nil,
				err:       nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Quantity: 15,
				Amount:   100,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 5, Amount: 50},
					{TimeUnix: 2, Quantity: 5, Amount: 50},
				},
			},
			output: output{
				Inventory: nil,
				err:       fmt.Errorf("ErrInsufficientQuantityInInventory : You want to withdraw quantity = 15 but you do not have enough quantity because your total quantity = 10"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Quantity: 10,
				Amount:   150,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 5, Amount: 50},
					{TimeUnix: 2, Quantity: 5, Amount: 50},
				},
			},
			output: output{
				Inventory: nil,
				err:       fmt.Errorf("ErrInsufficientAmountInInventory : You want to withdraw amount = 150 but you do not have enough amount because your total amount = 100"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Quantity: 7,
				Amount:   70,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 5, Amount: 50},
					{TimeUnix: 2, Quantity: 5, Amount: 50},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 2, Quantity: 3, Amount: 30},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Quantity: 10,
				Amount:   100,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 100},
				},
			},
			output: output{
				Inventory: nil,
				err:       nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Quantity:  5,
				Amount:    50,
				Inventory: Inventory{},
			},
			output: output{
				Inventory: nil,
				err:       fmt.Errorf("ErrInventoryIsEmpty : inventory is empty"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Quantity: 10,
				Amount:   100,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 5, Amount: 50},
					{TimeUnix: 2, Quantity: 5, Amount: 60},
				},
			},
			output: output{
				Inventory: nil,
				err:       fmt.Errorf("ErrAmountMismatch : amount mismatch: expected to enter amount = 110 but got = 100"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Quantity: 3,
				Amount:   30,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 5, Amount: 50},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 2, Amount: 20},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Quantity: 3,
				Amount:   30,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 0, Amount: 50},
				},
			},
			output: output{
				Inventory: nil,
				err:       fmt.Errorf("ErrInsufficientQuantityInInventory : You want to withdraw quantity = 3 but you do not have enough quantity because your total quantity = 0"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Quantity: 10,
				Amount:   100,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 5, Amount: 50},
					{TimeUnix: 2, Quantity: 20, Amount: 200},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 2, Quantity: 15, Amount: 150},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Quantity: 10,
				Amount:   75,
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 5, Amount: 50},
					{TimeUnix: 2, Quantity: 20, Amount: 100},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 2, Quantity: 15, Amount: 75},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		var output output
		output.Inventory, output.err = decreaseInventory(tt.input.Quantity, tt.input.Amount, tt.input.Inventory)
		output.err = goerrors.NormalizeTheError(output.err)
		testutils.TestCase("+v", tt.line, tt.input, tt.output, output)
	}
}

func Test_removeZeros(t *testing.T) {
	type input struct {
		Inventory Inventory
	}
	type output struct {
		Inventory Inventory
	}
	tests := []struct {
		line   string
		input  input
		output output
	}{
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 100},
					{TimeUnix: 2, Quantity: 20, Amount: 200},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 100},
					{TimeUnix: 2, Quantity: 20, Amount: 200},
				},
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 0, Amount: 0},
					{TimeUnix: 2, Quantity: 20, Amount: 200},
					{TimeUnix: 3, Quantity: 0, Amount: 0},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 2, Quantity: 20, Amount: 200},
				},
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: Inventory{},
			},
			output: output{
				Inventory: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 0, Amount: 0},
					{TimeUnix: 2, Quantity: 0, Amount: 0},
				},
			},
			output: output{
				Inventory: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 0},
					{TimeUnix: 2, Quantity: 0, Amount: 100},
					{TimeUnix: 3, Quantity: 0, Amount: 0},
				},
			},
			output: output{
				Inventory: Inventory{
					{TimeUnix: 1, Quantity: 10, Amount: 0},
					{TimeUnix: 2, Quantity: 0, Amount: 100},
				},
			},
		},
	}
	for _, tt := range tests {
		var output output
		output.Inventory = removeZeros(tt.input.Inventory)

		testutils.TestCase("+v", tt.line, tt.input, tt.output, output)
	}
}

type myDB struct {
	myInv         AccountIDAndInventory
	myEntries     []AccountingEntry
	lastEntryTime TimeUnix
	i             int
}

func (s *myDB) GetInventory(key AccountID) (Inventory, error) {
	inv, ok := s.myInv[key]
	if !ok {
		return Inventory{}, nil
	}
	return inv, nil
}
func (s *myDB) SetInventory(key AccountID, value Inventory) error {
	s.myInv[key] = value
	return nil
}
func (s *myDB) GetLastEntryTime() (TimeUnix, error) {
	return s.lastEntryTime, nil
}
func (s *myDB) SetEntry(value AccountingEntry) error {
	s.lastEntryTime = value.TimeUnix
	s.myEntries = append(s.myEntries, value)
	return nil
}
func (s *myDB) IterOnJournal() (AccountingEntry, bool, error) {
	if len(s.myEntries) == s.i {
		return AccountingEntry{}, false, nil
	}

	a := s.myEntries[s.i]
	s.i++
	return a, true, nil
}

func Test_AddToJournal(t *testing.T) {
	var kk myDB
	kk.myInv = make(AccountIDAndInventory)

	type input struct {
		AccountingEntry AccountingEntry
	}
	type output struct {
		err error
	}
	tests := []struct {
		line   string
		input  input
		output output
	}{
		{
			line: testutils.GetLine(),
			input: input{
				AccountingEntry: AccountingEntry{
					TimeUnix: 1000,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
						{INFLOW, -1, 10, 100},
					},
				},
			},
			output: output{
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				AccountingEntry: AccountingEntry{
					TimeUnix: 900,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
						{WAC, -1, 10, 100},
					},
				},
			},
			output: output{
				err: fmt.Errorf("ErrTimeShouldBeBigger : time should be bigger"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
						{WAC, 1, 10, 100},
					},
				},
			},
			output: output{
				err: fmt.Errorf("ErrDuplicateAccountInEntry : duplicate account ID 1 in entry"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
					},
				},
			},
			output: output{
				err: fmt.Errorf("ErrDebitNotEqualCredit : debit not equal credit and debit = 100 , credit = 0 and debit-credit = 100"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
						{WAC, -1, 10, 90},
					},
				},
			},
			output: output{
				err: fmt.Errorf("ErrDebitNotEqualCredit : debit not equal credit and debit = 190 , credit = 0 and debit-credit = 190"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{INFLOW, 99, 10, 100},
						{WAC, -1, 10, 100},
					},
				},
			},
			output: output{
				err: fmt.Errorf("ErrDebitNotEqualCredit : debit not equal credit and debit = 200 , credit = 0 and debit-credit = 200"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				AccountingEntry: AccountingEntry{
					TimeUnix: 1001,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
						{INFLOW, -1, 10, 100},
					},
				},
			},
			output: output{
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				AccountingEntry: AccountingEntry{
					TimeUnix: 1002,
					DoubleEntry: DoubleEntry{
						{9, 1, 10, 100},
						{INFLOW, -1, 10, 100},
					},
				},
			},
			output: output{
				err: fmt.Errorf("ErrTheCostFlowTypeIsWrong : the cost flow type is wrong"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				AccountingEntry: AccountingEntry{
					TimeUnix: 1003,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, -10, 100},
						{INFLOW, -1, 10, 100},
					},
				},
			},
			output: output{
				err: fmt.Errorf("ErrTheQuantityAndAmountShouldBeBothPositive : the quantity and amount should be both positive for account ID 1"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				AccountingEntry: AccountingEntry{
					TimeUnix: 1004,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
						{WAC, 2, 10, 100},
					},
				},
			},
			output: output{
				err: fmt.Errorf("ErrInsufficientQuantityInInventory : You want to withdraw quantity = 10 but you do not have enough quantity because your total quantity = 0"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				AccountingEntry: AccountingEntry{
					TimeUnix: 1005,
					DoubleEntry: DoubleEntry{
						{WAC, 1, 5, 50},
						{WAC, -1, 5, 50},
					},
				},
			},
			output: output{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		var output output
		output.err = AddToJournal(tt.input.AccountingEntry, &kk)
		output.err = goerrors.NormalizeTheError(output.err)
		testutils.TestCase("+v", tt.line, tt.input, tt.output, output)
	}

	fTest(kk.myInv, AccountIDAndInventory{
		-1: Inventory{{TimeUnix: 1005, Quantity: 15, Amount: 150}},
		1:  Inventory{{TimeUnix: 1005, Quantity: 15, Amount: 150}},
	})
	myInvExpected := kk.myInv
	err := CheckAllTheJournal(&kk)
	fTest(err, nil)
	fTest(kk.myInv, myInvExpected)
}

func Test_CheckAndProcessDoubleEntry(t *testing.T) {
	type input struct {
		LastTimeUnix          TimeUnix
		AccountingEntry       AccountingEntry
		AccountIDAndInventory AccountIDAndInventory
	}
	type output struct {
		AccountIDAndInventory AccountIDAndInventory
		err                   error
	}
	tests := []struct {
		line   string
		input  input
		output output
	}{
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 900,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
						{WAC, -1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{0, 0, 0}},
					-1: Inventory{{0, 0, 0}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrTimeShouldBeBigger : time should be bigger"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{0, 0, 0}},
					-1: Inventory{{0, 0, 0}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrDebitNotEqualCredit : debit not equal credit and debit = 100 , credit = 0 and debit-credit = 100"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{99, 1, 10, 100},
						{WAC, -1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{0, 0, 0}},
					-1: Inventory{{0, 0, 0}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrTheCostFlowTypeIsWrong : the cost flow type is wrong"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, -10, 100},
						{WAC, -1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{0, 0, 0}},
					-1: Inventory{{0, 0, 0}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrTheQuantityAndAmountShouldBeBothPositive : the quantity and amount should be both positive for account ID 1"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
						{WAC, 1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{0, 0, 0}},
					-1: Inventory{{0, 0, 0}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrDuplicateAccountInEntry : duplicate account ID 1 in entry"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
						{WAC, -1, 10, 50},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{0, 0, 0}},
					-1: Inventory{{0, 0, 0}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrDebitNotEqualCredit : debit not equal credit and debit = 150 , credit = 0 and debit-credit = 150"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
						{WAC, 2, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1: Inventory{{0, 0, 0}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrInventoryNotFoundForAccountID : inventory not found for account ID 2"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 100},
						{INFLOW, -1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{0, 0, 0}},
					-1: Inventory{{0, 0, 0}},
				},
			},
			output: output{
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1100, 10, 100}},
					-1: Inventory{{1100, 10, 100}},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{WAC, 1, 10, -100},
						{WAC, -1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{0, 0, 0}},
					-1: Inventory{{0, 0, 0}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrTheQuantityAndAmountShouldBeBothPositive : the quantity and amount should be both positive for account ID 1"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{WAC, 1, 10, 100},
						{WAC, -1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrAmountMismatch : amount mismatch: expected to enter amount = 83.83064516129032 but got = 100"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{FIFO, 1, 10, 100},
						{FIFO, -1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrAmountMismatch : amount mismatch: expected to enter amount = 18 but got = 100"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{LIFO, 1, 10, 100},
						{LIFO, -1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrAmountMismatch : amount mismatch: expected to enter amount = 128.83116883116884 but got = 100"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{HIFO, 1, 10, 100},
						{HIFO, -1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrAmountMismatch : amount mismatch: expected to enter amount = 972.4155844155844 but got = 100"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{LOFO, 1, 10, 100},
						{LOFO, -1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrAmountMismatch : amount mismatch: expected to enter amount = 1.6363636363636362 but got = 100"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{NONE, 1, 10, 100},
						{NONE, -1, 10, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1100, 238, 1979}},
					-1: Inventory{{1100, 238, 1979}},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{WAC, 1, 0, 0},
						{WAC, -1, 0, 0},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrQuantityAndAmountAreZero : you can't enter both quantity and amount as zeros for account ID 1"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 0, 100},
						{INFLOW, -1, 0, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1100, 248, 2179}},
					-1: Inventory{{1100, 248, 2179}},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{INFLOW, 1, 10, 0},
						{INFLOW, -1, 10, 0},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1100, 258, 2079}},
					-1: Inventory{{1100, 258, 2079}},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{WAC, 1, 0, 100},
						{WAC, -1, 0, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrYouShouldUseCostFlowTypeNONEIfYouHaveQuantityOrAmountZero : you should to use cost flow type NONE because your quantity or amount is zero for account ID 1"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{WAC, 1, 10, 0},
						{WAC, -1, 10, 0},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: nil,
				err:                   fmt.Errorf("ErrYouShouldUseCostFlowTypeNONEIfYouHaveQuantityOrAmountZero : you should to use cost flow type NONE because your quantity or amount is zero for account ID 1"),
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{NONE, 1, 0, 100},
						{NONE, -1, 0, 100},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1100, 248, 1979}},
					-1: Inventory{{1100, 248, 1979}},
				},
				err: nil,
			},
		},
		{
			line: testutils.GetLine(),
			input: input{
				LastTimeUnix: 1000,
				AccountingEntry: AccountingEntry{
					TimeUnix: 1100,
					DoubleEntry: DoubleEntry{
						{NONE, 1, 10, 0},
						{NONE, -1, 10, 0},
					},
				},
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
					-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				},
			},
			output: output{
				AccountIDAndInventory: AccountIDAndInventory{
					1:  Inventory{{1100, 238, 2079}},
					-1: Inventory{{1100, 238, 2079}},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		var output output
		output.AccountIDAndInventory, output.err = CheckAndProcessDoubleEntry(tt.input.LastTimeUnix, tt.input.AccountingEntry, tt.input.AccountIDAndInventory)
		output.err = goerrors.NormalizeTheError(output.err)
		testutils.TestCase("+v", tt.line, tt.input, tt.output, output)
	}
}
