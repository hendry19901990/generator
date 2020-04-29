package main

import (

	"os"
	"bufio"
	"strings"
	"fmt"
    "math/big"
    "crypto/sha256"
    "encoding/json"

	"github.com/vsergeev/btckeygenie/btckey"
	"github.com/jinzhu/gorm"

)

const (
	DB_TYPE       = "mysql"
    MYSQL_CONNECT = "test:test@/meme?charset=utf8&parseTime=True&loc=Local"
)

// OPTIONS
const (
   OPTION_COUNTRIES     = "countries"
   OPTION_PRIME_NUMBERS = "prime"
   OPTION_ONLY_PRIME_NUMBERS = "only_prime"
   OPTION_LINE = "line"
   OPTION_BIT = "bits"
   OPTION_DICTIONARY = "dictionary"
   OPTION_ALL = "all"
   OPTION_LIST = "list"
   OPTION_PERIODIC_TABLE = "periodic_table"
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
         only_prime    read list of prime numbers:    ./generator only_prime path_of_file  \n
         bits          for loop between bits size:  ./generator bits \n
         dictionary    read list of possible password of force brute: ./generator dictionary path_of_file\n
         list          read list of names of force brute: ./generator list path_of_file\n
         all           read list of DB and call api to prove if it has balance: ./generator all\n
         periodic_table   read file of any_type: ./generator path_of_file\n

`

const str_code = `{"BD": "880", "BE": "32", "BF": "226", "BG": "359", "BA": "387", "BB": "+1-246", "WF": "681", "BL": "590", "BM": "+1-441", "BN": "673", "BO": "591", "BH": "973", "BI": "257", "BJ": "229", "BT": "975", "JM": "+1-876", "BV": "", "BW": "267", "WS": "685", "BQ": "599", "BR": "55", "BS": "+1-242", "JE": "+44-1534", "BY": "375", "BZ": "501", "RU": "7", "RW": "250", "RS": "381", "TL": "670", "RE": "262", "TM": "993", "TJ": "992", "RO": "40", "TK": "690", "GW": "245", "GU": "+1-671", "GT": "502", "GS": "", "GR": "30", "GQ": "240", "GP": "590", "JP": "81", "GY": "592", "GG": "+44-1481", "GF": "594", "GE": "995", "GD": "+1-473", "GB": "44", "GA": "241", "SV": "503", "GN": "224", "GM": "220", "GL": "299", "GI": "350", "GH": "233", "OM": "968", "TN": "216", "JO": "962", "HR": "385", "HT": "509", "HU": "36", "HK": "852", "HN": "504", "HM": " ", "VE": "58", "PR": "+1-787 and 1-939", "PS": "970", "PW": "680", "PT": "351", "SJ": "47", "PY": "595", "IQ": "964", "PA": "507", "PF": "689", "PG": "675", "PE": "51", "PK": "92", "PH": "63", "PN": "870", "PL": "48", "PM": "508", "ZM": "260", "EH": "212", "EE": "372", "EG": "20", "ZA": "27", "EC": "593", "IT": "39", "VN": "84", "SB": "677", "ET": "251", "SO": "252", "ZW": "263", "SA": "966", "ES": "34", "ER": "291", "ME": "382", "MD": "373", "MG": "261", "MF": "590", "MA": "212", "MC": "377", "UZ": "998", "MM": "95", "ML": "223", "MO": "853", "MN": "976", "MH": "692", "MK": "389", "MU": "230", "MT": "356", "MW": "265", "MV": "960", "MQ": "596", "MP": "+1-670", "MS": "+1-664", "MR": "222", "IM": "+44-1624", "UG": "256", "TZ": "255", "MY": "60", "MX": "52", "IL": "972", "FR": "33", "IO": "246", "SH": "290", "FI": "358", "FJ": "679", "FK": "500", "FM": "691", "FO": "298", "NI": "505", "NL": "31", "NO": "47", "NA": "264", "VU": "678", "NC": "687", "NE": "227", "NF": "672", "NG": "234", "NZ": "64", "NP": "977", "NR": "674", "NU": "683", "CK": "682", "XK": "", "CI": "225", "CH": "41", "CO": "57", "CN": "86", "CM": "237", "CL": "56", "CC": "61", "CA": "1", "CG": "242", "CF": "236", "CD": "243", "CZ": "420", "CY": "357", "CX": "61", "CR": "506", "CW": "599", "CV": "238", "CU": "53", "SZ": "268", "SY": "963", "SX": "599", "KG": "996", "KE": "254", "SS": "211", "SR": "597", "KI": "686", "KH": "855", "KN": "+1-869", "KM": "269", "ST": "239", "SK": "421", "KR": "82", "SI": "386", "KP": "850", "KW": "965", "SN": "221", "SM": "378", "SL": "232", "SC": "248", "KZ": "7", "KY": "+1-345", "SG": "65", "SE": "46", "SD": "249", "DO": "+1-809 and 1-829", "DM": "+1-767", "DJ": "253", "DK": "45", "VG": "+1-284", "DE": "49", "YE": "967", "DZ": "213", "US": "1", "UY": "598", "YT": "262", "UM": "1", "LB": "961", "LC": "+1-758", "LA": "856", "TV": "688", "TW": "886", "TT": "+1-868", "TR": "90", "LK": "94", "LI": "423", "LV": "371", "TO": "676", "LT": "370", "LU": "352", "LR": "231", "LS": "266", "TH": "66", "TF": "", "TG": "228", "TD": "235", "TC": "+1-649", "LY": "218", "VA": "379", "VC": "+1-784", "AE": "971", "AD": "376", "AG": "+1-268", "AF": "93", "AI": "+1-264", "VI": "+1-340", "IS": "354", "IR": "98", "AM": "374", "AL": "355", "AO": "244", "AQ": "", "AS": "+1-684", "AR": "54", "AU": "61", "AT": "43", "AW": "297", "IN": "91", "AX": "+358-18", "AZ": "994", "IE": "353", "ID": "62", "UA": "380", "QA": "974", "MZ": "258"}`

var richMap map[string]struct{}

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

   richMap = FillMap(conn)

	option := os.Args[1]
	if option == OPTION_COUNTRIES{
		read_countries(getCountries())
	}else if option == OPTION_PRIME_NUMBERS {
		file_path := os.Args[2]
		read_prime_numbers(file_path)
	}else if option == OPTION_ONLY_PRIME_NUMBERS {
		file_path := os.Args[2]
		read_only_prime_numbers(file_path)
	}else if option == OPTION_LINE {
		file_path := os.Args[2]
		read_prime_numbers_extended(file_path)
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
	} else if option == OPTION_PERIODIC_TABLE {
    	file_path := os.Args[2]
      read_periodic_table(file_path)
	}else{
       fmt.Println("Option Default")
       getCodes()
	}

}

func getCodes(){
	list := getCountries()
	map_list := make(map[string]string)

	json.Unmarshal([]byte(str_code), &map_list)
	fmt.Println(len(map_list))
	for _, country := range list{
		phone_code, ok := map_list[country.Alpha2Code]
		if ok {
		   if len(phone_code) > 0 && phone_code[0:1] != "+"{
		   	  phone_code = "+" + phone_code
		   }
           //conn.Exec(fmt.Sprintf("INSERT INTO org_country(name, alpha_2_code, phone_code) VALUES('%s','%s','%s')",
          //           country.Name, country.Alpha2Code, phone_code))
		   fmt.Println(phone_code)
		}
	}
}

func default_func(){
	oka,   bi_a := getInt("45408662446006351146498425493603101118929405751231593740963758434475737141474", 10)
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

func read_periodic_table(file_path string){

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

	_, value := getInt("1766847064778384329583297500742918515827483896875618958121606201292273638", 10)

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

	n_set_sizes := []string{
		"72653859913610161834397480789764961790287049201970549985542013495161179381920",
		"86276458647412067178347008437845892125965870927340028107831141025503900516030",
		"90817324892012702292996850987206202237858811502463187481927516868951474227400",
		"95358191136613337407646693536566512349751752077586346856023892712399047938770",
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		list := strings.Split(line, "\t")

		for _, input := range list {
			for _, n_set := range n_set_sizes {
				part_1 := n_set[0:33]
		        part_2 := n_set[45:len(n_set)]
				oka, bi_a := getInt(part_1 + input + part_2, 10)

				if oka {
					fmt.Print(input, " ")
					execute(conn, bi_a)
				}

				okb, bi_b := getInt(input, 10)
				if okb {
					execute(conn, bi_b)
				}
			}
		}
		fmt.Print("\n")

    }

    if scanner.Err() != nil {
        fmt.Printf(" > Failed!: %v\n", scanner.Err())
    }
}

func read_only_prime_numbers(file_path string) {
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
			oka, bi_a := getInt(input, 10)
			if oka {
				fmt.Print(input, "  ")
				execute(conn, bi_a)
			}
		}
		fmt.Println("")
    }

    if scanner.Err() != nil {
        fmt.Printf(" > Failed!: %v\n", scanner.Err())
    }
}

func read_prime_numbers_extended(file_path string) {
    file, err := os.Open(file_path)
	if err != nil {
		fmt.Printf(" > Failed!: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		str_sha256 := SHA256(line)

		list := strings.Split(line, "\t")

        okb, bi_b := getInt(str_sha256, 16)

        if okb {
			fmt.Print(list[0], " ", bi_b.String(), "\n")
			execute(conn, bi_b)
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

	_, ok1 := richMap[address]
	if ok1 {
        fmt.Println(" YES address")
        addressDB := AddressRust{Private: priv.ToWIF(), Public: address}
        addressDB.Save(conn)
 	}

    _, ok2 := richMap[address_compressed]
	if ok2 {
        fmt.Println(" YES address_compressed")
        addressDB := AddressRust{Private: priv.ToWIFC(), Public: address_compressed}
        addressDB.Save(conn)
 	}
}

func executeOLD(conn *gorm.DB, bi *big.Int){
    priv               := btckey.NewPrivateKey(bi)
	address            := priv.PublicKey.ToAddressUncompressed()
	address_compressed := priv.PublicKey.ToAddress()

	if Exist(address, conn){
        fmt.Println(" YES address")
        addressDB := AddressRust{Private: priv.ToWIF(), Public: address}
        addressDB.Save(conn)
 	}

	if Exist(address_compressed, conn){
        fmt.Println(" YES address_compressed")
        addressDB := AddressRust{Private: priv.ToWIFC(), Public: address_compressed}
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
