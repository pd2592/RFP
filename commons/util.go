package commons

import "fmt"

func CheckErr(err error) {

	if err != nil {
		fmt.Println(err)
		//log.Fatal(err)
		//os.Exit(1)
	}
}
