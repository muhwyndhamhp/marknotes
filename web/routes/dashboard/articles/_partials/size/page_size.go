package size

import (
	"fmt"
	pubvariables "github.com/muhwyndhamhp/marknotes/pub/variables"
)

func Dropdown(page, pageSize int) pubvariables.DropdownVM {
	var arrays []pubvariables.DropdownItem
	item := 0
	for i := range []int{0, 1, 2} {
		size := (i + 1) * 10
		arrays = append(arrays, pubvariables.DropdownItem{
			Label:  fmt.Sprintf("%d", size),
			Path:   fmt.Sprintf("/dashboard/articles?page=%d&pageSize=%d&source=source-partial", page, size),
			Target: "#articles",
		})
		if size == pageSize {
			item = i
		}
	}
	return pubvariables.DropdownVM{
		Items:    arrays,
		Selected: item,
	}
}
