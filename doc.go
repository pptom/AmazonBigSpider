// a go spider!
package dudu

import (
	"github.com/hunterhug/GoSpider/util"
	"fmt"
)

var Dir = util.CurDir()
var Local = true

func init() {
	if util.FileExist(Dir + "/远程.txt") {
		Local = false
		fmt.Println("远程方式！！！")
	}
}
