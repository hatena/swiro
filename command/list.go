package command

import (
	"github.com/urfave/cli"
	"fmt"
	"os"
	//"github.com/taku-k/swiro/aws"
)

func CmdList(c *cli.Context) error {
	fmt.Fprintln(os.Stderr, "Not implemented")
	fmt.Println(c.String("force"))
	//ts, err := aws.NewRouteTables()
	//if err != nil {
	//	return err
	//}
	//for _, t := range ts {
	//	fmt.Println(t)
	//}
	return nil
}
