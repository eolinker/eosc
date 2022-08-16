package eocontext

import (
	"fmt"
)

type PrintFilter int

func (p PrintFilter) DoFilter(ctx EoContext, next IChain) (err error) {
	fmt.Print(">", p, "-")
	if next != nil {

		err = next.DoChain(ctx)

	}
	fmt.Print(p, "==>")
	return
}

func (p PrintFilter) Destroy() {
	fmt.Printf(":%d->", p)
}

func ExampleFiltersDoChain() {
	Count := 10

	filters := make(Filters, Count)
	for i := range filters {
		filters[i] = PrintFilter(i)
	}
	err := filters.DoChain(nil)

	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println()
	filters.Destroy()
	// output:>0->1->2->3->4->5->6->7->8->9-9==>8==>7==>6==>5==>4==>3==>2==>1==>0==>
	// :0->:1->:2->:3->:4->:5->:6->:7->:8->:9->

}
