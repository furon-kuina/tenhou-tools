package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("./examples/test.mjlog")
	if err != nil {
		log.Fatalln(err)
	}
	err = ParseMJLog(f)
	if err != nil {
		fmt.Println(err)
	}
}

type MJLogXML struct {
	XMLName  xml.Name `xml:"mjloggm"`
	Version  string   `xml:"ver,attr"`
	XMLNodes []Node   `xml:",any"`
}

type Node struct {
	XMLName xml.Name
	Attr    []xml.Attr `xml:",any,attr"`
	Nodes   []Node     `xml:",any"`
	Content string     `xml:",chardata"`
}

func ParseMJLog(r io.Reader) error {
	var mjlog MJLogXML
	err := xml.NewDecoder(r).Decode(&mjlog)
	if err != nil {
		return err
	}
	if len(mjlog.XMLNodes) < 4 {
		return fmt.Errorf(`not enough nodes in mjlog: %+v`, mjlog.XMLNodes)
	}
	shuffleNode := mjlog.XMLNodes[0]
	goNode := mjlog.XMLNodes[1]
	unNode := mjlog.XMLNodes[2]
	taikyokuNode := mjlog.XMLNodes[3]
	if shuffleNode.XMLName.Local != "SHUFFLE" {
		return fmt.Errorf(`expected "SHUFFLE" node, got %s`, shuffleNode.XMLName.Local)
	}
	if goNode.XMLName.Local != "GO" {
		return fmt.Errorf(`expected "GO" node, got %s`, shuffleNode.XMLName.Local)
	}
	if unNode.XMLName.Local != "UN" {
		return fmt.Errorf(`expected "UN" node, got %s`, shuffleNode.XMLName.Local)
	}
	if taikyokuNode.XMLName.Local != "TAIKYOKU" {
		return fmt.Errorf(`expected "TAIKYOKU" node, got %s`, shuffleNode.XMLName.Local)
	}

	game := Game{
		Seed:    "",
		Cpu:     false,
		Red:     false,
		Samma:   !containsAttr(unNode.Attr, "n3"),
		Tonnan:  false,
		Fast:    false,
		Level:   0,
		Lobby:   "",
		Players: []Player{},
		Rounds:  []Round{},
		Scores:  []int{},
		Ranking: []int{},
	}

	// parsing player info
	if game.Samma {
		game.Players = make([]Player, 3)
	} else {
		game.Players = make([]Player, 4)
	}
	n0Attr, err := getAttr(unNode.Attr, "n0")
	if err != nil {
		return fmt.Errorf("parsing unNode: %w", err)
	}
	n1Attr, err := getAttr(unNode.Attr, "n1")
	if err != nil {
		return fmt.Errorf("parsing unNode: %w", err)
	}
	n2Attr, err := getAttr(unNode.Attr, "n2")
	if err != nil {
		return fmt.Errorf("parsing unNode: %w", err)
	}
	game.Players[0].Name = n0Attr.Value
	game.Players[1].Name = n1Attr.Value
	game.Players[2].Name = n2Attr.Value
	if !game.Samma {
		n3Attr, err := getAttr(unNode.Attr, "n3")
		if err != nil {
			return fmt.Errorf("parsing unNode: %w", err)
		}
		game.Players[3].Name = n3Attr.Value
	}
	fmt.Println(game.Players)

	game.Scores = make([]int, len(game.Players))
	scores := []int{250, 250, 250, 250}
	for _, node := range mjlog.XMLNodes {
		switch node.XMLName.Local {
		case "AGARI":
			fmt.Println("AGARI node:", node)
			scIdx := slices.IndexFunc(node.Attr, func(attr xml.Attr) bool { return attr.Name.Local == "sc" })
			if scIdx == -1 {
				return fmt.Errorf(`attribute "sc" not found in %+v`, node.Attr)
			}
			scAttr := node.Attr[scIdx]
			before, delta, _ := parseScoreDelta(scAttr.Value)
			for i := range len(before) {
				scores[i] += delta[i]
			}
			kyoutaku, err := getKyoutaku(node.Attr)
			fmt.Println("kyoutaku:T", kyoutaku)
			if err != nil {
				return fmt.Errorf("getKyoutaku(%+v): %w", node.Attr, err)
			}
			agariPlayer, err := getAttr(node.Attr, "who")
			if err != nil {
				return fmt.Errorf("searching Agari Player: %w", err)
			}
			agariPlayerId, err := strconv.Atoi(agariPlayer.Value)
			if err != nil {
				return fmt.Errorf("searching Agari Player: %w", err)
			}
			scores[agariPlayerId] += kyoutaku
			fmt.Printf("before: %d, delta: %d, after: %d\n", before, delta, scores)
			owariIdx := slices.IndexFunc(node.Attr, func(attr xml.Attr) bool { return attr.Name.Local == "owari" })
			if owariIdx == -1 {
				continue
			}
			gameResultRaw := strings.Split(node.Attr[owariIdx].Value, ",")
			for i := range len(game.Players) {
				game.Scores[i], err = strconv.Atoi(gameResultRaw[2*i])
			}
			if err != nil {
				return fmt.Errorf(`strconv.Atoi(%v): %w`, gameResultRaw, err)
			}

		case "RYUUKYOKU":
			fmt.Println("RYUUKYOKU node:", node)
			scIdx := slices.IndexFunc(node.Attr, func(attr xml.Attr) bool { return attr.Name.Local == "sc" })
			if scIdx == -1 {
				return fmt.Errorf(`attribute "sc" not found in %+v`, node.Attr)
			}
			scAttr := node.Attr[scIdx]
			before, delta, _ := parseScoreDelta(scAttr.Value)
			for i := range len(before) {
				scores[i] += delta[i]
			}
			owariIdx := slices.IndexFunc(node.Attr, func(attr xml.Attr) bool { return attr.Name.Local == "owari" })
			if owariIdx == -1 {
				continue
			}
			gameResultRaw := strings.Split(node.Attr[owariIdx].Value, ",")
			for i := range len(game.Players) {
				game.Scores[i], err = strconv.Atoi(gameResultRaw[2*i])
			}
			if err != nil {
				return fmt.Errorf(`strconv.Atoi(%+v): %w`, gameResultRaw, err)
			}
			fmt.Printf("before: %d, delta: %d, after: %d\n", before, delta, scores)

		default:

		}
	}
	fmt.Println("game.Scores:", game.Scores)
	return nil
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

type Game struct {
	Seed    string
	Cpu     bool
	Red     bool
	Samma   bool
	Tonnan  bool
	Fast    bool
	Level   int
	Lobby   string
	Players []Player
	Rounds  []Round
	Scores  []int
	Ranking []int
}

type Player struct {
	Name string
	Dan  int
	Rate float64
}

type Round struct {
	Stage    int //East1=0, ..., West4=11
	Homba    int
	Kyoutaku int
	Dora     int
	Oya      int
	Scores   []int
}

type Ryuukyoku struct {
	ScoreDelta []int
	PostScore  []int
}

type Agari struct {
	ScoreDelta      []int
	PostScore       []int
	GameResultScore []int
	GameResultPoint []int
}
