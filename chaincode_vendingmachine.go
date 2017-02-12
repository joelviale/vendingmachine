package main

import (
	"errors"
	"fmt"
	//"golang.org/pkg/strconv"
	"strconv"
	"strings"
	//"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//type Products struct {
//	Product []string `json:"product"`
//}

// Separator
const SEPARATOR string = "##"

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// +----------------------------+
// | Init resets all the things |
// +----------------------------+
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. Initial Balance")
	}

	stub.PutState("Total_Balance", []byte(args[0]))

	return nil, nil
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// +----------------------------------------------------------+
// | Invoke is our entry point to invoke a chaincode function |
// +----------------------------------------------------------+
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "addVMC" {
		return t.addVMC(stub, args)
	} else if function == "removeVMC" {
		return t.removeVMC(stub, args)
	} else if function == "addCSP" {
		return t.addCSP(stub, args)
	} else if function == "removeCSP" {
		return t.removeCSP(stub, args)
	} else if function == "addSupplier" {
		return t.addSupplier(stub, args)
	} else if function == "removeSupplier" {
		return t.removeSupplier(stub, args)
	} else if function == "resetBalance" {
		return t.resetBalance(stub, args)
	} else if function == "updatePercentage" {
		return t.updatePercentage(stub, args)
	} else if function == "recordTransaction" {
		return t.recordTransaction(stub, args)
	} else if function == "addESIM" {
		return t.addESIM(stub, args)
	} else if function == "activateESIM" {
		return t.activateESIM(stub, args)
	} else if function == "removeESIM" {
		return t.removeESIM(stub, args)
	} else if function == "deactivateESIM" {
		return t.deactivateESIM(stub, args)
	} else if function == "createProduct" {
		return t.createProduct(stub, args)
	} else if function == "removeProduct" {
		return t.removeProduct(stub, args)
	} else if function == "updateInventory" {
		return t.updateInventory(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// +---------------------------------------------------------+
// | createProduct - invoke function to create a new Product |
// +---------------------------------------------------------+
func (t *SimpleChaincode) createProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var productId, productName, productImg, productPrice, productQRCode, entityId string
	
	productId = args[0]
	entityId = args[1]
	productName = args[2]
	productImg = args[3]
	productPrice = args[4]
	productQRCode = args[5]
	
	// Create all the key/value pairs in the ledger
	// The first key is necessary to list all the products
	stub.PutState("Product_" + productId, []byte(productId))
	stub.PutState(productId + "_Entity", []byte(entityId))
	stub.PutState(productId + "_Name", []byte(productName))
	stub.PutState(productId + "_Image", []byte(productImg))
	stub.PutState(productId + "_Price", []byte(productPrice))
	stub.PutState(productId + "_QRCode", []byte(productQRCode))

	fmt.Println("running createProduct()")

	//var ledgerKey = "products"
	//var newProduct = Product{make([]string, 1)}
	//newProduct.Product[0] = stateJSON

	return nil, nil
}

// +-----------------------------------------------------+
// | removeProduct - invoke function to remove a Product |
// | Params - productId                                  |
// +-----------------------------------------------------+
func (t *SimpleChaincode) removeProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var productId string
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	productId = args[0]

	// Delete all the key/value pairs to the ledger
	stub.DelState("Product_" + productId)
	stub.DelState(productId + "_Entity")
	stub.DelState(productId + "_Name")
	stub.DelState(productId + "_Image")
	stub.DelState(productId + "_Price")
	stub.DelState(productId + "_QRCode")

	fmt.Println("running removeProduct()")

	return nil, nil
}

// +------------------------------------------------------------------------+
// | updateInventory - invoke function to update the inventory of an entity |
// +------------------------------------------------------------------------+
func (t *SimpleChaincode) updateInventory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// Location Id corresponds to the location within the vending machine
	// Two quantities are maintained in the ledger:
	// - the quantity for one specific location within the vending machine
	// - the total quantity of the product in the vending machine, for all locations
	var entityId, locationId, productId, quantityString string
	var currentQuantity, currentTotalQuantity, newQuantity, newTotalQuantity int
	var err error
	
	currentQuantity = 0
	currentTotalQuantity = 0

	fmt.Println("running updateInventory()")
	
	entityId = args[0]
	locationId = args[1]
	productId = args[2]
	quantityString = args[3]
	
	// Can be positive (add to inventory) or negative (remove from inventory)
	deltaQuantity, err := strconv.Atoi(quantityString)

	// Retrieve current quantity for this location and product
	// Check if there is an existing quantity
	currentQuantityBytes, err := stub.GetState("InventoryByLocation" + SEPARATOR + entityId + SEPARATOR + locationId + SEPARATOR + productId);
	//if err != nil {
	if len(currentQuantityBytes) <= 0 {
		// No current quantity
		// Would need to check if deltaQuantity is positive
		newQuantity = deltaQuantity
	} else {
		// Product is already in this location
		currentQuantity, err = strconv.Atoi(string(currentQuantityBytes))
		newQuantity = currentQuantity + deltaQuantity
	}
	
	// Do the same for total quantity
	currentTotalQuantityBytes, err := stub.GetState("InventoryByProduct" + SEPARATOR + entityId + SEPARATOR + productId);
	//if err != nil {
	if len(currentTotalQuantityBytes) <= 0 {
		// No total quantity
		// Would need to check if deltaQuantity is positive
		newTotalQuantity = deltaQuantity
	} else {
		// Product is already in this location
		currentTotalQuantity, err = strconv.Atoi(string(currentTotalQuantityBytes))
		newTotalQuantity = currentTotalQuantity + deltaQuantity
	}
	
	// Store the quantities back to the ledger or delete the entry if new quantity is zero
	// Delete the entry if the new quantity is zero
	if newQuantity <= 0 {
		stub.DelState("InventoryByLocation" + SEPARATOR + entityId + SEPARATOR + locationId + SEPARATOR + productId);
	} else {
		stub.PutState("InventoryByLocation" + SEPARATOR + entityId + SEPARATOR + locationId + SEPARATOR + productId, []byte(strconv.Itoa(newQuantity)))
	}
	// Delete the entry if the new quantity is zero
	if newTotalQuantity <= 0 {
		stub.DelState("InventoryByProduct" + SEPARATOR + entityId + SEPARATOR + productId);
	} else {
		stub.PutState("InventoryByProduct" + SEPARATOR + entityId + SEPARATOR + productId, []byte(strconv.Itoa(newTotalQuantity)))
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// +-------------------------------------------+
// | addVMC - invoke function to add a new VMC |
// +-------------------------------------------+
func (t *SimpleChaincode) addVMC(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var VMCName string
	var initialBalance float64
	var err error
	
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	VMCName = args[0]
	initialBalance, err = strconv.ParseFloat(args[1], 64)

	// Create all the key/value pairs to the ledger
	stub.PutState(VMCName + "_Balance", []byte(strconv.FormatFloat(initialBalance, 'f', -1, 64)))

	fmt.Println("running addVMC()")

	if err != nil {
		return nil, err
	}
	return nil, nil
}

// +-------------------------------------------------------+
// | removeVMC - invoke function to remove an existing VMC |
// +-------------------------------------------------------+
func (t *SimpleChaincode) removeVMC(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var VMCName string
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	VMCName = args[0]

	// Delete all the key/value pairs from the ledger
	stub.DelState(VMCName + "_Balance")

	fmt.Println("running removeVMC()")

	return nil, nil
}

// +-------------------------------------------+
// | addCSP - invoke function to add a new CSP |
// +-------------------------------------------+
func (t *SimpleChaincode) addCSP(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var CSPName string
	var percentage, initialBalance float64
	var err error
	
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}
	
	CSPName = args[0]
	percentage, err = strconv.ParseFloat(args[1], 64)
	initialBalance, err = strconv.ParseFloat(args[2], 64)

	// Create all the key/value pairs to the ledger
	stub.PutState(CSPName + "_Percentage", []byte(strconv.FormatFloat(percentage, 'f', -1, 64)))
	stub.PutState(CSPName + "_Balance", []byte(strconv.FormatFloat(initialBalance, 'f', -1, 64)))

	fmt.Println("running addCSP()")

	if err != nil {
		return nil, err
	}
	return nil, nil
}

// +-------------------------------------------------------+
// | removeCSP - invoke function to remove an existing CSP |
// +-------------------------------------------------------+
func (t *SimpleChaincode) removeCSP(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var CSPName string
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	CSPName = args[0]

	// Delete all the key/value pairs from the ledger
	stub.DelState(CSPName + "_Percentage")
	stub.DelState(CSPName + "_Balance")

	fmt.Println("running removeCSP()")

	return nil, nil
}

// +-----------------------------------------------------+
// | addSupplier - invoke function to add a new supplier |
// +-----------------------------------------------------+
func (t *SimpleChaincode) addSupplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var supplierName string
	var percentage, initialBalance float64
	var err error
	
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}
	
	supplierName = args[0]
	percentage, err = strconv.ParseFloat(args[1], 64)
	initialBalance, err = strconv.ParseFloat(args[2], 64)

	// Create all the key/value pairs to the ledger
	stub.PutState(supplierName + "_Percentage", []byte(strconv.FormatFloat(percentage, 'f', -1, 64)))
	stub.PutState(supplierName + "_Balance", []byte(strconv.FormatFloat(initialBalance, 'f', -1, 64)))

	fmt.Println("running addSupplier()")

	if err != nil {
		return nil, err
	}
	return nil, nil
}

// +-----------------------------------------------------------------+
// | removeSupplier - invoke function to remove an existing supplier |
// +-----------------------------------------------------------------+
func (t *SimpleChaincode) removeSupplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var supplierName string
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	supplierName = args[0]

	// Delete all the key/value pairs from the ledger
	stub.DelState(supplierName + "_Percentage")
	stub.DelState(supplierName + "_Balance")

	fmt.Println("running removeSupplier()")

	return nil, nil
}

// +------------------------------------------------------------------+
// | resetBalance - invoke function to reset the balance of a company |
// +------------------------------------------------------------------+
func (t *SimpleChaincode) resetBalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var companyName string
	var balance float64
	var err error
	
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	companyName = args[0]
	balance, err = strconv.ParseFloat(args[1], 64)

	stub.PutState(companyName + "_Balance", []byte(strconv.FormatFloat(balance, 'f', -1, 64)))

	fmt.Println("running resetBalance()")

	if err != nil {
		return nil, err
	}
	return nil, nil
}

// +--------------------------------------------------------------------------+
// | updatePercentage - invoke function to update the percentage of a company |
// +--------------------------------------------------------------------------+
func (t *SimpleChaincode) updatePercentage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var companyName string
	var percentage float64
	var err error
	
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	companyName = args[0]
	percentage, err = strconv.ParseFloat(args[1], 64)

	stub.PutState(companyName + "_Percentage", []byte(strconv.FormatFloat(percentage, 'f', -1, 64)))

	fmt.Println("running updatePercentage()")

	if err != nil {
		return nil, err
	}
	return nil, nil
}

// +-------------------------------------------------------------------------------------------------------------+
// | recordTransaction - invoke function to record the transaction and update the companies balances accordingly |
// +-------------------------------------------------------------------------------------------------------------+

func (t *SimpleChaincode) recordTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var supplierName, CSPName, VMCName, transactionId string
	var amountval float64
	var CSPval, VMCval, Supplierval, Totalval, CSPPercentage, SupplierPercentage float64
	var CSPAdd, VMCAdd, SupplierAdd float64
	var err error
	var json string
	//var jsonResp string

	fmt.Println("running recordTransaction()")

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5. Transaction Id, Amount and names of the 3 companies")
	}
		
	// 0. Get the amount and company names from the parameters
	transactionId = args[0]
	amountval, err = strconv.ParseFloat(args[1], 64)
	supplierName = args[2]
	CSPName = args[3]
	VMCName = args[4]
	
	// 1. Retrieve the current balances and percentages from the ledger
	CSPvalbytes, err := stub.GetState(CSPName + "_Balance")
	VMCvalbytes, err := stub.GetState(VMCName + "_Balance")
	Suppliervalbytes, err := stub.GetState(supplierName + "_Balance")
	Totalvalbytes, err := stub.GetState("Total_Balance")

	CSPval, err = strconv.ParseFloat(string(CSPvalbytes),64)
	VMCval, err = strconv.ParseFloat(string(VMCvalbytes),64)
	Supplierval, err = strconv.ParseFloat(string(Suppliervalbytes),64)
	Totalval, err = strconv.ParseFloat(string(Totalvalbytes),64)

	CSPPercentagebytes, err := stub.GetState(CSPName + "_Percentage")
	SupplierPercentagebytes, err := stub.GetState(supplierName + "_Percentage")

	CSPPercentage, err = strconv.ParseFloat(string(CSPPercentagebytes),64)
	SupplierPercentage, err = strconv.ParseFloat(string(SupplierPercentagebytes),64)
	
	// 2. Calculate the amounts that needs to be added for each company
	CSPAdd = amountval*CSPPercentage
	SupplierAdd = amountval*SupplierPercentage
	VMCAdd = (amountval - CSPAdd) - SupplierAdd
	
	// 3. Update all the balances from the new amount
	Totalval = Totalval + amountval
	CSPval = CSPval + CSPAdd
	Supplierval = Supplierval + SupplierAdd
	VMCval = VMCval + VMCAdd
	
	// 4. Write the update balances back to the ledger
	stub.PutState(CSPName + "_Balance_" + transactionId, []byte(strconv.FormatFloat(CSPval, 'f', -1, 64)))
	stub.PutState(VMCName + "_Balance_" + transactionId, []byte(strconv.FormatFloat(VMCval, 'f', -1, 64)))
	stub.PutState(supplierName + "_Balance_" + transactionId, []byte(strconv.FormatFloat(Supplierval, 'f', -1, 64)))
	
	// 5. Write the update balances back to the ledger
	stub.PutState(CSPName + "_Balance", []byte(strconv.FormatFloat(CSPval, 'f', -1, 64)))
	stub.PutState(VMCName + "_Balance", []byte(strconv.FormatFloat(VMCval, 'f', -1, 64)))
	stub.PutState(supplierName + "_Balance", []byte(strconv.FormatFloat(Supplierval, 'f', -1, 64)))
	stub.PutState("Total_Balance", []byte(strconv.FormatFloat(Totalval, 'f', -1, 64)))

	// 5. Store all the new balances associated with the transactions
	json = "\"balances\":[{\"companyName\":\"" + supplierName + "\",\"balance\":" + strconv.FormatFloat(Supplierval, 'f', -1, 64) + "},"
	json += "{\"companyName\":\"" + CSPName + "\",\"balance\":" + strconv.FormatFloat(CSPval, 'f', -1, 64) + "},"
	json += "{\"companyName\":\"" + VMCName + "\",\"balance\":" + strconv.FormatFloat(VMCval, 'f', -1, 64) + "}]"
	
	fmt.Println("recordTransaction.json stored = " + json)
	stub.PutState(transactionId, []byte(json))
		
	// 5. Return the new balances -- CANNOT!
	//jsonResp = "{\"" + supplierName + "_Balance\":\"" + strconv.FormatFloat(Supplierval, 'f', -1, 64) + "\","
	//jsonResp += "\"" + CSPName + "_Balance\":\"" + strconv.FormatFloat(CSPval, 'f', -1, 64) + "\","
	//jsonResp += "\"" + VMCName + "_Balance\":\"" + strconv.FormatFloat(VMCval, 'f', -1, 64) + "\"}"
	
	//fmt.Println("recordTransaction.jsonResp = " + jsonResp)

	if err != nil {
		fmt.Println("recordTransaction.error")
		return nil, err
	}

	return nil, nil
	//return []byte(jsonResp), nil
}

// +---------------------------------------------+
// | addESIM - invoke function to add a new eSIM |
// | Params - eSIMId, Status, Manufacturer       |
// +---------------------------------------------+
func (t *SimpleChaincode) addESIM(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var manufacturer, eSIMId, status string
	var err error
	
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}
	
	eSIMId = args[0]
	status = args[1]
	manufacturer = args[2]

	// Create all the key/value pairs to the ledger
	stub.PutState(eSIMId + "_Status", []byte(status))
	stub.PutState(eSIMId + "_Manufacturer", []byte(manufacturer))

	fmt.Println("running addESIM()")

	if err != nil {
		return nil, err
	}
	return nil, nil
}

// +-------------------------------------------------------+
// | activateESIM - invoke function to activate an eSIM    |
// | Params - eSIMId, CSPName, EndUserId, IoTId, IoTSecret |
// +-------------------------------------------------------+
func (t *SimpleChaincode) activateESIM(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eSIMId, CSPName, endUserId, IoTId, IoTSecret string
	var err error
	
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}
	
	eSIMId = args[0]
	CSPName = args[1]
	endUserId = args[2]
	IoTId = args[3]
	IoTSecret = args[4]

	// Create all the key/value pairs to the ledger
	stub.PutState(eSIMId + "_Status", []byte("Active"))
	stub.PutState(eSIMId + "_CSP", []byte(CSPName))
	stub.PutState(eSIMId + "_EndUser", []byte(endUserId))
	stub.PutState(eSIMId + "_IoTId", []byte(IoTId))
	stub.PutState(eSIMId + "_IoTSecret", []byte(IoTSecret))

	fmt.Println("running activateESIM()")

	if err != nil {
		return nil, err
	}
	return nil, nil
}

// +---------------------------------------------------------+
// | deactivateESIM - invoke function to de-activate an eSIM |
// | Params - eSIMId                                         |
// +---------------------------------------------------------+
func (t *SimpleChaincode) deactivateESIM(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eSIMId string
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	eSIMId = args[0]

	// Delete all the key/value pairs to the ledger
	stub.PutState(eSIMId + "_Status", []byte("Inactive"))
	stub.DelState(eSIMId + "_CSP")
	stub.DelState(eSIMId + "_EndUser")
	stub.DelState(eSIMId + "_IoTId")
	stub.DelState(eSIMId + "_IoTSecret")

	fmt.Println("running deactivateESIM()")

	return nil, nil
}

// +------------------------------------------------+
// | removeESIM - invoke function to remove an eSIM |
// | Params - eSIMId                                |
// +------------------------------------------------+
func (t *SimpleChaincode) removeESIM(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eSIMId string
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	eSIMId = args[0]

	// Delete all the key/value pairs to the ledger
	stub.DelState(eSIMId + "_Status")
	stub.DelState(eSIMId + "_Manufacturer")

	fmt.Println("running removeESIM()")

	return nil, nil
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// +--------------------------------------+
// | Query is our entry point for queries |
// +--------------------------------------+
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "getTransaction" {
		return t.getTransaction(stub, args)
	} else if function == "getBalance" {
		return t.getBalance(stub, args)
	} else if function == "getBalanceWithTransaction" {
		return t.getBalanceWithTransaction(stub, args)
	} else if function == "getESIM" {
		return t.getESIM(stub, args)
	} else if function == "readProduct" {
		return t.readProduct(stub, args)
	} else if function == "readAllProducts" {
		return t.readAllProducts(stub, args)
	} else if function == "getInventoryByEntityAndProduct" {
		return t.getInventoryByEntityAndProduct(stub, args)
	} else if function == "getInventoryByEntityAndLocation" {
		return t.getInventoryByEntityAndLocation(stub, args)
	} else if function == "getAllInventoryByEntity" {
		return t.getAllInventoryByEntity(stub, args)
	} else if function == "getAllInventory" {
		return t.getAllInventory(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// +------------------------------------------------------------------------------------+
// | getTransaction - query function to read the balances associated with a transaction |
// +------------------------------------------------------------------------------------+
func (t *SimpleChaincode) getTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp, transactionId string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting transaction id")
	}
	
	transactionId = args[0]
	valAsbytes, err := stub.GetState(transactionId)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// +----------------------------------------------------------------+
// | getBalance - query function to read the balance of the company |
// +----------------------------------------------------------------+
func (t *SimpleChaincode) getBalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp, companyName string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the company to get the balance")
	}
	
	companyName = args[0]
	key = companyName + "_Balance"
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// +----------------------------------------------------------------------------------------------------------------+
// | getBalanceWithTransaction - query function to read the balance of the company associated with a transaction Id |
// +----------------------------------------------------------------------------------------------------------------+
func (t *SimpleChaincode) getBalanceWithTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp, companyName, transactionId string
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting transaction Id and name of the company to get the balance")
	}
	
	transactionId = args[0]
	companyName = args[1]
	key = companyName + "_Balance_" + transactionId
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// +------------------------------------------------------------+
// | getESIM - query function to read the parameters of an eSIM |
// +------------------------------------------------------------+
func (t *SimpleChaincode) getESIM(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eSIMId, status, manufacturer, CSPName, endUserId, IoTId, IoTSecret string
	var jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	eSIMId = args[0]
	
	manufacturerBytes, err := stub.GetState(eSIMId + "_Manufacturer")
	statusBytes, err := stub.GetState(eSIMId + "_Status")
	CSPNameBytes, err := stub.GetState(eSIMId + "_CSP")
	endUserIdBytes, err := stub.GetState(eSIMId + "_EndUser")
	IoTIdBytes, err := stub.GetState(eSIMId + "_IoTId")
	IoTSecretBytes, err := stub.GetState(eSIMId + "_IoTSecret")

	manufacturer = string(manufacturerBytes)
	status = string(statusBytes)
	CSPName = string(CSPNameBytes)
	endUserId = string(endUserIdBytes)
	IoTId = string(IoTIdBytes)
	IoTSecret = string(IoTSecretBytes)
	
	jsonResp = "{\"eSIMId\":\"" + eSIMId + "\",\"status\":\"" + status + "\",\"CSP\":\"" + CSPName + "\",\"manufacturer\":\"" + manufacturer
	jsonResp += "\",\"EndUser\":\"" + endUserId + "\",\"IoTId\":\"" + IoTId + "\",\"IoTSecret\":\"" + IoTSecret + "\"}";

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get products infos\"}"
		return nil, errors.New(jsonResp)
	}

	return []byte(jsonResp), nil
}

// +---------------------------------------------+
// | readProduct - read a product in the catalog |
// +---------------------------------------------+
func (t *SimpleChaincode) readProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var productId, productName, productImg, productPrice, productQRCode, entity string
	var jsonResp string
	var err error
	
	productId = args[0]
	
	// Read attributes from the ledger
	productEntityBytes, err := stub.GetState(productId + "_Entity")
	productNameBytes, err := stub.GetState(productId + "_Name")
	productImgBytes, err := stub.GetState(productId + "_Image")
	productPriceBytes, err := stub.GetState(productId + "_Price")
	productQRCodeBytes, err := stub.GetState(productId + "_QRCode")
	
	entity = string(productEntityBytes)
	productName = string(productNameBytes)
	productImg = string(productImgBytes)
	productPrice = string(productPriceBytes)
	productQRCode = string(productQRCodeBytes)
		
	jsonResp = "["
	jsonResp += "{\"productId\":\"" + productId + "\",\"relatedEntity\":\"" + entity + "\",\"productName\":\"" + productName + "\",\"productImg\":\"" + productImg
	jsonResp += "\",\"productPrice\":\"" + productPrice + "\",\"productQRCode\":\"" + productQRCode + "\"}";
	jsonResp += "]"
	
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get product info\"}"
		return nil, errors.New(jsonResp)
	}

	return []byte(jsonResp), nil
}

// +----------------------------------------------------------------------+
// | readAllProducts - query function to read all products in the catalog |
// +----------------------------------------------------------------------+
func (t *SimpleChaincode) readAllProducts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var productId, productName, productImg, productPrice, productQRCode, entity string
	var jsonResp string
	var err error
	
	// RangeQueryState function can be invoked by a chaincode to query of a range
	// of keys in the state. Assuming the startKey and endKey are in lexical order,
	// an iterator will be returned that can be used to iterate over all keys
	// between the startKey and endKey, inclusive. The order in which keys are
	// returned by the iterator is random.

	iter, err := stub.RangeQueryState("Product_", "Product_}")
	
	jsonResp = "["
	
	var i int
	i = 0
	
	for iter.HasNext() {
		ledgerKey, productIdBytes, err := iter.Next()
		productId = string(productIdBytes)
		fmt.Println("readAllProducts found product: " + productId + "\n and ledge key: " + ledgerKey)
		if err != nil {
			err = fmt.Errorf("readAllProducts iter.Next() failed: %s", err)
			//log.Error(err)
			return nil, err
		}
		// log.Debug("readAllProducts found product: " + productId + "\n")
		//err = json.Unmarshal(assetBytes, &state)
		
		// Read attributes from the ledger
		productEntityBytes, err := stub.GetState(productId + "_Entity")
		productNameBytes, err := stub.GetState(productId + "_Name")
		productImgBytes, err := stub.GetState(productId + "_Image")
		productPriceBytes, err := stub.GetState(productId + "_Price")
		productQRCodeBytes, err := stub.GetState(productId + "_QRCode")
		
		entity = string(productEntityBytes)
		productName = string(productNameBytes)
		productImg = string(productImgBytes)
		productPrice = string(productPriceBytes)
		productQRCode = string(productQRCodeBytes)
		
		if(i>0) {
			jsonResp += ","
		}
		
		jsonResp += "{\"productId\":\"" + productId + "\",\"relatedEntity\":\"" + entity + "\",\"productName\":\"" + productName + "\",\"productImg\":\"" + productImg
		jsonResp += "\",\"productPrice\":\"" + productPrice + "\",\"productQRCode\":\"" + productQRCode + "\"}";

		i ++
	}

	jsonResp += "]"
	
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get eSIM infos\"}"
		return nil, errors.New(jsonResp)
	}

	return []byte(jsonResp), nil
}

// +---------------------------------------------------------------------------------------+
// | getInventoryByEntityAndProduct - retrieve the quantity for the entity and the product |
// +---------------------------------------------------------------------------------------+
func (t *SimpleChaincode) getInventoryByEntityAndProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var entityId, productId, quantity string
	var jsonResp string
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	entityId = args[0]
	productId = args[1]
	
	quantityBytes, err := stub.GetState("InventoryByProduct" + SEPARATOR + entityId + SEPARATOR + productId)

	quantity = string(quantityBytes)
	
	jsonResp = "{\"quantity\":\"" + quantity + "\"}";

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get product infos\"}"
		return nil, errors.New(jsonResp)
	}

	return []byte(jsonResp), nil
}

// +------------------------------------------------------------------------------------------------+
// | getInventoryByEntityAndLocation - retrieve the product and quantity for an entity and location |
// +------------------------------------------------------------------------------------------------+
func (t *SimpleChaincode) getInventoryByEntityAndLocation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var entityId, locationId, productId, quantity, keyPrefix string
	var jsonResp string
	var q int
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	entityId = args[0]
	locationId = args[1]
	
	keyPrefix = "InventoryByLocation" + SEPARATOR + entityId + SEPARATOR + locationId + SEPARATOR
	l := len(keyPrefix)
	iter, err := stub.RangeQueryState(keyPrefix, keyPrefix + "{")
	
	jsonResp = "["
	
	var i int
	i = 0
	
	for iter.HasNext() {
		ledgerKey, quantityBytes, err := iter.Next()
		if err != nil {
			err = fmt.Errorf("iter.Next() failed: %s", err)
			return nil, err
		}
		productId = ledgerKey[l:len(ledgerKey)]
		quantity = string(quantityBytes)
		fmt.Println("getInventoryByEntityAndLocation found product: " + productId + "\n and quantity: " + quantity)
		
		q, err = strconv.Atoi(quantity)
		if q > 0 {
			if i > 0 {
				jsonResp += ","
			}
			jsonResp += "{\"productId\":\"" + productId + "\",\"quantity\":\"" + quantity + "\"}";
			i ++
		}
	}

	jsonResp += "]"

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get infos\"}"
		return nil, errors.New(jsonResp)
	}

	return []byte(jsonResp), nil
}

// +----------------------------------------------------------------------------------+
// | getAllInventoryByEntity - retrieve all products and quantities for each location |
// +----------------------------------------------------------------------------------+
func (t *SimpleChaincode) getAllInventoryByEntity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var entityId, locationId, productId, quantity, keyPrefix string
	var jsonResp string
	var q int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	entityId = args[0]
	
	// Format Inventory_EntityId_LocationId_ProductId
	keyPrefix = "InventoryByLocation" + SEPARATOR + entityId + SEPARATOR
	l := len(keyPrefix)
	iter, err := stub.RangeQueryState(keyPrefix, keyPrefix + "{")
	
	jsonResp = "["
	
	var i int
	i = 0
	
	for iter.HasNext() {
		ledgerKey, quantityBytes, err := iter.Next()
		
		fmt.Println("LedgerKey = " + ledgerKey);
		if err != nil {
			err = fmt.Errorf("iter.Next() failed: %s", err)
			return nil, err
		}
		
		// Retrieve locationId from the ledger key
		locationAndProduct := ledgerKey[l:len(ledgerKey)]
		fmt.Println("locationAndProduct = " + locationAndProduct);
		j := strings.Index(locationAndProduct, SEPARATOR)
		locationId = locationAndProduct[0:j]
		fmt.Println("locationId = " + locationId);
		productId = locationAndProduct[j+len(SEPARATOR):len(locationAndProduct)]
		fmt.Println("productId = " + productId);
		quantity = string(quantityBytes)
		fmt.Println("quantity = " + quantity);
		
		fmt.Println("getAllInventoryByEntity found product: " + productId + " in location " + locationId + " with quantity: " + quantity)

		q, err = strconv.Atoi(quantity)
		if q > 0 {
			if i > 0 {
				jsonResp += ","
			}
			jsonResp += "{\"productId\":\"" + productId + "\",\"locationId\":\"" + locationId + "\",\"quantity\":\"" + quantity + "\"}";
			i ++
		}
	}

	jsonResp += "]"

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get infos\"}"
		return nil, errors.New(jsonResp)
	}

	return []byte(jsonResp), nil
}

// +-------------------------------------------------------------------------+
// | getAllInventory - retrieve all products and quantities for all entities |
// +-------------------------------------------------------------------------+
func (t *SimpleChaincode) getAllInventory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}


// +----------------------------------------------+
// | read - query function to read key/value pair |
// +----------------------------------------------+
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
