package main

import (
	"a-library-for-others/csvparser"
	"fmt"
	"io"
	"os"
)

func main() {
	content := `Name,Age,Location
John Doe,30,New York
"Jane, Smith",25,"Los Angeles"
Alice Cooper,35,San Francisco
Bob O'Conner,40,Chicago
"Sam O'Neil",28,"Austin"`
	CreateFile(content)
	file, err := os.Open("output.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	csvParser := csvparser.CSVParserData{}
	idx := 0
	for {
		line, err := csvParser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading line:", err)
			os.Exit(1)
		}
        fmt.Println("Line", idx, ":", line)
		idx++
		numberOfFields := csvParser.GetNumberOfFields()
		for i := 0; i < numberOfFields; i++ {
			field, err := csvParser.GetField(i)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("    Field", i+1, ":", field)
		}
	}
}

func CreateFile(content string) {
	file, err := os.Create("output.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
