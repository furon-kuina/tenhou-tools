package cmd

import (
	"testing"
)

// func TestDownload(t *testing.T) {
// 	history := `23:56 | 13 | 三鳳南喰赤－ | <a href="http://tenhou.net/0/?log=2022121923gm-00b9-0000-62eb14e5">牌譜</a> | 砂吹きん(+44.3) 冬灯(+8.5) らんちぇすた(-52.8)<br>
// 	23:56 | 28 | 四鳳南喰赤－ | <a href="http://tenhou.net/0/?log=2022121923gm-00a9-0000-848364ac">牌譜</a> | 百野さん(+53.8) ジミー(+10.5) tk-benzo(-18.7) 紫の鏡(-45.6)<br>
// 	23:58 | 16 | 三鳳南喰赤－ | <a href="http://tenhou.net/0/?log=2022121923gm-00b9-0000-c6d2ff99">牌譜</a> | moly(+53.6) Hectorり(+18.1) 蒼月★.com(-71.7)<br>`

// 	p := newParser(strings.NewReader(history))
// 	res := p.getUrls()

// 	for _, url := range res {
// 		fmt.Println(url)
// 		resp, err := http.Get(url)
// 		if err != nil {
// 			t.Errorf("GET %s failed: %v", url, err)
// 		}
// 		game, err := parseMJLog(resp.Body)
// 		if err != nil {
// 			t.Errorf("failed to parse MJLog")
// 		}
// 		fmt.Printf("game data: %+v", game)
// 	}
// }

// func TestParseMJLog(t *testing.T) {
// 	cases := []struct {
// 		logFilePath string
// 		want        Game
// 	}{
// 		{"./testdata/test.mjlog", Game{Players: []Player{
// 			{url.QueryEscape("段原瑠々"), 3, 16, 2109.46},
// 			{url.QueryEscape("たなか教祖"), 4, 17, 2189.20},
// 			{url.QueryEscape("すれいか"), 2, 16, 2097.98},
// 			{url.QueryEscape("ふろん"), 1, 16, 2177.45},
// 		}},
// 		}}

// 	for _, tt := range cases {
// 		f, err := os.Open(tt.logFilePath)
// 		if err != nil {
// 			t.Fatalf("os.Open(%q): %v", tt.logFilePath, err)
// 		}
// 		game, err := parseMJLog(f)
// 		if err != nil {
// 			t.Fatalf("parseMJLog: %v", err)
// 		}
// 		if !cmp.Equal(*game, tt.want) {
// 			t.Errorf("diff: %s", cmp.Diff(*game, tt.want))
// 		}

// 	}

// }

func TestParseGameType(t *testing.T) {
	cases := []struct {
		give int
		want GameType
	}{
		{169, GameType{false, true, false, false, false, false, houou}},
	}
	for _, tt := range cases {
		res := parseGameType(tt.give)
		if tt.want != res {
			t.Errorf("want %v, got %v", tt.want, res)
		}
	}
}
