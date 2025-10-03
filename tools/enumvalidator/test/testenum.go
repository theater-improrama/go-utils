package test

//go:generate go run ../main.go -type=TestEnum,AnotherEnum,PrivateErrorEnum

type TestEnum int

const (
	TestEnum_1 TestEnum = iota
	TestEnum_2
	TestEnum_3
)

type AnotherEnum string

const (
	AnotherEnum_1 AnotherEnum = "1"
	AnotherEnum_2 AnotherEnum = "2"
	AnotherEnum_3 AnotherEnum = "3"
)

type PrivateErrorEnum float32

const (
	PrivateErrorEnum_1 PrivateErrorEnum = 1
	PrivateErrorEnum_2 PrivateErrorEnum = 2
	PrivateErrorEnum_3 PrivateErrorEnum = 3
)
