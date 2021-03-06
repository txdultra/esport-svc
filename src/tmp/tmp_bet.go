package tmp

type BET_COMPLETETION_TYPE int

const (
	BET_COMPLETETION_TYPE_VS     BET_COMPLETETION_TYPE = 1
	BET_COMPLETETION_TYPE_SINGLE BET_COMPLETETION_TYPE = 2
)

type BET_COMPLETETION_STATUS int

const (
	BET_COMPLETETION_STATUS_WAIT     BET_COMPLETETION_STATUS = 1
	BET_COMPLETETION_STATUS_MATCHING BET_COMPLETETION_STATUS = 2
	BET_COMPLETETION_STATUS_MATCHEND BET_COMPLETETION_STATUS = 3
	BET_COMPLETETION_STATUS_FAULT    BET_COMPLETETION_STATUS = 4
	BET_COMPLETETION_STATUS_CLOSED   BET_COMPLETETION_STATUS = 5
)

type BET_OBJ_POSITION int

const (
	BET_OBJ_POSITION_LEFT    BET_OBJ_POSITION = 1
	BET_OBJ_POSITION_CENTER  BET_OBJ_POSITION = 2
	BET_OBJ_POSITION_RIGHT   BET_OBJ_POSITION = 3
	BET_OBJ_POSITION_BYORDER BET_OBJ_POSITION = 4
)

type BET_TYPE int

const (
	BET_TYPE_WIN    BET_TYPE = 1
	BET_TYPE_LET    BET_TYPE = 2
	BET_TYPE_SINGLE BET_TYPE = 3
	BET_TYPE_MULTI  BET_TYPE = 4
)

type BET_STATE int

const (
	BET_STATE_AUDITING BET_TYPE  = 0
	BET_STATE_OPEN     BET_STATE = 1
	BET_STATE_CLOSE    BET_STATE = 2
	BET_STATE_SETTLING BET_STATE = 9
	BET_STATE_SETTLED  BET_STATE = 10
)
