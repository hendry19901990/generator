package main

import (

  "log"

	"github.com/jinzhu/gorm"
  _ "github.com/go-sql-driver/mysql"
    
)

type Address struct{
    Id                 int64      `json:"id" gorm:"primary_key"`
    Private            string     `gorm:"column:private"`
    PrivateCompressed  string     `gorm:"column:private_compressed"`
    Public             string     `gorm:"column:public"`
    PublicCompressed   string     `gorm:"column:public_compressed"`
}

type AddressRust struct{
    Id                 int64      `json:"id" gorm:"primary_key"`
    Private            string     `gorm:"column:wif"`
    Public             string     `gorm:"column:addr"`
}

func (AddressRust) TableName() string {
  return "address_rust"
}

type RichList struct{
    Id                 int64      `json:"id" gorm:"primary_key"`
    Address            string     `gorm:"column:address"`
}

func (RichList) TableName() string {
  return "rich_list"
}


func GetConnection() (connection *gorm.DB, err error){
   connection, err = gorm.Open(DB_TYPE, MYSQL_CONNECT)

   if err != nil {
   	 return
   }

   connection.SingularTable(true)
   connection.LogMode(false)
   connection.DB().SetConnMaxLifetime(0)

   return
}

func (add *Address) Save(connection *gorm.DB)  {
	connection.NewRecord(add)
	connection.Create(add)
}

func (add *AddressRust) Save(connection *gorm.DB)  {
  connection.NewRecord(add)
  connection.Create(add)
}

func NewAddressDB( priv, priv_comp, pub, pub_comp string) (add *Address){
    add = &Address{ Private: priv, PrivateCompressed: priv_comp, Public: pub, PublicCompressed: pub_comp}
    return
}

func All(conn *gorm.DB) []Address {
  list := make([]Address,0)

  conn.
    Where("id > 1500").
    Find(&list)

  return list
}

func Exist(address string, conn *gorm.DB) bool {
   var rich_list RichList

   conn.
    Where("address = ?", address).
    Find(&rich_list)

  return rich_list.Id != 0
}

func FillMap(conn *gorm.DB) map[string]struct{} {
  list := make(map[string]struct{})

  rows, _ := conn.Table("rich_list").
       Select("address as address").
       Rows()

  for rows.Next() {
     var address string
     if err := rows.Scan(&address); err != nil {
       log.Println(err)
     }else{
       list[address] = struct{}{}
     }
  }

  return list
} 