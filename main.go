package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Product struct {
	ProductID   string `json:"productID"`
	ProductName string `json:"productName"`
	Owner       string `json:"owner"`
	Status      string `json:"status"`
}

func (s *SmartContract) CreateProduct(ctx contractapi.TransactionContextInterface, productID, productName, owner string) error {
	exist, err := ctx.GetStub().GetState(productID)

	if err != nil {
		return err
	}

	if exist != nil {
		return fmt.Errorf("Product with id %v already exists", productID)
	}

	newProduct := Product{
		ProductID:   productID,
		ProductName: productName,
		Owner:       owner,
		Status:      "Factory",
	}

	byteData, err := json.Marshal(newProduct)

	return ctx.GetStub().PutState(productID, byteData)
}

func (s *SmartContract) TransferProduct(ctx contractapi.TransactionContextInterface, productID, newOwner string) error {
	exist, err := ctx.GetStub().GetState(productID)

	if err != nil {
		return err
	}

	if exist == nil {
		return fmt.Errorf("Product with id %v does not exists", productID)
	}

	var newProduct Product

	err = json.Unmarshal(exist, &newProduct)

	if err != nil {
		return err
	}

	newProduct.Owner = newOwner
	newProduct.Status = "Transferred"

	byteData, err := json.Marshal(newProduct)

	return ctx.GetStub().PutState(productID, byteData)
}

func (s *SmartContract) GetProduct(ctx contractapi.TransactionContextInterface, productID string) (*Product, error) {
	exist, err := ctx.GetStub().GetState(productID)

	if err != nil {
		return nil, err
	}

	if exist == nil {
		return nil, fmt.Errorf("Product with id %v does not exists", productID)
	}

	var newProduct Product

	err = json.Unmarshal(exist, &newProduct)

	if err != nil {
		return nil, err
	}

	return &newProduct, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(&SmartContract{})

	if err != nil {
		fmt.Errorf("Error while creating new chaincode")
	}

	err = chaincode.Start()

	if err != nil {
		fmt.Errorf("Error while starting chaincode")
	}

}
