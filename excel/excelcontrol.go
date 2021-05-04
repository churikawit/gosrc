package excel

import (
	"fmt"
	_ "fmt"
	_ "log"
	_ "math/rand"
	"os"
	"runtime"
	"strings"
	"time"
	_ "time"

	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/spreadsheet"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

// -----------------------------------------------------------------------------------------------------------

type ExcelControl struct {
	m_ExcelApp     int
	m_OpenFile     bool
	m_OpenFilename string
	m_Workbook     *spreadsheet.Workbook
	m_Worksheet    *spreadsheet.Sheet
}

func New() *ExcelControl {
	output := &ExcelControl{m_ExcelApp: 1, m_OpenFile: false, m_OpenFilename: ""}
	return output
}

// RowCol : row = [1,n], column = [1,256]
func RowCol(row, column int) (output string) {
	output = ""
	if column <= 0 {
		panic(fmt.Sprintf("RowCol : invalid column:%d\n", column))
	}

	var o, r int
	var temp int = column - 1
	var c byte = byte('A')

	for temp >= 0 {
		o = temp / 26
		r = temp % 26

		output = string(c+byte(r)) + output
		temp = o - 1
	}
	output = output + fmt.Sprintf("%d", row)
	return
}

func (ec *ExcelControl) OpenOrCreate(filename string) {
	if IsFileExists(filename) {
		ec.open(filename)
	} else {
		CreateDirectory(filename)
		ec.create(filename)
	}
}

func (ec *ExcelControl) create(filename string) {
	ec.m_OpenFilename = filename
	ec.m_Workbook = spreadsheet.New()
	sheet := ec.m_Workbook.AddSheet()
	ec.m_Worksheet = &sheet
	ec.m_Worksheet.Cell("A1").SetString("") // create empty cell
}

func (ec *ExcelControl) open(filename string) {
	ec.m_OpenFilename = filename
	var err error
	ec.m_Workbook, err = spreadsheet.Open(filename)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}
	ec.m_Worksheet = ec.GetSheet(1)
}

func (ec *ExcelControl) SaveAs(filename string) {
	ec.m_Workbook.SaveToFile(filename)
}

func (ec *ExcelControl) CountSheet() int {
	return ec.m_Workbook.SheetCount()
}

// RemoveSheet : sheet_no [1,n], can remove only last sheet
func (ec *ExcelControl) RemoveSheet(sheet_no int) {
	if 0 < sheet_no && sheet_no <= ec.CountSheet() {
		ec.m_Workbook.RemoveSheet(sheet_no - 1)
	}
}

func (ec *ExcelControl) GetSheet(sheet_no int) *spreadsheet.Sheet {
	if 0 < sheet_no && sheet_no <= ec.CountSheet() {
		return &ec.m_Workbook.Sheets()[sheet_no-1]
	}
	return nil
}

func (ec *ExcelControl) GetSheetName(sheet_no int) string {
	ss := ec.GetSheet(sheet_no)
	if ss != nil {
		return ss.Name()
	}
	return ""
}

// SetSheetName: sheet_no = [1,n], sheetname must differ from the other sheets
func (ec *ExcelControl) SetSheetName(sheetname string, sheet_no int) {
	ss := ec.GetSheet(sheet_no)
	if ss != nil {
		ss.SetName(sheetname)
	}
}

// SetActiveSheet: sheet_no = [1,n]
func (ec *ExcelControl) SetActiveSheet(sheet_no int) {
	var count int = ec.CountSheet()
	if count <= 0 {
		return
	}
	if sheet_no <= 0 {
		sheet_no = 1
	}
	if sheet_no > count {
		for j := count; j < sheet_no; j++ {
			ec.m_Workbook.AddSheet()
		}
	}

	ec.m_Worksheet = &(ec.m_Workbook.Sheets()[sheet_no-1])
}

func (ec *ExcelControl) SetPassword() {

}

// =================================== Data ==========================================

// WriteHeader: first_row = [1,n], first_col = [1,n]
func (ec *ExcelControl) WriteHeader(header []string, first_row int, first_col int) {
	l := len(header)

	if first_row <= 0 || first_col <= 0 {
		return
	}
	for i := 0; i < l; i++ {
		cellID := RowCol(first_row, first_col+i)
		ec.m_Worksheet.Cell(cellID).SetString(header[i])
	}
}

// WriteRow: first_row = [1,n], first_col = [1,n]
func (ec *ExcelControl) WriteRow(datarow []interface{}, first_row int, first_col int) {
	len := len(datarow)

	if first_row <= 0 || first_col <= 0 {
		return
	}
	for i := 0; i < len; i++ {
		cellID := RowCol(first_row, first_col+i)
		switch s2 := (datarow[i]).(type) {
		case bool:
			ec.m_Worksheet.Cell(cellID).SetBool(s2)
			// fmt.Printf("%d:bool\n", i+1)
		case time.Time:
			ec.m_Worksheet.Cell(cellID).SetTime(s2)
			// fmt.Printf("%d:time\n", i+1)
		case string:
			ec.m_Worksheet.Cell(cellID).SetString(s2)
			// fmt.Printf("%d:string\n", i+1)
		case int32:
			ec.m_Worksheet.Cell(cellID).SetNumber(float64(s2))
			// fmt.Printf("%d:int32\n", i+1)
		case int64:
			ec.m_Worksheet.Cell(cellID).SetNumber(float64(s2))
			// fmt.Printf("%d:int64\n", i+1)
		case float64:
			ec.m_Worksheet.Cell(cellID).SetNumber(s2)
			// fmt.Printf("%d:float64\n", i+1)
		default:
			ec.m_Worksheet.Cell(cellID).SetString("")
		}

	}
}

func (ec *ExcelControl) WriteString(txt string, first_row int, first_col int) {
	if strings.Index(txt, "=") >= 0 {
		ec.WriteFormula(txt, first_row, first_col)
	} else {
		cellID := RowCol(first_row, first_col)
		ec.m_Worksheet.Cell(cellID).SetString(txt)
	}
}
func (ec *ExcelControl) WriteFormula(txt string, first_row, first_col int) {
	index := strings.LastIndex(txt, "=")
	if index >= 0 {
		txt = txt[index+1:]
	}
	cellID := RowCol(first_row, first_col)
	ec.m_Worksheet.Cell(cellID).SetFormulaRaw(txt)
}

// =================================== Util ==========================================

// CreateDirectory : return false if error occur
func CreateDirectory(file_path string) bool {
	var pathSeparator string = "/"
	if runtime.GOOS == "windows" {
		pathSeparator = "\\"
	}
	index := strings.LastIndex(file_path, pathSeparator)

	if index != -1 {
		var dir_path string = file_path[0:index]
		err := os.MkdirAll(dir_path, 0755)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func IsFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
