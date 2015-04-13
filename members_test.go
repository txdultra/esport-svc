package main

/*
import (
	"./libs"
	"./models"
	"./utils"
	"fmt"
	"regexp"
	"testing"
	"time"
)

func TestCreateMember(t *testing.T) {
	member := new(models.Member)
	member.UserName = utils.MakeRndStrs(6)

	fmt.Println(member.UserName)
	if matched, err := regexp.MatchString(libs.USR_NAME_REGEX, member.UserName); err != nil || !matched {
		fmt.Println("表达式错误")
	}
	member.NickName = utils.MakeRndStrs(4)
	member.Email = utils.MakeRndStrs(5) + "@neotv.cn"
	member.Password = "123123"
	member.MobileIdentifier = utils.MakeRndStrs(32)
	member.MemberIdentifier = utils.MakeRndStrs(32)
	member.CreateTime = time.Now().Unix()
	member.CreateIP = utils.IpToInt("192.168.1.2")
	provider := libs.NewMemberProvider()
	_, err := provider.Create(member, 0, 0)
	if err != nil {
		t.Fatal(err.ErrorDescription)
	}
}

func BenchmarkCreateMember(b *testing.B) {
	for i := 0; i < b.N; i++ {
		member := new(models.Member)
		member.UserName = utils.MakeRndStrs(6)
		member.NickName = utils.MakeRndStrs(4)
		member.Email = utils.MakeRndStrs(5) + "@neotv.cn"
		member.Password = "123123"
		member.MobileIdentifier = utils.MakeRndStrs(32)
		member.MemberIdentifier = utils.MakeRndStrs(32)
		member.CreateTime = time.Now().Unix()
		member.CreateIP = utils.IpToInt("192.168.1.2")
		provider := libs.NewMemberProvider()
		_, err := provider.Create(member, 0, 0)
		if err != nil {
			fmt.Println(member.UserName)
			b.Fatal(err.ErrorDescription)
		}
	}
}
*/
