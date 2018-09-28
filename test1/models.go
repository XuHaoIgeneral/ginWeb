// Copyright 2014 Unknown
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	_ "github.com/lib/pq"

	"github.com/go-xorm/xorm"
)

// 银行账户
type Account struct {
	Id      int64
	Name    string `xorm:"unique"`
	Balance float64
	Version int `xorm:"version"` // 乐观锁
}

// ORM 引擎
var x *xorm.Engine

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbName   = "pg_test"

	prompt = `Please enter number of operation:
1. Create new account
2. Show detail of account
3. Deposit
4. Withdraw
5. Make transfer
6. List exist accounts by Id
7. List exist accounts by balance
8. Delete account
9. Exit`
)

func init() {
	// 创建 ORM 引擎与数据库
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	x, err = xorm.NewEngine("postgres", psqlInfo)
	if err != nil {
		glog.Fatalf("Fail to create engine: %v\n", err)
	}

	// 同步结构体与数据表
	if err = x.Sync(new(Account)); err != nil {
		glog.Fatalf("Fail to sync database: %v\n", err)
	}
}

// 创建新的账户
func newAccount(name string, balance float64) error {
	// 对未存在记录进行插入
	_, err := x.Insert(&Account{Name: name, Balance: balance})
	return err
}

// 获取账户信息
func getAccount(id int64) (*Account, error) {
	a := &Account{}
	// 直接操作 ID 的简便方法
	has, err := x.Id(id).Get(a)
	// 判断操作是否发生错误或对象是否存在
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("Account does not exist")
	}
	return a, nil
}

// 用户存款
func makeDeposit(id int64, deposit float64) (*Account, error) {
	a, err := getAccount(id)
	if err != nil {
		return nil, err
	}
	a.Balance += deposit
	// 对已有记录进行更新
	_, err = x.Update(a)
	return a, err
}

// 用户取款
func makeWithdraw(id int64, withdraw float64) (*Account, error) {
	a, err := getAccount(id)
	if err != nil {
		return nil, err
	}
	if a.Balance < withdraw {
		return nil, errors.New("Not enough balance")
	}
	a.Balance -= withdraw
	_, err = x.Update(a)
	return a, err
}

// 用户转账
func makeTransfer(id1, id2 int64, balance float64) error {
	a1, err := getAccount(id1)
	if err != nil {
		return err
	}

	a2, err := getAccount(id2)
	if err != nil {
		return err
	}

	if a1.Balance < balance {
		return errors.New("Not enough balance")
	}

	// 下面代码存在问题，需要采用事务回滚来改进（课时 2）
	a1.Balance -= balance
	a2.Balance += balance
	if _, err = x.Update(a1); err != nil {
		return err
	} else if _, err = x.Update(a2); err != nil {
		return err
	}

	return nil
}

// 按照 ID 正序排序返回所有账户
func getAccountsAscId() (as []Account, err error) {
	// 使用 Find 方法批量获取记录
	err = x.Find(&as)
	return as, err
}

// 按照存款倒序排序返回所有账户
func getAccountsDescBalance() (as []Account, err error) {
	// 使用 Desc 方法使结果呈倒序排序
	err = x.Desc("balance").Find(&as)
	return as, err
}

// 删除账户
func deleteAccount(id int64) error {
	// 通过 Delete 方法删除记录
	_, err := x.Delete(&Account{Id: id})
	return err
}

func main() {
	fmt.Println("Welcome bank of xorm!")
Exit:
	for {
		fmt.Println(prompt)

		// 获取用户输入并转换为数字
		var num int
		fmt.Scanf("%d\n", &num)

		// 判断用户选择
		switch num {
		case 1:
			fmt.Println("Please enter <name> <balance>:")
			var name string
			var balance float64
			fmt.Scanf("%s %f\n", &name, &balance)
			if err := newAccount(name, balance); err != nil {
				fmt.Println("Fail to create new account:", err)
			} else {
				fmt.Println("New account has been created")
			}
		case 2:
			fmt.Println("Please enter <id>:")
			var id int64
			fmt.Scanf("%d\n", &id)
			a, err := getAccount(id)
			if err != nil {
				fmt.Println("Fail to get account:", err)
			} else {
				fmt.Printf("%#v\n", a)
			}
		case 3:
			fmt.Println("Please enter <id> <deposit>:")
			var id int64
			var deposit float64
			fmt.Scanf("%d %f\n", &id, &deposit)
			a, err := makeDeposit(id, deposit)
			if err != nil {
				fmt.Println("Fail to deposit:", err)
			} else {
				fmt.Printf("%#v\n", a)
			}
		case 4:
			fmt.Println("Please enter <id> <withdraw>:")
			var id int64
			var withdraw float64
			fmt.Scanf("%d %f\n", &id, &withdraw)
			a, err := makeWithdraw(id, withdraw)
			if err != nil {
				fmt.Println("Fail to withdraw:", err)
			} else {
				fmt.Printf("%#v\n", a)
			}
		case 5:
			fmt.Println("Please enter <id> <balance> <id>:")
			var id1, id2 int64
			var balance float64
			fmt.Scanf("%d %f %d\n", &id1, &balance, &id2)
			if err := makeTransfer(id1, id2, balance); err != nil {
				fmt.Println("Fail to transfer:", err)
			} else {
				fmt.Println("Transfer has been made")
			}
		case 6:
			as, err := getAccountsAscId()
			if err != nil {
				fmt.Println("Fail to get accounts:", err)
			} else {
				for i, a := range as {
					fmt.Printf("%d: %#v\n", i+1, a)
				}
			}
		case 7:
			as, err := getAccountsDescBalance()
			if err != nil {
				fmt.Println("Fail to get accounts:", err)
			} else {
				for i, a := range as {
					fmt.Printf("%d: %#v\n", i+1, a)
				}
			}
		case 8:
			fmt.Println("Please enter <id>:")
			var id int64
			fmt.Scanf("%d\n", &id)
			if err := deleteAccount(id); err != nil {
				fmt.Println("Fail to delete account:", err)
			} else {
				fmt.Println("Account has been deleted")
			}
		case 9:
			fmt.Println("Thank you! Hope see you again soon!")
			break Exit
		default:
			fmt.Println("Unknown operation number:", num)
		}
		fmt.Println()
	}
}
