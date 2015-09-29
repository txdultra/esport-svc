package passport

import (
	"dbs"
	"fmt"
	"libs/passport/service"
	"libs/stat"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/thrift"
)

var default_username_minlen, default_username_maxlen, authorization_access_token_expries_seconds, Authorization_access_token_expries_refresh int
var MemberPasswordMinLen, MemberPasswordMaxLen int
var friend_limit_counts int
var openid_qq_key, openid_qq_consumer, use_ssdb_passport_db string
var credit_service_host, jing_service_host string

//搜索参数
var search_member_server string
var search_member_port, search_member_timeout int

const (
	MOD_NAME     = "passport"
	relation_db  = "user_relation"
	relation_tbl = "user_relation"
)

func init() {
	orm.RegisterModel(new(AccessToken), new(Member), new(MemberState), new(OpenIDOAuth),
		new(MemberRole), new(Role), new(Authority),
		new(MemberNickName), new(MemberProfile),
		new(ManageMemberMap), new(MemberConfig),
	)

	default_username_minlen, _ = beego.AppConfig.Int("member.username.minlen")
	default_username_maxlen, _ = beego.AppConfig.Int("member.username.maxlen")
	MemberPasswordMinLen, _ = beego.AppConfig.Int("member.password.minlen")
	MemberPasswordMaxLen, _ = beego.AppConfig.Int("member.password.maxlen")
	if MemberPasswordMinLen <= 0 {
		MemberPasswordMinLen = 6
	}
	if MemberPasswordMaxLen <= 0 {
		MemberPasswordMaxLen = 16
	}
	_authorization_access_token_expries_seconds, err := beego.AppConfig.Int("authorization_access_token_expries_seconds")
	if err != nil {
		authorization_access_token_expries_seconds = 30 * 3600 * 24
	} else {
		authorization_access_token_expries_seconds = _authorization_access_token_expries_seconds
	}
	_authorization_access_token_expries_refresh, err := beego.AppConfig.Int("authorization_access_token_expries_refresh")
	if err != nil {
		Authorization_access_token_expries_refresh = 10 * 3600 * 24
	} else {
		Authorization_access_token_expries_refresh = _authorization_access_token_expries_refresh
	}
	_friend_limit_counts, _ := beego.AppConfig.Int("friend.limit.counts")
	if _friend_limit_counts <= 0 {
		friend_limit_counts = 200
	} else {
		friend_limit_counts = _friend_limit_counts
	}

	openid_qq_key = beego.AppConfig.String("openid.qq.key")
	openid_qq_consumer = beego.AppConfig.String("openid.qq.consumer")
	use_ssdb_passport_db = beego.AppConfig.String("ssdb.passport.db")

	//search
	search_member_server = beego.AppConfig.String("search.member.server")
	search_member_port, _ = beego.AppConfig.Int("search.member.port")
	search_member_timeout, _ = beego.AppConfig.Int("search.member.timeout")

	//db
	register_relation_db()

	//注册计数器
	stat.RegisterCounter(MOD_NAME, &MemberProvider{})

	//积分系统地址
	credit_service_host = beego.AppConfig.String("credit.host")
	jing_service_host = beego.AppConfig.String("jings.host")

	//启动账号服务
	passport_service_run := beego.AppConfig.DefaultBool("passport.service.run", false)
	passport_service_port := beego.AppConfig.DefaultInt("passport.service.port", 20002)
	if passport_service_run {
		go runPassportServer(passport_service_port)
	}
}

func runPassportServer(port int) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	host := fmt.Sprintf("0.0.0.0:%d", port)
	serverTransport, err := thrift.NewTServerSocket(host)
	if err != nil {
		panic(err)
	}

	handler := &PassportServiceImpl{}
	processor := service.NewPassportServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("passport service server in " + host)
	server.Serve()
}

func register_relation_db() {
	relation_db_user := beego.AppConfig.String("db.relation.user")
	relation_db_pwd := beego.AppConfig.String("db.relation.pwd")
	relation_db_host := beego.AppConfig.String("db.relation.host")
	relation_db_port, _ := beego.AppConfig.Int("db.relation.port")
	relation_db_name := beego.AppConfig.String("db.relation.name")
	relation_db_charset := beego.AppConfig.String("db.relation.charset")
	relation_db_protocol := beego.AppConfig.String("db.relation.protocol")
	relation_db_time_local := beego.AppConfig.String("db.relation.time_local")
	relation_db_maxconns, _ := beego.AppConfig.Int("db.relation.maxconns")
	relation_db_maxidels, _ := beego.AppConfig.Int("db.relation.maxidles")
	if relation_db_maxconns <= 0 {
		relation_db_maxconns = 1500
	}
	if relation_db_maxidels <= 0 {
		relation_db_maxidels = 100
	}

	relation_db_addr := relation_db_host + ":" + strconv.Itoa(relation_db_port)
	connection_url := relation_db_user + ":" + relation_db_pwd + "@" + relation_db_protocol + "(" + relation_db_addr + ")/" + relation_db_name + "?charset=" + relation_db_charset + "&loc=" + relation_db_time_local
	dbs.LoadDb(relation_db, "mysql", connection_url, relation_db_maxidels, relation_db_maxconns)
}
