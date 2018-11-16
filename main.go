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

type User struct {
	Name  string
	Email string
	Dates []string
}

type Meeting struct {
	Title string
	Date  string
	Users []string
}

func readInput(p string) (int, error) {
	reader := bufio.NewReader(os.Stdin)

	if p == "" {
		pc, fn, line, _ := runtime.Caller(1)
		fmt.Printf("No string provided for readInput in %s %s:%d \n", runtime.FuncForPC(pc).Name(), fn, line)
		os.Exit(0)
	}

	fmt.Print(p)

	raw, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Please enter a valid integer and press Enter")
		return readInput(p)
	}

	data := strings.Trim(raw, "\n")

	res, err := strconv.Atoi(data)
	if err != nil {
		fmt.Println("Please enter a valid integer and press Enter")
		return readInput(p)
	}

	return res, err
}

func getCalendar() ([]string, error) {
	dayCounter, err := readInput("Enter a number of days you wish to put statuses for: ")
	if err != nil {
		return nil, err
	}

	if dayCounter <= 10 && dayCounter > 0 {
		date := time.Now().Local()
		calendar := []string{}
		dateAvailable := false
		for i := 0; i < dayCounter; i++ {
			humanDate := date.Format("2006-01-02")
			status, err := readInput("Enter a status (e.g. 1 = free or 0 = busy) for " + humanDate + "\n")
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
			return getCalendar()
		}

		return calendar, nil
	} else {
		fmt.Println("Please enter a number that is bigger than 0, but less than 10")
		return getCalendar()
	}
}

func addBooking(c []string) (string, error) {
	if len(c) == 0 {
		fmt.Print("It seems there are no dates provided")
		getCalendar()
	}

	fmt.Println("Available dates for meeting: ")
	i := 1
	for _, v := range c {
		fmt.Printf("%d: %s \n", i, v)
		i++
	}

	bookingNumber, err := readInput("Select date you wish to schedule meeting at: ")

	if err != nil {
		return "", err
	}

	bookingCount := len(c)

	if bookingNumber >= 1 && bookingNumber <= bookingCount {
		bookingNumber--
	} else {
		fmt.Println("The number you have input, does not correspond to any available date, so let's try again")
		bookingNumber = 0
		bookingCount = 0
		return addBooking(c)
	}

	return c[bookingNumber], err
}

func DBConnect() (sess *mgo.Session, err error) {
	uri := "mongodb://meetplango:secret1@ds227243.mlab.com:27243/heroku_5n3jl0d2"

	sess, err = mgo.Dial(uri)
	if err != nil {
		return
	}

	// defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	return
}

func DBWrite(c []string, b string, db *mgo.Database) {
	userC := db.C("user")
	meetingC := db.C("meeting")

	err := userC.Insert(&User{"Bogdan", "lazbog@tuta.io", c})
	if err != nil {
		log.Fatal("Problem inserting userC data: ", err)
		return
	} else {
		fmt.Println("userC data is there")
	}

	err = meetingC.Insert(&Meeting{"Coding session", b, []string{"Bogdan", "Dima"}})
	if err != nil {
		log.Fatal("Problem inserting meetingC data: ", err)
		return
	} else {
		fmt.Println("meetingC data is there")
	}

}

func DBRead(db *mgo.Database) {
	userC := db.C("user")
	meetingC := db.C("meeting")

	var resU []User
	err := userC.Find(nil).All(&resU)

	if err != nil {
		log.Fatal("Problem finding the userC data: ", err)
		return
	} else {
		fmt.Println("Results All Users: ", resU)
	}

	var resM []Meeting
	err = meetingC.Find(nil).All(&resM)

	if err != nil {
		log.Fatal("Problem finding the meetingC data: ", err)
		return
	} else {
		fmt.Println("Results All Meeting: ", resM)
	}

}

func main() {
	sess, err := DBConnect()
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		return
	}
	db := sess.DB("")
	defer sess.Close()

	fmt.Println("Connection Established")

	calendar, err := getCalendar()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	result, err := addBooking(calendar)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	DBWrite(calendar, result, db)
	DBRead(db)

	fmt.Printf("Booking is scheduled sucessfully at %s \n", result)
}
