package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

/*
	Functions that need to be supported:
		- help : displays a list of commands that the user can interact with, as well as the idea of time and location in this app

		- location : prints the current location value that weth will report data for. Default is current location server.

		settime functions.
			These functions will mutate the current internal time weth reports weather data for. If the time value is successfully changed,
			it will be printed.

			Otherwise, an error message will be printed.

		- settime --military=true
				prints hourly information using military time if set to true,
				uses AM/PM mode if false.

		- now : displays very specific weather data in the current location, at the current time.

		- hours <NUMBER> : displays specific weather daya for the next <NUMBER> of hours at and after *TIME*, at *LOCATION*
			"hours 0" or "hours 1" is equivalent to typing 'now'

		- days <NUMBER> : displays general weather data for the next <NUMBER> days at and after *TIME*, in *LOCATION*
		-
*/

var internalTime time.Time
var validMonthCodes = map[string]int{"january": 1, "february": 2, "march": 3, "april": 4, "may": 5, "june": 6, "july": 7, "august": 8, "september": 9, "october": 10, "november": 11, "december": 12, "jan": 1, "feb": 2, "mar": 3, "apr": 4, "jun": 6, "jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12}

func setTime(args []string) string {

	if len(args) == 0 {
		internalTime = time.Now()
		return "  set time to " + internalTime.Format(time.DateOnly) + " Hour: " + strconv.Itoa(internalTime.Hour())
	}

	var stateValues = map[string]int{"Hour": internalTime.Hour(), "Day": internalTime.Day(), "Month": int(internalTime.Month()), "Year": internalTime.Year()}
	var stateNames = [...]string{"Hour", "Day", "Month", "Year"}

	var bound = min(len(stateNames), len(args))

	for i := 0; i < bound; i++ {

		if args[i] == "*" {
			continue
		}

		runeValue, width := utf8.DecodeRuneInString(args[i])

		if runeValue == '\\' {

			relNum, error := strconv.Atoi(args[i][width:])
			if error != nil {
				return "  Error: Expected a number for " + stateNames[i] + ", got " + args[i][width:]
			}
			stateValues[stateNames[i]] += relNum
			continue
		}

		if stateNames[i] != "Month" {
			relNum, error := strconv.Atoi(args[i])
			if error != nil {
				return "  Error: Expected a number for " + stateNames[i] + ", got " + args[i]
			}
			stateValues[stateNames[i]] += relNum

		} else {

			monthNum, error := strconv.Atoi(args[i])
			if error == nil {
				if monthNum < 1 || monthNum > 12 {
					return "  Error: Expected Month number in range 1-12, got " + strconv.Itoa(monthNum)
				}
				stateValues[stateNames[i]] = monthNum
			}

			copy := strings.ToLower(args[i])

			if validMonthCodes[copy] != 0 {
				stateValues[stateNames[i]] = validMonthCodes[copy]
			}

			return "  Error: Expected a valid month code. Got " + args[i]
		}

	}

	// TODO: when we add support for locations, we need this last parameter to be the timezone associated with the
	// current standing location.
	internalTime = time.Date(stateValues["Year"], time.Month(stateValues["Month"]), stateValues["Day"], stateValues["Hour"], 0, 0, 0, time.Local)
	return "  set time to " + internalTime.Format(time.DateOnly) + " Hour: " + strconv.Itoa(internalTime.Hour())
}

func getTime([]string) string {
	return internalTime.Format(time.DateOnly) + " Hour: " + strconv.Itoa(internalTime.Hour())
}

func main() {

	fmt.Println("Welcome to the weth REPL! Type 'help' to print a list of commands")

	reader := bufio.NewReader(os.Stdin)

	var command2func = make(map[string]func([]string) string)
	internalTime = time.Now()

	command2func["settime"] = setTime
	command2func["time"] = getTime

	for {

		fmt.Print("-> ")

		line, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		line = strings.Trim(line, " \n")

		//if utf8.RuneCountInString(line) == 0 {
		//  continue;
		//}

		arguments := strings.Split(line, " ")

		if len(arguments) == 0 {
			continue
		}

		if command2func[arguments[0]] == nil {
			fmt.Printf("  %s: command not found\n", arguments[0])
			continue
		}

		output := command2func[arguments[0]](arguments[1:])
		fmt.Printf("  %s\n", output)
	}

}
