package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(20);not null;unique" comment:"用户名"`
	Password string `json:"password" gorm:"type:varchar(100);not null" comment:"密码"`
	Role     int    `json:"role" gorm:"type:int;default:0" comment:"0:普通用户,1:管理员"`
	Gender   int    `json:"gender" gorm:"type:int;default:0" comment:"0:女,1:男"`
	Phone    string `json:"phone" gorm:"type:varchar(20)" comment:"手机号"`
	Email    string `json:"email" gorm:"type:varchar(50)" comment:"邮箱"`
	Address  string `json:"address" gorm:"type:varchar(255)" comment:"地址"`
	RealName string `json:"real_name" gorm:"type:varchar(20)" comment:"真实姓名"`
	//冻结
	IsFrozen bool `json:"is_frozen" gorm:"default:false" comment:"是否冻结"`
}
type Bank struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20);not null;unique" comment:"银行名称"`
	Description string `gorm:"type:varchar(255);" comment:"银行描述"`
}
type Bank struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20);not null;unique" comment:"银行名称"`
	Description string `gorm:"type:varchar(255);" comment:"银行描述"`
}
type AccountType struct { //银行的名称 招商银行一卡通、牡丹行
	gorm.Model
	Name            string  `gorm:"type:varchar(20);not null;unique"` // 账户类型名称，如"招商银行一卡通"
	Description     string  `gorm:"type:varchar(255);"`               // 账户类型描述
	OverdraftPolicy bool    `gorm:"default:false"`                    // 是否允许透支
	InterestRate    float64 `gorm:"type:decimal(5,2);default:0"`      // 对应的利率，适用于贷款或透支
	BankID          uint    `gorm:"index"`                            // 外键，指向银行
	Bank            Bank    `gorm:"foreignKey:BankID"`
}

// Account 表示用户账户信息 进行存款
type Account struct {
	gorm.Model                          // 添加ID, CreatedAt, UpdatedAt, DeletedAt字段
	UserID                  uint        `gorm:"index"`                                  // 用户ID，索引以加速查询
	User                    User        `gorm:"foreignKey:UserID"`                      // 外键，指向 User
	AccountNumber           string      `gorm:"type:varchar(20);not null;unique;index"` // 账号，设置为唯一和索引
	AccountTypeID           uint        `gorm:"index"`                                  // 外键，指向 AccountType
	AccountType             AccountType `gorm:"foreignKey:AccountTypeID"`
	PasswordHash            string      `gorm:"type:varchar(255);not null"`      // 存储加密后的密码
	Balance                 float64     `gorm:"type:decimal(10,2);default:0"`    // 账户余额，默认值为0
	OverdraftLimit          float64     `gorm:"type:decimal(10,2);default:1000"` // 透支限额
	CreditRating            int         `gorm:"default:0"`                       // 信用等级
	IsOverdraftLimitReached bool        `gorm:"default:false"`                   // 是否达到透支限额
}

// Transaction 表示账户的交易记录
type Transaction struct {
	gorm.Model
	AccountID uint    `gorm:"index"` // 关联的账户ID
	Account   Account `gorm:"foreignKey:AccountID"`
	// 交易类型，1 是转账 2是存款 3是取款 4 是贷款 5是还款 6 是透支 7是透支还款 8是回滚
	TransactionType int       `gorm:"type:int;not null;default:1"`
	Amount          float64   `gorm:"type:decimal(10,2);not null"` // 交易金额
	TransactionDate time.Time // 交易日期
	Status          string    `gorm:"type:varchar(20)"` // 交易状态，如"success", "failed"
	//交易的卡号
	CardNumber string `gorm:"type:varchar(20);not null"`
	//接受方的卡号 银行卡号为00000000
	ToCardNumber string `gorm:"type:varchar(20);not null"`
}

// Loan 表示贷款信息
type Loan struct {
	gorm.Model
	AccountID       uint      `gorm:"index"` // 关联的账户ID
	Account         Account   `gorm:"foreignKey:AccountID"`
	AmountBorrowed  float64   `gorm:"type:decimal(10,2);not null"` // 借款金额
	LoanDate        time.Time // 贷款日期
	DueDate         time.Time // 还款到期日期
	InterestRate    float64   `gorm:"type:decimal(5,2);not null"`   // 利率
	InterestAccrued float64   `gorm:"type:decimal(10,2);default:0"` // 利息
	Status          bool      `gorm:"default:false"`                // 是否已还款
}

// Overdrafts 表存储透支还款记录
type Overdraft struct {
	gorm.Model
	AccountID        uint      `gorm:"index"`
	Account          Account   `gorm:"foreignKey:AccountID"`
	Amount           float64   `gorm:"type:decimal(10,2);not null"`
	RepaymentDueDate time.Time // 还款截止日期
	Repaid           bool      `gorm:"default:false"` // 是否已还款
}
