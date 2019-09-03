package main

import (

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


func GetConnection() (connection *gorm.DB, err error){
   connection, err = gorm.Open(DB_TYPE, MYSQL_CONNECT)

   if err != nil {
   	 return
   }

   connection.SingularTable(true)
   connection.LogMode(true)
   connection.DB().SetConnMaxLifetime(0)

   return
}

func (add *Address) Save(connection *gorm.DB)  {
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
    //Where("id > 500").
    Find(&list)

  return list
}