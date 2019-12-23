package main

import (
	"fmt"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type MainChannelCC struct {}

type Request struct {
  Name  string  `json:"name"`
  Content string  `json:"content"`
  Difference string `json:"difference"`
  Date  string  `json:"date"`
}

type Difficulty struct {
  Request string `json:"request"`
  Uploader string `json:"uploader"`
  Content string `json:content`
  Difference string `json:difference`
  Date string `json:"date"`
}

type Winner struct {
  Request string  `json:request`
  Name string  `json:"name"`
  Date string  `json:"date"`
}

type Reward struct {
  RequestName string `json:"requestName"`
  DataName string `json:"dataName"`
  RewardContent string `json:"rewardContent"`
}

type SubChannel struct {
  Name string `json:"name"`
  Owner string `json:"owner"`
  Request string `json:request`
  Date string `json:date`
}

type Data struct {
  Request string `json:"request"`
  Name string `json:"name"`
  Content string `json:"content"`
  Date string `json:"date"`
}

func requestUpload(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 4 {
    return shim.Error("invalid args length")
  }

  name := args[0]
  content := args[1]
  difference := args[2]
  date := args[3]

  request := &Request{
    Name: name,
    Content: content,
    Difference: difference,
    Date: date,
  }

  requestBytes, err := json.Marshal(request)
  if err != nil {
    return shim.Error(fmt.Sprintf("marshal request error %s", err))
  }

  requestKey, err := stub.CreateCompositeKey("request", []string{
    name,
  })

  if err != nil {
    return shim.Error(fmt.Sprintf("create key error: %s", err))
  }

  if err := stub.PutState(requestKey, requestBytes); err != nil {
    return shim.Error(fmt.Sprintf("save request error: %s", err))
  }

  return shim.Success(nil)
}

func requestQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 1 {
    return shim.Error("invalid args length")
  }

  name := args[0]

  requestKey, err := stub.CreateCompositeKey("request", []string{
    name,
  })

  if err != nil {
    return shim.Error(fmt.Sprintf("create key error: %s", err))
  }

  requestBytes, err := stub.GetState(requestKey)
  if err != nil || len(requestBytes) == 0 {
    return shim.Error(fmt.Sprintf("get %s request erro: %s", name, err))
  }

  return shim.Success(requestBytes)
}

func requestQueryArr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 0 {
    return shim.Error("invalid args length")
  }

  requestResult, err := stub.GetStateByPartialCompositeKey("request", []string{})
  if err != nil {
    return shim.Error(fmt.Sprintf("get request arr error: %s", err))
  }

  defer requestResult.Close()

  requestArr := make([]*Request, 0)
  for requestResult.HasNext() {
    requestVal, err := requestResult.Next()
    if err != nil {
      return shim.Error(fmt.Sprintf("query error: %s", err))
    }

    request := new(Request)

    if err := json.Unmarshal(requestVal.GetValue(), request); err != nil {
      return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
    }

    requestArr = append(requestArr, request)
  }

  requestArrBytes, err := json.Marshal(requestArr)
  if err != nil {
    return shim.Error(fmt.Sprintf("marshal error: %s", err))
  }

  return shim.Success(requestArrBytes)
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

  request := args[0]
  name := args[1]
  date := args[2]

  winner := &Winner{
    Request: request,
    Name: name,
    Date: date,
  }

  winnerBytes, err := json.Marshal(winner)
  if err != nil {
    return shim.Error(fmt.Sprintf("marshal request error %s", err))
  }

  winnerKey, err := stub.CreateCompositeKey("winner", []string{
    request,
    name,
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

  if len(args) != 2 {
    return shim.Error("invalid args length")
  }

  request := args[0]
  name := args[1]

  winnerKey, err := stub.CreateCompositeKey("winner", []string{
    request,
    name,
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

func winnerQueryArr(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 1 {
    return shim.Error("invalid args length")
  }

  request := args[0]

  winnerResult, err := stub.GetStateByPartialCompositeKey("winner", []string{
    request,
  })
  if err != nil {
    return shim.Error(fmt.Sprintf("query %s request winners arr error: %s", request, err))
  }

  defer winnerResult.Close()

  winnerArr := make([]*Winner, 0)
  for winnerResult.HasNext() {
    winnerVal, err := winnerResult.Next()
    if err != nil {
      return shim.Error(fmt.Sprintf("query error: %s", err))
    }

    winner := new(Winner)
    if err := json.Unmarshal(winnerVal.GetValue(), winner); err != nil {
      return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
    }

    winnerArr = append(winnerArr, winner)
  }

  winnerArrBytes, err := json.Marshal(winnerArr)
  if err != nil {
    return shim.Error(fmt.Sprintf("marshal error: %s", err))
  }

  return shim.Success(winnerArrBytes)
}

func subChannelUpload(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 4 {
    return shim.Error("invalid args length")
  }

  name := args[0]
  owner := args[1]
  request := args[2]
  date := args[3]

  subChannel := &SubChannel{
    Name: name,
    Owner: owner,
    Request: request,
    Date: date,
  }

  subChannelBytes, err := json.Marshal(subChannel)
  if err != nil {
    return shim.Error(fmt.Sprintf("marshal subChannel error %s", err))
  }

  subChannelKey, err := stub.CreateCompositeKey("subChannel", []string{
    owner,
    request,
  })

  if err != nil {
    return shim.Error(fmt.Sprintf("create key error: %s", err))
  }

  if err := stub.PutState(subChannelKey, subChannelBytes); err != nil {
    return shim.Error(fmt.Sprintf("save reward error: %s", err))
  }

  return shim.Success(nil)
}

func subChannelQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 2 {
    return shim.Error(fmt.Sprintf("invalid args length"))
  }

  owner := args[0]
  request := args[1]

  subChannelKey, err := stub.CreateCompositeKey("subChannel", []string{
    owner,
    request,
  })

  if err != nil {
    return shim.Error(fmt.Sprintf("create key error: %s", err))
  }

  subChannelBytes, err := stub.GetState(subChannelKey)
  if err != nil || len(subChannelBytes) == 0 {
    return shim.Error(fmt.Sprintf("get %s owner %s request subchannel erro: %s ", owner, request, err))
  }

  return shim.Success(subChannelBytes)
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

  dataBytes, err := json.Marshal(data)
  if err != nil {
    return shim.Error(fmt.Sprintf("marshal data error %s", err))
  }

  dataKey, err := stub.CreateCompositeKey("data", []string{
    request,
    name,
  })

  if err != nil {
    return shim.Error(fmt.Sprintf("create key error %s", err))
  }

  if err := stub.PutState(dataKey, dataBytes); err != nil {
    return shim.Error(fmt.Sprintf("save data error: %s", err))
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
    return shim.Error(fmt.Sprintf("get %s data error: %s", name, err))
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
    return shim.Error(fmt.Sprintf("query %s request data arr error: %s", request, err))
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

func rewardsUpload(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 3 {
    return shim.Error("invalid args length")
  }

  requestName := args[0]
  dataName := args[1]
  rewardContent := args[2]

  reward := &Reward{
    RequestName: requestName,
    DataName: dataName,
    RewardContent: rewardContent,
  }

  rewardBytes, err := json.Marshal(reward)
  if err != nil {
    return shim.Error(fmt.Sprintf("marshal reward error %s", err))
  }

  rewardKey, err := stub.CreateCompositeKey("reward", []string{
    requestName,
    dataName,
  })

  if err != nil {
    return shim.Error(fmt.Sprintf("create key error: %s", err))
  }

  if err := stub.PutState(rewardKey, rewardBytes); err != nil {
    return shim.Error(fmt.Sprintf("save reward error: %s", err))
  }

  return shim.Success(nil)
}

func rewardsReceive(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 2 {
    return shim.Error("invalid args length")
  }

  requestName := args[0]
  dataName := args[1]

  rewardKey, err := stub.CreateCompositeKey("reward", []string{
    requestName,
    dataName,
  })

  if err != nil {
    return shim.Error(fmt.Sprintf("create key error: %s", err))
  }

  rewardBytes, err := stub.GetState(rewardKey)
  if err != nil || len(rewardBytes) == 0 {
    return shim.Error(fmt.Sprintf("get %s request %s reward erro: %s", requestName, dataName, err))
  }

  return shim.Success(rewardBytes)
}

// Init is called during Instantiate transaction after the chaincode container
// has been established for the first time, allowing the chaincode to
// initialize its internal data
func (c *MainChannelCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke is called to update or query the ledger in a proposal transaction.
// Updated state variables are not committed to the ledger until the
// transaction is committed.
func (c *MainChannelCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
  funcName, args := stub.GetFunctionAndParameters()

  switch funcName {
    // 任务上传
  case "requestUpload":
    return requestUpload(stub, args)
    // 查询任务
  case "requestQuery":
    return requestQuery(stub, args)
    // 查询全部任务
  case "requestQueryArr":
    return requestQueryArr(stub, args)
    // 难度值上传
  case "difficultyUpload":
    return difficultyUpload(stub, args)
    // 难度值查询
  case "difficultyQuery":
    return difficultyQuery(stub, args)
    // 难度值统一查询
  case "difficultyQueryArr":
    return difficultyQueryArr(stub, args)
    // 判断胜利者
  case "winnerUpload":
    return winnerUpload(stub, args)
    // 查询胜利者
  case "winnerQuery":
    return winnerQuery(stub, args)
    // 查询全部胜利者
  case "winnerQueryArr":
    return winnerQueryArr(stub, args)
    // 子channel上传
  case "subChannelUpload":
    return subChannelUpload(stub, args)
    // 子channel查询
  case "subChannelQuery":
    return subChannelQuery(stub, args)
    // 数据上传
  case "dataUpload":
    return dataUpload(stub, args)
    // 查询数据
  case "dataQuery":
    return dataQuery(stub, args)
    // 数据统一查询
  case "dataQueryArr":
    return dataQueryArr(stub, args)
    // 奖励发放
  case "rewardsUpload":
    return rewardsUpload(stub, args)
    // 奖励获取
  case "rewardsReceive":
    return rewardsReceive(stub, args)
  }

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(MainChannelCC))
	if err != nil {
		fmt.Printf("Error starting AssertsExchange chaincode: %s", err)
	}
}
