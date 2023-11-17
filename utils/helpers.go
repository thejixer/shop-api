package utils

import "fmt"

func Typeof(a any) {
	fmt.Printf("the variable %+v is of type of %T \n", a, a)
}
