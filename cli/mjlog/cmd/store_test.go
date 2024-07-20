package cmd

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func TestAddGame(t *testing.T) {
	ctx := context.Background()
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		t.Fatalf("env.Parse(%+v): %v", *cfg, err)
	}
	db, close, err := New(ctx, *cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer close()
	history := `23:56 | 13 | 三鳳南喰赤－ | <a href="http://tenhou.net/0/?log=2022121923gm-00b9-0000-62eb14e5">牌譜</a> | 砂吹きん(+44.3) 冬灯(+8.5) らんちぇすた(-52.8)<br>
23:56 | 28 | 四鳳南喰赤－ | <a href="http://tenhou.net/0/?log=2022121923gm-00a9-0000-848364ac">牌譜</a> | 百野さん(+53.8) ジミー(+10.5) tk-benzo(-18.7) 紫の鏡(-45.6)<br>
23:34 | 21 | 四上南喰赤－ | noxusan(+48.9) <br>(+12.6) |||||(-19.5) DDT-X(-42.0)
23:34 | 09 | 三上南喰赤－ | Kizao(+56.1) バッチ(-12.4) ＿(-43.7)
14:22 | 23 | 四般南－－－ | kamezawa(+62.8) 青獅子丸(+13.6) idayu2(-25.9) KT072511(-50.5)`
	p := newParser(strings.NewReader(history))
	gi := p.games()
	if p.err != nil {
		t.Fatalf("parser returned error: %v", p.err)
	}
	for _, game := range gi {
		err := AddGame(ctx, game, db)
		if err != nil {
			t.Errorf("AddGame: %v", err)
		}
	}
}

func TestListGames(t *testing.T) {
	ctx := context.Background()
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		t.Fatalf("env.Parse(%+v): %v", *cfg, err)
	}
	db, close, err := New(ctx, *cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer close()

	prepareGames(ctx, t, db)

	yonmaGames, err := ListYonmaGame(ctx, db)
	if err != nil {
		t.Fatalf("ListYonmaGame: %v", err)
	}
	for _, game := range yonmaGames {
		fmt.Printf("game :%+v", *game)
	}

	sanmaGames, err := ListSanmaGame(ctx, db)
	if err != nil {
		t.Fatalf("ListSanmaGame: %v", err)
	}
	for _, game := range sanmaGames {
		fmt.Printf("games:%+v", *game)
	}
}

func prepareGames(ctx context.Context, t *testing.T, db *sqlx.DB) {
	t.Helper()

	sql1 := `DELETE FROM sanma_game;`
	if _, err := db.ExecContext(ctx, sql1); err != nil {
		t.Logf("failed to initialize task: %v", err)
	}
	sql2 := `DELETE FROM yonma_game;`
	if _, err := db.ExecContext(ctx, sql2); err != nil {
		t.Logf("failed to initialize task: %v", err)
	}
	history := `23:56 | 13 | 三鳳南喰赤－ | <a href="http://tenhou.net/0/?log=2022121923gm-00b9-0000-62eb14e5">牌譜</a> | 砂吹きん(+44.3) 冬灯(+8.5) らんちぇすた(-52.8)<br>
23:56 | 28 | 四鳳南喰赤－ | <a href="http://tenhou.net/0/?log=2022121923gm-00a9-0000-848364ac">牌譜</a> | 百野さん(+53.8) ジミー(+10.5) tk-benzo(-18.7) 紫の鏡(-45.6)<br>
23:34 | 21 | 四上南喰赤－ | noxusan(+48.9) <br>(+12.6) |||||(-19.5) DDT-X(-42.0)
23:34 | 09 | 三上南喰赤－ | Kizao(+56.1) バッチ(-12.4) ＿(-43.7)
14:22 | 23 | 四般南－－－ | kamezawa(+62.8) 青獅子丸(+13.6) idayu2(-25.9) KT072511(-50.5)`
	p := newParser(strings.NewReader(history))
	gi := p.games()
	if p.err != nil {
		t.Fatalf("parser returned error: %v", p.err)
	}
	for _, game := range gi {
		err := AddGame(ctx, game, db)
		if err != nil {
			t.Fatalf("AddGame: %v", err)
		}
	}
}
