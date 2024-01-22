package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
)

type State struct {
	Atm_id     string
	State_date string
	State      string
	Mandant_id string
}

type ATM struct {
	ATM_id string
}

func main() {

	// DB Connection Start
	var db *sql.DB
	var st_date, st_date2 time.Time
	// var date_only string
	var atm_id []string
	var placeholder_data []string
	var err error
	var arg_first_date, arg_second_date, bank string
	var mandant int

	// location, _ := time.LoadLocation("Asia/Bangkok")

	// today := time.Now()
	yesterday := time.Now().AddDate(0, 0, -1)
	arg_first_date = yesterday.Format("2006-01-02")
	arg_second_date = yesterday.Format("2006-01-02")

	// if len(os.Args) == 4 {
	// 	arg_first_date = os.Args[1]
	// 	bank = os.Args[2]
	// 	bank = strings.ToUpper(bank)
	// 	// } else if len(os.Args) == 2 {
	// 	// 	arg_first_date = os.Args[1]
	// 	// 	arg_second_date = os.Args[1]
	// } else
	if len(os.Args) == 4 {
		bank = os.Args[1]
		bank = strings.ToUpper(bank)
		arg_first_date = os.Args[2]
		arg_second_date = os.Args[3]

		// } else if len(os.Args) == 2 {
		// 	arg_first_date = os.Args[1]
		// 	arg_second_date = os.Args[1]
	} else if len(os.Args) == 3 {
		bank = os.Args[1]
		bank = strings.ToUpper(bank)
		arg_first_date = os.Args[2]
		arg_second_date = os.Args[2]
	} else if len(os.Args) == 2 {
		bank = os.Args[1]
		bank = strings.ToUpper(bank)
	}

	if bank == "BTN" {
		mandant = 21
	} else if bank == "BNI" {
		mandant = 20
	} else if bank == "BRI" {
		mandant = 27
	} else if bank == "BMRI" {
		mandant = 24
	}

	fmt.Println("Command         : ", os.Args)
	fmt.Println("Length Argument : ", len(os.Args))
	fmt.Println("First Date      : ", arg_first_date)
	fmt.Println("Second Date     : ", arg_second_date)
	fmt.Println("Bank            : ", bank)
	fmt.Println("Mandant         : ", mandant)

	fmt.Println("------------------------")
	// fmt.Println("Reading Config File")
	conf, err := ini.Load("cfg.ini")
	if err != nil {
		log.Fatal(err)
		// os.Exit(1)
	}
	fmt.Println("Read Config File Success")

	// conf = conf
	// fmt.Println("App Mode:", conf.Section("").Key("app_mode").String())

	cfg := mysql.Config{
		User:                 conf.Section("db").Key("user").String(),
		Passwd:               conf.Section("db").Key("pass").String(),
		Net:                  conf.Section("db").Key("net").String(),
		Addr:                 conf.Section("db").Key("host").String(),
		DBName:               conf.Section("db").Key("db_name").String(),
		AllowNativePasswords: true,
	}

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("------------------------")
	log.Default()
	fmt.Println("Connected to DB!")
	fmt.Println("------------------------")
	// DB Connection End

	// Summary Workflow List Start
	var sts []State
	// rows, err := db.Query("SELECT atm_id, DATE, state FROM history_state WHERE state in (6,2) AND DATE_FORMAT(CONVERT_TZ(date,'+00:00','+00:00'), '%Y-%m-%d') = DATE_SUB(CURDATE(), INTERVAL 1 DAY) order by date asc")
	rows, err := db.Query("SELECT atm_id, DATE, history_state.state, mandant_id FROM history_state INNER JOIN atm ON history_state.atm_id = atm.id WHERE history_state.state in (6,2) AND DATE_FORMAT(CONVERT_TZ(date,'+00:00','+00:00'), '%Y-%m-%d') between ? and ? and mandant_id = ? order by date asc", arg_first_date, arg_second_date, mandant)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var st State
		if err := rows.Scan(&st.Atm_id, &st.State_date, &st.State, &st.Mandant_id); err != nil {
			log.Fatal(err)
		}
		sts = append(sts, st)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	// Summary Workflow List End

	// Summary Workflow List Start
	var atms []ATM
	atm_rows, err := db.Query("SELECT id from atm")
	if err != nil {
		log.Fatal(err)
	}
	defer atm_rows.Close()

	for atm_rows.Next() {
		var atm ATM
		if err := atm_rows.Scan(&atm.ATM_id); err != nil {
			log.Fatal(err)
		}
		atms = append(atms, atm)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	// Summary Workflow List End

	// Working With Workflow Data Start
	states := sts
	// fmt.Println(states)

	for i := 0; i < len(atms); i++ {
		atm_id = append(atm_id, atms[i].ATM_id)
	}

	// fmt.Println(atms)

	// reformat date time to local jakarta time for all workflows start
	layoutFormat := "2006-01-02 15:04:05"
	var total, diff time.Duration

	for k := 0; k < len(atm_id); k++ {
		// fmt.Println(atm_id[k])
		total = 0
		var check time.Duration
		for i := 0; i < len(states); i++ {
			if atm_id[k] == states[i].Atm_id {
				if states[i].State == "6" {
					if states[i+1].State == "2" {
						st_date, _ = time.Parse(layoutFormat, states[i].State_date)
						st_date2, _ = time.Parse(layoutFormat, states[i+1].State_date)
						diff = st_date2.Sub(st_date)

						// string_diff := diff.String()
						// out := time.Time{}.Add(diff)
						// diff_format_seconds, _ := time.ParseDuration(string_diff)
						// fmt.Println(diff_format_seconds.Seconds())
						// fmt.Println(diff)
						total = total + diff

					}
					// date_only = states[i].State_date[0:10]
					// fmt.Println(date_only, ";", states[i].Atm_id, ";", total)
					check = total
				}
			}
		}
		if check != 0 {
			sec_only := check.Seconds()
			string_sec_only := fmt.Sprint(sec_only)
			// str_data := date_only + ";" + atm_id[k] + ";" + string_sec_only
			str_data := atm_id[k] + ";" + string_sec_only
			placeholder_data = append(placeholder_data, str_data)
			// fmt.Println(date_only, ";", atm_id[k], ";", check)
			//fmt.Println(str_data)
		}
	}
	// fmt.Println(placeholder_data)
	fmt.Println("Count Of Data   : ", len(placeholder_data))

	filename_date, _ := time.Parse("2006-01-02", arg_first_date)
	filename_date_2, _ := time.Parse("2006-01-02", arg_second_date)
	currentTime := filename_date.Format("20060102_")
	currentTime2 := filename_date_2.Format("20060102_")

	fileName := currentTime + currentTime2 + bank + "_terminal_downtime_report.csv"
	fmt.Println("Success Create CSV")
	fmt.Println("File Created    : ", fileName)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("File does not exists or cannot be created")
		os.Exit(1)
	}
	defer file.Close()

	// writer := csv.NewWriter(file)
	// writer.Write(placeholder_data)
	w := bufio.NewWriter(file)

	fmt.Fprintf(w, "%v\n", "Report Downtime ATM Periode "+arg_first_date+" S/D "+arg_second_date)
	fmt.Fprintf(w, "%v\n", "")
	fmt.Fprintf(w, "%v\n", "ATM ID;Downtime (Seconds)")
	// fmt.Fprintf(w, "%v\n", "Date;ATM ID;Downtime (Seconds)")
	for i := 0; i < len(placeholder_data); i++ {
		fmt.Fprintf(w, "%v\n", placeholder_data[i])
	}

	w.Flush()

	// for i := 0; i < len(states); i++ {
	// 	// fmt.Println(st_date)
	// 	if states[i].State == "6" {
	// 		if states[i+1].State == "2" {
	// 			st_date, _ = time.Parse(layoutFormat, states[i].State_date)
	// 			st_date2, _ = time.Parse(layoutFormat, states[i+1].State_date)
	// 			diff = st_date2.Sub(st_date)
	// 			total = total + diff

	// 			date_only = states[i].State_date[0:10]
	// 			fmt.Println(date_only, ";", states[i].Atm_id, ";", total)
	// 		}
	// 	}
	// }

	// reformat date time to local jakarta time for all workflows end
	// reformat date time to local jakarta time for all failed workflows end
	// Working With Workflow Data end

	// Send workflow data to template email start
	// var tmplFile = "template.html"
	// tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	// if err != nil {
	// 	panic(err)
	// }

	// buf := new(bytes.Buffer)
	// if err = tmpl.Execute(buf, sts); err != nil {
	// 	fmt.Println(err)
	// }

	// result := buf.String()

	// fmt.Println(result)
	// Send workflow data to template email end
}

// func ReadConfig() Config {
// 	var configfile = flags.Configfile
// 	_, err := os.Stat(configfile)
// 	if err != nil {
// 		log.Fatal("Config file is missing: ", configfile)
// 	}

// 	var config Config
// 	if _, err := toml.DecodeFile(configfile, &config); err != nil {
// 		log.Fatal(err)
// 	}
// 	//log.Print(config.Index)
// 	return config
// }
