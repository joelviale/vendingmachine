package main

import (
	"errors"
	"fmt"
	//"golang.org/pkg/strconv"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Write all the initial balances to the ledger
	stub.PutState("CSP_Balance", []byte(args[0]))
	stub.PutState("VMC_Balance", []byte(args[1]))
	stub.PutState("Supplier_Balance", []byte(args[2]))
	stub.PutState("Total_Balance", []byte(args[3]))

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var amountval float64
	var CSPval, VMCval, Supplierval, Totalval float64
	
	// 1. Retrieve the current balances from the ledger
	CSPvalbytes, err := stub.GetState("CSP_Balance")
	VMCvalbytes, err := stub.GetState("VMC_Balance")
	Suppliervalbytes, err := stub.GetState("Supplier_Balance")
	Totalvalbytes, err := stub.GetState("Total_Balance")

	CSPval, err = strconv.ParseFloat(string(CSPvalbytes),64)
	VMCval, err = strconv.ParseFloat(string(VMCvalbytes),64)
	Supplierval, err = strconv.ParseFloat(string(Suppliervalbytes),64)
	Totalval, err = strconv.ParseFloat(string(Totalvalbytes),64)
		
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}
	
	// 2. Get the amount from the parameters
	amountval, err = strconv.ParseFloat(args[1], 64)
	
	// 3. Update all the balances from the new amount
	Totalval = Totalval + amountval
	CSPval = CSPval + (amountval*0.1)
	VMCval = VMCval + (amountval*0.4)
	Supplierval = Supplierval + (amountval*0.5)
	
	// 4. Write the update balances back to the ledger
	stub.PutState("CSP_Balance", []byte(strconv.FormatFloat(CSPval, 'f', -1, 64)))
	stub.PutState("VMC_Balance", []byte(strconv.FormatFloat(VMCval, 'f', -1, 64)))
	stub.PutState("Supplier_Balance", []byte(strconv.FormatFloat(Supplierval, 'f', -1, 64)))
	stub.PutState("Total_Balance", []byte(strconv.FormatFloat(Totalval, 'f', -1, 64)))

	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}


// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}
	
	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
