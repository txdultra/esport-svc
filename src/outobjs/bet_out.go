package outobjs

import (
	"time"
	"tmp"
)

type OutBetSetting struct {
	MaxRatios  int `json:"max_ratios"`
	BaseStakes int `json:"base_stakes"`
}

type OutBetGame struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Img    int64  `json:"img_id"`
	ImgUrl string `json:"img_url"`
}

type OutBetMatch struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	SubTitle  string `json:"subtitle"`
	Img       int64  `json:"img_id"`
	ImgUrl    string `json:"img_url"`
	CanWagers int    `json:"can_wagers"`
}

type OutBetItem struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Img         int64  `json:"img_id"`
	ImgUrl      string `json:"img_url"`
}

type OutBetCompletetion struct {
	Id             int64                       `json:"id"`
	GameId         int64                       `json:"game_id"`
	Game           *OutBetGame                 `json:"game"`
	MatchId        int64                       `json:"match_id"`
	Match          *OutBetMatch                `json:"match"`
	Title          string                      `json:"title"`
	Description    string                      `json:"description"`
	Img            int64                       `json:"img_id"`
	ImgUrl         string                      `json:"img_url"`
	ItemAId        int64                       `json:"item_a_id"`
	ItemA          *OutBetItem                 `json:"item_a"`
	ItemBId        int64                       `json:"item_b_id"`
	ItemB          *OutBetItem                 `json:"item_b"`
	AScore         int                         `json:"a_score"`
	BScore         int                         `json:"b_score"`
	ABetters       int                         `json:"a_betters"`
	BBetters       int                         `json:"b_betters"`
	Winner         int64                       `json:"winner"`
	MatchStart     time.Time                   `json:"match_start"`
	MatchEnd       time.Time                   `json:"match_end"`
	Stakes         int64                       `json:"stakes"`
	TotalBetters   int                         `json:"total_betters"`
	Status         tmp.BET_COMPLETETION_STATUS `json:"status"`
	T              tmp.BET_COMPLETETION_TYPE   `json:"t"`
	Bets           int                         `json:"bets"`
	ProgramId      int64                       `json:"program_id"`
	ObtainedStakes int64                       `json:"obtained_stakes"`
}

type OutBetCompletetionPagedList struct {
	CurrentPage int                   `json:"current_page"`
	List        []*OutBetCompletetion `json:"list"`
}

type OutBetMatchList struct {
	MostHot   *OutBetCompletetion `json:"most_hot"`
	BetMatchs []*OutBetMatch      `json:"bet_matchs"`
}

type OutBetItemObj struct {
	Id             int64                `json:"id"`
	BetId          string               `json:"bet_id"`
	CompId         int64                `json:"comp_id"`
	ItemId         int64                `json:"item_id"`
	Odds           float64              `json:"odds"`
	Position       tmp.BET_OBJ_POSITION `json:"position"`
	Title          string               `json:"title"`
	Betters        int                  `json:"betters"`
	Stakes         int64                `json:"stakes"`
	Enabled        bool                 `json:"enabled"`
	IsWin          bool                 `json:"is_win"`
	BettedStakes   int64                `json:"betted_stakes"`
	ObtainedStakes int64                `json:"obtained_stakes"`
}

type OutBetModel struct {
	BetId        string           `json:"bet_id"`
	GameId       string           `json:"game_id"`
	Game         *OutBetGame      `json:"game"`
	CompId       int64            `json:"comp_id"`
	MatchId      int64            `json:"match_id"`
	Match        *OutBetMatch     `json:"match"`
	Title        string           `json:"title"`
	BetStartTime time.Time        `json:"start_time"`
	BetEndTime   time.Time        `json:"end_time"`
	BetType      tmp.BET_TYPE     `json:"bet_type"`
	State        tmp.BET_STATE    `json:"state"`
	BetObjs      []*OutBetItemObj `json:"bet_objs"`
}

type OutBetWinLoseModel struct {
	BetId        string           `json:"bet_id"`
	GameId       string           `json:"game_id"`
	Game         *OutBetGame      `json:"game"`
	CompId       int64            `json:"comp_id"`
	MatchId      int64            `json:"match_id"`
	Match        *OutBetMatch     `json:"match"`
	Title        string           `json:"title"`
	BetStartTime time.Time        `json:"start_time"`
	BetEndTime   time.Time        `json:"end_time"`
	BetType      tmp.BET_TYPE     `json:"bet_type"`
	State        tmp.BET_STATE    `json:"state"`
	BetObjs      []*OutBetItemObj `json:"bet_objs"`
}

type OutBetSingleModel struct {
	BetId        string           `json:"bet_id"`
	GameId       string           `json:"game_id"`
	Game         *OutBetGame      `json:"game"`
	CompId       int64            `json:"comp_id"`
	MatchId      int64            `json:"match_id"`
	Match        *OutBetMatch     `json:"match"`
	Title        string           `json:"title"`
	BetStartTime time.Time        `json:"start_time"`
	BetEndTime   time.Time        `json:"end_time"`
	BetType      tmp.BET_TYPE     `json:"bet_type"`
	State        tmp.BET_STATE    `json:"state"`
	BetObjs      []*OutBetItemObj `json:"bet_objs"`
	SubTitle     string           `json:"subtitle"`
	Description  string           `json:"description"`
}

type OutBetLetModel struct {
	BetId        string           `json:"bet_id"`
	GameId       string           `json:"game_id"`
	Game         *OutBetGame      `json:"game"`
	CompId       int64            `json:"comp_id"`
	MatchId      int64            `json:"match_id"`
	Match        *OutBetMatch     `json:"match"`
	Title        string           `json:"title"`
	BetStartTime time.Time        `json:"start_time"`
	BetEndTime   time.Time        `json:"end_time"`
	BetType      tmp.BET_TYPE     `json:"bet_type"`
	State        tmp.BET_STATE    `json:"state"`
	BetObjs      []*OutBetItemObj `json:"bet_objs"`
	Lets         float64          `json:"lets"`
}

type OutBetMultiModel struct {
	BetId        string           `json:"bet_id"`
	GameId       string           `json:"game_id"`
	Game         *OutBetGame      `json:"game"`
	CompId       int64            `json:"comp_id"`
	MatchId      int64            `json:"match_id"`
	Match        *OutBetMatch     `json:"match"`
	Title        string           `json:"title"`
	BetStartTime time.Time        `json:"start_time"`
	BetEndTime   time.Time        `json:"end_time"`
	BetType      tmp.BET_TYPE     `json:"bet_type"`
	State        tmp.BET_STATE    `json:"state"`
	BetObjs      []*OutBetItemObj `json:"bet_objs"`
}

type OutBetUserStats struct {
	Uid    int64 `json:"uid"`
	Stakes int64 `json:"stakes"`
	Wins   int   `json:"wins"`
	Loses  int   `json:"loses"`
	Rank   int   `json:"rank"`
}

type OutBetExplain struct {
	Content string `json:"content"`
}
