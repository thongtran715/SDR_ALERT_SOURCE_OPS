package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// SourceOperatorInfo This struct consists of info holding the name of the operator as well as the info of the data counts
type SourceOperatorInfo struct {
	operatorName                   string
	totalMessagesReceived          int
	totalMessagesLessThan10mins    int
	totalMessagesLessThan1min      int
	totalMessagesLessThan1hour     int
	totalMessagesLessThan2hour     int
	totalMessagesDLRS              int
	totalMessagesLessThan10seconds int
}

// OperatorNode : Contains Source
type OperatorNode struct {
	leftOperator       *OperatorNode
	rightOperator      *OperatorNode
	sourceOperatorInfo SourceOperatorInfo
}

// SourceNode : contains the string and operator root tree
type SourceNode struct {
	leftSource       *SourceNode
	rightSource      *SourceNode
	sourceName       string
	operatorTreeRoot *OperatorNode
}

// SourceTree : root
type SourceTree struct {
	root *SourceNode
}

// ****************************************************************************************************************
// 												DISPLAY
func (node *OperatorNode) display() {
	if node == nil {
		return
	}
	fmt.Println("Name of operator: ", node.sourceOperatorInfo.operatorName)
	fmt.Println("Less than 1 minute: ", node.sourceOperatorInfo.totalMessagesLessThan1min)
	fmt.Println("Less than 10 minutes: ", node.sourceOperatorInfo.totalMessagesLessThan10mins)
	fmt.Println("Less than 1 hour: ", node.sourceOperatorInfo.totalMessagesLessThan1hour)
	fmt.Println("Total Message Less Than 2 hours: ", node.sourceOperatorInfo.totalMessagesLessThan2hour)
	fmt.Println("Total Message Less Than 10 seconds: ", node.sourceOperatorInfo.totalMessagesLessThan10seconds)
	fmt.Println("Number of messages that have status 4: ", node.sourceOperatorInfo.totalMessagesDLRS)
	fmt.Println("Total Message Received: ", node.sourceOperatorInfo.totalMessagesReceived)
	fmt.Println("*******************************************")
	node.leftOperator.display()
	node.rightOperator.display()
}
func (node *SourceNode) display() {
	if node == nil {
		return
	}
	fmt.Println("Source Name", node.sourceName)
	node.operatorTreeRoot.display()
	node.leftSource.display()
	node.rightSource.display()
}

func (tree *SourceTree) display() {
	if tree.root == nil {
		return
	}
	tree.root.display()
}

// ****************************************************************************************************************
// 									ADD FUNCTIONS

func (node *OperatorNode) addOperator(operator SourceOperatorInfo) bool {
	if node == nil {
		return false
	}
	if strings.Compare(node.sourceOperatorInfo.operatorName, operator.operatorName) == -1 {
		if node.leftOperator == nil {
			node.leftOperator = &OperatorNode{
				sourceOperatorInfo: operator,
			}
			return true
		}
		return node.leftOperator.addOperator(operator)
	} else if strings.Compare(node.sourceOperatorInfo.operatorName, operator.operatorName) == 1 {
		if node.rightOperator == nil {
			node.rightOperator = &OperatorNode{
				sourceOperatorInfo: operator,
			}
			return true
		}
		return node.rightOperator.addOperator(operator)
	} else {
		return false
	}
}

// This function is to add new source for tree Node
func (node *SourceNode) addSource(source string, operator SourceOperatorInfo) bool {
	if node == nil {
		fmt.Println(source)
		return false
	}
	compare := strings.Compare(node.sourceName, source)
	if compare == -1 {
		if node.leftSource == nil {
			node.leftSource = &SourceNode{
				sourceName: source,
				operatorTreeRoot: &OperatorNode{
					sourceOperatorInfo: operator,
				},
			}
			return true
		}
		return node.leftSource.addSource(source, operator)
	} else if compare == 1 {
		if node.rightSource == nil {
			node.rightSource = &SourceNode{
				sourceName: source,
				operatorTreeRoot: &OperatorNode{
					sourceOperatorInfo: operator,
				},
			}
			return true
		}
		return node.rightSource.addSource(source, operator)
	} else {
		return node.operatorTreeRoot.addOperator(operator)
	}
}

func (tree *SourceTree) addSource(source string, operator SourceOperatorInfo) bool {
	if tree.root == nil {
		tree.root = &SourceNode{
			sourceName: source,
			operatorTreeRoot: &OperatorNode{
				sourceOperatorInfo: operator,
			},
		}
		return true
	}
	return tree.root.addSource(source, operator)
}

// ****************************************************************************************************************
// 									FIND FUNCTIONS

func (node *OperatorNode) findOperator(name string) *OperatorNode {
	if node == nil {
		return nil
	}
	if strings.Compare(node.sourceOperatorInfo.operatorName, name) == 0 {
		return node
	} else if strings.Compare(node.sourceOperatorInfo.operatorName, name) == -1 {
		return node.leftOperator.findOperator(name)
	} else if strings.Compare(node.sourceOperatorInfo.operatorName, name) == 1 {
		return node.rightOperator.findOperator(name)
	} else {
		return nil
	}
}

func (node *SourceNode) findSource(name string) *SourceNode {
	if node == nil {
		return nil
	}
	if strings.Compare(node.sourceName, name) == 0 {
		return node
	} else if strings.Compare(node.sourceName, name) == -1 {
		return node.leftSource.findSource(name)
	} else {
		return node.rightSource.findSource(name)
	}
}

func (tree *SourceTree) findSource(name string) *SourceNode {
	if tree.root == nil {
		return nil
	}
	return tree.root.findSource(name)
}

// ****************************************************************************************************************
// 									CHECK IF EXISTED FUNCTIONS

func (tree *SourceTree) checkIfSourceAndOperatorExisted(source string, operator string) bool {
	if tree.root == nil {
		return false
	}
	sourceNode := tree.findSource(source)
	if sourceNode != nil {
		opearatorNode := sourceNode.operatorTreeRoot.findOperator(operator)
		if opearatorNode != nil {
			return true
		}
		return false
	}
	return false
}

// ****************************************************************************************************************
// 									TOTAL SOURCES FUNCTIONS
func (node *SourceNode) totalTrees() int {
	if node == nil {
		return 0
	}
	return node.leftSource.totalTrees() + node.rightSource.totalTrees() + 1
}
func (tree *SourceTree) totalTrees() int {
	if tree.root == nil {
		return 0
	}
	return tree.root.totalTrees()
}

// ****************************************************************************************************************
// 									FIND AND DISPLAY  FUNCTIONS

func (node *SourceNode) findAndDisplaySource(name string) {
	if node == nil {
		return
	}
	compare := strings.Compare(node.sourceName, name)
	if compare == 0 {
		node.operatorTreeRoot.display()
	} else if compare > 0 {
		node.rightSource.findAndDisplaySource(name)
	} else {
		node.leftSource.findAndDisplaySource(name)
	}
}
func (tree *SourceTree) findAndDisplaySource(name string) {
	if tree.root == nil {
		return
	}
	tree.root.findAndDisplaySource(name)
}

// ****************************************************************************************************************
// 									FIND AND INCREMENTS FUNCTIONS

func (node *OperatorNode) findAndIncrementOperator(typeIncrement, opsName string) bool {
	if node == nil {
		return false
	}
	compare := strings.Compare(node.sourceOperatorInfo.operatorName, opsName)
	if compare == 0 {
		switch typeIncrement {
		case "totalMessagesLessThan10mins":
			node.sourceOperatorInfo.totalMessagesLessThan10mins++
		case "totalMessagesLessThan1min":
			node.sourceOperatorInfo.totalMessagesLessThan1min++
		case "totalMessagesLessThan1hour":
			node.sourceOperatorInfo.totalMessagesLessThan1hour++
		case "totalMessagesLessThan2hour":
			node.sourceOperatorInfo.totalMessagesLessThan2hour++
		case "totalMessagesDLRS":
			node.sourceOperatorInfo.totalMessagesDLRS++
		case "totalMessagesLessThan10seconds":
			node.sourceOperatorInfo.totalMessagesLessThan10seconds++
		case "totalMessagesReceived":
			node.sourceOperatorInfo.totalMessagesReceived++
		default:
		}
		return true
	} else if compare == -1 {
		return node.leftOperator.findAndIncrementOperator(typeIncrement, opsName)
	} else {
		return node.rightOperator.findAndIncrementOperator(typeIncrement, opsName)
	}
}

func (node *SourceNode) findAndIncrementOperator(typeIncrement, opsName, sourceName string) bool {
	if node == nil {
		return false
	}
	compare := strings.Compare(node.sourceName, sourceName)
	if compare == 0 {
		return node.operatorTreeRoot.findAndIncrementOperator(typeIncrement, opsName)
	} else if compare < 0 {
		return node.leftSource.findAndIncrementOperator(typeIncrement, opsName, sourceName)
	} else {
		return node.rightSource.findAndIncrementOperator(typeIncrement, opsName, sourceName)
	}
}

func (tree *SourceTree) findAndIncrement(typeName, opsName, sourceName string) bool {
	if tree.root == nil {
		return false
	}
	return tree.root.findAndIncrementOperator(typeName, opsName, sourceName)
}

// ****************************************************************************************************************
// 									FIND THE MOST FUNCTIONS

func (node *OperatorNode) findTheMostOPs() int {
	if node == nil {
		return 0
	}
	return node.leftOperator.findTheMostOPs() + node.rightOperator.findTheMostOPs() + 1
}
func (node *SourceNode) findTheMostOps(sourceName *string, value *int) {
	if node == nil {
		return
	}
	val := node.operatorTreeRoot.findTheMostOPs()
	if *value < val {
		*value = val
		*sourceName = node.sourceName
	}
	node.leftSource.findTheMostOps(sourceName, value)
	node.rightSource.findTheMostOps(sourceName, value)
}

func (tree *SourceTree) findTheMos() (string, int) {
	if tree.root == nil {
		return "", 0
	}
	sourceName, val := "", 0
	tree.root.findTheMostOps(&sourceName, &val)
	return sourceName, val
}

// ****************************************************************************************************************
// 									SDR Alert Based

func sendAlertToSlack(payload map[string]interface{}) {
	webHookURL := "https://hooks.slack.com/services/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"
	bytesRepresentation, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(webHookURL, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
}

func (node *OperatorNode) calculateSDRPercentage(sourceName string) {
	if node == nil {
		return
	}
	averageOneMinute := (float64(node.sourceOperatorInfo.totalMessagesLessThan1min) / float64(node.sourceOperatorInfo.totalMessagesReceived)) * 100
	averageTenMinutes := (float64(node.sourceOperatorInfo.totalMessagesLessThan10mins) / float64(node.sourceOperatorInfo.totalMessagesReceived)) * 100
	averageOneHour := (float64(node.sourceOperatorInfo.totalMessagesLessThan1hour) / float64(node.sourceOperatorInfo.totalMessagesReceived)) * 100
	averageTwoHours := (float64(node.sourceOperatorInfo.totalMessagesLessThan2hour) / float64(node.sourceOperatorInfo.totalMessagesReceived)) * 100
	averageTenSeconds := (float64(node.sourceOperatorInfo.totalMessagesLessThan10seconds) / float64(node.sourceOperatorInfo.totalMessagesReceived)) * 100
	averageStatus := (float64(node.sourceOperatorInfo.totalMessagesDLRS) / float64(node.sourceOperatorInfo.totalMessagesReceived)) * 100
	payload := make(map[string]interface{})
	payload["username"] = "SDR BASED ALERT"
	var str bytes.Buffer
	str.WriteString("Source Name: " + sourceName + "\nOperator Name: " + node.sourceOperatorInfo.operatorName)
	if averageOneMinute > 60 {
		str.WriteString("\nPercentage Under 1 minute: " + strconv.FormatFloat(averageOneMinute, 'f', 6, 64))
	}
	if averageTenMinutes > 60 {
		str.WriteString("\nPercentage Under Ten Minutes: " + strconv.FormatFloat(averageTenMinutes, 'f', 6, 64))
	}
	if averageOneHour > 60 {
		str.WriteString("\nPercentage Under 1 Hour: " + strconv.FormatFloat(averageOneHour, 'f', 6, 64))
	}
	if averageTwoHours > 60 {
		str.WriteString("\nPercentage Under 2 Hours: " + strconv.FormatFloat(averageTwoHours, 'f', 6, 64))
	}
	if averageTenSeconds > 60 {
		str.WriteString("\nPercentage Under 10 seconds: " + strconv.FormatFloat(averageTenSeconds, 'f', 6, 64))
	}
	if averageStatus > 60 {
		str.WriteString("\nPercentage Status (4): " + strconv.FormatFloat(averageStatus, 'f', 6, 64))
	}
	payload["text"] = str.String()
	sendAlertToSlack(payload)
	node.leftOperator.calculateSDRPercentage(sourceName)
	node.rightOperator.calculateSDRPercentage(sourceName)
}

func (node *SourceNode) calculateSDRPercentage() {
	if node == nil {
		return
	}
	node.operatorTreeRoot.calculateSDRPercentage(node.sourceName)
	node.leftSource.calculateSDRPercentage()
	node.rightSource.calculateSDRPercentage()
}
func (tree *SourceTree) calculateSDRPercentage() {
	if tree.root == nil {
		return
	}
	tree.root.calculateSDRPercentage()
}
