package define

var (
	//// 密钥
	//Jwtkey = "d_gateway"
	//// token有效期，7天
	//TokenExpire = time.Now().Add(time.Second * 3600 * 24 * 7).Unix()
	//// 刷新token有效期14天
	//RefreshTokenExpire = time.Now().Add(time.Second * 3600 * 24 * 14).Unix()

	// DefaultSize 后台默认分页没有显示条数
	DefaultSize = 10
	DefaultPage = 1

	LoadTypeHTTP = 0
	LoadTypeTCP  = 1
	LoadTypeGRPC = 2

	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1

	RedisFlowDayKey  = "flow_day_count"
	RedisFlowHourKey = "flow_hour_count"

	FlowTotal         = "flow_total"
	FlowServicePrefix = "flow_service_"
	FlowAppPrefix     = "flow_app_"
)

type UserClaim struct {
	Id      uint
	Name    string
	IsAdmin bool // 是否超管
	RoleId  uint // 所属角色
}

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeTCP:  "TCP",
		LoadTypeGRPC: "GRPC",
	}
)

// PageHelper 分页参数处理
func DefaultPageNum(inputPageNo, inputPageSize int) (outpageNo, outPageSize int) {
	pageNo := inputPageNo
	pageSize := inputPageSize

	if pageNo <= 0 {
		pageNo = DefaultPage
	}

	if pageSize <= 0 {
		pageSize = DefaultSize
	}

	return pageNo, pageSize
}
