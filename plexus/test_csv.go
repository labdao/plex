package main
 
import (
    "encoding/csv"
    "fmt"
    "os"
)
 
func main() {
    f, e := os.Create("./People.csv")
    if e != nil {
        fmt.Println(e)
    }
 
    writer := csv.NewWriter(f)
    var data = [][]string{
        {"Name", "Age", "Occupation"},
        {"Sally", "22", "Nurse"},
        {"Joe", "43", "Sportsman"},
        {"Louis", "39", "Author"},
    }
	fmt.Println(data)
 
    e = writer.WriteAll(data)
    if e != nil {
        fmt.Println(e)
    }
}