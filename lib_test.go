package accounting

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/HashemJaafar7/testutils"
)

func fTest[t any](actual t, expected t) {
	testutils.Test(true, false, true, 10, "v", actual, expected)
}

func fTest1[t any](actual, expected t) {
	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf("actual:\n%v\nexpected:\n%v\n", actual, expected)
	}
}

func Test_calculateTotalInventory(t *testing.T) {
	{
		a1, a2 := calculateTotalInventory([]InventoryRecord{
			{EntryNumber: 0, Quantity: 10, Amount: 100},
			{EntryNumber: 0, Quantity: 20, Amount: 200},
			{EntryNumber: 0, Quantity: 30, Amount: 300},
		})
		var e1 Quantity = 60
		var e2 Amount = 600
		fTest(a1, e1)
		fTest(a2, e2)
	}
	{
		a1, a2 := calculateTotalInventory([]InventoryRecord{
			{EntryNumber: 0, Quantity: 5, Amount: 50},
			{EntryNumber: 0, Quantity: 15, Amount: 150},
		})
		var e1 Quantity = 20
		var e2 Amount = 200
		fTest(a1, e1)
		fTest(a2, e2)
	}
	{
		a1, a2 := calculateTotalInventory([]InventoryRecord{})
		var e1 Quantity = 0
		var e2 Amount = 0
		fTest(a1, e1)
		fTest(a2, e2)
	}
	{
		a1, a2 := calculateTotalInventory([]InventoryRecord{
			{EntryNumber: 0, Quantity: 0, Amount: 0},
			{EntryNumber: 0, Quantity: 0, Amount: 0},
		})
		var e1 Quantity = 0
		var e2 Amount = 0
		fTest(a1, e1)
		fTest(a2, e2)
	}
}

func Test_GetStatus(t *testing.T) {
	// Case 1: INFLOW and positive accountAddress
	{
		a := GetStatus(INFLOW, AccountAddress(1))
		e := IsDebit(true)
		fTest(a, e)
	}

	// Case 2: INFLOW and negative accountAddress
	{
		a := GetStatus(INFLOW, AccountAddress(-1))
		e := IsDebit(false)
		fTest(a, e)
	}

	// Case 3: Non-INFLOW and positive accountAddress
	{
		a := GetStatus(WAC, AccountAddress(1))
		e := IsDebit(false)
		fTest(a, e)
	}

	// Case 4: Non-INFLOW and negative accountAddress
	{
		a := GetStatus(WAC, AccountAddress(-1))
		e := IsDebit(true)
		fTest(a, e)
	}
}

func Test_sortInventoryByEntryNumber(t *testing.T) {
	// Case 1: Inventory with different entry numbers
	{
		inv := Inventory{
			{EntryNumber: 3, Quantity: 10, Amount: 100},
			{EntryNumber: 1, Quantity: 20, Amount: 200},
			{EntryNumber: 2, Quantity: 30, Amount: 300},
		}
		sortInventoryByEntryNumber(inv)
		expected := Inventory{
			{EntryNumber: 1, Quantity: 20, Amount: 200},
			{EntryNumber: 2, Quantity: 30, Amount: 300},
			{EntryNumber: 3, Quantity: 10, Amount: 100},
		}
		fTest(inv, expected)
	}

	// Case 2: Inventory with identical entry numbers
	{
		inv := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 100},
			{EntryNumber: 1, Quantity: 20, Amount: 200},
		}
		sortInventoryByEntryNumber(inv)
		expected := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 100},
			{EntryNumber: 1, Quantity: 20, Amount: 200},
		}
		fTest(inv, expected)
	}

	// Case 3: Empty inventory
	{
		inv := Inventory{}
		sortInventoryByEntryNumber(inv)
		expected := Inventory{}
		fTest(inv, expected)
	}

	// Case 4: Inventory with a single element
	{
		inv := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 100},
		}
		sortInventoryByEntryNumber(inv)
		expected := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 100},
		}
		fTest(inv, expected)
	}
}

func Test_addAmountOnInventory(t *testing.T) {
	// Case 1: Empty inventory
	{
		entryNumberVariable := EntryNumber(1)
		quantityVariable := Quantity(100)
		inventoryVariable := Inventory{}
		expected := Inventory{{entryNumberVariable, quantityVariable, 0}}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, quantityVariable, Amount(0), inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 2: Inventory with zero total amount
	{
		entryNumberVariable := EntryNumber(2)
		quantityVariable := Quantity(50)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 100, Amount: 0},
			{EntryNumber: 2, Quantity: 200, Amount: 0},
		}
		expected := Inventory{{entryNumberVariable, 350, 0}}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, quantityVariable, Amount(0), inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 3: Normal case with non-zero amount and quantity
	{
		entryNumberVariable := EntryNumber(3)
		quantityVariable := Quantity(60)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 100},
			{EntryNumber: 2, Quantity: 20, Amount: 200},
		}
		expected := Inventory{
			{EntryNumber: 3, Quantity: 90, Amount: 300},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, quantityVariable, Amount(0), inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 4: Negative quantity that would result in negative total
	{
		entryNumberVariable := EntryNumber(4)
		quantityVariable := Quantity(-150)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 100, Amount: 100},
		}
		var expected Inventory
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, quantityVariable, Amount(0), inventoryVariable)
		expectedErr := "ErrInsufficientQuantityInInventory : You want to withdraw quantity = 150 but you do not have enough quantity because your total quantity = 100"
		fTest(actual, expected)
		fTest(actualErr.Error(), expectedErr)
	}

	// Case 5: Single record in inventory
	{
		entryNumberVariable := EntryNumber(5)
		quantityVariable := Quantity(30)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 50, Amount: 500},
		}
		expected := Inventory{
			{EntryNumber: 5, Quantity: 80, Amount: 500},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, quantityVariable, Amount(0), inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 6: Zero quantity addition
	{
		entryNumberVariable := EntryNumber(6)
		quantityVariable := Quantity(0)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 100},
		}
		expected := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 100},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, quantityVariable, Amount(0), inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 1: Successful decrease with sufficient quantity and amount
	{
		qty := Quantity(-5)
		amt := Amount(-50)
		inventoryVariable := Inventory{{EntryNumber: 1, Quantity: 10, Amount: 100}}
		expected := Inventory{{EntryNumber: 1, Quantity: 5, Amount: 50}}
		actual, actualErr := addQuantityAndAmountOnInventory(1, qty, amt, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 2: Insufficient quantity
	{
		qty := Quantity(-15)
		amt := Amount(-50)
		inventoryVariable := Inventory{{EntryNumber: 1, Quantity: 10, Amount: 100}}
		var expected Inventory
		actual, actualErr := addQuantityAndAmountOnInventory(1, qty, amt, inventoryVariable)
		expectedErr := fmt.Errorf("ErrInsufficientQuantityInInventory : You want to withdraw quantity = 15 but you do not have enough quantity because your total quantity = 10")
		fTest(actual, expected)
		fTest(actualErr.Error(), expectedErr.Error())
	}

	// Case 3: Insufficient amount
	{
		qty := Quantity(-5)
		amt := Amount(-150)
		inventoryVariable := Inventory{{EntryNumber: 1, Quantity: 10, Amount: 100}}
		var expected Inventory
		actual, actualErr := addQuantityAndAmountOnInventory(1, qty, amt, inventoryVariable)
		expectedErr := fmt.Errorf("ErrInsufficientAmountInInventory : You want to withdraw amount = 150 but you do not have enough amount because your total amount = 100")
		fTest(actual, expected)
		fTest(actualErr.Error(), expectedErr.Error())
	}

	// Case 4: Exact match of quantity and amount
	{
		qty := Quantity(-10)
		amt := Amount(-100)
		inventoryVariable := Inventory{{EntryNumber: 1, Quantity: 10, Amount: 100}}
		expected := Inventory{{EntryNumber: 1, Quantity: 0, Amount: 0}}
		actual, actualErr := addQuantityAndAmountOnInventory(1, qty, amt, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 5: Zero quantity and amount decrease
	{
		qty := Quantity(-0)
		amt := Amount(-0)
		inventoryVariable := Inventory{{EntryNumber: 1, Quantity: 10, Amount: 100}}
		expected := Inventory{{EntryNumber: 1, Quantity: 10, Amount: 100}}
		actual, actualErr := addQuantityAndAmountOnInventory(1, qty, amt, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 1: Empty inventory
	{
		entryNumberVariable := EntryNumber(1)
		amountVariable := Amount(100)
		inventoryVariable := Inventory{}
		expected := Inventory{{entryNumberVariable, 0, amountVariable}}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 2: Inventory with zero total quantity
	{
		entryNumberVariable := EntryNumber(2)
		amountVariable := Amount(50)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 0, Amount: 100},
			{EntryNumber: 2, Quantity: 0, Amount: 200},
		}
		expected := Inventory{{entryNumberVariable, 0, 350}}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 3: Inventory with non-zero total quantity
	{
		entryNumberVariable := EntryNumber(3)
		amountVariable := Amount(60)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 100},
			{EntryNumber: 2, Quantity: 20, Amount: 200},
		}
		expected := Inventory{
			{EntryNumber: 3, Quantity: 30, Amount: 360},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 4: Single record in inventory
	{
		entryNumberVariable := EntryNumber(4)
		amountVariable := Amount(30)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 5, Amount: 50},
		}
		expected := Inventory{
			{EntryNumber: 4, Quantity: 5, Amount: 80},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 5: Inventory with mixed quantities
	{
		entryNumberVariable := EntryNumber(5)
		amountVariable := Amount(90)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 0, Amount: 0},
			{EntryNumber: 2, Quantity: 10, Amount: 100},
			{EntryNumber: 3, Quantity: 20, Amount: 200},
		}
		expected := Inventory{
			{EntryNumber: 5, Quantity: 30, Amount: 390},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 6: Inventory with zero quantities
	{
		entryNumberVariable := EntryNumber(6)
		amountVariable := Amount(90)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 0, Amount: 0},
		}
		expected := Inventory{
			{EntryNumber: 6, Quantity: 0, Amount: 90},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 7: Inventory with negative amount
	{
		entryNumberVariable := EntryNumber(7)
		amountVariable := Amount(-90)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 0, Amount: 0},
			{EntryNumber: 2, Quantity: 10, Amount: 100},
			{EntryNumber: 3, Quantity: 20, Amount: 200},
		}
		expected := Inventory{
			{EntryNumber: 7, Quantity: 30, Amount: 210},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 8: Subtracting a negative amount from an inventory with non-zero quantities
	{
		entryNumberVariable := EntryNumber(8)
		amountVariable := Amount(-50)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 100},
			{EntryNumber: 2, Quantity: 20, Amount: 200},
		}
		expected := Inventory{
			{EntryNumber: 8, Quantity: 30, Amount: 250},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 9: Subtracting a negative amount from an inventory with zero quantities
	{
		entryNumberVariable := EntryNumber(9)
		amountVariable := Amount(-50)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 0, Amount: 100},
			{EntryNumber: 2, Quantity: 0, Amount: 200},
		}
		expected := Inventory{
			{EntryNumber: 9, Quantity: 0, Amount: 250},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 10: Subtracting a negative amount from an inventory with mixed quantities
	{
		entryNumberVariable := EntryNumber(10)
		amountVariable := Amount(-100)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 0, Amount: 50},
			{EntryNumber: 2, Quantity: 10, Amount: 100},
			{EntryNumber: 3, Quantity: 20, Amount: 200},
		}
		expected := Inventory{
			{EntryNumber: 10, Quantity: 30, Amount: 250},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 11: Subtracting a negative amount from an inventory with a single record
	{
		entryNumberVariable := EntryNumber(11)
		amountVariable := Amount(-30)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 5, Amount: 50},
		}
		expected := Inventory{
			{EntryNumber: 11, Quantity: 5, Amount: 20},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 12: Subtracting a negative amount from an inventory with zero amounts
	{
		entryNumberVariable := EntryNumber(12)
		amountVariable := Amount(-50)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 0},
			{EntryNumber: 2, Quantity: 20, Amount: 0},
		}
		var expected Inventory
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		expectedErr := "ErrInsufficientAmountInInventory : You want to withdraw amount = 50 but you do not have enough amount because your total amount = 0"
		fTest(actual, expected)
		fTest(actualErr.Error(), expectedErr)
	}

	// Case 13: Subtracting a negative amount from an empty inventory
	{
		entryNumberVariable := EntryNumber(13)
		amountVariable := Amount(-50)
		inventoryVariable := Inventory{}
		var expected Inventory
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		expectedErr := "ErrInsufficientAmountInInventory : You want to withdraw amount = 50 but you do not have enough amount because your total amount = 0"
		fTest(actual, expected)
		fTest(actualErr.Error(), expectedErr)
	}

	// Case 14: Subtracting a negative amount from an inventory with mixed quantities
	{
		entryNumberVariable := EntryNumber(10)
		amountVariable := Amount(-100)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 0, Amount: 50},
			{EntryNumber: 2, Quantity: 1000, Amount: 100},
			{EntryNumber: 3, Quantity: 25, Amount: 200},
			{EntryNumber: 4, Quantity: 25, Amount: 1},
		}
		expected := Inventory{
			{EntryNumber: 10, Quantity: 1050, Amount: 251},
		}
		actual, actualErr := addQuantityAndAmountOnInventory(entryNumberVariable, Quantity(0), amountVariable, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

}

func Test_decreaseInventory(t *testing.T) {
	// Case 1: Sufficient inventory for the requested quantity and amount
	{
		qty := Quantity(10)
		amt := Amount(100)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 5, Amount: 50},
			{EntryNumber: 2, Quantity: 5, Amount: 50},
		}
		var expected Inventory
		actual, actualErr := decreaseInventory(qty, amt, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 2: Insufficient quantity in inventory
	{
		qty := Quantity(15)
		amt := Amount(100)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 5, Amount: 50},
			{EntryNumber: 2, Quantity: 5, Amount: 50},
		}
		var expected Inventory
		actual, actualErr := decreaseInventory(qty, amt, inventoryVariable)
		expectedErr := "ErrInsufficientQuantityInInventory : You want to withdraw quantity = 15 but you do not have enough quantity because your total quantity = 10"
		fTest(actual, expected)
		fTest(actualErr.Error(), expectedErr)
	}

	// Case 3: Insufficient amount in inventory
	{
		qty := Quantity(10)
		amt := Amount(150)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 5, Amount: 50},
			{EntryNumber: 2, Quantity: 5, Amount: 50},
		}
		var expected Inventory
		actual, actualErr := decreaseInventory(qty, amt, inventoryVariable)
		expectedErr := "ErrInsufficientAmountInInventory : You want to withdraw amount = 150 but you do not have enough amount because your total amount = 100"
		fTest(actual, expected)
		fTest(actualErr.Error(), expectedErr)
	}

	// Case 4: Partial consumption of inventory
	{
		qty := Quantity(7)
		amt := Amount(70)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 5, Amount: 50},
			{EntryNumber: 2, Quantity: 5, Amount: 50},
		}
		expected := Inventory{
			{EntryNumber: 2, Quantity: 3, Amount: 30},
		}
		actual, actualErr := decreaseInventory(qty, amt, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 5: Exact match of quantity and amount
	{
		qty := Quantity(10)
		amt := Amount(100)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 100},
		}
		var expected Inventory
		actual, actualErr := decreaseInventory(qty, amt, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 6: Empty inventory
	{
		qty := Quantity(5)
		amt := Amount(50)
		inventoryVariable := Inventory{}
		var expected Inventory
		actual, actualErr := decreaseInventory(qty, amt, inventoryVariable)
		var expectedErr error = fmt.Errorf("ErrInventoryIsEmpty : inventory is empty")
		fTest(actual, expected)
		fTest(actualErr.Error(), expectedErr.Error())
	}

	// Case 7: Quantity mismatch after processing
	{
		qty := Quantity(10)
		amt := Amount(100)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 5, Amount: 50},
			{EntryNumber: 2, Quantity: 5, Amount: 60},
		}
		var expected Inventory
		actual, actualErr := decreaseInventory(qty, amt, inventoryVariable)
		var expectedErr error = fmt.Errorf("ErrAmountMismatch : amount mismatch: expected to enter amount = 110 but got = 100")
		fTest(actual, expected)
		fTest(actualErr.Error(), expectedErr.Error())
	}

	// Case 8: Single record with partial consumption
	{
		qty := Quantity(3)
		amt := Amount(30)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 5, Amount: 50},
		}
		expected := Inventory{
			{EntryNumber: 1, Quantity: 2, Amount: 20},
		}
		actual, actualErr := decreaseInventory(qty, amt, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 9: Single record with zero quantity
	{
		qty := Quantity(3)
		amt := Amount(30)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 0, Amount: 50},
		}
		var expected Inventory
		actual, actualErr := decreaseInventory(qty, amt, inventoryVariable)
		var expectedErr error = fmt.Errorf("ErrInsufficientQuantityInInventory : You want to withdraw quantity = 3 but you do not have enough quantity because your total quantity = 0")
		fTest(actual, expected)
		fTest(actualErr.Error(), expectedErr.Error())
	}

	// Case 10: Sufficient inventory for the requested quantity and amount
	{
		qty := Quantity(10)
		amt := Amount(100)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 5, Amount: 50},
			{EntryNumber: 2, Quantity: 20, Amount: 200},
		}
		var expected Inventory = Inventory{
			{EntryNumber: 2, Quantity: 15, Amount: 150},
		}
		actual, actualErr := decreaseInventory(qty, amt, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}

	// Case 11: Sufficient inventory for the requested quantity and amount
	{
		qty := Quantity(10)
		amt := Amount(75)
		inventoryVariable := Inventory{
			{EntryNumber: 1, Quantity: 5, Amount: 50},
			{EntryNumber: 2, Quantity: 20, Amount: 100},
		}
		var expected Inventory = Inventory{
			{EntryNumber: 2, Quantity: 15, Amount: 75},
		}
		actual, actualErr := decreaseInventory(qty, amt, inventoryVariable)
		var expectedErr error
		fTest(actual, expected)
		fTest(actualErr, expectedErr)
	}
}

func Test_removeZeros(t *testing.T) {
	// Case 1: Inventory with no zero records
	{
		inv := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 100},
			{EntryNumber: 2, Quantity: 20, Amount: 200},
		}
		expected := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 100},
			{EntryNumber: 2, Quantity: 20, Amount: 200},
		}
		actual := removeZeros(inv)
		fTest(actual, expected)
	}

	// Case 2: Inventory with some zero records
	{
		inv := Inventory{
			{EntryNumber: 1, Quantity: 0, Amount: 0},
			{EntryNumber: 2, Quantity: 20, Amount: 200},
			{EntryNumber: 3, Quantity: 0, Amount: 0},
		}
		expected := Inventory{
			{EntryNumber: 2, Quantity: 20, Amount: 200},
		}
		actual := removeZeros(inv)
		fTest(actual, expected)
	}

	// Case 3: Empty inventory
	{
		inv := Inventory{}
		var expected Inventory
		actual := removeZeros(inv)
		fTest(actual, expected)
	}

	// Case 4: All zero records
	{
		inv := Inventory{
			{EntryNumber: 1, Quantity: 0, Amount: 0},
			{EntryNumber: 2, Quantity: 0, Amount: 0},
		}
		var expected Inventory
		actual := removeZeros(inv)
		fTest(actual, expected)
	}

	// Case 5: Mixed zero and non-zero values
	{
		inv := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 0},
			{EntryNumber: 2, Quantity: 0, Amount: 100},
			{EntryNumber: 3, Quantity: 0, Amount: 0},
		}
		expected := Inventory{
			{EntryNumber: 1, Quantity: 10, Amount: 0},
			{EntryNumber: 2, Quantity: 0, Amount: 100},
		}
		actual := removeZeros(inv)
		fTest(actual, expected)
	}
}

func Test_AddToJournal(t *testing.T) {
	var myInv AccountAddressAndInventory
	var myEntries map[EntryNumber]AccountingEntry

	myInv = make(AccountAddressAndInventory)
	myEntries = make(map[EntryNumber]AccountingEntry)

	getInv := func(key AccountAddress) (Inventory, error) {
		inv, ok := myInv[key]
		if !ok {
			// return nil, fmt.Errorf("address %v not found", key)
		}
		return inv, nil
	}
	setInv := func(key AccountAddress, value Inventory) error {
		myInv[key] = value
		return nil
	}
	getEnt := func() (AccountingEntry, error) {
		var max AccountingEntry
		for k, v := range myEntries {
			if k > max.EntryNumber {
				max = v
			}
		}

		return max, nil
	}
	setEntry := func(value AccountingEntry) error {
		myEntries[value.EntryNumber] = value
		return nil
	}

	myInv[1] = Inventory{{0, 0, 0}}
	myInv[-1] = Inventory{{0, 0, 0}}

	// Case 1: Basic valid entry
	{
		entry := AccountingEntry{
			EntryNumber: 1,
			TimeUnix:    1000,
			DoubleEntry: DoubleEntry{
				{INFLOW, 1, 10, 100},
				{INFLOW, -1, 10, 100},
			},
		}

		actErr := AddToJournal(entry, getInv, setInv, getEnt, setEntry)
		var expErr error
		fTest(actErr, expErr)
	}

	// Case 2: Invalid time (earlier than previous entry)
	{
		entry := AccountingEntry{
			EntryNumber: 2,
			TimeUnix:    900, // Earlier than previous entry
			DoubleEntry: DoubleEntry{
				{INFLOW, 1, 10, 100},
				{WAC, -1, 10, 100},
			},
		}

		actErr := AddToJournal(entry, getInv, setInv, getEnt, setEntry)
		var expErr error = fmt.Errorf("ErrTimeShouldBeBigger : time should be bigger")
		fTest(actErr.Error(), expErr.Error())
	}

	// Case 3: Invalid entry number sequence
	{
		entry := AccountingEntry{
			EntryNumber: 4, // Skips entry number 3
			TimeUnix:    1100,
			DoubleEntry: DoubleEntry{
				{INFLOW, 1, 10, 100},
				{WAC, -1, 10, 100},
			},
		}

		actErr := AddToJournal(entry, getInv, setInv, getEnt, setEntry)
		var expErr error = fmt.Errorf("ErrEntryNumberShouldBeBiggerByOne : entry number should be bigger by one from the last entry number")
		fTest(actErr.Error(), expErr.Error())
	}

	// Case 4: Duplicate accounts in entry
	{
		entry := AccountingEntry{
			EntryNumber: 2,
			TimeUnix:    1100,
			DoubleEntry: DoubleEntry{
				{INFLOW, 1, 10, 100},
				{WAC, 1, 10, 100}, // Same account
			},
		}

		actErr := AddToJournal(entry, getInv, setInv, getEnt, setEntry)
		var expErr error = fmt.Errorf("ErrDuplicateAccountInEntry : duplicate account address 1 in entry")
		fTest(actErr.Error(), expErr.Error())
	}

	// Case 5: Single entry (invalid)
	{
		entry := AccountingEntry{
			EntryNumber: 2,
			TimeUnix:    1100,
			DoubleEntry: DoubleEntry{
				{INFLOW, 1, 10, 100},
			},
		}

		actErr := AddToJournal(entry, getInv, setInv, getEnt, setEntry)
		var expErr error = fmt.Errorf("ErrEntryMustHaveAtLeast_2Entries : entry must have at least 2 entries")
		fTest(actErr.Error(), expErr.Error())
	}

	// Case 6: Unbalanced debits and credits
	{
		entry := AccountingEntry{
			EntryNumber: 2,
			TimeUnix:    1100,
			DoubleEntry: DoubleEntry{
				{INFLOW, 1, 10, 100},
				{WAC, -1, 10, 90}, // Different amounts
			},
		}

		actErr := AddToJournal(entry, getInv, setInv, getEnt, setEntry)
		var expErr error = fmt.Errorf("ErrSumOfAmountsIsNotZeroAndDebitMoreThanCredit : debit not equal credit and debit = 190 , credit = 0 and debit-credit = 190")
		fTest(actErr.Error(), expErr.Error())
	}

	// Case 7: Account not found
	{
		entry := AccountingEntry{
			EntryNumber: 2,
			TimeUnix:    1100,
			DoubleEntry: DoubleEntry{
				{INFLOW, 99, 10, 100}, // Non-existent account
				{WAC, -1, 10, 100},
			},
		}

		actErr := AddToJournal(entry, getInv, setInv, getEnt, setEntry)
		var expErr error = fmt.Errorf("ErrSumOfAmountsIsNotZeroAndDebitMoreThanCredit : debit not equal credit and debit = 200 , credit = 0 and debit-credit = 200")
		fTest(actErr.Error(), expErr.Error())
	}

	// Case 8:
	{
		entry := AccountingEntry{
			EntryNumber: 1,
			TimeUnix:    1000,
			DoubleEntry: DoubleEntry{
				{INFLOW, 1, -10, 100},
				{INFLOW, -1, 10, 100},
			},
		}

		actErr := AddToJournal(entry, getInv, setInv, getEnt, setEntry)
		var expErr error = fmt.Errorf("ErrEntryNumberShouldBeBiggerByOne : entry number should be bigger by one from the last entry number")
		fTest(actErr.Error(), expErr.Error())
	}

	// Case 9:
	{
		entry := AccountingEntry{
			EntryNumber: 2,
			TimeUnix:    1000,
			DoubleEntry: DoubleEntry{
				{INFLOW, 1, 10, 100},
				{INFLOW, -1, 10, 100},
			},
		}

		actErr := AddToJournal(entry, getInv, setInv, getEnt, setEntry)
		var expErr error
		fTest(actErr, expErr)
	}

	// Case 10:
	{
		entry := AccountingEntry{
			EntryNumber: 3,
			TimeUnix:    1000,
			DoubleEntry: DoubleEntry{
				{9, 1, 10, 100},
				{INFLOW, -1, 10, 100},
			},
		}

		actErr := AddToJournal(entry, getInv, setInv, getEnt, setEntry)
		var expErr error = fmt.Errorf("ErrTheCostFlowTypeIsWrong : the cost flow type is wrong")
		fTest(actErr.Error(), expErr.Error())
	}

	// Case 11:
	{
		entry := AccountingEntry{
			EntryNumber: 3,
			TimeUnix:    1000,
			DoubleEntry: DoubleEntry{
				{INFLOW, 1, -10, 100},
				{INFLOW, -1, 10, 100},
			},
		}

		actErr := AddToJournal(entry, getInv, setInv, getEnt, setEntry)
		var expErr error = fmt.Errorf("ErrTheQuantityAndAmountShouldBeBothPositive : the quantity and amount should be both positive for account address 1")
		fTest(actErr.Error(), expErr.Error())
	}

	// Case 12:
	{
		entry := AccountingEntry{
			EntryNumber: 3,
			TimeUnix:    1000,
			DoubleEntry: DoubleEntry{
				{INFLOW, 1, 10, 100},
				{WAC, 2, 10, 100},
			},
		}

		actErr := AddToJournal(entry, getInv, setInv, getEnt, setEntry)
		var expErr error = fmt.Errorf("ErrInsufficientQuantityInInventory : You want to withdraw quantity = 10 but you do not have enough quantity because your total quantity = 0")
		fTest(actErr.Error(), expErr.Error())
	}

	i := EntryNumber(1)
	iterOnJournalFunction := func() (AccountingEntry, bool, error) {
		a, ok := myEntries[i]
		i++
		return a, ok, nil
	}

	myInvExpected := myInv
	err := CheckAllTheJournal(setInv, iterOnJournalFunction)
	fTest(err, nil)
	fTest(myInv, myInvExpected)
}

func TestCheckAndProcessDoubleEntry(t *testing.T) {
	tests := []struct {
		name              string
		i_lastEntryNumber EntryNumber
		i_lastTimeUnix    TimeUnix
		i_entry           AccountingEntry
		i_inv             AccountAddressAndInventory
		o_inv             AccountAddressAndInventory
		o_err             error
	}{
		{
			name:              "ErrEntryNumberShouldBeBiggerByOne",
			i_lastEntryNumber: EntryNumber(0),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 0,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{INFLOW, 1, 10, 100},
					{WAC, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{0, 0, 0}},
				-1: Inventory{{0, 0, 0}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrEntryNumberShouldBeBiggerByOne : entry number should be bigger by one from the last entry number"),
		},
		{
			name:              "ErrEntryNumberShouldBeBiggerByOne",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 3,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{INFLOW, 1, 10, 100},
					{WAC, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{0, 0, 0}},
				-1: Inventory{{0, 0, 0}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrEntryNumberShouldBeBiggerByOne : entry number should be bigger by one from the last entry number"),
		},
		{
			name:              "ErrTimeShouldBeBigger",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    900,
				DoubleEntry: DoubleEntry{
					{INFLOW, 1, 10, 100},
					{WAC, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{0, 0, 0}},
				-1: Inventory{{0, 0, 0}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrTimeShouldBeBigger : time should be bigger"),
		},
		{
			name:              "ErrEntryMustHaveAtLeast_2Entries",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{INFLOW, 1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{0, 0, 0}},
				-1: Inventory{{0, 0, 0}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrEntryMustHaveAtLeast_2Entries : entry must have at least 2 entries"),
		},
		{
			name:              "ErrTheCostFlowTypeIsWrong",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{CostFlowType(99), 1, 10, 100},
					{WAC, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{0, 0, 0}},
				-1: Inventory{{0, 0, 0}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrTheCostFlowTypeIsWrong : the cost flow type is wrong"),
		},
		{
			name:              "ErrTheQuantityAndAmountShouldBeBothPositive",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{INFLOW, 1, -10, 100},
					{WAC, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{0, 0, 0}},
				-1: Inventory{{0, 0, 0}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrTheQuantityAndAmountShouldBeBothPositive : the quantity and amount should be both positive for account address 1"),
		},
		{
			name:              "ErrDuplicateAccountInEntry",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{INFLOW, 1, 10, 100},
					{WAC, 1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{0, 0, 0}},
				-1: Inventory{{0, 0, 0}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrDuplicateAccountInEntry : duplicate account address 1 in entry"),
		},
		{
			name:              "ErrSumOfAmountsIsNotZeroAndDebitMoreThanCredit",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{INFLOW, 1, 10, 100},
					{WAC, -1, 10, 50},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{0, 0, 0}},
				-1: Inventory{{0, 0, 0}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrSumOfAmountsIsNotZeroAndDebitMoreThanCredit : debit not equal credit and debit = 150 , credit = 0 and debit-credit = 150"),
		},
		{
			name:              "ErrInventoryNotFoundForAccountAddress",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{INFLOW, 1, 10, 100},
					{WAC, 2, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1: Inventory{{0, 0, 0}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrInventoryNotFoundForAccountAddress : inventory not found for account address 2"),
		},
		{
			name:              "valid",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{INFLOW, 1, 10, 100},
					{INFLOW, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{0, 0, 0}},
				-1: Inventory{{0, 0, 0}},
			},
			o_inv: AccountAddressAndInventory{
				1:  Inventory{{2, 10, 100}},
				-1: Inventory{{2, 10, 100}},
			},
			o_err: nil,
		},
		{
			name:              "ErrTheQuantityAndAmountShouldBeBothPositive",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{WAC, 1, 10, -100},
					{WAC, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{0, 0, 0}},
				-1: Inventory{{0, 0, 0}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrTheQuantityAndAmountShouldBeBothPositive : the quantity and amount should be both positive for account address 1"),
		},
		{
			name:              "ErrAmountMismatch",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{WAC, 1, 10, 100},
					{WAC, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrAmountMismatch : amount mismatch: expected to enter amount = 83.83064516129032 but got = 100"),
		},
		{
			name:              "ErrAmountMismatch",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{FIFO, 1, 10, 100},
					{FIFO, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrAmountMismatch : amount mismatch: expected to enter amount = 18 but got = 100"),
		},
		{
			name:              "ErrAmountMismatch",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{LIFO, 1, 10, 100},
					{LIFO, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrAmountMismatch : amount mismatch: expected to enter amount = 128.83116883116884 but got = 100"),
		},
		{
			name:              "ErrAmountMismatch",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{HIFO, 1, 10, 100},
					{HIFO, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrAmountMismatch : amount mismatch: expected to enter amount = 972.4155844155844 but got = 100"),
		},
		{
			name:              "ErrAmountMismatch",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{LOFO, 1, 10, 100},
					{LOFO, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrAmountMismatch : amount mismatch: expected to enter amount = 1.6363636363636362 but got = 100"),
		},
		{
			name:              "valid",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{NONE, 1, 10, 100},
					{NONE, -1, 10, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: AccountAddressAndInventory{
				1:  Inventory{{2, 238, 1979}},
				-1: Inventory{{2, 238, 1979}},
			},
			o_err: nil,
		},
		{
			name:              "ErrQuantityAndAmountShouldBothBeDebitOrCredit",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{WAC, 1, 0, 0},
					{WAC, -1, 0, 0},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrQuantityAndAmountShouldBothBeDebitOrCredit : quantity and amount should both be debit or credit for account address 1"),
		},
		{
			name:              "valid",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{INFLOW, 1, 0, 100},
					{INFLOW, -1, 0, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: AccountAddressAndInventory{
				1:  Inventory{{2, 248, 2179}},
				-1: Inventory{{2, 248, 2179}},
			},
			o_err: nil,
		},
		{
			name:              "valid",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{INFLOW, 1, 10, 0},
					{INFLOW, -1, 10, 0},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: AccountAddressAndInventory{
				1:  Inventory{{2, 258, 2079}},
				-1: Inventory{{2, 258, 2079}},
			},
			o_err: nil,
		},
		{
			name:              "ErrYouShouldUseCostFlowTypeNONEIfYouHaveQuantityOrAmountZero",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{WAC, 1, 0, 100},
					{WAC, -1, 0, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrYouShouldUseCostFlowTypeNONEIfYouHaveQuantityOrAmountZero : you should to use cost flow type NONE because your quantity or amount is zero for account address 1"),
		},
		{
			name:              "ErrYouShouldUseCostFlowTypeNONEIfYouHaveQuantityOrAmountZero",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{WAC, 1, 10, 0},
					{WAC, -1, 10, 0},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: nil,
			o_err: fmt.Errorf("ErrYouShouldUseCostFlowTypeNONEIfYouHaveQuantityOrAmountZero : you should to use cost flow type NONE because your quantity or amount is zero for account address 1"),
		},
		{
			name:              "valid",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{NONE, 1, 0, 100},
					{NONE, -1, 0, 100},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: AccountAddressAndInventory{
				1:  Inventory{{2, 248, 1979}},
				-1: Inventory{{2, 248, 1979}},
			},
			o_err: nil,
		},
		{
			name:              "valid",
			i_lastEntryNumber: EntryNumber(1),
			i_lastTimeUnix:    TimeUnix(1000),
			i_entry: AccountingEntry{
				EntryNumber: 2,
				TimeUnix:    1100,
				DoubleEntry: DoubleEntry{
					{NONE, 1, 10, 0},
					{NONE, -1, 10, 0},
				},
			},
			i_inv: AccountAddressAndInventory{
				1:  Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
				-1: Inventory{{1, 50, 90}, {2, 5, 908}, {3, 61, 80}, {7, 77, 992}, {6, 55, 9}},
			},
			o_inv: AccountAddressAndInventory{
				1:  Inventory{{2, 238, 2079}},
				-1: Inventory{{2, 238, 2079}},
			},
			o_err: nil,
		},
	}
	for _, tt := range tests {
		got, err := CheckAndProcessDoubleEntry(tt.i_lastEntryNumber, tt.i_lastTimeUnix, tt.i_entry, tt.i_inv)

		fmt.Println(tt.name)
		fTest1(got, tt.o_inv)
		if err != nil {
			fTest1(err.Error(), tt.o_err.Error())
		} else {
			fTest1(err, tt.o_err)
		}
	}
}
