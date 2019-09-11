package main

import (
	"fmt"
	"path/filepath"
	"strings"

	arg "github.com/alexflint/go-arg"
	"github.com/mitchellh/go-homedir"
	"github.com/nosarthur/todo/db"
)

// AddCmd is a subcommand
type AddCmd struct {
	Task []string `arg:"positional, required"`
}

// ListCmd lists all todos
type ListCmd struct {
}

// RmCmd removes a todo
type RmCmd struct {
	Numbers []int `arg:"positional, required"`
}

type args struct {
	Add  *AddCmd  `arg:"subcommand:add" help:"Add a task"`
	List *ListCmd `arg:"subcommand:list"`
	Rm   *RmCmd   `arg:"subcommand:rm"`
}

func (args) Version() string {
	return "todo 0.1.0"
}
func main() {
	d, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	dbPath := filepath.Join(d, "my_task.db")
	db.MustInit(dbPath)

	var args args
	arg.MustParse(&args)

	switch {
	case args.Add != nil:
		if err := db.Add(strings.Join(args.Add.Task, " ")); err != nil {
			fmt.Printf("%s", err)
			panic("Fail to add task!")
		}
	case args.Rm != nil:
		nums := map[int]struct{}{}
		for _, x := range args.Rm.Numbers {
			nums[x] = struct{}{}
		}
		db.Rm(nums)
	case args.List != nil:
		fallthrough
	default:
		if err := db.List(); err != nil {
			fmt.Printf("%s", err)
			panic("cannot list")
		}
	}

}
