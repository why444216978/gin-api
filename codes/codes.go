package codes

//1XXX，参数相关
const ERRNO_MISS_ORIGIN_CUSTOMER_ID = 1000
const ERRNO_MISS_TYPE = 1001
const ERRNO_MISS_PRICE_ID = 1002
const ERRNO_MISS_PRICE_LIST = 1003
const ERRNO_MISS_DESC_LIST = 1004
const ERRNO_WRONG_PRICE_SPEC = 1005
const ERRNO_WRONG_PRICE_NUM = 1006
const ERRNO_WRONG_DESC = 1007
const ERRNO_WRONG_OVERTIME = 1008
const ERRNO_WRONG_BREED_ID = 1009
const ERRNO_MISS_BREED_ID = 1010
const ERRNO_WRONG_TYPE = 1011
const ERRNO_MISS_PRODUCT_ID = 1012
const ERRNO_PARAMS_EMPTY = 1013

//2XXX，业务验证相关
const ERRNO_REPEAT_ADD_BREED = 2000
const ERRNO_CUSTOMER_NOT_HAS_BREED = 2001
const ERRNO_PRICE_ID_NOT_LAST = 2002
const ERRNO_MAX_PRODUCT = 2003

//3XXX，权限相关
const NO_AUTHORIZE_SPY = 3001

//5XXX，服务器错误相关
const SERVER_ERROR = 5000
const ERRNO_DATA_ERR = 5001

var ErrorMsg = map[int]string{
	//1XXX
	ERRNO_MISS_ORIGIN_CUSTOMER_ID: "缺少情报员ID",
	ERRNO_MISS_TYPE:               "缺少type参数",
	ERRNO_MISS_PRICE_ID:           "缺少报价ID",
	ERRNO_MISS_PRICE_LIST:         "缺少价格列表",
	ERRNO_MISS_DESC_LIST:          "缺少描述列表",
	ERRNO_WRONG_PRICE_SPEC:        "缺少规格",
	ERRNO_WRONG_PRICE_NUM:         "报价不合法",
	ERRNO_WRONG_DESC:              "描述不合法",
	ERRNO_WRONG_OVERTIME:          "超时不允许修改",
	ERRNO_WRONG_BREED_ID:          "ID错误，该情报员不存在该品类",
	ERRNO_MISS_BREED_ID:           "缺少breed_id",
	ERRNO_WRONG_TYPE:              "type不合法",
	ERRNO_MISS_PRODUCT_ID:         "缺少product_id",
	ERRNO_PARAMS_EMPTY:            "缺少参数",

	//2XXX
	ERRNO_REPEAT_ADD_BREED:       "重复提交",
	ERRNO_CUSTOMER_NOT_HAS_BREED: "该情报员未添加该品类",
	ERRNO_PRICE_ID_NOT_LAST:      "所传price_id不是最新报价",
	ERRNO_MAX_PRODUCT:            "最多只能添加5个品类",

	//3XXX
	NO_AUTHORIZE_SPY: "不是情报员",

	//5XXX
	SERVER_ERROR:   "服务器错误",
	ERRNO_DATA_ERR: "数据错误",
}

var ErrorUserMsg = map[int]string{

	//1XXX
	ERRNO_MISS_ORIGIN_CUSTOMER_ID: "请求参数错误",
	ERRNO_MISS_TYPE:               "请求参数错误",
	ERRNO_MISS_PRICE_ID:           "请求参数错误",
	ERRNO_MISS_PRICE_LIST:         "请求参数错误",
	ERRNO_MISS_DESC_LIST:          "请求参数错误",
	ERRNO_WRONG_PRICE_SPEC:        "请求参数错误",
	ERRNO_WRONG_PRICE_NUM:         "请求参数错误",
	ERRNO_WRONG_DESC:              "请求参数错误",
	ERRNO_WRONG_OVERTIME:          "请求参数错误",
	ERRNO_WRONG_BREED_ID:          "请求参数错误",
	ERRNO_MISS_BREED_ID:           "请求参数错误",
	ERRNO_WRONG_TYPE:              "请求参数错误",
	ERRNO_MISS_PRODUCT_ID:         "请求参数错误",
	ERRNO_PARAMS_EMPTY:            "请求参数错误",

	//2XXX
	ERRNO_REPEAT_ADD_BREED:       "重复提交",
	ERRNO_CUSTOMER_NOT_HAS_BREED: "请求参数错误",
	ERRNO_PRICE_ID_NOT_LAST:      "请求参数错误",
	ERRNO_MAX_PRODUCT:            "您最多只能添加5个品类",

	//3XXX
	NO_AUTHORIZE_SPY: "您不是情报员，请申请成为情报员",

	//5XXX
	SERVER_ERROR:   "服务器暂时有点小问题，稍后再试",
	ERRNO_DATA_ERR: "服务器暂时有点小问题，稍后再试",
}
