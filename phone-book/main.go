package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

const (
	filenameDefault = "phonebook.json"
	filenameEnv     = "PHONEBOOK_PATH"
)

type (
	Person struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	Records struct {
		Persons []Person
	}
)

func (x *Records) AddPerson(person Person) {
	x.Persons = append(x.Persons, person)
}

func (x *Records) DeletePerson(name string) {
	for i, person := range x.Persons {
		if person.Name == name {
			x.Persons = append(x.Persons[:i], x.Persons[i+1:]...)
			break
		}
	}
}

func (x *Records) GetPersons(name string) []Person {
	if name == "" {
		return x.Persons
	}

	persons := []Person{}

	for _, person := range x.Persons {
		if person.Name == name {
			persons = append(persons, person)
		}
	}

	return persons
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <subcommand> <options>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Subcommands:\n")
	fmt.Fprintf(os.Stderr, "  add - Add a new phone book record\n")
	fmt.Fprintf(os.Stderr, "     Options:\n")
	fmt.Fprintf(os.Stderr, "       -name - Name of the person\n")
	fmt.Fprintf(os.Stderr, "       -phone - Phone number of the person\n")
	fmt.Fprintf(os.Stderr, "  delete - Delete an existing phone book record\n")
	fmt.Fprintf(os.Stderr, "     Options:\n")
	fmt.Fprintf(os.Stderr, "       -name - Name of the person\n")
	fmt.Fprintf(os.Stderr, "  get - Get one or all phone book records\n")
	fmt.Fprintf(os.Stderr, "     Options:\n")
	fmt.Fprintf(os.Stderr, "       -name - Name of the person\n")
	fmt.Fprintf(os.Stderr, "Global options:\n")

	flag.PrintDefaults()
}

func main() {
	var name string
	var phone string
	var filename string

	flag.StringVar(&filename, "filename", filenameDefault, "Specify a path where file will be stored")

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addCmd.StringVar(&name, "name", "", "Name of the person")
	addCmd.StringVar(&phone, "phone", "", "Phone number of the person")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteCmd.StringVar(&name, "name", "", "Name of the person")

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getCmd.StringVar(&name, "name", "", "Name of the person")

	if len(os.Args) < 2 {
		fmt.Println("Expected a subcommand -- add or delete or get")
		usage()
		os.Exit(1)
	}

	flag.Parse()

	if envvar := os.Getenv(filenameEnv); envvar != "" {
		fmt.Printf("Using phone book path %s from environment variable %s\n", envvar, filenameEnv)
		flag.Set("filename", envvar)
	}

	persons, err := os.ReadFile(filename)

	if err != nil {
		if os.IsNotExist(err) {
			persons = []byte("[]")
			os.WriteFile(filename, persons, 0644)
		} else {
			panic(err)
		}
	}

	var records Records
	json.Unmarshal(persons, &records.Persons)

	switch flag.Arg(0) {
	case "add":
		addCmd.Parse(flag.Args()[1:])

		if addCmd.NFlag() < 2 {
			fmt.Println("You must supply both --name and --phone options")
			addCmd.Usage()
			os.Exit(1)
		}

		person := Person{name, phone}
		records.AddPerson(person)

	case "delete":
		deleteCmd.Parse(flag.Args()[1:])

		if deleteCmd.NFlag() < 1 {
			fmt.Println("You must supply both --name option")
			deleteCmd.Usage()
			os.Exit(1)
		}

		records.DeletePerson(name)

	case "get":
		getCmd.Parse(flag.Args()[1:])

		if getCmd.NFlag() < 1 {
			fmt.Println("You must supply both --name option")
			getCmd.Usage()
			os.Exit(1)
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
		fmt.Fprintln(writer, "Name\tPhone\t")

		for _, person := range records.GetPersons(name) {
			fmt.Fprintln(writer, person.Name, "\t", person.Phone, "\t")
		}

		writer.Flush()

	default:
		fmt.Println("Invalid subcommand specified")
		usage()
		os.Exit(1)
	}

	persons, _ = json.MarshalIndent(records.Persons, "", " ")
	err = os.WriteFile(filename, persons, 0644)

	if err != nil {
		fmt.Printf("Error writing to file %s\n", filename)
		os.Exit(1)
	}
}
