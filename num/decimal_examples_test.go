package num

import (
	"fmt"
)

// ExampleNewDecFromString demonstrates creating a decimal from a string
func ExampleNewDecFromString() {
	// Create decimal from string
	price, err := NewDecFromString("123.456")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Price: %s\n", price.String())
	fmt.Printf("Is set: %t\n", price.Isset())

	// Try with invalid string
	_, err = NewDecFromString("invalid")
	fmt.Printf("Invalid string error: %v\n", err != nil)

	// Output:
	// Price: 123.456
	// Is set: true
	// Invalid string error: true
}

// ExampleMustDecFromString demonstrates creating a decimal from a string (panics on error)
func ExampleMustDecFromString() {
	// Create decimal from valid string
	price := MustDecFromString("999.99")

	fmt.Printf("Price: %s\n", price.String())
	if f, ok := price.Float64(); ok {
		fmt.Printf("As float: %.2f\n", f)
	}

	// Output:
	// Price: 999.99
}

// ExampleNewDecFromInt demonstrates creating a decimal from an integer
func ExampleNewDecFromInt() {
	// Create decimal from integer
	quantity := NewDecFromInt(42)

	fmt.Printf("Quantity: %s\n", quantity.String())
	fmt.Printf("As int: %d\n", quantity.IntPart())
	fmt.Printf("Is zero: %t\n", quantity.IsZero())

	// Output:
	// Quantity: 42
	// As int: 42
	// Is zero: false
}

// ExampleNewDecFromFloat demonstrates creating a decimal from a float
func ExampleNewDecFromFloat() {
	// Create decimal from float
	temperature := NewDecFromFloat(36.5)

	fmt.Printf("Temperature: %s\n", temperature.String())
	if f, ok := temperature.Float64(); ok {
		fmt.Printf("As float: %.1f\n", f)
	}
	fmt.Printf("Number of decimals: %d\n", temperature.NumberOfDecimals())

	// Output:
	// Temperature: 36.5
	// As float: 36.5
	// Number of decimals: 1
}

// ExampleMustDecFromAny demonstrates creating a decimal from any supported type
func ExampleMustDecFromAny() {
	// Create decimals from different types
	fromInt := MustDecFromAny(100)
	fromFloat := MustDecFromAny(3.14159)
	fromString := MustDecFromAny("42.42")

	fmt.Printf("From int: %s\n", fromInt.String())
	fmt.Printf("From float: %s\n", fromFloat.String())
	fmt.Printf("From string: %s\n", fromString.String())

	// Output:
	// From int: 100
	// From float: 3.14159
	// From string: 42.42
}

// ExampleDec_Add demonstrates decimal addition
func ExampleDec_Add() {
	// Perform decimal addition
	price1 := MustDecFromString("123.45")
	price2 := MustDecFromString("67.89")

	total := price1.Add(price2)

	fmt.Printf("Price 1: %s\n", price1.String())
	fmt.Printf("Price 2: %s\n", price2.String())
	fmt.Printf("Total: %s\n", total.String())

	// Output:
	// Price 1: 123.45
	// Price 2: 67.89
	// Total: 191.34
}

// ExampleDec_Sub demonstrates decimal subtraction
func ExampleDec_Sub() {
	// Perform decimal subtraction
	balance := MustDecFromString("1000.00")
	withdrawal := MustDecFromString("750.25")

	remaining := balance.Sub(withdrawal)

	fmt.Printf("Initial balance: %s\n", balance.String())
	fmt.Printf("Withdrawal: %s\n", withdrawal.String())
	fmt.Printf("Remaining: %s\n", remaining.String())

	// Output:
	// Initial balance: 1000
	// Withdrawal: 750.25
	// Remaining: 249.75
}

// ExampleDec_Mul demonstrates decimal multiplication
func ExampleDec_Mul() {
	// Calculate area
	length := MustDecFromString("12.5")
	width := MustDecFromString("8.4")

	area := length.Mul(width)

	fmt.Printf("Length: %s\n", length.String())
	fmt.Printf("Width: %s\n", width.String())
	fmt.Printf("Area: %s\n", area.String())

	// Output:
	// Length: 12.5
	// Width: 8.4
	// Area: 105
}

// ExampleDec_Div demonstrates decimal division
func ExampleDec_Div() {
	// Calculate unit price
	totalCost := MustDecFromString("127.50")
	quantity := MustDecFromString("25")

	unitPrice := totalCost.Div(quantity)

	fmt.Printf("Total cost: %s\n", totalCost.String())
	fmt.Printf("Quantity: %s\n", quantity.String())
	fmt.Printf("Unit price: %s\n", unitPrice.String())

	// Output:
	// Total cost: 127.5
	// Quantity: 25
	// Unit price: 5.1
}

// ExampleDec_Percent demonstrates percentage calculations
func ExampleDec_Percent() {
	// Calculate percentage
	value := MustDecFromString("850.00")
	percentage := MustDecFromString("15") // 15%

	percentValue := value.ApplyPercent(percentage)
	finalValue := value.AddPercent(percentage)

	fmt.Printf("Original value: %s\n", value.String())
	fmt.Printf("Percentage: %s%%\n", percentage.String())
	fmt.Printf("Percentage amount: %s\n", percentValue.String())
	fmt.Printf("Final value (with percentage): %s\n", finalValue.String())

	// Output:
	// Original value: 850
	// Percentage: 15%
	// Percentage amount: 127.5
	// Final value (with percentage): 977.5
}

// ExampleDec_Round demonstrates decimal rounding
func ExampleDec_Round() {
	// Round to different precision levels
	value := MustDecFromString("123.456789")

	rounded2 := value.RoundTo(2)
	rounded4 := value.RoundTo(4)
	rounded0 := value.RoundTo(0)

	fmt.Printf("Original: %s\n", value.String())
	fmt.Printf("Rounded to 2 decimals: %s\n", rounded2.String())
	fmt.Printf("Rounded to 4 decimals: %s\n", rounded4.String())
	fmt.Printf("Rounded to 0 decimals: %s\n", rounded0.String())

	// Output:
	// Original: 123.456789
	// Rounded to 2 decimals: 123.46
	// Rounded to 4 decimals: 123.4568
	// Rounded to 0 decimals: 123
}

// ExampleDec_LessThan demonstrates decimal comparison
func ExampleDec_LessThan() {
	// Compare decimal values
	price1 := MustDecFromString("99.99")
	price2 := MustDecFromString("100.00")
	price3 := MustDecFromString("99.99")

	fmt.Printf("Price 1: %s\n", price1.String())
	fmt.Printf("Price 2: %s\n", price2.String())
	fmt.Printf("Price 3: %s\n", price3.String())

	fmt.Printf("Price 1 < Price 2: %t\n", price1.LessThan(price2))
	fmt.Printf("Price 1 > Price 2: %t\n", price1.GreaterThan(price2))
	fmt.Printf("Price 1 == Price 3: %t\n", price1.Equal(price3))
	fmt.Printf("Price 1 <= Price 2: %t\n", price1.LessThanOrEqual(price2))

	// Output:
	// Price 1: 99.99
	// Price 2: 100
	// Price 3: 99.99
	// Price 1 < Price 2: true
	// Price 1 > Price 2: false
	// Price 1 == Price 3: true
	// Price 1 <= Price 2: true
}

// ExampleDec_IsZero demonstrates zero checking
func ExampleDec_IsZero() {
	// Check if decimal is zero
	zero := Zero
	nonZero := MustDecFromString("0.01")
	alsoZero := MustDecFromString("0.00")

	fmt.Printf("Zero value: %s, IsZero: %t\n", zero.String(), zero.IsZero())
	fmt.Printf("Non-zero value: %s, IsZero: %t\n", nonZero.String(), nonZero.IsZero())
	fmt.Printf("Also zero: %s, IsZero: %t\n", alsoZero.String(), alsoZero.IsZero())

	// Output:
	// Zero value: 0, IsZero: true
	// Non-zero value: 0.01, IsZero: false
	// Also zero: 0, IsZero: true
}

// ExampleAbs demonstrates absolute value calculation
func ExampleAbs() {
	// Calculate absolute values for different types
	intVal := -42
	floatVal := -3.14

	absInt := Abs(intVal)
	absFloat := Abs(floatVal)

	fmt.Printf("Original int: %d, Absolute: %d\n", intVal, absInt)
	fmt.Printf("Original float: %.2f, Absolute: %.2f\n", floatVal, absFloat)

	// Output:
	// Original int: -42, Absolute: 42
	// Original float: -3.14, Absolute: 3.14
}
