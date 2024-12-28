package errs

const (
	// CommonInvalidInput  任何模块均可使用
	CommonInvalidInput        = 400001
	CommonInternalServerError = 500001
)

// 用户模块
const (
	// UserInvalidInput 用户模块输入错误，这是一个含糊的错误
	UserInvalidInput        = 401001
	UserInternalServerError = 501001
	// UserInvalidOrPassword 用户不存在或密码错误,防止别人攻击
	UserInvalidOrPassword = 401002
)

const (
	ArticleInvalidInput        = 402001
	ArticleInternalServerError = 502002
)

var (
	UserInvalidInputV1 = Code{
		Number: 401001,
		Msg:    "用户输入有误",
	}
)

type Code struct {
	Number int
	Msg    string
}
