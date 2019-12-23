package main

import (
	"fmt"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SubChannelCC struct {}

type Data struct {
  Request string `json:"request"`
  Name string `json:"name"`
  Content string `json:"content"`
  Date string `json:"date"`
}

type Difficulty struct {
  Request string `json:"request"`
  Uploader string `json:"uploader"`
  Content string `json:content`
  Difference string `json:difference`
  Date string `json:"date"`
}

type Winner struct {
  Name string  `json:"name"`
  Date string  `json:"date"`
  Request string  `json:request`
}

func dataUpload(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 4 {
    return shim.Error("invalid args length")
  }

  request := args[0]
  name := args[1]
  content := args[2]
  date := args[3]

  data := &Data{
    Request: request,
    Name: name,
    Content: content,
    Date: date,
  }

  dataKey, err := stub.CreateCompositeKey("data", []string{
    request,
    name,
  })

  dataBytes, err := json.Marshal(data)
  if err != nil {
    return shim.Error(fmt.Sprintf("marshal error: %s", err))
  }

  if err := stub.PutState(dataKey, dataBytes); err != nil {
    return shim.Error(fmt.Sprintf("save %s request %s name data error: %s", request, name, err))
  }

	eventPayLoad := "Upload data name " + name + "for request " +  request + "at " + date
	payloadAsBytes := []byte(eventPayLoad)
	eventErr := stub.SetEvent("dataUploadEvent", payloadAsBytes)
	if (eventErr != nil) {
		return shim.Error(fmt.Sprintf("failed to emit event dataUploadEvent"))
	}

  return shim.Success(nil)
}

func dataQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 2 {
    return shim.Error("invalid args length")
  }

  request := args[0]
  name := args[1]

  dataKey, err := stub.CreateCompositeKey("data", []string{
    request,
    name,
  })
  if err != nil {
    return shim.Error(fmt.Sprintf("create key error: %s", err))
  }

  dataBytes, err := stub.GetState(dataKey)
  if err != nil || len(dataBytes) == 0 {
    return shim.Error(fmt.Sprintf("get %s request data error: %s", request, err))
  }

  return shim.Success(dataBytes)
}

func dataQueryArr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 1 {
    return shim.Error("invalid args length")
  }

  request := args[0]

  dataResult, err := stub.GetStateByPartialCompositeKey("data", []string{
    request,
  })

  if err != nil {
    return shim.Error(fmt.Sprintf("get %s request data arr error: %s", request, err))
  }

  defer dataResult.Close()

  dataArr := make([]*Data, 0)
  for dataResult.HasNext() {
    dataVal, err := dataResult.Next()
    if err != nil {
      return shim.Error(fmt.Sprintf("query error: %s", err))
    }

    data := new(Data)
    if err := json.Unmarshal(dataVal.GetValue(), data); err != nil {
      return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
    }

    dataArr = append(dataArr, data)
  }

  dataBytes, err := json.Marshal(dataArr)
  if err != nil {
    return shim.Error(fmt.Sprintf("marshal error: %s", err))
  }

  return shim.Success(dataBytes)
}

func difficultyUpload(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 5 {
    return shim.Error(fmt.Sprintf("invalid args length"))
  }

  request := args[0]
  uploader := args[1]
  content := args[2]
  difference := args[3]
  date := args[4]

  difficulty := &Difficulty{
    Request: request,
    Uploader: uploader,
    Content: content,
    Difference: difference,
    Date: date,
  }

  difficultyBytes, err := json.Marshal(difficulty)
  if err != nil {
    return shim.Error(fmt.Sprintf("marshal difficulty error: %s", err))
  }

  difficultyKey, err := stub.CreateCompositeKey("difficulty", []string{
    request,
    difference,
    uploader,
  })

  if err != nil {
    return shim.Error(fmt.Sprintf("create key error: %s", err))
  }

  if err := stub.PutState(difficultyKey, difficultyBytes); err != nil {
    return shim.Error(fmt.Sprintf("save %s request %s uploader difficulty error: %s", request, uploader, err))
  }

  return shim.Success(nil)
}

func difficultyQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 3 {
    return shim.Error("invalid args length")
  }

  request := args[0]
  uploader := args[1]
  difference := args[2]

  difficultyKey, err := stub.CreateCompositeKey("difficulty", []string{
    request,
    difference,
    uploader,
  })

  if err != nil {
    return shim.Error(fmt.Sprintf("create key error: %s", err))
  }

  difficultyBytes, err := stub.GetState(difficultyKey)
  if err != nil || len(difficultyBytes) == 0 {
    return shim.Error(fmt.Sprintf("get %s request %s uploader difficulty error: %s", request, uploader, err))
  }

  return shim.Success(difficultyBytes)
}

func difficultyQueryArr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 2 {
    return shim.Error("invalid args length")
  }

  request := args[0]
  difference := args[1]

  difficultyResult, err := stub.GetStateByPartialCompositeKey("difficulty", []string{
    request,
    difference,
  })
  if err != nil {
    return shim.Error(fmt.Sprintf("get %s request difficulty arr error: %s", request, err))
  }

  defer difficultyResult.Close()

  difficultyArr := make([]*Difficulty, 0)
  for difficultyResult.HasNext() {
    difficultyVal, err := difficultyResult.Next()
    if err != nil {
      return shim.Error(fmt.Sprintf("query error: %s", err))
    }

    difficulty := new(Difficulty)
    if err := json.Unmarshal(difficultyVal.GetValue(), difficulty); err != nil {
      return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
    }

    difficultyArr = append(difficultyArr, difficulty)
  }

  difficultyBytes, err := json.Marshal(difficultyArr)
  if err != nil {
    return shim.Error(fmt.Sprintf("marshal error: %s", err))
  }

  return shim.Success(difficultyBytes)
}


func winnerUpload(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 3 {
    return shim.Error("invalid args length")
  }

  name := args[0]
  date := args[1]
  request := args[2]

  winner := &Winner{
    Name: name,
    Date: date,
    Request: request,
  }

  winnerBytes, err := json.Marshal(winner)
  if err != nil {
    return shim.Error(fmt.Sprintf("marshal request error %s", err))
  }

  winnerKey, err := stub.CreateCompositeKey("winner", []string{
    request,
  })

  if err != nil {
    return shim.Error(fmt.Sprintf("create key error: %s", err))
  }

  if err := stub.PutState(winnerKey, winnerBytes); err != nil {
    return shim.Error(fmt.Sprintf("save winner error: %s", err))
  }

  return shim.Success(nil)
}

func winnerQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 1 {
    return shim.Error("invalid args length")
  }

  request := args[0]

  winnerKey, err := stub.CreateCompositeKey("winner", []string{
    request,
  })

  if err != nil {
    return shim.Error(fmt.Sprintf("create key error: %s", err))
  }

  winnerBytes, err := stub.GetState(winnerKey)
  if err != nil || len(winnerBytes) == 0 {
    return shim.Error(fmt.Sprintf("get %s winner erro: %s", request, err))
  }

  return shim.Success(winnerBytes)
}

// Init is called during Instantiate transaction after the chaincode container
// has been established for the first time, allowing the chaincode to
// initialize its internal data
func (c *SubChannelCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke is called to update or query the ledger in a proposal transaction.
// Updated state variables are not committed to the ledger until the
// transaction is committed.
func (c *SubChannelCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
  funcName, args := stub.GetFunctionAndParameters()

  switch funcName {
    // 数据上传
  case "dataUpload":
    return dataUpload(stub, args)
    // 数据查询
  case "dataQuery":
    return dataQuery(stub, args)
    // 全部数据查询
  case "dataQueryArr":
    return dataQueryArr(stub, args)
    // 难度值上传
  case "difficultyUpload":
    return difficultyUpload(stub, args)
    // 难度值查询
  case "difficultyQuery":
    return difficultyQuery(stub, args)
    // 难度值范围查询
  case "difficultyQueryArr":
    return difficultyQueryArr(stub, args)
    // 判断胜利者
  case "winnerUpload":
    return winnerUpload(stub, args)
    // 查询胜利者
  case "winnerQuery":
    return winnerQuery(stub, args)
  }

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SubChannelCC))
	if err != nil {
		fmt.Printf("Error starting AssertsExchange chaincode: %s", err)
	}
}
