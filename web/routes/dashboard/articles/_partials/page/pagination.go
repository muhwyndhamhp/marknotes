package page

import (
	"fmt"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
)

func PageDropdown(page, pageSize, count int) pub_variables.DropdownVM {
	arrays := []pub_variables.DropdownItem{}
	item := 0
	for i := 0; (i)*pageSize <= count; i++ {
		currentPage := i + 1
		size := pageSize
		arrays = append(arrays, pub_variables.DropdownItem{
			Label:  fmt.Sprintf("%d", currentPage),
			Path:   fmt.Sprintf("/dashboard/articles?page=%d&pageSize=%d&source=source-partial", currentPage, size),
			Target: "#articles",
		})
		if currentPage == page {
			item = i
		}
	}

	return pub_variables.DropdownVM{
		Items:    arrays,
		Selected: item,
	}
}
