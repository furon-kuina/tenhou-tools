package cmd

import (
	"encoding/xml"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"time"
)

type MJLogXML struct {
	XMLName  xml.Name `xml:"mjloggm"`
	Version  string   `xml:"ver,attr"`
	XMLNodes []node   `xml:",any"`
}

type node struct {
	XMLName xml.Name
	Attr    []xml.Attr `xml:",any,attr"`
	Nodes   []node     `xml:",any"`
	Content string     `xml:",chardata"`
}

func parseMJLog(r io.Reader) (*GameDetail, error) {
	var mjlog MJLogXML
	err := xml.NewDecoder(r).Decode(&mjlog)
	if err != nil {
		return nil, err
	}
	if len(mjlog.XMLNodes) < 4 {
		return nil, fmt.Errorf(`not enough nodes in mjlog: %+v`, mjlog.XMLNodes)
	}
	shuffleNode := mjlog.XMLNodes[0]
	goNode := mjlog.XMLNodes[1]
	unNode := mjlog.XMLNodes[2]
	taikyokuNode := mjlog.XMLNodes[3]
	if shuffleNode.XMLName.Local != "SHUFFLE" {
		return nil, fmt.Errorf(`expected "SHUFFLE" node, got %s`, shuffleNode.XMLName.Local)
	}
	if goNode.XMLName.Local != "GO" {
		return nil, fmt.Errorf(`expected "GO" node, got %s`, shuffleNode.XMLName.Local)
	}
	if unNode.XMLName.Local != "UN" {
		return nil, fmt.Errorf(`expected "UN" node, got %s`, shuffleNode.XMLName.Local)
	}
	if taikyokuNode.XMLName.Local != "TAIKYOKU" {
		return nil, fmt.Errorf(`expected "TAIKYOKU" node, got %s`, shuffleNode.XMLName.Local)
	}

	var game GameDetail

	game.Gametype.Sanma = !containsAttr(unNode.Attr, "n3")
	playerNum := 4
	if game.Gametype.Sanma {
		playerNum = 3
	}
	// parsing player info
	n0Attr, err := getAttr(unNode.Attr, "n0")
	if err != nil {
		return nil, fmt.Errorf("parsing unNode: %w", err)
	}
	n1Attr, err := getAttr(unNode.Attr, "n1")
	if err != nil {
		return nil, fmt.Errorf("parsing unNode: %w", err)
	}
	n2Attr, err := getAttr(unNode.Attr, "n2")
	if err != nil {
		return nil, fmt.Errorf("parsing unNode: %w", err)
	}
	if game.Gametype.Sanma {
		game.Players = make([]Player, 3)
	} else {
		game.Players = make([]Player, 4)
	}
	game.Players[0].Name = n0Attr.Value
	game.Players[1].Name = n1Attr.Value
	game.Players[2].Name = n2Attr.Value
	if game.Gametype.Sanma {
		n3Attr, err := getAttr(unNode.Attr, "n3")
		if err != nil {
			return nil, fmt.Errorf("parsing unNode: %w", err)
		}
		game.Players[3].Name = n3Attr.Value
	}

	danAttr, err := getAttr(unNode.Attr, "dan")
	if err != nil {
		return nil, fmt.Errorf("unNode: %v", err)
	}
	splitDanAttr := strings.Split(danAttr.Value, ",")
	if len(splitDanAttr) != playerNum {
		return nil, fmt.Errorf("dan attr and player number doesn't match")
	}
	for i, danStr := range splitDanAttr {
		dan, err := strconv.Atoi(danStr)
		if err != nil {
			return nil, fmt.Errorf("strconv.Atoi(%q): %v", danStr, err)
		}
		game.Players[i].Dan = dan
	}
	rateAttr, err := getAttr(unNode.Attr, "rate")
	if err != nil {
		return nil, fmt.Errorf("unNode: %v", err)
	}
	splitRateAttr := strings.Split(rateAttr.Value, ",")
	if len(splitDanAttr) != playerNum {
		return nil, fmt.Errorf("rate attr and player number doesn't match")
	}
	for i, rateStr := range splitRateAttr {
		rate, err := strconv.ParseFloat(rateStr, 32)
		if err != nil {
			return nil, fmt.Errorf("strconv.Atoi(%q): %v", rateStr, err)
		}
		game.Players[i].Rate = float32(rate)
	}

	return &game, nil
}

func getAttr(attrs []xml.Attr, target string) (xml.Attr, error) {
	targetIdx := slices.IndexFunc(attrs, func(attr xml.Attr) bool { return attr.Name.Local == target })
	if targetIdx == -1 {
		return xml.Attr{}, fmt.Errorf("findAttr(%+v, %q): not found", attrs, target)
	}
	return attrs[targetIdx], nil
}

func containsAttr(attrs []xml.Attr, target string) bool {
	targetIdx := slices.IndexFunc(attrs, func(attr xml.Attr) bool { return attr.Name.Local == target })
	return targetIdx != -1
}

func getKyoutaku(attrs []xml.Attr) (int, error) {
	ba, err := getAttr(attrs, "ba")
	if err != nil {
		return 0, fmt.Errorf(`findAttr(%+v, "ba"): %w`, attrs, err)
	}
	baSpl := strings.Split(ba.Value, ",")
	tmp, err := strconv.Atoi(baSpl[1])
	if err != nil {
		return 0, fmt.Errorf(`strconv.Atoi(%s): %v`, baSpl[1], err)
	}
	return tmp * 10, nil
}

// Score delta of each round is recorded as
// "92,-3,108,-3,564,-5,236,11" in mjlog,
// where 2i'th field denotes the score of i'th player at
// the start of the round
// 2i+1'th field denotes the delta of the score
func parseScoreDelta(s string) (before []int, delta []int, err error) {
	spl := strings.Split(s, ",")
	for i, str := range spl {
		tmp, err := strconv.Atoi(str)
		if err != nil {
			return nil, nil, fmt.Errorf("strconv.Atoi(%q)", str)
		}
		if i%2 == 0 {
			before = append(before, tmp)
		} else {
			delta = append(delta, tmp)
		}
	}
	return before, delta, nil
}

const (
	_ = iota
	Ippan
	Joukyuu
	Tokujou
	Houou
)

const (
	_ = iota
	NashiNashi
	AriAri
)

type GameDetail struct {
	Time     time.Time
	Gametype GameType
	Players  []Player
}

const (
	ippan = iota
	joukyuu
	tokujou
	houou
)

type GameType struct {
	Cpu        bool
	Aka        bool
	NashiNashi bool
	Tonpuu     bool
	Sanma      bool
	Fast       bool
	Table      int
}

type Player struct {
	Name string
	Dan  int
	Rate float32
}

func parseGameType(typeInt int) GameType {
	gt := GameType{}
	gt.Cpu = (typeInt & 0x1) == 0
	gt.Aka = (typeInt & 0x2) == 0
	gt.NashiNashi = (typeInt & 0x4) != 0
	gt.Tonpuu = (typeInt & 0x8) == 0
	gt.Sanma = (typeInt & 0x10) != 0
	gt.Fast = (typeInt & 0x40) != 0

	x := typeInt & 0x20
	y := typeInt & 0x80
	if x == 0 && y == 0 {
		gt.Table = ippan
	}
	if x == 0 && y != 0 {
		gt.Table = tokujou
	}
	if x != 0 && y == 0 {
		gt.Table = joukyuu
	}
	if x != 0 && y != 0 {
		gt.Table = houou
	}
	return gt
}
