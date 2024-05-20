package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(20);not null" comment:"用户名"`
	Password string `json:"password" gorm:"type:varchar(100);not null" comment:"密码"`
	Role     int    `json:"role" gorm:"type:int;default:0" comment:"0:普通用户,1:管理员"`
	Gender   int    `json:"gender" gorm:"type:int;default:0" comment:"0:女,1:男"`
}
type AccountType struct { //银行的名称 招商银行一卡通、牡丹行
	gorm.Model
	Name            string  `gorm:"type:varchar(20);not null;unique"` // 账户类型名称，如"招商银行一卡通"
	Description     string  `gorm:"type:varchar(255);"`               // 账户类型描述
	OverdraftPolicy bool    `gorm:"default:false"`                    // 是否允许透支
	InterestRate    float64 `gorm:"type:decimal(5,2);default:0"`      // 对应的利率，适用于贷款或透支
}

// Account 表示用户账户信息 进行存款
type Account struct {
	gorm.Model                 // 添加ID, CreatedAt, UpdatedAt, DeletedAt字段
	UserID         uint        `gorm:"index"`                                  // 用户ID，索引以加速查询
	AccountNumber  string      `gorm:"type:varchar(20);not null;unique;index"` // 账号，设置为唯一和索引
	AccountTypeID  uint        `gorm:"index"`                                  // 外键，指向 AccountType
	AccountType    AccountType `gorm:"foreignKey:AccountTypeID"`
	PasswordHash   string      `gorm:"type:varchar(255);not null"`      // 存储加密后的密码
	Balance        float64     `gorm:"type:decimal(10,2);default:0"`    // 账户余额，默认值为0
	OverdraftLimit float64     `gorm:"type:decimal(10,2);default:1000"` // 透支限额
	CreditRating   int         `gorm:"default:0"`                       // 信用等级
}

// Transaction 表示账户的交易记录
type Transaction struct {
	gorm.Model
	AccountID       uint      `gorm:"index"` // 关联的账户ID
	Account         Account   `gorm:"foreignKey:AccountID"`
	Type            string    `gorm:"type:varchar(20);not null"`   // 交易类型，如"deposit", "withdrawal"
	Amount          float64   `gorm:"type:decimal(10,2);not null"` // 交易金额
	TransactionDate time.Time // 交易日期
	Status          string    `gorm:"type:varchar(20)"` // 交易状态，如"success", "failed"
}

// Loan 表示贷款信息
type Loan struct {
	gorm.Model
	AccountID      uint      `gorm:"index"` // 关联的账户ID
	Account        Account   `gorm:"foreignKey:AccountID"`
	AmountBorrowed float64   `gorm:"type:decimal(10,2);not null"` // 借款金额
	LoanDate       time.Time // 贷款日期
	DueDate        time.Time // 还款到期日期
	InterestRate   float64   `gorm:"type:decimal(5,2);not null"` // 利率
	Status         string    `gorm:"type:varchar(20);not null"`  // 贷款状态，如"active", "closed"
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
