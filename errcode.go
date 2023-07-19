package goerr

/*
ErrCode
*

	项目组代号:10
	服务代号:01
	模块代号:0~99
	错误码：0~99

	| 错误标识	              | 错误码	 | HTTP状态码  | 描述							  |
	| ----------------------- | -------- | ---------- | ----------------------------- |
	| ErrNo                   | 10010000 | 200        | OK                            |
	| ErrInternalServer       | 10010001 | 500        | Internal server error         |
	| ErrParams               | 10010002 | 400        | Illegal params                |
	| ErrValidateParams       | 10010012 | 400        | validator.ValidationErrors    |
	| ErrAuthentication       | 10010003 | 401        | Authentication failed         |
	| ErrNotFound             | 10010004 | 404        | Page not found                |
	| ErrAuthenticationHeader | 10010005 | 401        | Authentication header Illegal |
	| ErrAppKey               | 10010006 | 401        | Invalid app key               |
	| ErrSecretKey            | 10010007 | 401        | Invalid secret key            |
	| ErrPermission           | 10010008 | 403        | Permission denied             |
	| ErrTooManyRequests      | 10010013 | 429        | Too Many Requests             |
	| ErrInvalidJson          | 10010009 | 500        | Invalid Json                  |
	| ErrTimeout              | 10010010 | 504        | Server response timeout       |
	| ErrElasticsearchServer  | 10010101 | 500        | Elasticsearch server error    |
	| ErrElasticsearchDSL     | 10010102 | 500        | Elasticsearch  DSL error      |
	| ErrMysqlServer          | 10010201 | 500        | Mysql server error            |
	| ErrMysqlSQL             | 10010202 | 500        | Illegal SQL                   |
	| ErrMongoServer          | 10010301 | 500        | MongoDB server error          |
	| ErrMongoDSL             | 10010302 | 500        | MongoDB DSL error             |
	| ErrRedisServer          | 10010401 | 500        | Redis server error            |
	| ErrKafkaServer          | 10010501 | 500        | Kafka server error            |
*/
type ErrCode struct {
	Code     int
	HTTPCode int
	Desc     string
}

var (
	ErrNo = ErrCode{
		Code:     0,
		HTTPCode: 200,
		Desc:     "OK",
	}
	ErrRequestFail = ErrCode{
		Code:     101,
		HTTPCode: 200,
		Desc:     "请求失败",
	}
	ErrInternalServer = ErrCode{
		Code:     10010001,
		HTTPCode: 500,
		Desc:     "Internal server error",
	}
	ErrParams = ErrCode{
		Code:     10010002,
		HTTPCode: 422,
		Desc:     "Illegal params",
	}
	ErrValidateParams = ErrCode{
		Code:     10010012,
		HTTPCode: 422,
		Desc:     "validator.ValidationErrors",
	}
	ErrAuthentication = ErrCode{
		Code:     10010003,
		HTTPCode: 401,
		Desc:     "Authentication failed",
	}
	ErrNotFound = ErrCode{
		Code:     10010004,
		HTTPCode: 404,
		Desc:     "Page not found",
	}
	ErrAuthenticationHeader = ErrCode{
		Code:     10010005,
		HTTPCode: 403,
		Desc:     "Authentication header Illegal",
	}
	ErrAppKey = ErrCode{
		Code:     10010006,
		HTTPCode: 403,
		Desc:     "Invalid app key",
	}
	ErrSign = ErrCode{
		Code:     10010007,
		HTTPCode: 401,
		Desc:     "Invalid secret key",
	}
	ErrPermission = ErrCode{
		Code:     10010008,
		HTTPCode: 403,
		Desc:     "Permission denied",
	}
	ErrTooManyRequests = ErrCode{
		Code:     10010013,
		HTTPCode: 429,
		Desc:     "Too Many Requests",
	}
	ErrInvalidJson = ErrCode{
		Code:     10010009,
		HTTPCode: 500,
		Desc:     "Invalid Json",
	}
	ErrTimeout = ErrCode{
		Code:     10010010,
		HTTPCode: 408,
		Desc:     "Server response timeout",
	}
	ErrAuthExpired = ErrCode{
		Code:     10010011,
		HTTPCode: 504,
		Desc:     "Authentication expired",
	}
	ErrElasticsearchServer = ErrCode{
		Code:     10010101,
		HTTPCode: 500,
		Desc:     "Elasticsearch server error",
	}
	ErrElasticsearchDSL = ErrCode{
		Code:     10010102,
		HTTPCode: 500,
		Desc:     "Elasticsearch  DSL error",
	}
	ErrMysqlServer = ErrCode{
		Code:     10010201,
		HTTPCode: 500,
		Desc:     "Mysql server error",
	}
	ErrMysqlSQL = ErrCode{
		Code:     10010202,
		HTTPCode: 500,
		Desc:     "Illegal SQL",
	}
	ErrMongoServer = ErrCode{
		Code:     10010301,
		HTTPCode: 500,
		Desc:     "MongoDB server error",
	}
	ErrMongoDSL = ErrCode{
		Code:     10010302,
		HTTPCode: 500,
		Desc:     "MongoDB DSL error",
	}
	ErrRedisServer = ErrCode{
		Code:     10010401,
		HTTPCode: 500,
		Desc:     "Redis server error",
	}
	ErrKafkaServer = ErrCode{
		Code:     10010501,
		HTTPCode: 500,
		Desc:     "Kafka server error",
	}
	ErrSearch = ErrCode{
		Code:     10010601,
		HTTPCode: 500,
		Desc:     "search error",
	}
)
