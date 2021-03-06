package main

import (
	"github.com/spy16/sabre"
)

func main() {
	// Accept business rules from file, command-line, http request etc.
	// These rules can change as per business requirements and your
	// application doesn't have to change.
	ruleSrc := `(and (regular-user? current-user)
					 (not-blacklisted? current-user))`

	shouldDiscount, err := runDiscountingRule(ruleSrc, "bob")
	if err != nil {
		panic(err)
	}

	if shouldDiscount {
		// apply discount for the order
	} else {
		// don't apply discount
	}
}

func runDiscountingRule(rule string, user string) (bool, error) {
	// Define a scope with no bindings. (not even special forms)
	scope := sabre.NewScope(nil)

	// Define and expose your rules which ideally should have no
	// side effects.
	scope.BindGo("and", and)
	scope.BindGo("or", or)
	scope.BindGo("regular-user?", isRegularUser)
	scope.BindGo("minimum-cart-price?", isMinCartPrice)
	scope.BindGo("not-blacklisted?", isNotBlacklisted)

	// Bind current user name
	scope.BindGo("current-user", user)

	shouldDiscount, err := sabre.ReadEvalStr(scope, rule)
	return isTruthy(shouldDiscount), err
}

func isTruthy(v sabre.Value) bool {
	if v == nil || v == (sabre.Nil{}) {
		return false
	}
	if b, ok := v.(sabre.Bool); ok {
		return bool(b)
	}
	return true
}

func isNotBlacklisted(user string) bool {
	return user != "joe"
}

func isMinCartPrice(price float64) bool {
	return price >= 100
}

func isRegularUser(user string) bool {
	return user == "bob"
}

func and(rest ...bool) bool {
	if len(rest) == 0 {
		return true
	}
	result := rest[0]
	for _, r := range rest {
		result = result && r
		if !result {
			return false
		}
	}
	return true
}

func or(rest ...bool) bool {
	if len(rest) == 0 {
		return true
	}
	result := rest[0]
	for _, r := range rest {
		if result {
			return true
		}
		result = result || r
	}
	return false
}
