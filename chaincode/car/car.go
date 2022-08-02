/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// Car describes basic details of what makes up a car
type Car struct {
	CarID   string `json:"carId"`
	CarType string `json:"cartype"`
	Make    string `json:"make"`
	Model   string `json:"model"`
	Colour  string `json:"colour"`
	Dealer  string `json:"dealer"`
	Owner   string `json:"owner"`
	Status  string `json:"status"`
}

// ManufactureCar adds a new car to the world state with given details
func (s *SmartContract) ManufactureCar(ctx contractapi.TransactionContextInterface, CarData string) (string, error) {

	// Check manufacturer authorization - this sample assumes Org1 is the manufacturer with privilege to manufacture new car
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return "error", fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != "Org1MSP" {
		return "error", fmt.Errorf("client is not authorized to manufacture new car")
	}

	if len(CarData) == 0 {
		return "error", fmt.Errorf("please pass the correct Car data")
	}

	var car Car

	error := json.Unmarshal([]byte(CarData), &car)
	if error != nil {
		return "error", fmt.Errorf("failed while unmarshling car data %s", error.Error())
	}

	carAsBytes, err := json.Marshal(car)
	if err != nil {
		return "error", fmt.Errorf("failed while marshling proposal records %s", err.Error())
	}

	exists, error := s.CarExists(ctx, car.CarID)
	if error != nil {
		return "error", error
	}
	if exists {
		return "nil", fmt.Errorf("Car with CarId %s already exists", car.CarID)
	}

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(car.CarID, carAsBytes)

}

//Update Car Dealer information and Status
func (s *SmartContract) UpdateDealer(ctx contractapi.TransactionContextInterface, carId string, dealer string) (string, error) {

	// Check manufacturer authorization - this sample assumes Org1 is the manufacturer with privilege to update Dealer information
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return "error", fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != "Org1MSP" {
		return "error", fmt.Errorf("client is not authorized to update dealer information")
	}

	if len(carId) == 0 {
		return "error", fmt.Errorf("please pass the correct Car ID")
	}

	carAsBytes, err := ctx.GetStub().GetState(carId)

	if err != nil {
		return "error", fmt.Errorf("failed to get car records %s", err.Error())
	}

	if carAsBytes == nil {
		return "error", fmt.Errorf("the car %s does not exist", carId)
	}

	// overwriting original data with new data
	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	car.Dealer = dealer
	car.Status = "READY_FOR_SALE"

	carAsBytes, err = json.Marshal(car)
	if err != nil {
		return "", fmt.Errorf("failed marshal %s", err.Error())
	}

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(car.CarID, carAsBytes)

}

//Sell Car with owner information
func (s *SmartContract) SellCar(ctx contractapi.TransactionContextInterface, carId string, owner string) (string, error) {

	// Check dealer authorization - this sample assumes Org2 is the dealer with privilege to sell the car
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return "error", fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != "Org2MSP" {
		return "error", fmt.Errorf("client is not authorized to sell car")
	}

	if len(carId) == 0 {
		return "error", fmt.Errorf("please pass the correct Car ID")
	}

	carAsBytes, err := ctx.GetStub().GetState(carId)

	if err != nil {
		return "", fmt.Errorf("failed to get car records %s", err.Error())
	}

	if carAsBytes == nil {
		return "", fmt.Errorf("the car %s does not exist", carId)
	}

	// overwriting original data with new data
	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	if car.Status != "READY_FOR_SALE" {
		return "error", fmt.Errorf("Car is not on sale, please contact dealer")
	}

	if car.Status == "SOLD" {
		return "error", fmt.Errorf("car already sold. Please try purchasing other car")
	}

	car.Owner = owner
	car.Status = "SOLD"

	carAsBytes, err = json.Marshal(car)
	if err != nil {
		return "", fmt.Errorf("failed marshal %s", err.Error())
	}

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(car.CarID, carAsBytes)

}

// ChangeCarOwner updates the owner field of car with given id in world state
func (s *SmartContract) ChangeCarOwner(ctx contractapi.TransactionContextInterface, carId string, currentOwner string, newOwner string) (string, error) {
	car, err := s.QueryCar(ctx, carId)

	if err != nil {
		return "error", fmt.Errorf("error while querying car: %v", err)
	}
	if car.Owner != currentOwner {
		return "error", fmt.Errorf("current owner does not match")
	}

	car.Owner = newOwner

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(carId, carAsBytes)
}

// Car Exists returns true when Car with given ID exists in world state
func (s *SmartContract) CarExists(ctx contractapi.TransactionContextInterface, CarId string) (bool, error) {
	proposalJSON, err := ctx.GetStub().GetState(CarId)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return proposalJSON != nil, nil
}

// getHistoryforCar based on CarId
func (s *SmartContract) GetHistoryForCar(ctx contractapi.TransactionContextInterface, carId string) (string, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(carId)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}
		if bArrayMemberAlreadyWritten {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return string(buffer.Bytes()), nil
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, carId string) (*Car, error) {
	carAsBytes, err := ctx.GetStub().GetState(carId)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carId)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create car chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting car chaincode: %s", err.Error())
	}
}
