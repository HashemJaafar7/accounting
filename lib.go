package accounting

import (
	"math"
	"slices"

	"github.com/HashemJaafar7/goerrors"
)

// errors
const (
	ErrEntryNumberShouldBeBiggerByOne                            = "ErrEntryNumberShouldBeBiggerByOne"
	ErrTimeShouldBeBigger                                        = "ErrTimeShouldBeBigger"
	ErrEntryMustHaveAtLeast_2Entries                             = "ErrEntryMustHaveAtLeast_2Entries"
	ErrDuplicateAccountInEntry                                   = "ErrDuplicateAccountInEntry"
	ErrSumOfAmountsIsNotZeroAndDebitMoreThanCredit               = "ErrSumOfAmountsIsNotZeroAndDebitMoreThanCredit"
	ErrInventoryNotFoundForAccountAddress                        = "ErrInventoryNotFoundForAccountAddress"
	ErrQuantityAndAmountShouldBothBeDebitOrCredit                = "ErrQuantityAndAmountShouldBothBeDebitOrCredit"
	ErrTheCostFlowTypeIsWrong                                    = "ErrTheCostFlowTypeIsWrong"
	ErrInventoryIsEmpty                                          = "ErrInventoryIsEmpty"
	ErrInsufficientQuantityInInventory                           = "ErrInsufficientQuantityInInventory"
	ErrAmountMismatch                                            = "ErrAmountMismatch"
	ErrInsufficientAmountInInventory                             = "ErrInsufficientAmountInInventory"
	ErrTheQuantityAndAmountShouldBeBothPositive                  = "ErrTheQuantityAndAmountShouldBeBothPositive"
	ErrYouShouldUseCostFlowTypeNONEIfYouHaveQuantityOrAmountZero = "ErrYouShouldUseCostFlowTypeNONEIfYouHaveQuantityOrAmountZero"
)

// error functions
func fErrQuantityAndAmountShouldBothBeDebitOrCredit(address AccountAddress) error {
	return goerrors.Errorf(ErrQuantityAndAmountShouldBothBeDebitOrCredit, "quantity and amount should both be debit or credit for account address %v", address)
}
func fErrYouShouldUseCostFlowTypeNONEIfYouHaveQuantityOrAmountZero(address AccountAddress) error {
	return goerrors.Errorf(ErrYouShouldUseCostFlowTypeNONEIfYouHaveQuantityOrAmountZero, "you should to use cost flow type NONE because your quantity or amount is zero for account address %v", address)
}
func fErrInsufficientQuantityInInventory(inputQuantity, totalQuantity Quantity) error {
	return goerrors.Errorf(ErrInsufficientQuantityInInventory, "You want to withdraw quantity = %v but you do not have enough quantity because your total quantity = %v", math.Abs(float64(inputQuantity)), totalQuantity)
}
func fErrInsufficientAmountInInventory(inputAmount, totalAmount Amount) error {
	return goerrors.Errorf(ErrInsufficientAmountInInventory, "You want to withdraw amount = %v but you do not have enough amount because your total amount = %v", math.Abs(float64(inputAmount)), totalAmount)
}

const (
	INFLOW CostFlowType = iota
	WAC
	FIFO
	LIFO
	HIFO
	LOFO
	NONE
	theNumberOfCostFlowTypes
)

type IsDebit bool
type CostFlowType uint8
type AccountAddress int64
type Quantity float64
type Amount float64
type EntryNumber uint64
type TimeUnix = int64

type SingleEntry struct {
	CostFlowType
	AccountAddress
	Quantity
	Amount
}

type DoubleEntry []SingleEntry

type AccountingEntry struct {
	EntryNumber
	TimeUnix // the time in unix in seconds
	DoubleEntry
}

type InventoryRecord struct {
	EntryNumber
	Quantity
	Amount
}

type Inventory []InventoryRecord
type AccountAddressAndInventory map[AccountAddress]Inventory

type GetInventory func(key AccountAddress) (Inventory, error)
type SetInventory func(key AccountAddress, value Inventory) error
type GetLastEntry func() (AccountingEntry, error)
type SetEntry func(value AccountingEntry) error
type IterOnJournal func() (AccountingEntry, bool, error)

// IsNatureDebit determines if an account has a debit nature based on its address.
// A positive or zero account address indicates a debit nature account (assets, expenses),
// while a negative address indicates a credit nature account (liabilities, revenues, equity).
//
// Parameters:
//   - accountAddress: The address of the account to check
//
// Returns:
//   - isDebit: true if the account has a debit nature, false if credit nature
func IsNatureDebit(accountAddress AccountAddress) IsDebit {
	return accountAddress >= 0
}

// GetStatus determines if a account status a debit based on cost flow type and account address.
// It compares the cost flow direction (inflow/outflow) with the natural debit/credit state of the account.
// Returns true if the account status a debit, false if it a credit.
// Parameters:
//   - costFlowType: Indicates whether money/value is flowing in or out (INFLOW/WAC/FIFO/LIFO/HIFO/LOFO/NONE)
//   - accountAddress: The address/identifier of the account being affected
func GetStatus(costFlowType CostFlowType, accountAddress AccountAddress) IsDebit {
	return costFlowType == INFLOW == IsNatureDebit(accountAddress)
}

// CheckAndProcessDoubleEntry validates and processes a double-entry accounting transaction.
// It ensures the integrity of the accounting entry and updates the inventory records accordingly.
//
// Parameters:
//   - lastEntryNumber: The previous entry number for sequence validation
//   - lastTimeUnix: The timestamp of the last entry for chronological validation
//   - entry: The accounting entry to be processed
//   - accountAddressAndInventoryVariable: Current state of inventory records for all accounts
//
// Returns:
//   - AccountAddressAndInventory: Updated inventory records after processing the entry
//   - error: Error if any validation fails or processing encounters issues
//
// The function performs the following validations:
//   - Ensures entry number is sequential
//   - Verifies timestamp is after the last entry
//   - Checks minimum of 2 entries in double-entry
//   - Validates debit and credit balance
//   - Prevents duplicate accounts in single entry
//   - Verifies positive amounts and quantities
//   - Ensures valid cost flow types
//
// After validation, it processes inventory records according to various transaction scenarios:
//   - Handles inflow and outflow of quantities and amounts
//   - Manages inventory adjustments (zero quantity or amount cases)
//   - Applies cost flow accounting methods
//   - Removes zero-value inventory records
//
// The function handles different combinations of positive, negative, and zero values
// for both amounts and quantities, applying appropriate business rules for each case.
func CheckAndProcessDoubleEntry(lastEntryNumber EntryNumber, lastTimeUnix TimeUnix, entry AccountingEntry, accountAddressAndInventoryVariable AccountAddressAndInventory) (AccountAddressAndInventory, error) {
	if entry.EntryNumber != lastEntryNumber+1 {
		return nil, goerrors.Errorf(ErrEntryNumberShouldBeBiggerByOne, "entry number should be bigger by one from the last entry number")
	}

	if entry.TimeUnix < lastTimeUnix {
		return nil, goerrors.Errorf(ErrTimeShouldBeBigger, "time should be bigger")
	}

	if len(entry.DoubleEntry) < 2 {
		return nil, goerrors.Errorf(ErrEntryMustHaveAtLeast_2Entries, "entry must have at least 2 entries")
	}

	totalDebit := Amount(0)
	totalCredit := Amount(0)
	accounts := make(map[AccountAddress]bool)
	for _, single := range entry.DoubleEntry {
		if single.CostFlowType >= theNumberOfCostFlowTypes {
			return nil, goerrors.Errorf(ErrTheCostFlowTypeIsWrong, "the cost flow type is wrong")
		}
		if single.Amount < 0 || single.Quantity < 0 {
			return nil, goerrors.Errorf(ErrTheQuantityAndAmountShouldBeBothPositive, "the quantity and amount should be both positive for account address %v", single.AccountAddress)
		}
		if _, exists := accounts[single.AccountAddress]; exists {
			return nil, goerrors.Errorf(ErrDuplicateAccountInEntry, "duplicate account address %v in entry", single.AccountAddress)
		}
		accounts[single.AccountAddress] = true

		if GetStatus(single.CostFlowType, single.AccountAddress) {
			totalDebit += single.Amount
		} else {
			totalCredit += single.Amount
		}

	}

	if totalDebit != totalCredit {
		return nil, goerrors.Errorf(ErrSumOfAmountsIsNotZeroAndDebitMoreThanCredit, "debit not equal credit and debit = %v , credit = %v and debit-credit = %v", totalDebit, totalCredit, totalDebit-totalCredit)
	}

	for _, single := range entry.DoubleEntry {
		address := single.AccountAddress

		inventoryVariable, ok := accountAddressAndInventoryVariable[address]
		if !ok {
			return nil, goerrors.Errorf(ErrInventoryNotFoundForAccountAddress, "inventory not found for account address %v", address)
		}

		qty := single.Quantity
		amt := single.Amount

		if single.CostFlowType != INFLOW {
			qty = -qty
			amt = -amt
		}

		var err error
		// i should to deal with amt == 0 and qty != 0 because i deal with amt != 0 and qty == 0 before and that will make the amount zero and quantity not zero
		switch {
		case amt > 0 && qty > 0:
			inventoryVariable = append(inventoryVariable, InventoryRecord{entry.EntryNumber, qty, amt})
		case amt > 0 && qty == 0: // not sure: but it cuse to adjust the inventory: like feeding sheep
			inventoryVariable, err = addQuantityAndAmountOnInventory(entry.EntryNumber, qty, amt, inventoryVariable)
		case amt > 0 && qty < 0:
			panic(ErrQuantityAndAmountShouldBothBeDebitOrCredit)
		case amt == 0 && qty > 0: // not sure: like gift but i dont want this to happen because it will lead to decrease the quantity without any amount and that will make some entry verbose
			inventoryVariable, err = addQuantityAndAmountOnInventory(entry.EntryNumber, qty, amt, inventoryVariable)
		case amt == 0 && qty == 0:
			return nil, fErrQuantityAndAmountShouldBothBeDebitOrCredit(address)
		case amt == 0 && qty < 0: // not sure: it happens when the account is dont have any amount in the balance because it came from gifts or we make the amount 0
			if single.CostFlowType != NONE {
				return nil, fErrYouShouldUseCostFlowTypeNONEIfYouHaveQuantityOrAmountZero(address)
			}
			inventoryVariable, err = addQuantityAndAmountOnInventory(entry.EntryNumber, qty, amt, inventoryVariable)
		case amt < 0 && qty > 0:
			panic(ErrQuantityAndAmountShouldBothBeDebitOrCredit)
		case amt < 0 && qty == 0: // not sure: but it cuse to adjust the inventory: like smashing a car or depreciation or market value
			if single.CostFlowType != NONE {
				return nil, fErrYouShouldUseCostFlowTypeNONEIfYouHaveQuantityOrAmountZero(address)
			}
			inventoryVariable, err = addQuantityAndAmountOnInventory(entry.EntryNumber, qty, amt, inventoryVariable)
		case amt < 0 && qty < 0:
			inventoryVariable, err = checkAndProcessCostOutFlow(entry.EntryNumber, single, inventoryVariable)
		}

		if err != nil {
			return nil, err
		}

		inventoryVariable = removeZeros(inventoryVariable)
		accountAddressAndInventoryVariable[address] = inventoryVariable
	}

	return accountAddressAndInventoryVariable, nil
}

func checkAndProcessCostOutFlow(entryNumberVariable EntryNumber, singleEntryVariable SingleEntry, inventoryVariable Inventory) (Inventory, error) {
	qty := singleEntryVariable.Quantity
	amt := singleEntryVariable.Amount

	switch singleEntryVariable.CostFlowType {
	case WAC:
		totalQuantity, totalAmount := calculateTotalInventory(inventoryVariable)
		inventoryVariable = Inventory{{entryNumberVariable, totalQuantity, totalAmount}}
		return decreaseInventory(qty, amt, inventoryVariable)
	case FIFO:
		sortInventoryByEntryNumber(inventoryVariable)
		return decreaseInventory(qty, amt, inventoryVariable)
	case LIFO:
		sortInventoryByEntryNumber(inventoryVariable)
		slices.Reverse(inventoryVariable)
		return decreaseInventory(qty, amt, inventoryVariable)
	case HIFO:
		sortInventoryByPrice(inventoryVariable)
		slices.Reverse(inventoryVariable)
		return decreaseInventory(qty, amt, inventoryVariable)
	case LOFO:
		sortInventoryByPrice(inventoryVariable)
		return decreaseInventory(qty, amt, inventoryVariable)
	case NONE:
		return addQuantityAndAmountOnInventory(entryNumberVariable, -qty, -amt, inventoryVariable)
	default:
		panic(ErrTheCostFlowTypeIsWrong)
	}
}

func sortInventoryByPrice(inventory Inventory) {
	slices.SortFunc(inventory, func(a, b InventoryRecord) int {
		price1 := a.Amount / Amount(a.Quantity)
		price2 := b.Amount / Amount(b.Quantity)
		switch {
		case price1 > price2:
			return 1
		case price1 < price2:
			return -1
		default:
			return 0
		}
	})
}

func sortInventoryByEntryNumber(inventory Inventory) {
	slices.SortFunc(inventory, func(a, b InventoryRecord) int {
		switch {
		case a.EntryNumber > b.EntryNumber:
			return 1
		case a.EntryNumber < b.EntryNumber:
			return -1
		default:
			return 0
		}
	})
}

func calculateTotalInventory(inventory Inventory) (Quantity, Amount) {
	var totalQuantity Quantity
	var totalAmount Amount
	for _, r := range inventory {
		totalQuantity += r.Quantity
		totalAmount += r.Amount
	}
	return totalQuantity, totalAmount
}

func decreaseInventory(qty Quantity, amt Amount, inventoryVariable Inventory) (Inventory, error) {
	if len(inventoryVariable) == 0 {
		return nil, goerrors.Errorf(ErrInventoryIsEmpty, "inventory is empty")
	}

	totalQty, totalAmt := calculateTotalInventory(inventoryVariable)

	if totalQty < qty {
		return nil, fErrInsufficientQuantityInInventory(qty, totalQty)
	}

	if totalAmt < amt {
		return nil, fErrInsufficientAmountInInventory(amt, totalAmt)
	}

	// Create resultInventory slice
	var resultInventory Inventory

	var qtyAccumulator Quantity
	var amtAccumulator Amount

	remainingQty := qty

	// Process FIFO
	for _, record := range inventoryVariable {
		if record.Quantity <= remainingQty {
			// Take entire record
			remainingQty -= record.Quantity

			qtyAccumulator += record.Quantity
			amtAccumulator += record.Amount
		} else {
			// Take partial record
			price := float64(record.Amount) / float64(record.Quantity)
			newQty := record.Quantity - remainingQty
			newAmount := Amount(float64(newQty) * price)

			resultInventory = append(resultInventory, InventoryRecord{
				EntryNumber: record.EntryNumber,
				Quantity:    newQty,
				Amount:      newAmount,
			})

			qtyAccumulator += remainingQty
			amtAccumulator += Amount(float64(remainingQty) * price)

			remainingQty = 0
		}
	}

	if amtAccumulator != amt {
		return nil, goerrors.Errorf(ErrAmountMismatch, "amount mismatch: expected to enter amount = %v but got = %v", amtAccumulator, amt)
	}

	return resultInventory, nil
}

func addQuantityAndAmountOnInventory(entryNumberVariable EntryNumber, qty Quantity, amt Amount, inventoryVariable Inventory) (Inventory, error) {
	if amt == 0 && qty == 0 {
		return inventoryVariable, nil
	}

	totalQty, totalAmt := calculateTotalInventory(inventoryVariable)

	if totalAmt+amt < 0 {
		return nil, fErrInsufficientAmountInInventory(amt, totalAmt)
	}

	if totalQty+qty < 0 {
		return nil, fErrInsufficientQuantityInInventory(qty, totalQty)
	}

	return Inventory{{entryNumberVariable, totalQty + qty, totalAmt + amt}}, nil
}

func removeZeros(inventoryVariable Inventory) Inventory {
	var resultInventory Inventory
	for _, record := range inventoryVariable {
		if !(record.Amount == 0 && record.Quantity == 0) {
			resultInventory = append(resultInventory, record)
		}
	}
	return resultInventory
}

// AddToJournal adds a new accounting entry to the journal while maintaining double-entry accounting principles.
// It takes the following parameters:
//   - entry: The AccountingEntry to be added to the journal
//   - getInventoryFunction: A function to retrieve the current inventory for an account address
//   - setInventoryFunction: A function to update the inventory for an account address
//   - getLastEntryFunction: A function to get the last entry from the journal
//   - setEntryFunction: A function to save a new entry to the journal
//
// The function performs the following steps:
// 1. Retrieves current inventory for all accounts involved in the entry
// 2. Gets the last journal entry for reference
// 3. Validates and processes the double-entry accounting rules
// 4. Saves the new entry to the journal
// 5. Updates the inventory for all affected accounts
//
// Returns an error if any operation fails during the process.
func AddToJournal(entry AccountingEntry,
	getInventoryFunction GetInventory,
	setInventoryFunction SetInventory,
	getLastEntryFunction GetLastEntry,
	setEntryFunction SetEntry,
) error {

	addressAndInventory := make(AccountAddressAndInventory)
	for _, singleEntryVariable := range entry.DoubleEntry {
		inv, err := getInventoryFunction(singleEntryVariable.AccountAddress)
		if err != nil {
			return err
		}
		addressAndInventory[singleEntryVariable.AccountAddress] = inv
	}

	lastEntry, err := getLastEntryFunction()
	if err != nil {
		return err
	}

	addressAndInventory, err = CheckAndProcessDoubleEntry(lastEntry.EntryNumber, lastEntry.TimeUnix, entry, addressAndInventory)
	if err != nil {
		return err
	}

	err = setEntryFunction(entry)
	if err != nil {
		return err
	}

	for address, inv := range addressAndInventory {
		err := setInventoryFunction(address, inv)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckAllTheJournal iterates through journal entries and processes double-entry accounting
// by updating account inventories. It takes two function parameters:
//
// setInventoryFunction: A function that updates the inventory for a given address
// iterOnJournalFunction: A function that iterates through journal entries
//
// The function processes entries sequentially, checking and validating double-entry accounting rules.
// For each processed entry, it updates an in-memory map of address inventories.
// Finally, it persists all updated inventories using the setInventoryFunction.
//
// Returns an error if any operation fails during journal processing or inventory updates.
func CheckAllTheJournal(setInventoryFunction SetInventory, iterOnJournalFunction IterOnJournal) error {

	var lastEntry AccountingEntry
	addressAndInventory := make(AccountAddressAndInventory)
	for {
		entry, isFinish, err := iterOnJournalFunction()
		if err != nil {
			return err
		}

		if isFinish {
			break
		}

		addressAndInventory, err = CheckAndProcessDoubleEntry(lastEntry.EntryNumber, lastEntry.TimeUnix, entry, addressAndInventory)
		if err != nil {
			return err
		}

		lastEntry = entry
	}

	for address, inv := range addressAndInventory {
		err := setInventoryFunction(address, inv)
		if err != nil {
			return err
		}
	}

	return nil
}
