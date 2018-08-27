package outobjs

import "libs/matchrace"

type OutMatchPlayer struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Img    int64  `json:"img"`
	ImgUrl string `json:"img_url"`
}

type OutMatchPlayerPagedList struct {
	Total       int               `json:"total"`
	TotalPage   int               `json:"pages"`
	CurrentPage int               `json:"current_page"`
	Size        int               `json:"size"`
	Lists       []*OutMatchPlayer `json:"lists"`
}

type OutMatchMode struct {
	Id        int64               `json:"id"`
	MatchId   int64               `json:"match_id"`
	ModeType  matchrace.MODE_TYPE `json:"mode_type"`
	Title     string              `json:"title"`
	IsWebView bool                `json:"is_webview"`
	WebUrl    string              `json:"web_url"`
	Groups    []*OutMatchGroup    `json:"groups"`
	Recents   []*OutMatchRecent   `json:"recents"`
	Elimins   []*OutMatchEliminMs `json:"elimins"`
}

type OutMatchModeForAdmin struct {
	Id           int64               `json:"id"`
	MatchId      int64               `json:"match_id"`
	ModeType     matchrace.MODE_TYPE `json:"mode_type"`
	Title        string              `json:"title"`
	DisplayOrder int                 `json:"displayorder"`
	IsView       bool                `json:"is_view"`
}

type OutMatchGroup struct {
	Id      int64                  `json:"id"`
	Title   string                 `json:"title"`
	Players []*OutMatchGroupPlayer `json:"players"`
}

type OutMatchGroupPlayer struct {
	Id           int64           `json:"id"`
	GroupId      int64           `json:"group_id"`
	PlayerId     int64           `json:"player_id"`
	Player       *OutMatchPlayer `json:"player"`
	Wins         int16           `json:"wins"`
	Pings        int16           `json:"pings"`
	Loses        int16           `json:"loses"`
	Points       int             `json:"points"`
	Outlet       bool            `json:"outlet"`
	DisplayOrder int             `json:"displayorder"`
	Vss          []*OutMatchVs   `json:"vss"`
}

type OutMatchRecent struct {
	Id       int64           `json:"id"`
	PlayerId int64           `json:"player_id"`
	Player   *OutMatchPlayer `json:"player"`
	M1       matchrace.WLP   `json:"m1"`
	M2       matchrace.WLP   `json:"m2"`
	M3       matchrace.WLP   `json:"m3"`
	M4       matchrace.WLP   `json:"m4"`
	M5       matchrace.WLP   `json:"m5"`
}

type OutMatchRecentForAdmin struct {
	Id           int64           `json:"id"`
	ModeId       int64           `json:"mode_id"`
	PlayerId     int64           `json:"player_id"`
	Player       *OutMatchPlayer `json:"player"`
	M1           matchrace.WLP   `json:"m1"`
	M2           matchrace.WLP   `json:"m2"`
	M3           matchrace.WLP   `json:"m3"`
	M4           matchrace.WLP   `json:"m4"`
	M5           matchrace.WLP   `json:"m5"`
	DisplayOrder int             `json:"displayorder"`
	Disabled     bool            `json:"disabled"`
}

type OutMatchEliminMs struct {
	Id           int64                   `json:"id"`
	Title        string                  `json:"title"`
	Icon         int64                   `json:"icon"`
	IconUrl      string                  `json:"icon_url"`
	T            matchrace.ELIMIN_MSTYPE `json:"t"`
	DisplayOrder int                     `json:"displayorder"`
	Evs          []*OutMatchEliminVs     `json:"evs"`
}

type OutMatchEliminVs struct {
	Id   int64       `json:"id"`
	VsId int64       `json:"vsid"`
	MsId int64       `json:"msid"`
	Vs   *OutMatchVs `json:"vs"`
}

type OutMatchVs struct {
	Id      int64  `json:"id"`
	A       int64  `json:"a"`
	AName   string `json:"a_name"`
	AImg    int64  `json:"a_img"`
	AImgUrl string `json:"a_img_url"`
	AScore  int16  `json:"a_score"`
	AOutlet bool   `json:"a_outlet"`
	B       int64  `json:"b"`
	BName   string `json:"b_name"`
	BImg    int64  `json:"b_img"`
	BImgUrl string `json:"b_img_url"`
	BScore  int16  `json:"b_score"`
	BOutlet bool   `json:"b_outlet"`
	MatchId int64  `json:"match_id"`
	ModeId  int64  `json:"mode_id"`
	Refid   int64  `json:"ref_id"`
}

func GetOutMatchPlayer(player *matchrace.MatchPlayer) *OutMatchPlayer {
	if player == nil {
		return nil
	}
	return &OutMatchPlayer{
		Id:     player.Id,
		Name:   player.Name,
		Img:    player.Img,
		ImgUrl: file.GetFileUrl(player.Img),
	}
}

func GetOutMatchPlayerById(playerId int64) *OutMatchPlayer {
	mrs := &matchrace.MatchRaceService{}
	player := mrs.GetPlayer(playerId)
	return GetOutMatchPlayer(player)
}

func GetOutMatchVs(vs *matchrace.MatchVs) *OutMatchVs {
	if vs == nil {
		return nil
	}
	return &OutMatchVs{
		Id:      vs.Id,
		A:       vs.A,
		AName:   vs.AName,
		AImg:    vs.AImg,
		AImgUrl: file.GetFileUrl(vs.AImg),
		AScore:  vs.AScore,
		B:       vs.B,
		BName:   vs.BName,
		BImg:    vs.BImg,
		BImgUrl: file.GetFileUrl(vs.BImg),
		BScore:  vs.BScore,
		MatchId: vs.MatchId,
		ModeId:  vs.ModeId,
		Refid:   vs.RefId,
	}
}

func GetOutMatchEliminMs(mem *matchrace.MatchEliminMs) *OutMatchEliminMs {
	if mem == nil {
		return nil
	}
	mrs := &matchrace.MatchRaceService{}
	evs := mrs.GetEliminVss(mem.Id)
	out_vss := []*OutMatchEliminVs{}
	for _, ev := range evs {
		_vs := mrs.GetMatchVs(ev.VsId)
		if _vs != nil {
			out_vss = append(out_vss, GetOutMatchEliminVs(ev, _vs))
		}
	}
	return &OutMatchEliminMs{
		Id:           mem.Id,
		Title:        mem.Title,
		Icon:         mem.Icon,
		IconUrl:      file.GetFileUrl(mem.Icon),
		T:            mem.T,
		DisplayOrder: mem.DisplayOrder,
		Evs:          out_vss,
	}
}

func GetOutMatchEliminVs(mevs *matchrace.MatchEliminVs, vs *matchrace.MatchVs) *OutMatchEliminVs {
	if mevs == nil {
		return nil
	}
	outvs := GetOutMatchVs(vs)
	if outvs.A == mevs.OutletId {
		outvs.AOutlet = true
	}
	if outvs.B == mevs.OutletId {
		outvs.BOutlet = true
	}
	return &OutMatchEliminVs{
		Id:   mevs.Id,
		VsId: mevs.VsId,
		MsId: mevs.MsId,
		Vs:   outvs,
	}
}

func GetOutMatchRecent(recent *matchrace.MatchRecent) *OutMatchRecent {
	if recent == nil {
		return nil
	}
	return &OutMatchRecent{
		Id:       recent.Id,
		PlayerId: recent.Player,
		Player:   GetOutMatchPlayerById(recent.Player),
		M1:       recent.M1,
		M2:       recent.M2,
		M3:       recent.M3,
		M4:       recent.M4,
		M5:       recent.M5,
	}
}

func GetOutMatchRecentForAdmin(recent *matchrace.MatchRecent) *OutMatchRecentForAdmin {
	if recent == nil {
		return nil
	}
	return &OutMatchRecentForAdmin{
		Id:           recent.Id,
		ModeId:       recent.ModeId,
		PlayerId:     recent.Player,
		Player:       GetOutMatchPlayerById(recent.Player),
		M1:           recent.M1,
		M2:           recent.M2,
		M3:           recent.M3,
		M4:           recent.M4,
		M5:           recent.M5,
		DisplayOrder: recent.DisplayOrder,
		Disabled:     recent.Disabled,
	}
}

func GetOutMatchGroupPlayer(mgp *matchrace.MatchGroupPlayer, vss []*matchrace.MatchVs) *OutMatchGroupPlayer {
	if mgp == nil {
		return nil
	}
	out := &OutMatchGroupPlayer{
		Id:           mgp.Id,
		GroupId:      mgp.GroupId,
		PlayerId:     mgp.Player,
		Player:       GetOutMatchPlayerById(mgp.Player),
		Wins:         mgp.Wins,
		Pings:        mgp.Pings,
		Loses:        mgp.Loses,
		Points:       mgp.Points,
		Outlet:       mgp.Outlet,
		DisplayOrder: mgp.DisplayOrder,
		Vss:          []*OutMatchVs{},
	}
	for _, vs := range vss {
		out_vs := GetOutMatchVs(vs)
		if out_vs != nil {
			out.Vss = append(out.Vss, out_vs)
		}
	}
	return out
}

func GetOutMatchGroup(matchId int64, group *matchrace.MatchGroup, players []*matchrace.MatchGroupPlayer) *OutMatchGroup {
	if group == nil {
		return nil
	}
	out := &OutMatchGroup{
		Id:      group.Id,
		Title:   group.Title,
		Players: []*OutMatchGroupPlayer{},
	}
	mrs := &matchrace.MatchRaceService{}
	for _, player := range players {
		vss := mrs.GetMatchVss(player.Id, matchId)
		out_player := GetOutMatchGroupPlayer(player, vss)
		if out_player != nil {
			out.Players = append(out.Players, out_player)
		}
	}
	return out
}

func GetOutMatchModeForAdmin(race *matchrace.RaceMode) *OutMatchModeForAdmin {
	return &OutMatchModeForAdmin{
		Id:           race.Id,
		MatchId:      race.MatchId,
		ModeType:     race.ModeType,
		Title:        race.Title,
		DisplayOrder: race.DisplayOrder,
		IsView:       race.IsView,
	}
}

func GetOutMatchMode(race *matchrace.RaceMode) *OutMatchMode {
	out := &OutMatchMode{
		Id:        race.Id,
		MatchId:   race.MatchId,
		ModeType:  race.ModeType,
		Title:     race.Title,
		IsWebView: false,
		WebUrl:    "",
		Groups:    []*OutMatchGroup{},
		Recents:   []*OutMatchRecent{},
		Elimins:   []*OutMatchEliminMs{},
	}
	mrs := &matchrace.MatchRaceService{}
	switch race.ModeType {
	case matchrace.MODE_TYPE_GROUP:
		groups := mrs.GetMatchGroups(race.Id)
		for _, group := range groups {
			players := mrs.GetMatchGroupPlayers(group.Id)
			out_g := GetOutMatchGroup(race.MatchId, group, players)
			if out_g != nil {
				out.Groups = append(out.Groups, out_g)
			}
		}
		break
	case matchrace.MODE_TYPE_RECENT:
		recents := mrs.GetMatchRecents(race.Id)
		for _, recent := range recents {
			out_r := GetOutMatchRecent(recent)
			if out_r != nil {
				out.Recents = append(out.Recents, out_r)
			}
		}
		break
	case matchrace.MODE_TYPE_ELIMIN:
		elimins := mrs.GetEliminMss(race.Id)
		for _, elimin := range elimins {
			out_eli := GetOutMatchEliminMs(elimin)
			if out_eli != nil {
				out.Elimins = append(out.Elimins, out_eli)
			}
		}
		break
	}
	return out
}
