//

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

type Meeting struct {
	Title string
	Date  string
	Users []string
}

func contains(u []string, e string) bool {
	for _, v := range u {
		if v == e {
			return true
		}
	}
	return false
}

func readInput(p string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	if p == "" {
		pc, fn, line, _ := runtime.Caller(1)
		fmt.Printf("No string provided for readInput in %s %s:%d \n", runtime.FuncForPC(pc).Name(), fn, line)
		os.Exit(0)
	}

	fmt.Print(p)

	raw, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Please enter a valid string and press Enter")
		return readInput(p)
	}

	data := strings.Trim(raw, "\n")

	return data, err
}

func createMeeting() (string, []string, error) {
	meeting, err := readInput("Enter a title for meeting: ")
	if err != nil {
		return "", nil, err
	}

	rawCount, err := readInput("Enter a number of participants (from 2 to 5): ")
	if err != nil {
		return "", nil, err
	}

	pCount, err := strconv.Atoi(rawCount)

	if err != nil {
		return "", nil, err
	}

	participants := []string{}

	for i := 0; i < pCount; i++ {
		participantName, err := readInput("Enter a name of participant: ")
		if err != nil {
			return "", nil, err
		}
		participants = append(participants, participantName)
	}

	return meeting, participants, nil
}

func getCalendar(p string) ([]string, error) {
	rawCounter, err := readInput(p + ", enter a number of days you wish to put statuses for: ")
	if err != nil {
		return nil, err
	}

	dayCounter, err := strconv.Atoi(rawCounter)

	if err != nil {
		return nil, err
	}

	if dayCounter <= 10 && dayCounter > 0 {
		date := time.Now().Local()
		calendar := []string{}
		dateAvailable := false
		for i := 0; i < dayCounter; i++ {
			humanDate := date.Format("2006-01-02")
			rawStatus, err := readInput("Enter a status (e.g. 1 = free or 0 = busy) for " + humanDate + "\n")
			if err != nil {
				return nil, err
			}

			status, err := strconv.Atoi(rawStatus)

			if err != nil {
				return nil, err
			}

			if status == 1 {
				calendar = append(calendar, humanDate)
				dateAvailable = true
			}
			date = date.AddDate(0, 0, 1)
		}
		if !dateAvailable {
			fmt.Println("Please set at least one date as availabe by entering 1 as status in front of it")
			return getCalendar(p)
		}

		return calendar, nil
	} else {
		fmt.Println("Please enter a number that is bigger than 0, but less than 10")
		return getCalendar(p)
	}
}

func addBooking(c map[string][]string) (bool, string, error) {
	if len(c) == 0 {
		return false, "It seems there are no dates provided", nil
	}

	res := make(map[string][]string)

	for usr, dates := range c {
		for _, date := range dates {
			if vals, ok := res[date]; ok {
				if contains(vals, usr) == false {
					res[date] = append(res[date], usr)
				}
			} else {
				res[date] = []string{usr}
			}
		}
	}

	for k, v := range res {
		if len(v) == len(c) {
			return true, k, nil
		}
	}

	return false, "Not found any availabe date", nil

}

func DBConnect() (sess *mgo.Session, err error) {
	uri := "mongodb://meetplango:secret1@ds227243.mlab.com:27243/heroku_5n3jl0d2"

	sess, err = mgo.Dial(uri)
	if err != nil {
		return
	}

	sess.SetSafe(&mgo.Safe{})

	return
}

func DBWrite(t string, d string, u []string, db *mgo.Database) {
	meetingC := db.C("meeting")

	err := meetingC.Insert(&Meeting{t, d, u}) //{"title, date, []string{"Bogdan", "Dima"}})
	if err != nil {
		log.Fatal("Problem inserting meetingC data: ", err)
		return
	} else {
		fmt.Println("meetingC data is there")
	}
}

func DBRead(db *mgo.Database) {
	meetingC := db.C("meeting")

	var resM []Meeting
	err := meetingC.Find(nil).All(&resM)

	if err != nil {
		log.Fatal("Problem finding the meetingC data: ", err)
		return
	} else {
		fmt.Println("Results All Meeting: ", resM)
	}
}

func main() {
	conn, err := DBConnect()
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		return
	}
	db := conn.DB("")
	defer conn.Close()

	fmt.Println("Connection Established")

	m, p, err := createMeeting()
	if err != nil {
		fmt.Printf("Can't create meeting %v\n", err)
		return
	}

	uC := make(map[string][]string)

	for i := 0; i < len(p); i++ {
		calendar, err := getCalendar(p[i])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		uC[p[i]] = calendar
	}

	check, booking, err := addBooking(uC)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if check {
		DBWrite(m, booking, p, db)
		DBRead(db)

		fmt.Println("Meeting %s is scheduled sucessfully at %s", m, booking)
	} else {
		fmt.Println(booking)
	}
}
