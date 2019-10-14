package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"pes2020/pkg/util"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func main() {
	util.RenewFile(resultFile)
	// 30级总评>=80
	fetchData("overall", 80, 0, false)
	// 满级总评>=85
	//fetchData("potential", 0, 85, true)
	// 30级总评>=85,满级总评>=90
	fetchData("elite", 85, 90, true)
}

var cookies = map[string]string{
	"columns": "pos%2Cname%2Cclub_team%2Cleague%2Cnationality%2Cheight%2Cweight%2Cage%2Cfoot%2Coffensive_awareness%2Cball_control%2Cdribbling" +
		"%2Ctight_possession%2Clow_pass%2Clofted_pass%2Cfinishing%2Cheading%2Cplace_kicking%2Ccurl%2Cspeed%2Cacceleration%2Ckicking_power%2Cjump" +
		"%2Cphysical_contact%2Cbalance%2Cstamina%2Cdefensive_awareness%2Cball_winning%2Caggression%2Cgk_awareness%2Cgk_catching%2Cgk_clearing" +
		"%2Cgk_reflexes%2Cgk_reach%2Cweak_foot_usage%2Cweak_foot_accuracy%2Cform%2Cinjury_resistance%2Coverall_rating%2Cmax_level%2Coverall_at_max_level",
}
var headers = map[string]string{
	"User-Agent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.89 Safari/537.36",
}

var footMap = map[string]string{"Right foot": "右脚", "Left foot": "左脚"}
var posMap = map[string]string{
	"CF": "中锋", "RWF": "右边锋", "LWF": "左边锋", "SS": "影锋", "AMF": "前腰", "CMF": "中场", "RMF": "右前卫",
	"LMF": "左前卫", "DMF": "后腰", "CB": "中后卫", "RB": "右后卫", "LB": "左后卫", "GK": "门将",
}
var growthMap = readConfig("growth.txt")
var clubTeamMap = readConfig("club_team.txt")
var nationalityMap = readConfig("nationality.txt")
var leagueMap = readConfig("league.txt")
var resultFile = "result.txt"

func readConfig(fileName string) map[string]string {
	var lines = util.ReadLines(fileName)
	configMap := make(map[string]string)
	for _, line := range lines {
		keyValuePair := strings.Split(line, "\t")
		configMap[keyValuePair[0]] = keyValuePair[1]
	}
	return configMap
}

func fetchData(mode string, minRatingLevel30, minRatingLevelMax int, isMaxLevel bool) {
	page := 1
	playerId := (page - 1) * 32
	finished := false
	keySet := make(map[string]bool)
	for !finished {
		url := ""
		if mode == "overall" {
			url = "http://pesdb.net/pes2020/?page=" + strconv.Itoa(page)
		} else if mode == "potential" {
			url = "http://pesdb.net/pes2020/?sort=overall_at_max_level&page=" + strconv.Itoa(page)
		} else if mode == "elite" {
			url = "http://pesdb.net/pes2020/?sort=overall_at_max_level&overall_rating=" + strconv.Itoa(minRatingLevel30) + "-99&page=" + strconv.Itoa(page)
		}
		html := util.GetRequest(url, headers, cookies)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		playersTable := doc.Find("table[class=players]")
		playerList := playersTable.Find("tr")

		playerList.Each(func(i int, player *goquery.Selection) {
			playerData := player.Find("td")
			if playerData.Eq(0).Text() == "" {
				return
			}
			p := NewPlayer(playerData)
			if keySet[p.Key] == true {
				return
			} else {
				keySet[p.Key] = true
			}
			if mode == "overall" {
				if p.OverallRating < minRatingLevel30 {
					finished = true
					return
				}
			} else if mode == "potential" || mode == "elite" {
				if p.OverallAtMaxLevel < minRatingLevelMax {
					finished = true
					return
				}
			}
			if isMaxLevel {
				p = p.getMaxLevelData()
			}
			playerId++
			p.printData(playerId)
		})
		time.Sleep(time.Duration(3)*time.Second)
		page++
	}
}

func NewPlayer(playerData *goquery.Selection) Player {
	player := Player{
		Pos:                playerData.Eq(0).Text(),
		Name:               playerData.Eq(1).Text(),
		ClubTeam:           playerData.Eq(2).Text(),
		League:             playerData.Eq(3).Text(),
		Nationality:        playerData.Eq(4).Text(),
		Height:             playerData.Eq(5).Text(),
		Weight:             playerData.Eq(6).Text(),
		Age:                playerData.Eq(7).Text(),
		Foot:               playerData.Eq(8).Text(),
		OffensiveAwareness: util.ParseStrToInt(playerData.Eq(9).Text()),
		BallControl:        util.ParseStrToInt(playerData.Eq(10).Text()),
		Dribbling:          util.ParseStrToInt(playerData.Eq(11).Text()),
		TightPossession:    util.ParseStrToInt(playerData.Eq(12).Text()),
		LowPass:            util.ParseStrToInt(playerData.Eq(13).Text()),
		LoftedPass:         util.ParseStrToInt(playerData.Eq(14).Text()),
		Finishing:          util.ParseStrToInt(playerData.Eq(15).Text()),
		Heading:            util.ParseStrToInt(playerData.Eq(16).Text()),
		PlaceKicking:       util.ParseStrToInt(playerData.Eq(17).Text()),
		Curl:               util.ParseStrToInt(playerData.Eq(18).Text()),
		Speed:              util.ParseStrToInt(playerData.Eq(19).Text()),
		Acceleration:       util.ParseStrToInt(playerData.Eq(20).Text()),
		KickingPower:       util.ParseStrToInt(playerData.Eq(21).Text()),
		Jump:               util.ParseStrToInt(playerData.Eq(22).Text()),
		PhysicalContact:    util.ParseStrToInt(playerData.Eq(23).Text()),
		Balance:            util.ParseStrToInt(playerData.Eq(24).Text()),
		Stamina:            util.ParseStrToInt(playerData.Eq(25).Text()),
		DefensiveAwareness: util.ParseStrToInt(playerData.Eq(26).Text()),
		BallWinning:        util.ParseStrToInt(playerData.Eq(27).Text()),
		Aggression:         util.ParseStrToInt(playerData.Eq(28).Text()),
		GkAwareness:        util.ParseStrToInt(playerData.Eq(29).Text()),
		GkCatching:         util.ParseStrToInt(playerData.Eq(30).Text()),
		GkClearing:         util.ParseStrToInt(playerData.Eq(31).Text()),
		GkReflexes:         util.ParseStrToInt(playerData.Eq(32).Text()),
		GkReach:            util.ParseStrToInt(playerData.Eq(33).Text()),
		WeakFootUsage:      util.ParseStrToInt(playerData.Eq(34).Text()),
		WeakFootAccuracy:   util.ParseStrToInt(playerData.Eq(35).Text()),
		Form:               util.ParseStrToInt(playerData.Eq(36).Text()),
		InjuryResistance:   util.ParseStrToInt(playerData.Eq(37).Text()),
		OverallRating:      util.ParseStrToInt(playerData.Eq(38).Text()),
		MaxLevel:           util.ParseStrToInt(playerData.Eq(39).Text()),
		OverallAtMaxLevel:  util.ParseStrToInt(playerData.Eq(40).Text()),
		Key:                playerData.Eq(1).Text() + "_" + posMap[playerData.Eq(0).Text()] + "_" + playerData.Eq(38).Text(),
	}
	return player
}

func (p Player) printData(playerId int) {
	format := strings.Repeat("%s\t", 11) + strings.Repeat("%d\t", 32) + "%s\n"
	output := fmt.Sprintf(format, strconv.Itoa(playerId), p.Name, posMap[p.Pos], "", clubTeamMap[p.ClubTeam], leagueMap[p.League],
		nationalityMap[p.Nationality], p.Height, p.Weight, p.Age, footMap[p.Foot], p.OverallRating, p.MaxLevel, p.OverallAtMaxLevel,
		p.OffensiveAwareness, p.BallControl, p.Dribbling, p.TightPossession, p.LowPass, p.LoftedPass, p.Finishing, p.Heading, p.PlaceKicking,
		p.Curl, p.Speed, p.Acceleration, p.KickingPower, p.Jump, p.PhysicalContact, p.Balance, p.Stamina, p.DefensiveAwareness, p.BallWinning,
		p.Aggression, p.GkAwareness, p.GkCatching, p.GkClearing, p.GkReflexes, p.GkReach, p.WeakFootUsage, p.WeakFootAccuracy, p.Form,
		p.InjuryResistance, p.Key, )
	fmt.Print(output)
	util.AppendLine(resultFile, output)
}

func (p Player) getMaxLevelData() Player {
	for level := 2; level < p.MaxLevel+1; level++ {
		growthList := strings.Split(growthMap[strconv.Itoa(level)], ",")
		for _, attribute := range growthList {
			field := reflect.ValueOf(&p).Elem().FieldByName(attribute)
			attrValue := field.Int()
			if attrValue < 99 {
				field.SetInt(attrValue + 1)
			}
		}
	}
	if p.Pos != "GK" {
		p.GkAwareness, p.GkCatching, p.GkClearing, p.GkReflexes, p.GkReach = 40, 40, 40, 40, 40
	}
	return p
}

type Player struct {
	Pos, Name, ClubTeam, League, Nationality, Foot, Height, Weight, Age, Key string
	OffensiveAwareness, BallControl, Dribbling, TightPossession, LowPass, LoftedPass, Finishing, Heading, PlaceKicking, Curl, Speed,
	Acceleration, KickingPower, Jump, PhysicalContact, Balance, Stamina, DefensiveAwareness, BallWinning, Aggression, GkAwareness, GkCatching,
	GkClearing, GkReflexes, GkReach, WeakFootUsage, WeakFootAccuracy, Form, InjuryResistance, OverallRating, MaxLevel, OverallAtMaxLevel int
}
