package main

import (
	"fmt"
	"net/http"

	"github.com/xuri/excelize/v2"
)

func main() {

}

type WorkbookJson struct {
	Sheets []SheetJson `json:"sheets"`
}

type SheetJson struct {
	Name  string     `json:"name"`
	Cells []CellJson `json:"cells"`
}

type CellJson struct {
	Index string `json:"index"`
	Value string `json:"value"`
	Color string `json:"color"`
}

func writeToCellHandler(w http.ResponseWriter, r *http.Request) {
	// Get the file from the request
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Open the file
	f, err := excelize.OpenReader(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the value from the cell
	cell := f.GetCellValue("Sheet1", "B2")
	fmt.Println(cell)

	// Set the value of the cell
	f.SetCellValue("Sheet1", "B2", 200)
	// Save the file
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
