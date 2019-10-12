package main

import (
	"fmt"
	"pes2020/pkg/util"
	"strings"
)

func main() {
	couchList := initCouchList()
	printOppositeTacticsCouch(couchList)
}

var tacticsNum = 7

// 初始化教练列表
func initCouchList() [] Couch {

	lines := util.ReadLines("couch.txt")
	var couchList []Couch
	for _, line := range lines {
		couchDataArr := strings.Split(line, "\t")
		// 初始化数据
		offenceCouch := NewCouch(couchDataArr, true)
		defenceCouch := NewCouch(couchDataArr, false)
		couchList = append(couchList, offenceCouch, defenceCouch)
	}
	return couchList
}

func printOppositeTacticsCouch(couchList []Couch) {
	couchSize := len(couchList)
	for i := 0; i < couchSize; i++ {
		for j := i + 1; j < couchSize; j++ {
			sameTacticsCount := 0
			for z := 0; z < tacticsNum; z++ {
				if couchList[i].ActiveTactics[z] == couchList[j].ActiveTactics[z] {
					sameTacticsCount++
				}
			}
			if sameTacticsCount == 0 {
				fmt.Printf("%s_%s_%s/%s_%s_%s\n", couchList[i].Name, couchList[i].TacticsName, couchList[i].ActiveFormation, couchList[j].Name, couchList[j].TacticsName, couchList[j].ActiveFormation)
			}
		}
	}
}

func NewCouch(couchDataArr []string, isOffence bool) Couch {
	var couch = Couch{
		Name:             couchDataArr[0],
		OffenceFormation: couchDataArr[7],
		OffenceTactics:   []string{couchDataArr[8], couchDataArr[9], couchDataArr[10], couchDataArr[11], couchDataArr[12], couchDataArr[13], couchDataArr[14]},
		DefenceFormation: couchDataArr[15],
		DefenceTactics:   []string{couchDataArr[16], couchDataArr[17], couchDataArr[18], couchDataArr[19], couchDataArr[20], couchDataArr[21], couchDataArr[22]},
	}
	if isOffence {
		couch.ActiveFormation = couch.OffenceFormation
		couch.ActiveTactics = couch.OffenceTactics
		couch.TacticsName = "进攻"
	} else {
		couch.ActiveFormation = couch.DefenceFormation
		couch.ActiveTactics = couch.DefenceTactics
		couch.TacticsName = "防守"
	}
	return couch
}

type Couch struct {
	Name             string
	OffenceFormation string
	OffenceTactics   []string
	DefenceFormation string
	DefenceTactics   []string
	ActiveFormation  string
	ActiveTactics    []string
	TacticsName      string
}
