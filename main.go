package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

type State struct {
	Atm_id     string
	State_date string
	State      string
}

type ATM struct {
	ATM_id string
}

func main() {

	// DB Connection Start
	var db *sql.DB
	var st_date, st_date2 time.Time
	var date_only string
	var atm_id []string
	var placeholder_data []string
	// location, _ := time.LoadLocation("Asia/Bangkok")

	cfg := mysql.Config{
		User:                 "appops",
		Passwd:               "J@l1n0ps123!",
		Net:                  "tcp",
		Addr:                 "10.133.5.64:6033",
		DBName:               "atm",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	// DB Connection End

	// Summary Workflow List Start
	var sts []State
	rows, err := db.Query("SELECT atm_id, DATE, state FROM history_state WHERE state in (6,2) AND DATE_FORMAT(CONVERT_TZ(date,'+00:00','+00:00'), '%Y-%m-%d') = DATE_SUB(CURDATE(), INTERVAL 1 DAY) order by date asc")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var st State
		if err := rows.Scan(&st.Atm_id, &st.State_date, &st.State); err != nil {
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

						string_diff := diff.String()
						// out := time.Time{}.Add(diff)
						diff_format_seconds, _ := time.ParseDuration(string_diff)
						fmt.Println(diff_format_seconds.Seconds())
						// fmt.Println(diff)
						total = total + diff

					}
					date_only = states[i].State_date[0:10]
					// fmt.Println(date_only, ";", states[i].Atm_id, ";", total)
					check = total
				}
			}
		}
		if check != 0 {
			str_data := date_only + ";" + atm_id[k] + ";" + check.String()
			placeholder_data = append(placeholder_data, str_data)
			// fmt.Println(date_only, ";", atm_id[k], ";", check)
			//fmt.Println(str_data)
		}
	}
	fmt.Println(placeholder_data)

	currentTime := time.Now().AddDate(0, 0, -1).Format("20060102_")
	fileName := currentTime + "terminal_downtime_report.csv"
	fmt.Println(fileName)
	file, err := os.OpenFile(currentTime+"terminal_downtime_report.csv", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("File does not exists or cannot be created")
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write(placeholder_data)
	w := bufio.NewWriter(file)

	fmt.Fprintf(w, "%v\n", "Date;ATM ID;Downtime Duration")
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
