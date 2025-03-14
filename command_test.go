package ishell_test

import (
	"fmt"
	"testing"
	"github.com/ryupatterson/ishell"
	"github.com/stretchr/testify/assert"
)

func newCmd(name string, help string) *ishell.Cmd {
	return &ishell.Cmd{
		Name: name,
		Help: help,
	}
}

func TestAddCommand(t *testing.T) {
	cmd := newCmd("root", "")
	assert.Equal(t, len(cmd.Children()), 0, "should be empty")
	cmd.AddCmd(newCmd("child", ""))
	assert.Equal(t, len(cmd.Children()), 1, "should include one child command")
}

func TestDeleteCommand(t *testing.T) {
	cmd := newCmd("root", "")
	cmd.AddCmd(newCmd("child", ""))
	assert.Equal(t, len(cmd.Children()), 1, "should include one child command")
	cmd.DeleteCmd("child")
	assert.Equal(t, len(cmd.Children()), 0, "should be empty")
}

func TestFindCmd(t *testing.T) {
	cmd := newCmd("root", "")
	cmd.AddCmd(newCmd("child1", ""))
	cmd.AddCmd(newCmd("child2", ""))
	res, err := cmd.FindCmd([]string{"child1"})
	if err != nil {
		t.Fatal("finding should work")
	}
	assert.Equal(t, res.Name, "child1")

	res, err = cmd.FindCmd([]string{"child2"})
	if err != nil {
		t.Fatal("finding should work")
	}
	assert.Equal(t, res.Name, "child2")

	res, err = cmd.FindCmd([]string{"child3"})
	if err == nil {
		t.Fatal("should not find this child!")
	}
	assert.Nil(t, res)
}

func TestFindAlias(t *testing.T) {
	cmd := newCmd("root", "")
	subcmd := newCmd("child1", "")
	subcmd.Aliases = []string{"alias1", "alias2"}
	cmd.AddCmd(subcmd)

	res, err := cmd.FindCmd([]string{"alias1"})
	if err != nil {
		t.Fatal("finding alias should work")
	}
	assert.Equal(t, res.Name, "child1")

	res, err = cmd.FindCmd([]string{"alias2"})
	if err != nil {
		t.Fatal("finding alias should work")
	}
	assert.Equal(t, res.Name, "child1")

	res, err = cmd.FindCmd([]string{"alias3"})
	if err == nil {
		t.Fatal("should not find this child!")
	}
	assert.Nil(t, res)
}

func TestHelpText(t *testing.T) {
	cmd := newCmd("root", "help for root command")
	cmd.AddCmd(newCmd("child1", "help for child1 command"))
	cmd.AddCmd(newCmd("child2", "help for child2 command"))
	res := cmd.HelpText()
	expected := "\nhelp for root command\n\nCommands:\n  child1      help for child1 command\n  child2      help for child2 command\n\n"
	assert.Equal(t, res, expected)
}

func TestChildrenSortedAlphabetically(t *testing.T) {
	cmd := newCmd("root", "help for root command")
	cmd.AddCmd(newCmd("child2", "help for child1 command"))
	cmd.AddCmd(newCmd("child1", "help for child2 command"))
	children := cmd.Children()
	assert.Equal(t, children[0].Name, "child1", "must be first")
	assert.Equal(t, children[1].Name, "child2", "must be second")
}

// test creation of new cmd args
func TestCmdArgs(t *testing.T) {
	_, err := ishell.NewCmdArg("-x", "--test-1", ishell.IntType, false, true)
	assert.NoError(t, err, "Error must be nil")

	// test flag param
	_, err = ishell.NewCmdArg(".", "--test_2", ishell.BoolType, false, false)
	assert.Error(t, err, "Flag parameter must error")

	// test longflag
	_, err = ishell.NewCmdArg("-x", "-test2", ishell.BoolType, false, false)
	assert.Error(t, err, "Longflag first char - test must err")

	_, err = ishell.NewCmdArg("-x", "--test.2", ishell.BoolType, false, false)
	assert.Error(t, err, "Longflag illegal char, test must err")

	// test typ param
	_, err = ishell.NewCmdArg("-x", "--test_3", 3, false, false)
	assert.Error(t, err, "Illegal typ value, test must err")

	// test positional
	_, err = ishell.NewCmdArg("", "test_3", ishell.StringType, false, false)
	if !assert.NoError(t, err, "Valid Positional argument") {
		fmt.Println(err)
	}
}

func TestCmdArgsParsing(t *testing.T) {
	arg1_flag := "-x"
	arg1_key := "--test1"
	arg1_type := ishell.IntType
	arg1, _ := ishell.NewCmdArg(arg1_flag, arg1_key, arg1_type, false, true)

	arg2_flag := "-y"
	arg2_key := "--test2"
	arg2_type := ishell.BoolType
	arg2, _ := ishell.NewCmdArg(arg2_flag, arg2_key, arg2_type, false, false)

	arg3_flag := "-z"
	arg3_key := "--test3"
	arg3_type := ishell.StringType
	arg3, _ := ishell.NewCmdArg(arg3_flag, arg3_key, arg3_type, true, false)

	cmd := ishell.Cmd{
		Name: "root",
		Help: "root help",
		Func: nil,
	}
	cmd.AddCmdArg(arg1)
	cmd.AddCmdArg(arg2)
	cmd.AddCmdArg(arg3)

	// test 1
	// basic case
	ex1 := []string{"-x", "1"}
	parsed1, err := cmd.ParseArgs(ex1)
	if assert.NoError(t, err, "CmdArgsParsing:Test1 should not error") {
		assert.Equal(t, 1, len(parsed1), "CmdArgsParsing:Test1 should have 1 argument")
		// check param x
		assert.Equal(t, 0, parsed1[0].Index, fmt.Sprintf("CmdArgsParsing:Test1 Key %d != %d", 0, parsed1[0].Index))
		assert.Equal(t, arg1_type, parsed1[0].Typ, fmt.Sprintf("CmdArgsParsing:Test1 Typ %d != %d", arg1_type, parsed1[0].Typ))
		assert.Equal(t, "1", parsed1[0].Value, fmt.Sprintf("CmdArgsParsing:Test1 Value %s != %s", "1", parsed1[0].Value))
	}


	// test 2
	// checking case where there is a boolean flag that is combined with a flag that takes a value, i.e. -yz
	ex2 := []string{"-x", "1", "-yz", "test"}
	parsed2, err := cmd.ParseArgs(ex2)
	if assert.NoError(t, err, "CmdArgsParsing:Test2 should not error") {
		assert.Equal(t, 3, len(parsed2), "CmdArgsParsing:Test2 should have 3 arguments")
		// check param x
		assert.Equal(t, 0, parsed2[0].Index, fmt.Sprintf("CmdArgsParsing:Test2 Key %d != %d", 0, parsed2[0].Index))
		assert.Equal(t, arg1_type, parsed2[0].Typ, fmt.Sprintf("CmdArgsParsing:Test2 Typ %d != %d", arg1_type, parsed2[0].Typ))
		assert.Equal(t, "1", parsed2[0].Value, fmt.Sprintf("CmdArgsParsing:Test2 Value %s != %s", "1", parsed2[0].Value))
		// check param y
		assert.Equal(t, 1, parsed2[1].Index, fmt.Sprintf("CmdArgsParsing:Test2 Key %d != %d", 1, parsed2[1].Index))
		assert.Equal(t, arg2_type, parsed2[1].Typ, fmt.Sprintf("CmdArgsParsing:Test2 Typ %d != %d", arg2_type, parsed2[1].Typ))
		// check param z
		assert.Equal(t, 2, parsed2[2].Index, fmt.Sprintf("CmdArgsParsing:Test2 Key %d != %d", 2, parsed2[2].Index))
		assert.Equal(t, arg3_type, parsed2[2].Typ, fmt.Sprintf("CmdArgsParsing:Test2 Typ %d != %d", arg3_type, parsed2[2].Typ))
		assert.Equal(t, "test", parsed2[2].Value, fmt.Sprintf("CmdArgsParsing:Test2 Value %s != %s", "test", parsed2[2].Value))
	}

	// test 3
	// checking case where a required arg is missing
	ex3 := []string{"-yz", "test"}
	_, err = cmd.ParseArgs(ex3)
	assert.Error(t, err, "CmdArgsParsing:Test3 should error due to a missing required arg")

	// test 4
	// checking case where param z is missing its value
	ex4 := []string{"-x", "1", "-yz"}
	_, err = cmd.ParseArgs(ex4)
	assert.Error(t, err, "Process should error due to a missing required arg")
}

func TestPositionalCmdArgsParsing(t *testing.T) {
	arg1_type := ishell.StringType
	arg2_type := ishell.StringType
	arg3_type := ishell.IntType

	arg1_key := "test1"
	arg2_key := "test2"
	arg3_key := "--test3"
	// positional argument
	arg1, err := ishell.NewCmdArg("", arg1_key, arg1_type, false, true)
	assert.NoError(t, err, "PositionalCmdArgsParsing:Arg1 arg creation should not error")
	arg2, err := ishell.NewCmdArg("", arg2_key, arg2_type, false, false)
	assert.NoError(t, err, "PositionalCmdArgsParsing:Arg2 arg creation should not error")
	arg3, err := ishell.NewCmdArg("-x", arg3_key, arg3_type, false, true)
	assert.NoError(t, err, "PositionalCmdArgsParsing:Arg3 arg creation should not error")

	cmd := ishell.Cmd{
		Name: "root",
		Help: "root help",
		Func: nil,
	}
	cmd.AddCmdArg(arg1)
	cmd.AddCmdArg(arg2)
	cmd.AddCmdArg(arg3)

	ex1 := []string{"-x", "1", "test"}
	parsed1, err := cmd.ParseArgs(ex1)
	if assert.NoError(t, err, "PositionalCmdArgsParsing:Test1 should not error") {
		assert.Equal(t, 2, len(parsed1), "PositionalCmdArgsParsing:Test1 should have 2 arguments")
		// check param x
		assert.Equal(t, 2, parsed1[0].Index, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Index %d != %d", 2, parsed1[0].Index))
		assert.Equal(t, arg3_key, parsed1[0].Key, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Key %s != %s", arg3_key, parsed1[0].Key))
		assert.Equal(t, arg3_type, parsed1[0].Typ, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Typ %d != %d", arg3_type, parsed1[0].Typ))
		assert.Equal(t, "1", parsed1[0].Value, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Value %s != %s", "1", parsed1[0].Value))
		// check test1
		assert.Equal(t, 0, parsed1[1].Index, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Index %d != %d", 0, parsed1[1].Index))
		assert.Equal(t, arg1_key, parsed1[1].Key, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Key %s != %s", arg1_key, parsed1[1].Key))
		assert.Equal(t, arg1_type, parsed1[1].Typ, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Typ %d != %d", arg1_type, parsed1[1].Typ))
		assert.Equal(t, "test", parsed1[1].Value, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Value %s != %s", "test", parsed1[1].Value))
	}

	ex2 := []string{"test1", "-x", "1", "test2"}
	parsed2, err := cmd.ParseArgs(ex2)
	if assert.NoError(t, err, "PositionalCmdArgsParsing:Test1 should not error") {
		assert.Equal(t, 3, len(parsed2), "PositionalCmdArgsParsing:Test1 should have 3 arguments")
		// check test1
		idx := 0
		assert.Equal(t, 0, parsed2[idx].Index, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Index %d != %d", 0, parsed2[idx].Index))
		assert.Equal(t, arg1_key, parsed2[idx].Key, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Key %s != %s", arg1_key, parsed2[idx].Key))
		assert.Equal(t, arg1_type, parsed2[idx].Typ, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Typ %d != %d", arg1_type, parsed2[idx].Typ))
		assert.Equal(t, "test1", parsed2[idx].Value, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Value %s != %s", "test1", parsed2[idx].Value))
		// check param -x
		idx = 1
		assert.Equal(t, 2, parsed2[idx].Index, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Index %d != %d", 2, parsed2[idx].Index))
		assert.Equal(t, arg3_key, parsed2[idx].Key, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Key %s != %s", arg3_key, parsed2[idx].Key))
		assert.Equal(t, arg3_type, parsed2[idx].Typ, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Typ %d != %d", arg3_type, parsed2[idx].Typ))
		assert.Equal(t, "1", parsed2[idx].Value, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Value %s != %s", "1", parsed2[idx].Value))
		// check test2
		idx = 2
		assert.Equal(t, 1, parsed2[idx].Index, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Index %d != %d", 1, parsed2[idx].Index))
		assert.Equal(t, arg2_key, parsed2[idx].Key, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Key %s != %s", arg2_key, parsed2[idx].Key))
		assert.Equal(t, arg2_type, parsed2[idx].Typ, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Typ %d != %d", arg2_type, parsed2[idx].Typ))
		assert.Equal(t, "test2", parsed2[idx].Value, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Value %s != %s", "test2", parsed2[idx].Value))
	}

}

// test canHaveMultiple with positionals
// it doesn't work well if there are multiple positional args
func TestPositionalCmdArgsParsing2(t *testing.T) {
	arg1_type := ishell.StringType
	arg2_type := ishell.StringType

	arg1_key := "test1"
	arg2_key := "test2"
	// positional argument
	arg1, err := ishell.NewCmdArg("", arg1_key, arg1_type, true, true)
	assert.NoError(t, err, "PositionalCmdArgsParsing:Arg1 arg creation should not error")
	arg2, err := ishell.NewCmdArg("", arg2_key, arg2_type, false, false)
	assert.NoError(t, err, "PositionalCmdArgsParsing:Arg2 arg creation should not error")

	cmd := ishell.Cmd{
		Name: "root",
		Help: "root help",
		Func: nil,
	}
	cmd.AddCmdArg(arg1)
	cmd.AddCmdArg(arg2)

	// because 'test1' can accept multiple, no args should be parsed as 'test2'
	ex1 := []string{"test1", "test1"}
	parsed1, err := cmd.ParseArgs(ex1)
	if assert.NoError(t, err, "PositionalCmdArgsParsing:Test1 should not error") {
		assert.Equal(t, 2, len(parsed1), "PositionalCmdArgsParsing2:Test1 should have 2 arguments")
		// check first test1
		idx := 0
		assert.Equal(t, 0, parsed1[idx].Index, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Index %d != %d", 0, parsed1[idx].Index))
		assert.Equal(t, arg1_key, parsed1[idx].Key, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Key %s != %s", arg1_key, parsed1[idx].Key))
		assert.Equal(t, arg1_type, parsed1[idx].Typ, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Typ %d != %d", arg1_type, parsed1[idx].Typ))
		assert.Equal(t, "test1", parsed1[idx].Value, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Value %s != %s", "test1", parsed1[idx].Value))
		// check second test1
		idx = 1
		assert.Equal(t, 0, parsed1[idx].Index, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Index %d != %d", 0, parsed1[idx].Index))
		assert.Equal(t, arg1_key, parsed1[idx].Key, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Key %s != %s", arg1_key, parsed1[idx].Key))
		assert.Equal(t, arg1_type, parsed1[idx].Typ, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Typ %d != %d", arg1_type, parsed1[idx].Typ))
		assert.Equal(t, "test1", parsed1[idx].Value, fmt.Sprintf("PositionalCmdArgsParsing:Test1 Value %s != %s", "test1", parsed1[idx].Value))
	}
}