package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// CDR Struct
type CDR struct {
	timestamp string
	state     string
	id        string
	src       string
	operator  string
	statusIND string
}

func rangeOfIndexDelay(delay int) int {
	if delay < 10 {
		return 0
	} else if delay < 60 {
		return 1
	} else if delay < 600 {
		return 2
	} else if delay < 3600 {
		return 3
	}
	return 4
}
func main() {
	// List of CDR that has under certain time
	// listDelayUnder10secs := []CDR{}
	// listDelayUnder60secs := []CDR{}
	// listDelayUnder10mins := []CDR{}
	// listDelayUnder1hour := []CDR{}
	// Need to create a hash table for easy look up
	cdrNTable := map[string]CDR{}
	cdrMTable := map[string]CDR{}

	// Need to create a hash table of hash table to store

	const longForm = "20060102150405"
	file, err := os.Open("../CDR_20180702_2.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), ";")
		// 1 -> SEQ_NUM
		// 2 -> ENTRY_TS
		// 3 -> STATE_TS
		// 4 -> TYPE
		// 5 -> STATE
		// 6 -> SRC_NAME
		// 7 -> SINK_NAME
		// 8 -> ROUTING_INDEX
		// 9 -> ROUTING_DEST
		// 10 -> ORIG_TON
		// 11 -> ORIG_NP
		// 12 -> ORIG
		// 13 -> DEST_IMSI
		// 14 -> DEST_MSC
		// 15 -> DEST_TON
		// 16 -> DEST_NP
		// 17 -> DEST_MSISDN
		// 18 -> DEST
		// 19 -> ID
		// 20 -> REF_ID
		// 21 -> ESME_ID
		// 22 -> ERROR
		// 23 -> EMSE_SPECIFIC_ERROR
		// 24 -> DR_ERROR_ORIGINAL
		// 25 -> DR_ERROR
		// 26 -> STATUS_IND
		// 27 -> ATTEMPTS
		// 28 -> DATA_LEN
		// 29 -> DCS
		// 39 -> PROTOCOL_ID
		// 31 -> RCPT_REQ
		// 32 -> USER_PROVIDED_ID

		operator := strings.Split(arr[9], ":")
		if arr[4] == "M" {
			if operator[0] != "null" {
				aCDR := CDR{
					timestamp: arr[2],
					state:     arr[5],
					id:        arr[19],
					src:       arr[6],
					operator:  operator[0],
				}
				_, ok := cdrMTable[arr[19]]
				if ok == false {
					cdrMTable[arr[19]] = aCDR
				}
			}
		} else {
			if operator[0] != "null" {
				aCDR := CDR{
					timestamp: arr[3],
					state:     arr[5],
					id:        arr[20],
					src:       arr[7],
					operator:  operator[0],
					statusIND: arr[27],
				}
				_, ok := cdrNTable[arr[20]]
				if ok == false {
					cdrNTable[arr[20]] = aCDR
				} else {
					cdr := cdrNTable[arr[20]]
					firstTimeStamp, _ := strconv.Atoi(cdr.timestamp)
					secondTimeStamp, _ := strconv.Atoi(aCDR.timestamp)
					if secondTimeStamp > firstTimeStamp {
						delete(cdrNTable, arr[20])
						cdrNTable[arr[20]] = aCDR
					}
				}
			}
		}
	}
	fmt.Println("Processing............................................")
	tree := SourceTree{
		root: nil,
	}

	for keyM := range cdrMTable {
		valueN, existed := cdrNTable[keyM]
		if existed {
			ops := SourceOperatorInfo{
				operatorName: valueN.operator,
			}
			tree.addSource(valueN.src, ops)
		}
	}

	fmt.Println("Processing Calculating ..................................")

	for keyM, valueM := range cdrMTable {
		valueN, existed := cdrNTable[keyM]
		if existed {
			status := valueN.statusIND
			timeM, _ := time.Parse(longForm, valueM.timestamp)
			timeN, _ := time.Parse(longForm, valueN.timestamp)
			delay := (timeN.Day()-timeM.Day())*84000 + (timeN.Hour()-timeM.Hour())*3600 + (timeN.Minute()-timeM.Minute())*60 + (timeN.Second() - timeM.Second())
			if status == "4" {
				tree.findAndIncrement("totalMessagesDLRS", valueN.operator, valueM.src)
			}
			if delay < 10 {
				tree.findAndIncrement("totalMessagesLessThan10seconds", valueN.operator, valueM.src)
			} else if delay < 60 {
				tree.findAndIncrement("totalMessagesLessThan1min", valueN.operator, valueM.src)
			} else if delay < 600 {
				tree.findAndIncrement("totalMessagesLessThan10mins", valueN.operator, valueM.src)
			} else if delay < 3600 {
				tree.findAndIncrement("totalMessagesLessThan1hour", valueN.operator, valueM.src)
			} else if delay < 7200 {
				tree.findAndIncrement("totalMessagesLessThan2hour", valueN.operator, valueM.src)
			}
			tree.findAndIncrement("totalMessagesReceived", valueN.operator, valueM.src)
		}
	}

	sourceName, val := tree.findTheMos()
	fmt.Println(sourceName, val)

	tree.findAndDisplaySource("VGE_C_ING_MING_7_mming")

	tree.calculateSDRPercentage()
	// // Remember how many of them are failing in certain interval of time
	// delayInterval := [5]int{0, 0, 0, 0, 0}
	// // Remember the logs for all different M that might catch later for 72 hours or 4 next files
	// var delayMList []string
	// // Backtrace the N table to find the matching pair with M table
	// for key, value := range cdrMTable {
	// 	cdr, ok := cdrNTable[key]
	// 	// if Ok it means if the M key maps to N in one CDR File
	// 	if ok {
	// 		// fmt.Println("ID: ", value.id)
	// 		timeM, _ := time.Parse(longForm, value.timestamp)
	// 		timeN, _ := time.Parse(longForm, cdr.timestamp)
	// 		delay := (timeN.Day()-timeM.Day())*84000 + (timeN.Hour()-timeM.Hour())*3600 + (timeN.Minute()-timeM.Minute())*60 + (timeN.Second() - timeM.Second())
	// 		interval := rangeOfIndexDelay(delay)
	// 		delayInterval[interval]++
	// 		switch interval {
	// 		case 0:
	// 			listDelayUnder10secs = append(listDelayUnder10secs, cdr)
	// 		case 1:
	// 			listDelayUnder60secs = append(listDelayUnder60secs, cdr)
	// 		case 2:
	// 			listDelayUnder10mins = append(listDelayUnder10mins, cdr)
	// 		case 3:
	// 			listDelayUnder1hour = append(listDelayUnder1hour, cdr)
	// 		default:
	// 		}
	// 	} else {
	// 		delayMList = append(delayMList, value.id)
	// 	}
	// }
	// totalMessage := 0
	// for i := 0; i < 5; i++ {
	// 	totalMessage += delayInterval[i]
	// }

	// sourceTree := SourceTree{
	// 	root: nil,
	// }

	// lengthUnder10mins := len(listDelayUnder10mins)
	// lengthUnder1min := len(listDelayUnder60secs)
	// lengthUnder1hour := len(listDelayUnder1hour)

	// // Processing 10 mins list
	// for i := 0; i < lengthUnder10mins; i++ {
	// 	sourceName := listDelayUnder10mins[i].src
	// 	operator := SourceOperatorInfo{
	// 		lessThan10mins:              1,
	// 		lessThan1min:                0,
	// 		lessThan1hour:               0,
	// 		numberOfMessagesHaveStatus4: 0,
	// 		operatorName:                listDelayUnder10mins[i].operator,
	// 	}
	// 	if strings.Compare(listDelayUnder10mins[i].statusIND, "4") == 0 {
	// 		operator.numberOfMessagesHaveStatus4 = 1
	// 	}
	// 	if sourceTree.findSource(listDelayUnder10mins[i].src) == nil {
	// 		sourceTree.addSource(sourceName, operator)
	// 	} else {
	// 		if sourceTree.root.operatorTreeRoot.findOperator(listDelayUnder10mins[i].operator) == nil {
	// 			sourceTree.root.operatorTreeRoot.addOperator(operator)
	// 		} else {

	// 			sourceTree.incrementValues(listDelayUnder10mins[i].statusIND, "10mins", sourceName, listDelayUnder10mins[i].operator)
	// 		}
	// 	}
	// }

	// // Processing 1 min list
	// for i := 0; i < lengthUnder1min; i++ {
	// 	sourceName := listDelayUnder60secs[i].src
	// 	operator := SourceOperatorInfo{
	// 		lessThan10mins:              0,
	// 		lessThan1min:                1,
	// 		lessThan1hour:               0,
	// 		numberOfMessagesHaveStatus4: 0,
	// 		operatorName:                listDelayUnder60secs[i].operator,
	// 	}
	// 	if strings.Compare(listDelayUnder60secs[i].statusIND, "4") == 0 {
	// 		operator.numberOfMessagesHaveStatus4 = 1
	// 	}
	// 	if sourceTree.findSource(listDelayUnder60secs[i].src) == nil {
	// 		sourceTree.addSource(sourceName, operator)
	// 	} else {
	// 		if sourceTree.root.operatorTreeRoot.findOperator(listDelayUnder60secs[i].operator) == nil {
	// 			sourceTree.root.operatorTreeRoot.addOperator(operator)
	// 		} else {
	// 			sourceTree.incrementValues(listDelayUnder60secs[i].statusIND, "1min", sourceName, listDelayUnder60secs[i].operator)
	// 		}
	// 	}
	// }
	// // Processing 1 hour test
	// for i := 0; i < lengthUnder1hour; i++ {
	// 	sourceName := listDelayUnder1hour[i].src
	// 	operator := SourceOperatorInfo{
	// 		lessThan10mins:              0,
	// 		lessThan1min:                0,
	// 		lessThan1hour:               1,
	// 		numberOfMessagesHaveStatus4: 0,
	// 		operatorName:                listDelayUnder1hour[i].operator,
	// 	}
	// 	if strings.Compare(listDelayUnder1hour[i].statusIND, "4") == 0 {
	// 		operator.numberOfMessagesHaveStatus4 = 1
	// 	}
	// 	if sourceTree.findSource(listDelayUnder1hour[i].src) == nil {
	// 		sourceTree.addSource(sourceName, operator)
	// 	} else {
	// 		if sourceTree.root.operatorTreeRoot.findOperator(listDelayUnder1hour[i].operator) == nil {
	// 			sourceTree.root.operatorTreeRoot.addOperator(operator)
	// 		} else {

	// 			sourceTree.incrementValues(listDelayUnder1hour[i].statusIND, "1hour", sourceName, listDelayUnder1hour[i].operator)
	// 		}
	// 	}
	// }

	// //sourceTree.display()

	// // fmt.Println("Total Messages Have Been Received", totalMessage)

	// // fmt.Printf("Percentage of Messages that has DLR under 10 seconds %.2f\n", (float64(delayInterval[0])/float64(totalMessage))*100)

	// // fmt.Printf("Percentage of Messages that has DLR under 1 minute %.2f\n", (float64(delayInterval[1])/float64(totalMessage))*100)

	// // fmt.Printf("Percentage of Messages that has DLR under 10 minutes %.2f\n", (float64(delayInterval[2])/float64(totalMessage))*100)

	// // fmt.Printf("Percentage of Messages that has DLR under 1 hour %.2f\n", (float64(delayInterval[3])/float64(totalMessage))*100)

	// sourceTree.findStatus4()

}
