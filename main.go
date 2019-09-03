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
   OPTION_LIST = "list"
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
         list          read list of names of force brute: ./generator list path_of_file\n 
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
    }else if option == OPTION_LIST {
    	file_path := os.Args[2]
        read_dictionary_bracket(file_path)
	}else{
       fmt.Println("Option Default")
       default_func()
	}

}

func default_func(){
	oka,   bi_a := getInt("45408662446006351146498425493603101118929405751231593740963758434475737113700", 10) 
    _, bi_end   := getInt("57669001306428065956053000376875938421040345304064124051023973211784186134399", 10)
 
    for bi_a.Cmp(bi_end) == -1 {
	    fmt.Println(bi_a)
	    
	    if oka {
			execute(conn, bi_a)
		}

		_, y := getInt("1", 10)
        bi_a = bi_a.Add(bi_a, y)
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

func read_dictionary_bracket(file_path string){

	file, err := os.Open(file_path)
	if err != nil {
		fmt.Printf(" > Failed!: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(file)
	line := ""
	for scanner.Scan() {
		line       = scanner.Text()
		i_i := 0
		change_    := false
		line_clean := line
		for i, r := range line {
			if string(r) == "(" {
				i_i = i-1
				change_ = true
				break
			}
        } 
        if change_ {
        	line_clean = strings.Trim(line[:i_i], " ")
        }		
		str_sha256_pass    :=  SHA256( line_clean )

   	    oka, bi_a := getInt(str_sha256_pass, 16)
	    if oka {
			fmt.Println(line_clean, " ", bi_a)
			execute(conn, bi_a)
		 }
        line = ""
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

