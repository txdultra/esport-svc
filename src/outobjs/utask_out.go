package outobjs

import "libs/utask"

type OutUserTaskGroup struct {
	GroupId      int64                 `json:"group_id"`
	GroupType    utask.TASK_GROUP_TYPE `json:"group_type"`
	TaskType     utask.TASK_TYPE       `json:"task_type"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	Icon         int64                 `json:"icon"`
	IconUrl      string                `json:"icon_url"`
	BgImg        int64                 `json:"bgimg"`
	BgImgUrl     string                `json:"bgimg_url"`
	TaskCount    int                   `json:"task_count"`
	DoneCount    int                   `json:"done_count"`
	DisplayOrder int                   `json:"displayorder"`
	Tasks        []*OutUserTask        `json:"tasks"`
	Ex1          string                `json:"extension1"`
	Ex2          string                `json:"extension2"`
	Ex3          string                `json:"extension3"`
	Ex4          string                `json:"extension4"`
	Ex5          string                `json:"extension5"`
}

type OutUserTask struct {
	TaskId       int64                  `json:"task_id"`
	TaskType     utask.TASK_TYPE        `json:"task_type"`
	GroupId      int64                  `json:"group_id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Icon         int64                  `json:"icon"`
	IconUrl      string                 `json:"icon_url"`
	Limit        int                    `json:"limits"`
	Dones        int                    `json:"dones"`
	Period       int                    `json:"period"`
	PeriodType   utask.TASK_PERIOD_TYPE `json:"periodtype"`
	Reward       utask.TASK_REWARD_TYPE `json:"reward"`
	Prize        int64                  `json:"prize"`
	DisplayOrder int                    `json:"displayorder"`
	ResetVar     int                    `json:"resetvar"`
	ResetTime    int                    `json:"resettime"`
	Ex1          string                 `json:"extension1"`
	Ex2          string                 `json:"extension2"`
	Ex3          string                 `json:"extension3"`
	Ex4          string                 `json:"extension4"`
	Ex5          string                 `json:"extension5"`
}
