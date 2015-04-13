package outobjs

import (
	"libs"
	"libs/passport"
)

var file libs.IFileStorage = libs.NewFileStorage()
var friendship *passport.FriendShips = &passport.FriendShips{}
