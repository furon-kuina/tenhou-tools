package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGameInfo(t *testing.T) {
	give := `23:56 | 13 | 三鳳南喰赤－ | <a href="http://tenhou.net/0/?log=2022121923gm-00b9-0000-62eb14e5">牌譜</a> | 砂吹きん(+44.3) 冬灯(+8.5) らんちぇすた(-52.8)<br>
23:56 | 28 | 四鳳南喰赤－ | <a href="http://tenhou.net/0/?log=2022121923gm-00a9-0000-848364ac">牌譜</a> | 百野さん(+53.8) ジミー(+10.5) tk-benzo(-18.7) 紫の鏡(-45.6)<br>
23:34 | 21 | 四上南喰赤－ | noxusan(+48.9) <br>(+12.6) |||||(-19.5) DDT-X(-42.0)
23:34 | 09 | 三上南喰赤－ | Kizao(+56.1) バッチ(-12.4) ＿(-43.7)
14:22 | 23 | 四般南－－－ | kamezawa(+62.8) 青獅子丸(+13.6) idayu2(-25.9) KT072511(-50.5)`
	want := []*Game{
		{
			Time:     time.Date(0, 1, 1, 23, 56, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
			PlayerID: []string{"砂吹きん", "冬灯", "らんちぇすた"},
			Scores:   []int{443, 85, -528},
			LogURL:   "http://tenhou.net/0/log/?2022121923gm-00b9-0000-62eb14e5",
			GameType: GameType{
				Aka:   true,
				Sanma: true,
				Table: houou,
			},
		},
		{
			Time:     time.Date(0, 1, 1, 23, 56, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
			PlayerID: []string{"百野さん", "ジミー", "tk-benzo", "紫の鏡"},
			Scores:   []int{538, 105, -187, -456},
			LogURL:   "http://tenhou.net/0/log/?2022121923gm-00a9-0000-848364ac",
			GameType: GameType{
				Aka:   true,
				Table: houou,
			},
		},
		{
			Time:     time.Date(0, 1, 1, 23, 34, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
			PlayerID: []string{"noxusan", "<br>", "|||||", "DDT-X"},
			Scores:   []int{489, 126, -195, -420},
			GameType: GameType{
				Aka:   true,
				Table: joukyuu,
			},
		}, {
			Time:     time.Date(0, 1, 1, 23, 34, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
			PlayerID: []string{"Kizao", "バッチ", "＿"},
			Scores:   []int{561, -124, -437},
			GameType: GameType{
				Aka:   true,
				Sanma: true,
				Table: joukyuu,
			},
		}, {
			Time:     time.Date(0, 1, 1, 14, 22, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
			PlayerID: []string{"kamezawa", "青獅子丸", "idayu2", "KT072511"},
			Scores:   []int{628, 136, -259, -505},
			GameType: GameType{
				NashiNashi: true,
				Table:      ippan,
			},
		},
	}
	p := newParser(strings.NewReader(give))
	res := p.games()
	if p.err != nil {
		t.Fatalf("parser returned error: %v", p.err)
	}
	if !cmp.Equal(want, res) {
		t.Errorf("unexpected result: %s", cmp.Diff(want, res))
	}
}
