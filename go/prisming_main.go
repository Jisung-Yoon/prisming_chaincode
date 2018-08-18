package main

import (
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}


type Donor struct {
	ObjectType     string      `json:"doctype"` // field for couchdb
	Id     string     `json:"id"`
	Name     string     `json:"name"`
	Phone     string	`json:"phone"`
	Credit     int     `json:"credit"`
	Assets_array []string `json:"assetArray"`
}


type Asset struct {
	ObjectType     string      `json:"doctype"` // field for couchdb
	Id     string     `json:"id"`
	Name  string `json:"name"`
	DonorId     string     `json:"donorid"`
	NPOId string `json:"npoid`
	Owner_history []OwnerRelation `json:"owner"`
	Status     string     `json:"status"`
	ProductType     string     `json:"producttype"`
	Picture     string     `json:"pichash"` // generated by hashing algorithm

}

type NPO struct {
	ObjectType     string      `json:"doctype"` // field for couchdb
	Id     string     `json:"id"`
	Name     string     `json:"name"`
	Assets_array []string `json:"assetsarray"`
	Needs []string `json:"needs"`
}

type Recipient struct {
	ObjectType     string      `json:"doctype"` // field for couchdb
	Id     string	`json:"id"`
	Name     string     `json:"name"`
	Types string `json:"type"`
	Asset_array []string `json:"assetarray"`
}


type OwnerRelation struct {
	Id         string `json:"id"`
	Username   string `json:"username"`    //this is mostly cosmetic/handy, the real relation is by Id not Username
	User_type   string `json:"user_type"`     //this is mostly cosmetic/handy, the real relation is by Id not Company
}


// Donation needs from NPO
type Need struct {
	Id string `json:"id"`
	NPOID string `json:"npoid`
	ProductType string `json:"producttype"`
	Name string `json:"name"`
	Status string `json:"status"`
	Total_count int `json:"totalcount"`
	Current_count int `json:"currentcount"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}


// Init function
// Activate when instantiating
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Marbles Is Starting Up")
	funcName, args := stub.GetFunctionAndParameters()
	var err error
	txId := stub.GetTxID()

	fmt.Println("Init() is running")
	fmt.Println("Transaction ID:", txId)
	fmt.Println("  GetFunctionAndParameters() function:", funcName)
	fmt.Println("  GetFunctionAndParameters() args count:", len(args))
	fmt.Println("  GetFunctionAndParameters() args found:", args)


	// showing the alternative argument shim function
	alt := stub.GetStringArgs()
	fmt.Println("  GetStringArgs() args count:", len(alt))
	fmt.Println("  GetStringArgs() args found:", alt)

	// store compatible marbles application version
	err = stub.PutState("marbles_ui", []byte("0.0.1"))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Ready for action")                          //self-test pass
	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)
	fmt.Println(args)

	if function == "query"{
		return t.query(stub, args)
	} else if function == "enroll_donor"{
		return t.enroll_donor(stub, args)
	} else if function == "enroll_npo" {
		return t.enroll_npo(stub, args)
	} else if function == "enroll_needs" {
		return t.enroll_needs(stub, args)
	} else if function == "propose_asset" {
		return t.propose_asset(stub, args)
	} else if function == "approve_asset" {
		return t.approve_asset(stub, args)
	} else if function == "delete_asset" {
		return t.delete_asset(stub, args)
	} else if function == "enroll_recipient" {
		return t.enroll_recipient(stub, args)
	} else if function == "borrow_asset" {
		return t.borrow_asset(stub, args)
	} else if function == "give_asset" {
		return t.give_asset(stub, args)
	} else if function == "get_back_asset" {
		return t.get_back_asset(stub, args)
	} else if function == "read_everything" {
		return t.read_everything(stub)
	} else if function == "get_history" {
		return t.get_history(stub, args)
	}

	// error out
	fmt.Println("Received unknown invoke function name - " + function)
	return shim.Error("Received unknown invoke function name - '" + function + "'")
}

func (t *SimpleChaincode) enroll_donor(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var temp_donor Donor  // Entities
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	temp_donor.ObjectType = "Donor"
	temp_donor.Id = args[0] // d0~d999999999
	temp_donor.Name = args[1]
	temp_donor.Phone = args[2]
	temp_donor.Credit = 0
	temp_donor.Assets_array = []string{}

	fmt.Println(temp_donor)

	donorAsBytes, _ := json.Marshal(temp_donor)
	fmt.Println("writing donor to state")
	fmt.Println(string(donorAsBytes))

	err = stub.PutState(temp_donor.Id, donorAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store donor")
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) enroll_npo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var temp_NPO NPO  // Entities
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}


	temp_NPO.ObjectType = "NPO"
	temp_NPO.Id = args[0]
	temp_NPO.Name = args[1]
	temp_NPO.Assets_array = []string{}
	temp_NPO.Needs = []string{}

	fmt.Println(temp_NPO)

	NPOAsBytes, _ := json.Marshal(temp_NPO)
	fmt.Println("writing NPO information to ledger")
	fmt.Println(string(NPOAsBytes))

	err = stub.PutState(temp_NPO.Id, NPOAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store NPO")
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) enroll_recipient(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var temp_rec Recipient  // Entities
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	temp_rec.ObjectType = "Recipient"
	temp_rec.Id = args[0]
	temp_rec.Name = args[1]
	temp_rec.Types = args[2]
	temp_rec.Asset_array = []string{}

	fmt.Println(temp_rec)

	RecAsBytes, _ := json.Marshal(temp_rec)
	fmt.Println("writing Recipient information to ledger")
	fmt.Println(string(RecAsBytes))

	err = stub.PutState(temp_rec.Id, RecAsBytes)  //store owner by its Id
	if err != nil {
		fmt.Println("Could not store Recipient")
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) enroll_needs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var temp_npo NPO
	temp_npo_id := args[1]
	temp_npo_by_byte, err := stub.GetState(temp_npo_id)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get npo state for\"}"
		return shim.Error(jsonResp)
	}
	if temp_npo_by_byte == nil {
		jsonResp := "{\"Error\":\"Nil amount npo state\"}"
		return shim.Error(jsonResp)
	}


	json.Unmarshal(temp_npo_by_byte, &temp_npo)



	var temp_need Need

	temp_need.Id = args[0]
	temp_need.NPOID = args[1]
	temp_need.Name = args[2]
	temp_need.ProductType = args[3]
	temp_need.Total_count,_ = strconv.Atoi(args[4])
	temp_need.Current_count = 0
	temp_need.Status = "I"

	fmt.Println(temp_need)

	temp_npo.Needs = append(temp_npo.Needs, temp_need.Id)

	fmt.Println(temp_npo)
	NPOAsBytes, _ := json.Marshal(temp_npo)
	fmt.Println("writing NPO information to ledger")
	fmt.Println(string(NPOAsBytes))

	err = stub.PutState(temp_npo.Id, NPOAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store NPO")
		return shim.Error(err.Error())
	}

	fmt.Println(temp_need)
	NeedsAsBytes, _ := json.Marshal(temp_need)
	fmt.Println("writing needs information to ledger")
	fmt.Println(string(NeedsAsBytes))

	err = stub.PutState(temp_need.Id, NeedsAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store needs")
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) propose_asset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var temp_asset Asset  // Entities
	var err error

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	temp_asset.ObjectType = "Asset"
	temp_asset.Id = args[0]
	temp_asset.Name = args[1]

	var temp_donor Donor
	temp_donor_by_byte, err := stub.GetState(args[2])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get donor state\"}"
		return shim.Error(jsonResp)
	}

	if temp_donor_by_byte == nil {
		jsonResp := "{\"Error\":\"Nil amount for \"}"
		return shim.Error(jsonResp)
	}

	json.Unmarshal(temp_donor_by_byte, &temp_donor)

	temp_asset.DonorId = temp_donor.Id

	var temp_npo NPO
	temp_npo_by_byte, err := stub.GetState(args[3])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get NPO state for \"}"
		return shim.Error(jsonResp)
	}

	if temp_donor_by_byte == nil {
		jsonResp := "{\"Error\":\"Nil amount for NPO state\"}"
		return shim.Error(jsonResp)
	}
	json.Unmarshal(temp_npo_by_byte, &temp_npo)

	temp_asset.NPOId = temp_npo.Id
	temp_asset.Owner_history = []OwnerRelation{}
	temp_asset.Status = "Proposed"
	temp_asset.ProductType = args[4]
	temp_asset.Picture = args[5]

	fmt.Println(temp_asset)

	AssetAsBytes, _ := json.Marshal(temp_asset)
	fmt.Println("writing Asset information to ledger")
	fmt.Println(string(AssetAsBytes))

	err = stub.PutState(temp_asset.Id, AssetAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store Asset")
		return shim.Error(err.Error())
	}


	temp_donor.Assets_array = append(temp_donor.Assets_array, temp_asset.Id)
	fmt.Println(temp_donor)
	DonorAsBytes, _ := json.Marshal(temp_donor)
	fmt.Println("Updating donor information to ledger")
	fmt.Println(string(DonorAsBytes))

	err = stub.PutState(temp_donor.Id, DonorAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not update donor")
		return shim.Error(err.Error())
	}

	temp_npo.Assets_array = append(temp_npo.Assets_array, temp_asset.Id)
	fmt.Println(temp_npo)
	NpoAsBytes, _ := json.Marshal(temp_npo)
	fmt.Println("Updating npo information to ledger")
	fmt.Println(string(NpoAsBytes))

	err = stub.PutState(temp_npo.Id,  NpoAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not update NPO")
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) approve_asset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var temp_asset Asset
	temp_asset_by_byte, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get Asset state\"}"
		return shim.Error(jsonResp)
	}
	json.Unmarshal(temp_asset_by_byte, &temp_asset)

	fmt.Println(temp_asset.NPOId)
	if temp_asset.NPOId != args[1]{
		jsonResp := "{\"Error\":\"Asset is not owned by given NPO\"}"
		return shim.Error(jsonResp)
	}

	var temp_npo NPO
	temp_npo_by_byte, err := stub.GetState(temp_asset.NPOId)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get npo state \"}"
		return shim.Error(jsonResp)
	}
	json.Unmarshal(temp_npo_by_byte, &temp_npo)

	fmt.Println(temp_npo)

	var temp_need Need
	check := false
	for _, v := range temp_npo.Needs {
		temp_npo_by_byte, err := stub.GetState(v)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get needs state \"}"
			return shim.Error(jsonResp)
		}
		json.Unmarshal(temp_npo_by_byte, &temp_need)
		if temp_need.Name == temp_asset.Name {
			temp_need.Current_count = temp_need.Current_count + 1
			if temp_need.Current_count == temp_need.Total_count{
				temp_need.Status = "C"
			}
			check = true
			break
		}
	}


	temp_asset.Status = "Approved"

	fmt.Println(temp_asset)

	AssetAsBytes, _ := json.Marshal(temp_asset)
	fmt.Println("writing Asset information to ledger")
	fmt.Println(string(AssetAsBytes))

	err = stub.PutState(temp_asset.Id, AssetAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store Asset")
		return shim.Error(err.Error())
	}

	if check == true {
		var temp_donor Donor
		temp_donor_id := temp_asset.DonorId
		temp_donor_by_byte, err := stub.GetState(temp_donor_id)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get Donor state\"}"
			return shim.Error(jsonResp)
		}
		json.Unmarshal(temp_donor_by_byte, &temp_donor)
		temp_donor.Credit = temp_donor.Credit + 1
		fmt.Println(temp_donor)

		DonorAsBytes, _ := json.Marshal(temp_donor)
		fmt.Println("writing Donor information to ledger")
		fmt.Println(string(DonorAsBytes))

		err = stub.PutState(temp_donor.Id, DonorAsBytes)                    //store owner by its Id
		if err != nil {
			fmt.Println("Could not store Donor")
			return shim.Error(err.Error())
		}

		NpoAsBytes, _ := json.Marshal(temp_npo)
		fmt.Println("writing Npo information to ledger")
		fmt.Println(string(NpoAsBytes))

		err = stub.PutState(temp_npo.Id, NpoAsBytes)                    //store owner by its Id
		if err != nil {
			fmt.Println("Could not store Npo")
			return shim.Error(err.Error())
		}
		NeedsAsBytes, _ := json.Marshal(temp_need)
		fmt.Println("writing Needs information to ledger")
		fmt.Println(string(NeedsAsBytes))

		err = stub.PutState(temp_need.Id, NeedsAsBytes)                    //store owner by its Id
		if err != nil {
			fmt.Println("Could not store Needs")
			return shim.Error(err.Error())
		}

	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) delete_asset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var temp_asset Asset
	temp_asset_by_byte, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get Asset state \"}"
		return shim.Error(jsonResp)
	}
	json.Unmarshal(temp_asset_by_byte, &temp_asset)

	if temp_asset.NPOId != args[1]{
		jsonResp := "{\"Error\":\"Asset is not owned by given NPO\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(temp_asset.Id)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store Asset")
		return shim.Error(err.Error())
	}

	var temp_npo NPO
	temp_npo_by_byte, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get Npo state \"}"
		return shim.Error(jsonResp)
	}
	json.Unmarshal(temp_npo_by_byte, &temp_npo)

	fmt.Println(temp_npo)
	for i, v := range temp_npo.Assets_array {
		if v == temp_asset.Id {
			temp_npo.Assets_array = append(temp_npo.Assets_array[:i], temp_npo.Assets_array[i+1:]...)
			break
		}
	}
	fmt.Println(temp_npo)

	NpoAsBytes, _ := json.Marshal(temp_npo)
	fmt.Println("writing NPO information to ledger")
	fmt.Println(string(NpoAsBytes))

	err = stub.PutState(temp_npo.Id, NpoAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store NPO")
		return shim.Error(err.Error())
	}

	var temp_donor Donor
	temp_donor_by_byte, err := stub.GetState(temp_asset.DonorId)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get Donor state for \"}"
		return shim.Error(jsonResp)
	}
	json.Unmarshal(temp_donor_by_byte, &temp_donor)

	fmt.Println(temp_donor)
	for i, v := range temp_donor.Assets_array {
		if v == temp_asset.Id {
			temp_donor.Assets_array = append(temp_donor.Assets_array[:i], temp_donor.Assets_array[i+1:]...)
			break
		}
	}
	fmt.Println(temp_donor)

	DonorAsBytes, _ := json.Marshal(temp_donor)
	fmt.Println("writing Donor information to ledger")
	fmt.Println(string(DonorAsBytes))

	err = stub.PutState(temp_donor.Id, DonorAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store Donor")
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) borrow_asset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}


	var temp_asset Asset
	temp_asset_by_byte, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get asset state\"}"
		return shim.Error(jsonResp)
	}

	if temp_asset_by_byte == nil {
		jsonResp := "{\"Error\":\"Nil amount asset state\"}"
		return shim.Error(jsonResp)
	}

	json.Unmarshal(temp_asset_by_byte, &temp_asset)

	var temp_rec Recipient
	temp_rec_by_byte, err := stub.GetState(args[1])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get rec state for \"}"
		return shim.Error(jsonResp)
	}

	if temp_asset_by_byte == nil {
		jsonResp := "{\"Error\":\"Nil amount rec information \"}"
		return shim.Error(jsonResp)
	}

	json.Unmarshal(temp_rec_by_byte, &temp_rec)
	temp_rec.Asset_array = append(temp_rec.Asset_array, temp_asset.Id)


	var temp_owner_relation OwnerRelation
	temp_owner_relation.Id = temp_rec.Id
	temp_owner_relation.Username = temp_rec.Name
	temp_owner_relation.User_type= temp_rec.Types
	fmt.Println(temp_owner_relation)

	temp_asset.Owner_history = append(temp_asset.Owner_history, temp_owner_relation)
	temp_asset.Status = "Borrowed"

	AssetAsBytes, _ := json.Marshal(temp_asset)
	fmt.Println("writing Asset information to ledger")
	fmt.Println(string(AssetAsBytes))

	err = stub.PutState(temp_asset.Id, AssetAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store Asset")
		return shim.Error(err.Error())
	}

	RecAsBytes, _ := json.Marshal(temp_rec)
	fmt.Println("writing Rec information to ledger")
	fmt.Println(string(RecAsBytes))

	err = stub.PutState(temp_rec.Id, RecAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store Rec")
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) give_asset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}


	var temp_asset Asset
	temp_asset_by_byte, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get asset state\"}"
		return shim.Error(jsonResp)
	}

	if temp_asset_by_byte == nil {
		jsonResp := "{\"Error\":\"Nil amount asset state\"}"
		return shim.Error(jsonResp)
	}

	json.Unmarshal(temp_asset_by_byte, &temp_asset)

	var temp_rec Recipient
	temp_rec_by_byte, err := stub.GetState(args[1])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get rec state for \"}"
		return shim.Error(jsonResp)
	}

	if temp_asset_by_byte == nil {
		jsonResp := "{\"Error\":\"Nil amount rec information \"}"
		return shim.Error(jsonResp)
	}

	json.Unmarshal(temp_rec_by_byte, &temp_rec)
	temp_rec.Asset_array = append(temp_rec.Asset_array, temp_asset.Id)


	var temp_owner_relation OwnerRelation
	temp_owner_relation.Id = temp_rec.Id
	temp_owner_relation.Username = temp_rec.Name
	temp_owner_relation.User_type= temp_rec.Types
	fmt.Println(temp_owner_relation)

	temp_asset.Owner_history = append(temp_asset.Owner_history, temp_owner_relation)
	temp_asset.Status = "Given"

	AssetAsBytes, _ := json.Marshal(temp_asset)
	fmt.Println("writing Asset information to ledger")
	fmt.Println(string(AssetAsBytes))

	err = stub.PutState(temp_asset.Id, AssetAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store Asset")
		return shim.Error(err.Error())
	}

	RecAsBytes, _ := json.Marshal(temp_rec)
	fmt.Println("writing Rec information to ledger")
	fmt.Println(string(RecAsBytes))

	err = stub.PutState(temp_rec.Id, RecAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store Rec")
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) get_back_asset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}


	var temp_asset Asset
	temp_asset_by_byte, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get asset state for\"}"
		return shim.Error(jsonResp)
	}

	if temp_asset_by_byte == nil {
		jsonResp := "{\"Error\":\"Nil amount asset state \"}"
		return shim.Error(jsonResp)
	}

	json.Unmarshal(temp_asset_by_byte, &temp_asset)

	var temp_rec Recipient
	temp_rec_by_byte, err := stub.GetState(args[1])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get rec state \"}"
		return shim.Error(jsonResp)
	}

	if temp_asset_by_byte == nil {
		jsonResp := "{\"Error\":\"Nil amount rec state\"}"
		return shim.Error(jsonResp)
	}

	json.Unmarshal(temp_rec_by_byte, &temp_rec)



	temp_asset.Status = "Approved"
	for i, v := range temp_rec.Asset_array {
		if v == temp_asset.Id {
			temp_rec.Asset_array = append(temp_rec.Asset_array[:i], temp_rec.Asset_array[i+1:]...)
			break
		}
	}

	AssetAsBytes, _ := json.Marshal(temp_asset)
	fmt.Println("writing Asset information to ledger")
	fmt.Println(string(AssetAsBytes))

	err = stub.PutState(temp_asset.Id, AssetAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store Asset")
		return shim.Error(err.Error())
	}

	RecAsBytes, _ := json.Marshal(temp_rec)
	fmt.Println("writing Rec information to ledger")
	fmt.Println(string(RecAsBytes))

	err = stub.PutState(temp_rec.Id, RecAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store Rec")
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// Query - General_query_function (needed to be specify)
// ============================================================================================================================
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}



func (t *SimpleChaincode) read_everything(stub shim.ChaincodeStubInterface) pb.Response {
	type Everything struct {
		Donors   []Donor
		NPOs  []NPO
		Recipients []Recipient
		Assets []Asset
		Needs []Need
	}
	var everything Everything

	// ---- Get All Assets ---- //
	assetsIterator, err := stub.GetStateByRange("a0", "a9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer assetsIterator.Close()

	for assetsIterator.HasNext() {
		aKeyValue, err := assetsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on asset id - ", queryKeyAsStr)
		var asset Asset
		json.Unmarshal(queryValAsBytes, &asset)                   //un stringify it aka JSON.parse()
		everything.Assets = append(everything.Assets, asset)
	}
	fmt.Println("asset array - ", everything.Assets)


	// ---- Get All Donors ---- //
	donorsIterator, err := stub.GetStateByRange("d0", "d9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer donorsIterator.Close()

	for donorsIterator.HasNext() {
		aKeyValue, err := donorsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on donors id - ", queryKeyAsStr)
		var donor Donor
		json.Unmarshal(queryValAsBytes, &donor)                  //un stringify it aka JSON.parse()
		everything.Donors = append(everything.Donors, donor)   //add this marble to the list
	}
	fmt.Println("donor array - ", everything.Donors)

	// ---- Get All NPOs ---- //
	nposIterator, err := stub.GetStateByRange("n0", "n9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer nposIterator.Close()

	for nposIterator.HasNext() {
		aKeyValue, err := nposIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on NPO id - ", queryKeyAsStr)
		var npo NPO
		json.Unmarshal(queryValAsBytes, &npo)                   //un stringify it aka JSON.parse()
		everything.NPOs = append(everything.NPOs, npo)
	}
	fmt.Println("NPO array - ", everything.NPOs)

	// ---- Get All recipient ---- //
	recsIterator, err := stub.GetStateByRange("r0", "r9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer recsIterator.Close()

	for recsIterator.HasNext() {
		aKeyValue, err := recsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on NPO id - ", queryKeyAsStr)
		var recipient Recipient
		json.Unmarshal(queryValAsBytes, &recipient)                   //un stringify it aka JSON.parse()
		everything.Recipients = append(everything.Recipients, recipient)
	}

	fmt.Println("Reciptents array - ", everything.Recipients)

	// ---- Get All recipient ---- //
	needsIterator, err := stub.GetStateByRange("e0", "e9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer needsIterator.Close()

	for needsIterator.HasNext() {
		aKeyValue, err := needsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on Needs id - ", queryKeyAsStr)
		var need Need
		json.Unmarshal(queryValAsBytes, &need)                   //un stringify it aka JSON.parse()
		everything.Needs = append(everything.Needs, need)
	}

	fmt.Println("Needs array - ", everything.Needs)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything)              //convert to array of bytes
	return shim.Success(everythingAsBytes)
}


func (t *SimpleChaincode) get_history(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId    string   `json:"txId"`
		Value   Asset   `json:"value"`
	}
	var history []AuditHistory;
	var temp_asset Asset

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	assetId := args[0]
	fmt.Printf("- start getHistoryForAseet: %s\n", assetId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(assetId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx AuditHistory
		tx.TxId = historyData.TxId                     //copy transaction id over
		json.Unmarshal(historyData.Value, &temp_asset)     //un stringify it aka JSON.parse()
		if historyData.Value == nil {                  //marble has been deleted
			var emptyAsset Asset
			tx.Value = emptyAsset                 //copy nil marble
		} else {
			json.Unmarshal(historyData.Value, &temp_asset) //un stringify it aka JSON.parse()
			tx.Value = temp_asset                      //copy marble over
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForAssets returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}
