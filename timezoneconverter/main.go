package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

type TimeInfo struct {
	Location string
	Time     string
	Offset   string
}

const timeLayout = "2006-01-02 15:04"

func getOffset(targetTime time.Time, timezoneTime time.Time) string {
	_, timezone1 := targetTime.Zone()
	_, timezone2 := timezoneTime.Zone()
	hours := float64(timezone1-timezone2) / 3600.0
	minutes := (hours - float64(int(hours))) * 60.0

	if hours < 0 {
		return fmt.Sprintf("%d hours %d minutes ahead", -int(hours), -int(minutes))
	} else {
		return fmt.Sprintf("%d hours %d minutes behind", int(hours), int(minutes))
	}
}

func convertTime(targetTime time.Time, timezones []string) error {
	converted := []TimeInfo{}
	converted = append(converted, TimeInfo{"Local", targetTime.Format(timeLayout), "0 hours 0 minutes"})

	for _, timezone := range timezones {
		location, err := time.LoadLocation(timezone)

		if err != nil {
			return fmt.Errorf("%s is not a valid timezone", timezone)
		}

		timezoneTime := targetTime.In(location)
		offset := getOffset(targetTime, timezoneTime)
		converted = append(converted, TimeInfo{timezone, timezoneTime.Format(timeLayout), offset})
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(writer, "Location\tTime\tOffset\t")

	for _, timeInfo := range converted {
		fmt.Fprintln(writer, timeInfo.Location, "\t", timeInfo.Time, "\t", timeInfo.Offset, "\t")
	}

	writer.Flush()
	return nil
}

func main() {
	targetTime := time.Now()
	timezones := []string{}

	flag.Func("targetTime", fmt.Sprintf("Date-time in %q format", timeLayout), func(s string) error {
		location, err := time.LoadLocation("Local")

		if err != nil {
			return fmt.Errorf("failed parsing local time")
		}

		parsed, err := time.ParseInLocation(timeLayout, s, location)

		if err != nil {
			return fmt.Errorf("failed parsing time at timezone %s", s)
		}

		if parsed.Weekday() == time.Saturday || parsed.Weekday() == time.Sunday {
			return fmt.Errorf("schedule meetings on a weekday")
		}

		targetTime = parsed
		return nil
	})

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [-targetTime] <timezones...>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("You must specify at least one timezone string")
		flag.Usage()
		os.Exit(1)
	}

	for _, timezone := range flag.Args() {
		_, err := time.LoadLocation(timezone)

		if err != nil {
			fmt.Printf("%s is not a valid timezone\n", timezone)
			os.Exit(1)
		}

		timezones = append(timezones, timezone)
	}

	convertTime(targetTime, timezones)
}
