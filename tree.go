package main

import (
	"fmt"
	"strings"
)

// SourceOperatorInfo This struct consists of info holding the name of the operator as well as the info of the data counts
type SourceOperatorInfo struct {
	lessThan1min                int
	lessThan10mins              int
	lessThan1hour               int
	operatorName                string
	numberOfMessagesHaveStatus4 int
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

func (node *OperatorNode) display() {
	if node == nil {
		return
	}
	fmt.Println("Name of operator: ", node.sourceOperatorInfo.operatorName)
	fmt.Println("Less than 1 minute: ", node.sourceOperatorInfo.lessThan1min)
	fmt.Println("Less than 10 minutes: ", node.sourceOperatorInfo.lessThan10mins)
	fmt.Println("Less than 1 hour: ", node.sourceOperatorInfo.lessThan1hour)
	fmt.Println("Number of messages that have status 4: ", node.sourceOperatorInfo.numberOfMessagesHaveStatus4)
	fmt.Println("*******************************************")
	fmt.Println("*******************************************")
	node.leftOperator.display()
	node.rightOperator.display()
}

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

// This function is to add new source for tree Node
func (node *SourceNode) addSource(source string, operator SourceOperatorInfo) bool {
	if node == nil {
		fmt.Println(source)
		return false
	}
	if strings.Compare(node.sourceName, source) == -1 {
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
	} else if strings.Compare(node.sourceName, source) == 1 {
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

func (tree *SourceTree) incrementValues(status, typeTime string, source string, operator string) bool {

	soureNode := tree.findSource(source)
	if soureNode != nil {
		operatorNode := soureNode.operatorTreeRoot.findOperator(operator)
		if operatorNode != nil {
			if strings.Compare(typeTime, "10mins") == 0 {
				operatorNode.sourceOperatorInfo.lessThan10mins++
			} else if strings.Compare(typeTime, "1min") == 0 {
				operatorNode.sourceOperatorInfo.lessThan1min++
			} else if strings.Compare(typeTime, "1hour") == 0 {
				operatorNode.sourceOperatorInfo.lessThan1hour++
			}
			if strings.Compare(status, "4") == 0 {
				operatorNode.sourceOperatorInfo.numberOfMessagesHaveStatus4++
			}
			return true
		}
		return false
	}
	return false
}

// func (node *SourceNode) findSourceAndOperatorMoreThanCertainTime(time int) bool {
// 	if node == nil {
// 		return false
// 	}

// }

// func (tree *SourceTree) findSourceAndOperatorMoreThanCertainTime(time int) bool {
// 	if tree.root == nil {
// 		fmt.Println("The tree is empty")
// 		return nil
// 	}

// }

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

func (node *OperatorNode) findStatus4() {
	if node == nil {
		return
	}

	if node.sourceOperatorInfo.numberOfMessagesHaveStatus4 >= 1 {
		fmt.Println(node.sourceOperatorInfo.operatorName)
	}
	node.leftOperator.findStatus4()
	node.rightOperator.findStatus4()
}
func (node *SourceNode) findStatus4() {
	if node == nil {
		return
	}
	node.operatorTreeRoot.findStatus4()
	node.leftSource.findStatus4()
	node.rightSource.findStatus4()
}

func (tree *SourceTree) findStatus4() {
	if tree.root == nil {
		return
	}
	tree.root.findStatus4()
}
