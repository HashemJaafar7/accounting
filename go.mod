module github.com/HashemJaafar7/accounting

go 1.24.3

// version should be update
require (
	github.com/HashemJaafar7/goerrors v0.1.8
	github.com/HashemJaafar7/testutils v0.1.8
)

replace (
	github.com/HashemJaafar7/goerrors => ../goerrors
	github.com/HashemJaafar7/testutils => ../testutils
)

require github.com/google/gofuzz v1.2.0 // indirect
