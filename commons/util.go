package commons

import (
	"fmt"
	"strings"
)

func CheckErr(err error) {

	if err != nil {
		fmt.Println(err)
		//log.Fatal(err)
		//os.Exit(1)
	}
}

func CreateCondStr(formCondVal []string, columnCondVal []string) string {
	conditionStr := " WHERE"
	i := 0
	for i = 0; i < len(columnCondVal); i++ {
		conditionStr = conditionStr + " and " + columnCondVal[i] + " = '" + formCondVal[i] + "'"
	}
	conditionStr = strings.Replace(conditionStr, "WHERE and", " WHERE", -2)
	return conditionStr
}
