package cmd

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/utf8string"
)

type historyParser struct {
	scanner bufio.Scanner
	err     error
}

func newParser(src io.Reader) *historyParser {
	return &historyParser{
		scanner: *bufio.NewScanner(src),
	}
}

func (p *historyParser) games() []*Game {
	if p.err != nil {
		return nil
	}
	var res []*Game
	for p.scanner.Scan() {
		res = append(res, p.parseHistoryLine())
	}
	return res
}

func (p *historyParser) parseHistoryLine() *Game {
	if p.err != nil {
		return nil
	}
	line := p.scanner.Text()

	isHouou := strings.Contains(line, "http://tenhou.net")

	fieldNum := 4
	if isHouou {
		fieldNum = 5
	}

	splitLine := strings.SplitN(strings.TrimSuffix(line, "<br>"), "|", fieldNum)
	if len(splitLine) != fieldNum {
		p.err = fmt.Errorf("a history line has unexpected format: %q", line)
		return nil
	}
	timeStr := strings.TrimSpace(splitLine[0])
	hour, err := strconv.Atoi(timeStr[:2])
	if err != nil {
		p.err = fmt.Errorf("error parsing time %q: %v", splitLine[0], err)
		return nil
	}
	minute, err := strconv.Atoi(timeStr[3:])
	if err != nil {
		p.err = fmt.Errorf("error parsing time %q: %v", splitLine[0], err)
	}

	res := &Game{
		Time:     time.Date(0, 1, 1, hour, minute, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
		GameType: p.parseGameType(strings.TrimSpace(splitLine[2])),
		PlayerID: make([]string, 4),
		Scores:   make([]int, 4),
	}
	if res.GameType.Sanma {
		res.PlayerID = make([]string, 3)
		res.Scores = make([]int, 3)
	}
	// 鳳凰卓の場合
	if isHouou {
		splitATag := strings.Split(splitLine[3], `"`)
		if len(splitATag) != 3 {
			p.err = fmt.Errorf("a tag has unexpected format: %q", splitLine[3])
			return nil
		}
		res.LogURL = strings.Replace(splitATag[1], "?log=", "log/?", 1)
	}

	nameAndScores := strings.Fields(splitLine[fieldNum-1])
	for i, field := range nameAndScores {
		scoreStart := strings.LastIndex(field, "(")
		if scoreStart == -1 {
			p.err = fmt.Errorf("a history line has unexpected format: %q", line)
			return nil
		}
		name := field[:scoreStart]
		scoreString := field[scoreStart+1 : len(field)-1]
		res.PlayerID[i] = name
		res.Scores[i], p.err = scoreStringToInt(scoreString)
	}
	return res
}

type Game struct {
	Time time.Time
	// 牌譜へのリンク 鳳凰卓以外は空文字列
	LogURL   string
	GameType GameType
	// 天鳳ID 順位順
	PlayerID []string
	// 終局時の点棒 順位順
	Scores []int
}

// (+|-) xy.z -> +/-xyz
func scoreStringToInt(scoreStr string) (int, error) {
	splitScoreStr := strings.Split(scoreStr, ".")
	sign := scoreStr[0]
	s1 := splitScoreStr[0][1:]
	s2 := splitScoreStr[1]
	v1, err := strconv.Atoi(s1)
	if err != nil {
		return 0, err
	}
	v2, err := strconv.Atoi(s2)
	if err != nil {
		return 0, err
	}
	res := v1*10 + v2
	if sign == '-' {
		res *= -1
	}
	return res, nil
}

// "四鳳東喰赤ー"形式の文字列をパースする
func (p *historyParser) parseGameType(gameTypeStr string) GameType {
	if p.err != nil {
		return GameType{}
	}
	gameTypeRunes := utf8string.NewString(gameTypeStr)
	if gameTypeRunes.RuneCount() != 6 {
		p.err = fmt.Errorf("unexpected gametype %q", gameTypeStr)
		return GameType{}
	}
	gt := GameType{}
	gt.Sanma = gameTypeRunes.At(0) == '三'
	gt.Table = getTableConstant(gameTypeRunes.At(1))
	gt.Tonpuu = gameTypeRunes.At(2) == '東'
	gt.NashiNashi = gameTypeRunes.At(3) != '喰'
	gt.Aka = gameTypeRunes.At(4) == '赤'

	return gt
}

// 般 -> 1
// 上 -> 2
// 特 -> 3
// 鳳 -> 4
func getTableConstant(ch rune) int {
	switch ch {
	case '般':
		return ippan
	case '上':
		return joukyuu
	case '特':
		return tokujou
	case '鳳':
		return houou
	default:
		return -1
	}
}
