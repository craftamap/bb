package bb

import (
	"fmt"
	"os"

	"github.com/craftamap/bb/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
