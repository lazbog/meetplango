package main

import (
	"bufio"
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

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

func getCalendar() (map[string]int, error) {
	dayCounter, err := readInput("Enter a number of days you wish to put statuses for: ")
	if err != nil {
		return nil, err
	}

	if dayCounter <= 10 && dayCounter > 0 {
		date := time.Now().Local()
		calendar := make(map[string]int)
		dateAvailable := false
		for i := 0; i < dayCounter; i++ {
			humanDate := date.Format("2006-01-02")
			status, err := readInput("Enter a status (e.g. 1 = free or 0 = busy) for " + humanDate + "\n")
			if err != nil {
				return nil, err
			}
			calendar[humanDate] = status
			if status == 1 {
				dateAvailable = true
			}
			date = date.AddDate(0, 0, 1)
		}
		if !dateAvailable {
			fmt.Println("Please set at least one date as availabe by entering 1 as status in front of it")
			return getCalendar()
		}

		return calendar, err
	} else {
		fmt.Println("Please enter a number that is bigger than 0, but less than 10")
		return getCalendar()
	}
}

func addBooking(c map[string]int) (string, error) {
	if len(c) == 0 {
		fmt.Print("It seems there are no dates provided")
		getCalendar()
	}
	booking := []string{}
	fmt.Println("Available dates for meeting: ")
	i := 1
	for k, v := range c {
		if v == 1 {
			booking = append(booking, k)
			fmt.Printf("%d: %s \n", i, k)
			i++
		}
	}

	bookingCount := len(booking)

	if bookingCount == 1 {
		return booking[0], nil
	}

	bookingNumber, err := readInput("Select date you wish to schedule meeting at: ")

	if err != nil {
		return "", err
	}

	if bookingNumber >= 1 && bookingNumber <= bookingCount {
		bookingNumber--
	} else {
		fmt.Println("The number you have input, does not correspond to any available date, so let's try again")
		bookingNumber = 0
		bookingCount = 0
		return addBooking(c)
	}

	return booking[bookingNumber], err
}

func main() {
	url := "ds227243.mlab.com:27243/heroku_5n3jl0d2";
	mgo.Dial(url)
	
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

	fmt.Printf("Booking is scheduled sucessfully at %s \n", result)
}
