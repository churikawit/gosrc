package excel

import "testing"

func TestGetSheetName(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"SheetName1", "SheetName1"},
		{"世界", "世界"},
		{"สวัสดี", "สวัสดี"},
	}

	ec := New()
	filename := "test.xlsx"
	ec.OpenOrCreate(filename)

	for _, c := range cases {
		ec.SetSheetName(c.in, 1)
		got := ec.GetSheetName(1)
		if got != c.want {
			t.Errorf("GetSheetName(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestSetActiveSheet(t *testing.T) {
	ec := New()
	filename := "test.xlsx"
	ec.OpenOrCreate(filename)

	want := 10
	in := 10
	ec.SetActiveSheet(in)
	got := ec.CountSheet()
	if got != want {
		t.Errorf("SetActiveSheet(%d) == %d, want %d", in, got, want)
	}
}

func TestRowCol(t *testing.T) {
	cases := []struct {
		in_row int
		in_col int
		want   string
	}{
		{1, 1, "A1"},
		{10, 10, "J10"},
		{2, 5, "E2"},
		{2, 16384, "XFD2"},
		{1048576, 16384, "XFD1048576"},
	}

	for _, c := range cases {
		got := RowCol(c.in_row, c.in_col)
		if got != c.want {
			t.Errorf("RowCol(%d,%d) == %q, want %q", c.in_row, c.in_col, got, c.want)
		}
	}
}
