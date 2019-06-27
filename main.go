package main

import (

	"os"
	"bufio"
	"strings"
	"fmt"
    "math/big"
    "crypto/sha256"

	"github.com/vsergeev/btckeygenie/btckey"
	
)

const (
	DB_TYPE = "mysql"
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

				priv_a :=  btckey.NewPrivateKey(bi_a)
				address_a            := priv_a.PublicKey.ToAddress()
				address_compressed_a := priv_a.PublicKey.ToAddress()

				fmt.Print(input)
				if Call(address_a) || Call(address_compressed_a){
			        fmt.Println(" YES")
			        addressDB_a := NewAddressDB(priv_a.ToWIF(), priv_a.ToWIFC(), address_a, address_compressed_a)
			        addressDB_a.Save(conn)
				}
			}

			 
			if okb {
                
				priv_b :=  btckey.NewPrivateKey(bi_b)
				address_b            := priv_b.PublicKey.ToAddress()
				address_compressed_b := priv_b.PublicKey.ToAddress()

				fmt.Print(" ", bi_b, "\n")
				if Call(address_b) || Call(address_compressed_b){
			        fmt.Println(" YES")
			        addressDB_b := NewAddressDB(priv_b.ToWIF(), priv_b.ToWIFC(), address_b, address_compressed_b)
			        addressDB_b.Save(conn)
				}
		    }
			
		}
 
    }

    if scanner.Err() != nil {
        fmt.Printf(" > Failed!: %v\n", scanner.Err())
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

