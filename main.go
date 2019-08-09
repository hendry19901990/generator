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

// OPTIONS
const (
   OPTION_COUNTRIES     = "countries"
   OPTION_PRIME_NUMBERS = "prime"
   OPTION_BIT = "bits"
   OPTION_DICTIONARY = "dictionary"
   OPTION_ALL = "all"
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

var conn *gorm.DB

const OPTIONS = `
      ./generator <option>  \n
      option: \n 
         countries     read list of contries from api rest: ./generator countries \n
         prime         read list of prime numbers:    ./generator prime path_of_file  \n
         bits          for loop between bits size:  ./generator bits \n 
         dictionary    read list of possible password of force brute: ./generator dictionary path_of_file\n 
         all           read list of DB and call api to prove if it has balance: ./generator all\n 

`

func main(){
   if len(os.Args) == 1 {
   	 fmt.Println(OPTIONS)
   	 return
   }

	var errDB error
	conn, errDB = GetConnection()
	if errDB != nil {
		fmt.Printf(" > Failed!: %v\n", errDB)
		return
	}

	option := os.Args[1]
	if option == OPTION_COUNTRIES{
		read_countries(getCountries())        
	}else if option == OPTION_PRIME_NUMBERS {
		file_path := os.Args[2]
		read_prime_numbers(file_path)
	}else if option == OPTION_BIT {
        read_every_bit()
    }else if option == OPTION_DICTIONARY {
    	file_path := os.Args[2]
        read_dictionary(file_path)
    }else if option == OPTION_ALL {
        read_all()
	}else{
       fmt.Println("Option Default")
       default_func()
	}

}

func default_func(){
	list := []string{
	 		"4669523849932130508876392554713407521319117239637943224980015676156491",
			"4906275427767802358357703730938087362176142642699093827933107888253709",
			"2409130781894986571956777721649968801511465915451196376269177305066867",
			"7595009151080016652449223792726748985452052945413160073645842090827711",
			"3822535632033509464266159811805197854872067042990716005808372194664933",
			"5885903965180586669073549360644800583458138238012033647539649735017287",
			"5850725702766829291491370712136286009948642125131436113342815786444567",
			"4237080979868607742750808600846638318022863593147774739556427943294937",
			"3773180816219384606784189538899553110499442295782576702222280384917551",
			"9547848065153773335707495885453566120069130270246768806790708393909999",
       }

    for _, v := range list{
    	oka, bi_a := getInt(v, 10)
	    if oka {
			execute(conn, bi_a)
		}

		str_sha256_v    :=  SHA256(v)
		okb, bi_b := getInt(str_sha256_v, 16)
		if okb {
			execute(conn, bi_b)
	    }
	    fmt.Println(bi_a)
    }
}

func read_all(){
	list := All(conn)
	for _, add := range list {
       if Call(add.Public, false){
          fmt.Println("This one has Balance", add.Public)
       }

       if Call(add.PublicCompressed, false){
          fmt.Println("This one has Balance", add.PublicCompressed)
       }
	}
	fmt.Println("read_all just executed")
}

func read_dictionary(file_path string){

	file, err := os.Open(file_path)
	if err != nil {
		fmt.Printf(" > Failed!: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		str_sha256_pass    :=  SHA256( line )

   	    oka, bi_a := getInt(str_sha256_pass, 16)
	    if oka {
			fmt.Println(line, " ", bi_a)
			execute(conn, bi_a)
		 }
	}

}

func read_every_bit(){

	_, value := getInt("1766847064778384329583297500742918515827483896875618958121606201292306265", 10)

	for value.BitLen() > 239 {
		fmt.Println(value.BitLen(), " = ", value.String())
        execute(conn, value)
        _, y := getInt("1", 10)
        value = value.Sub(value, y)
	}
}

func read_countries(list []Country){
   for _, country := range list {
   	   str_sha256_name    :=  SHA256(strings.ToUpper(country.Name))
   	   str_sha256_capital :=  SHA256(strings.ToUpper(country.Capital))

   	   oka, bi_a := getInt(str_sha256_name, 16)
	   okb, bi_b := getInt(str_sha256_capital, 16)

	   if oka {
			fmt.Print(country.Name, " ", bi_a.BitLen())
			execute(conn, bi_a)
		}

		if okb {
			fmt.Print(" ", str_sha256_capital, "\n")
			execute(conn, bi_b)
	    }
   }
}

func read_prime_numbers(file_path string) {
    file, err := os.Open(file_path)
	if err != nil {
		fmt.Printf(" > Failed!: %v\n", err)
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

	if Call(address, true) || Call(address_compressed, true){
        fmt.Println(" YES")
        addressDB := NewAddressDB(priv.ToWIF(), priv.ToWIFC(), address, address_compressed)
        addressDB.Save(conn)
        addressDB = nil
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

