package commands

import "fmt"

const reset = "\033[0m"
const redC = "\033[31m"
const greenC = "\033[32m"
const yellowC = "\033[33m"
const blueC = "\033[34m"
const purpleC = "\033[35m"
const cyanC = "\033[36m"
const grayC = "\033[37m"
const whiteC = "\033[97m"

func cyan(str string) string {
	return fmt.Sprintf("%s%s%s", cyanC, str, reset)
}
