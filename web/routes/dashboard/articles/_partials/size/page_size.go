package size

import (
	"fmt"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
)

func SizeDropdown(page, pageSize int) pub_variables.DropdownVM {
	arrays := []pub_variables.DropdownItem{}
	item := 0
	for i := range []int{0, 1, 2} {
		size := (i + 1) * 10
		arrays = append(arrays, pub_variables.DropdownItem{
			Label:  fmt.Sprintf("%d", size),
			Path:   fmt.Sprintf("/dashboard/articles?page=%d&pageSize=%d&source=source-partial", page, size),
			Target: "#articles",
		})
		if size == pageSize {
			item = i
		}
	}
	return pub_variables.DropdownVM{
		Items:    arrays,
		Selected: item,
	}
}
