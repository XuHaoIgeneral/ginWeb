package models

type User struct {
	Id            int64  `xorm:"pk autoincr"`
	WXUnionId     string `xorm:"unique" json:"union_id"`  //微信id
	NickName      string `xorm:"unique" json:"nick_name"` //昵称
	Avatar        string `json:"avatar"`                  //头像
	Adoption      int    `json:"adoption"`                //总采纳
	Totallncome   int    `json:"totallncome"`             //总收益
	IncomeAccount string `json:"income_account"`          //收益账户
	CashAccount   string `json:"cash_account"`            //可提现现账户
}

type Wx struct {
	Id         int64             `xorm:"pk autoincr" json:"id"`
	Openid     string            `xorm:"unique" json:"openid"`
	Nickname   string            `json:"nickname"`
	Sex        int               `json:"sex"`
	Province   string            `json:"province"`
	City       string            `json:"city"`
	Country    string            `json:"country"`
	Headimgurl string            `json:"headimgurl"`
	Privilege  map[string]string `json:"privilege"`
	Unionid    string            `json:"unionid"`
}

type TransactionDetail struct {
	Id         int64   `xorm:"pk autoincr"`
	UserId     int64   `xorm:"notnull"` //用户id
	Type       int     `xorm:"notnull"` //类型
	Account    string  `xorm:"notnull"` //账户
	QuestionId int64   `xorm:"notnull"` //题目id
	Money      float64 `xorm:"decimal"` //流水金额
	Commission float64 `xorm:"decimal"` //手续费
	Status     int     `xorm:"notnull"` //状态
}

type Question struct {
	Id                int64   `xorm:"pk autoincr"`
	Sponsor           int64   `xorm:"notnull"` //发起人id
	AnswerId          int64   `xorm:"notnull"` //fk 答案id
	Type              int     `xorm:"notnull"` //类型
	Title             string  `xorm:"notnull"` //问题标题
	Content           string  `xorm:"text"`    //问题主体
	Price             float64 `xorm:"decimal notnull"`
	CrowdFundingPrice float64 `xorm:"decimal notnull"`
	Tag               int     `xorm:"notnull"`
	RateUp            int
	RateDown          int
	Income            float64 `xorm:"decimal"`
	CreateTime        string  `xorm:"timestamp notnull"`
	UpdateTime        string  `xorm:"timestamp notnull"`
}
