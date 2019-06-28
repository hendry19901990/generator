package main

import (

	"os"
	"bufio"
	"strings"
	"fmt"
    "math/big"
    "crypto/sha256"

	"github.com/vsergeev/btckeygenie/btckey"
	"github.com/jinzhu/gorm"
	
)

const (
	DB_TYPE       = "mysql"
    MYSQL_CONNECT = "test:test@/wallet?charset=utf8&parseTime=True&loc=Local"
)

/*

CREATE DATABASE `wallet`;

CREATE TABLE `wallet`.`address` ( 
	`id` INT NOT NULL AUTO_INCREMENT , 
	`private` VARCHAR(100) NOT NULL ,
	`private_compressed` VARCHAR(100) NOT NULL ,
	`public` VARCHAR(100) NOT NULL ,
	`public_compressed` VARCHAR(100) NOT NULL ,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB;

*/


func main(){

	file_path := os.Args[1]
	file, err := os.Open(file_path)

	if err != nil {
		fmt.Printf(" > Failed!: %v\n", err)
		return
	}

	conn, errDB := GetConnection()
	if errDB != nil {
		fmt.Printf(" > Failed!: %v\n", errDB)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		list := strings.Split(line, "\t")

		for _, input := range list {

			str_sha256 := SHA256(input)

			oka, bi_a := getInt(input, 10)
			okb, bi_b := getInt(str_sha256, 16)

			if oka {
				fmt.Print(input)
				execute(conn, bi_a)
			}

			if okb {
				fmt.Print(" ", bi_b.String(), "\n")
				execute(conn, bi_b)
		    }
			
		}
 
    }

    if scanner.Err() != nil {
        fmt.Printf(" > Failed!: %v\n", scanner.Err())
    }


}

func execute(conn *gorm.DB, bi *big.Int){
    priv               := btckey.NewPrivateKey(bi)
	address            := priv.PublicKey.ToAddressUncompressed()
	address_compressed := priv.PublicKey.ToAddress()

	if Call(address) || Call(address_compressed){
        fmt.Println(" YES")
        addressDB := NewAddressDB(priv.ToWIF(), priv.ToWIFC(), address, address_compressed)
        addressDB.Save(conn)
	}
}

func getInt(input string, type_ int) (bool, *big.Int) {
	bi    := big.NewInt(0)
	_, ok := bi.SetString(input, type_)
	return ok, bi
}

func SHA256(str string) string {
    hash := sha256.Sum256([]byte(str))
    return fmt.Sprintf("%x", hash)
}

