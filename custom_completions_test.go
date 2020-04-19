package cobra

import (
	"bytes"
	"strings"
	"testing"
)

func validArgsFunc(cmd *Command, args []string, toComplete string) ([]string, ShellCompDirective) {
	if len(args) != 0 {
		return nil, ShellCompDirectiveNoFileComp
	}

	var completions []string
	for _, comp := range []string{"one\tThe first", "two\tThe second"} {
		if strings.HasPrefix(comp, toComplete) {
			completions = append(completions, comp)
		}
	}
	return completions, ShellCompDirectiveDefault
}

func validArgsFunc2(cmd *Command, args []string, toComplete string) ([]string, ShellCompDirective) {
	if len(args) != 0 {
		return nil, ShellCompDirectiveNoFileComp
	}

	var completions []string
	for _, comp := range []string{"three\tThe third", "four\tThe fourth"} {
		if strings.HasPrefix(comp, toComplete) {
			completions = append(completions, comp)
		}
	}
	return completions, ShellCompDirectiveDefault
}

func TestValidArgsFuncSingleCmd(t *testing.T) {
	rootCmd := &Command{
		Use:               "root",
		ValidArgsFunction: validArgsFunc,
		Run:               emptyRun,
	}

	// Test completing an empty string
	output, err := executeCommand(rootCmd, ShellCompNoDescRequestCmd, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		"one",
		"two",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Check completing with a prefix
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"two",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}

func TestValidArgsFuncSingleCmdInvalidArg(t *testing.T) {
	rootCmd := &Command{
		Use: "root",
		// If we don't specify a value for Args, this test fails.
		// This is only true for a root command without any subcommands, and is caused
		// by the fact that the __complete command becomes a subcommand when there should not be one.
		// The problem is in the implementation of legacyArgs().
		Args:              MinimumNArgs(1),
		ValidArgsFunction: validArgsFunc,
		Run:               emptyRun,
	}

	// Check completing with wrong number of args
	output, err := executeCommand(rootCmd, ShellCompNoDescRequestCmd, "unexpectedArg", "t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}

func TestValidArgsFuncChildCmds(t *testing.T) {
	rootCmd := &Command{Use: "root", Args: NoArgs, Run: emptyRun}
	child1Cmd := &Command{
		Use:               "child1",
		ValidArgsFunction: validArgsFunc,
		Run:               emptyRun,
	}
	child2Cmd := &Command{
		Use:               "child2",
		ValidArgsFunction: validArgsFunc2,
		Run:               emptyRun,
	}
	rootCmd.AddCommand(child1Cmd, child2Cmd)

	// Test completion of first sub-command with empty argument
	output, err := executeCommand(rootCmd, ShellCompNoDescRequestCmd, "child1", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		"one",
		"two",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test completion of first sub-command with a prefix to complete
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "child1", "t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"two",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Check completing with wrong number of args
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "child1", "unexpectedArg", "t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test completion of second sub-command with empty argument
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "child2", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"three",
		"four",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "child2", "t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"three",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Check completing with wrong number of args
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "child2", "unexpectedArg", "t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}

func TestValidArgsFuncAliases(t *testing.T) {
	rootCmd := &Command{Use: "root", Args: NoArgs, Run: emptyRun}
	child := &Command{
		Use:               "child",
		Aliases:           []string{"son", "daughter"},
		ValidArgsFunction: validArgsFunc,
		Run:               emptyRun,
	}
	rootCmd.AddCommand(child)

	// Test completion of first sub-command with empty argument
	output, err := executeCommand(rootCmd, ShellCompNoDescRequestCmd, "son", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		"one",
		"two",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test completion of first sub-command with a prefix to complete
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "daughter", "t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"two",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Check completing with wrong number of args
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "son", "unexpectedArg", "t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}

func TestValidArgsFuncInBashScript(t *testing.T) {
	rootCmd := &Command{Use: "root", Args: NoArgs, Run: emptyRun}
	child := &Command{
		Use:               "child",
		ValidArgsFunction: validArgsFunc,
		Run:               emptyRun,
	}
	rootCmd.AddCommand(child)

	buf := new(bytes.Buffer)
	rootCmd.GenBashCompletion(buf)
	output := buf.String()

	check(t, output, "has_completion_function=1")
}

func TestNoValidArgsFuncInBashScript(t *testing.T) {
	rootCmd := &Command{Use: "root", Args: NoArgs, Run: emptyRun}
	child := &Command{
		Use: "child",
		Run: emptyRun,
	}
	rootCmd.AddCommand(child)

	buf := new(bytes.Buffer)
	rootCmd.GenBashCompletion(buf)
	output := buf.String()

	checkOmit(t, output, "has_completion_function=1")
}

func TestCompleteCmdInBashScript(t *testing.T) {
	rootCmd := &Command{Use: "root", Args: NoArgs, Run: emptyRun}
	child := &Command{
		Use:               "child",
		ValidArgsFunction: validArgsFunc,
		Run:               emptyRun,
	}
	rootCmd.AddCommand(child)

	buf := new(bytes.Buffer)
	rootCmd.GenBashCompletion(buf)
	output := buf.String()

	check(t, output, ShellCompNoDescRequestCmd)
}

func TestCompleteNoDesCmdInFishScript(t *testing.T) {
	rootCmd := &Command{Use: "root", Args: NoArgs, Run: emptyRun}
	child := &Command{
		Use:               "child",
		ValidArgsFunction: validArgsFunc,
		Run:               emptyRun,
	}
	rootCmd.AddCommand(child)

	buf := new(bytes.Buffer)
	rootCmd.GenFishCompletion(buf, false)
	output := buf.String()

	check(t, output, ShellCompNoDescRequestCmd)
}

func TestCompleteCmdInFishScript(t *testing.T) {
	rootCmd := &Command{Use: "root", Args: NoArgs, Run: emptyRun}
	child := &Command{
		Use:               "child",
		ValidArgsFunction: validArgsFunc,
		Run:               emptyRun,
	}
	rootCmd.AddCommand(child)

	buf := new(bytes.Buffer)
	rootCmd.GenFishCompletion(buf, true)
	output := buf.String()

	check(t, output, ShellCompRequestCmd)
	checkOmit(t, output, ShellCompNoDescRequestCmd)
}

func TestFlagCompletionInGo(t *testing.T) {
	rootCmd := &Command{
		Use: "root",
		Run: emptyRun,
	}
	rootCmd.Flags().IntP("introot", "i", -1, "help message for flag introot")
	rootCmd.RegisterFlagCompletionFunc("introot", func(cmd *Command, args []string, toComplete string) ([]string, ShellCompDirective) {
		completions := []string{}
		for _, comp := range []string{"1\tThe first", "2\tThe second", "10\tThe tenth"} {
			if strings.HasPrefix(comp, toComplete) {
				completions = append(completions, comp)
			}
		}
		return completions, ShellCompDirectiveDefault
	})
	rootCmd.Flags().String("filename", "", "Enter a filename")
	rootCmd.RegisterFlagCompletionFunc("filename", func(cmd *Command, args []string, toComplete string) ([]string, ShellCompDirective) {
		completions := []string{}
		for _, comp := range []string{"file.yaml\tYAML format", "myfile.json\tJSON format", "file.xml\tXML format"} {
			if strings.HasPrefix(comp, toComplete) {
				completions = append(completions, comp)
			}
		}
		return completions, ShellCompDirectiveNoSpace | ShellCompDirectiveNoFileComp
	})

	// Test completing an empty string
	output, err := executeCommand(rootCmd, ShellCompNoDescRequestCmd, "--introot", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		"1",
		"2",
		"10",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Check completing with a prefix
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "--introot", "1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"1",
		"10",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test completing an empty string
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "--filename", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"file.yaml",
		"myfile.json",
		"file.xml",
		":6",
		"Completion ended with directive: ShellCompDirectiveNoSpace, ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Check completing with a prefix
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "--filename", "f")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"file.yaml",
		"file.xml",
		":6",
		"Completion ended with directive: ShellCompDirectiveNoSpace, ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}

func TestValidArgsFuncChildCmdsWithDesc(t *testing.T) {
	rootCmd := &Command{Use: "root", Args: NoArgs, Run: emptyRun}
	child1Cmd := &Command{
		Use:               "child1",
		ValidArgsFunction: validArgsFunc,
		Run:               emptyRun,
	}
	child2Cmd := &Command{
		Use:               "child2",
		ValidArgsFunction: validArgsFunc2,
		Run:               emptyRun,
	}
	rootCmd.AddCommand(child1Cmd, child2Cmd)

	// Test completion of first sub-command with empty argument
	output, err := executeCommand(rootCmd, ShellCompRequestCmd, "child1", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		"one\tThe first",
		"two\tThe second",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test completion of first sub-command with a prefix to complete
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "child1", "t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"two\tThe second",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Check completing with wrong number of args
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "child1", "unexpectedArg", "t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test completion of second sub-command with empty argument
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "child2", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"three\tThe third",
		"four\tThe fourth",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "child2", "t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"three\tThe third",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Check completing with wrong number of args
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "child2", "unexpectedArg", "t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}

func TestFlagCompletionInGoWithDesc(t *testing.T) {
	rootCmd := &Command{
		Use: "root",
		Run: emptyRun,
	}
	rootCmd.Flags().IntP("introot", "i", -1, "help message for flag introot")
	rootCmd.RegisterFlagCompletionFunc("introot", func(cmd *Command, args []string, toComplete string) ([]string, ShellCompDirective) {
		completions := []string{}
		for _, comp := range []string{"1\tThe first", "2\tThe second", "10\tThe tenth"} {
			if strings.HasPrefix(comp, toComplete) {
				completions = append(completions, comp)
			}
		}
		return completions, ShellCompDirectiveDefault
	})
	rootCmd.Flags().String("filename", "", "Enter a filename")
	rootCmd.RegisterFlagCompletionFunc("filename", func(cmd *Command, args []string, toComplete string) ([]string, ShellCompDirective) {
		completions := []string{}
		for _, comp := range []string{"file.yaml\tYAML format", "myfile.json\tJSON format", "file.xml\tXML format"} {
			if strings.HasPrefix(comp, toComplete) {
				completions = append(completions, comp)
			}
		}
		return completions, ShellCompDirectiveNoSpace | ShellCompDirectiveNoFileComp
	})

	// Test completing an empty string
	output, err := executeCommand(rootCmd, ShellCompRequestCmd, "--introot", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		"1\tThe first",
		"2\tThe second",
		"10\tThe tenth",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Check completing with a prefix
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "--introot", "1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"1\tThe first",
		"10\tThe tenth",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test completing an empty string
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "--filename", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"file.yaml\tYAML format",
		"myfile.json\tJSON format",
		"file.xml\tXML format",
		":6",
		"Completion ended with directive: ShellCompDirectiveNoSpace, ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Check completing with a prefix
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "--filename", "f")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"file.yaml\tYAML format",
		"file.xml\tXML format",
		":6",
		"Completion ended with directive: ShellCompDirectiveNoSpace, ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}

func TestFlagNameCompletionInGo(t *testing.T) {
	rootCmd := &Command{
		Use: "root",
		Run: emptyRun,
	}
	childCmd := &Command{
		Use: "childCmd",
		Run: emptyRun,
	}
	rootCmd.AddCommand(childCmd)

	rootCmd.Flags().IntP("first", "f", -1, "first flag")
	rootCmd.PersistentFlags().BoolP("second", "s", false, "second flag")
	childCmd.Flags().String("subFlag", "", "sub flag")

	// Test that flag names are not shown if the user has not given the '-' prefix
	output, err := executeCommand(rootCmd, ShellCompNoDescRequestCmd, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		"childCmd",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that flag names are completed
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "-")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--first",
		"--first=",
		"-f",
		"--second",
		"-s",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that flag names are completed when a prefix is given
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "--f")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--first",
		"--first=",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that flag names are completed in a sub-cmd
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "childCmd", "-")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--second",
		"-s",
		"--subFlag",
		"--subFlag=",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}

func TestFlagNameCompletionInGoWithDesc(t *testing.T) {
	rootCmd := &Command{
		Use: "root",
		Run: emptyRun,
	}
	childCmd := &Command{
		Use:   "childCmd",
		Short: "first command",
		Run:   emptyRun,
	}
	rootCmd.AddCommand(childCmd)

	rootCmd.Flags().IntP("first", "f", -1, "first flag")
	rootCmd.PersistentFlags().BoolP("second", "s", false, "second flag")
	childCmd.Flags().String("subFlag", "", "sub flag")

	// Test that flag names are not shown if the user has not given the '-' prefix
	output, err := executeCommand(rootCmd, ShellCompRequestCmd, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		"childCmd\tfirst command",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that flag names are completed
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "-")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--first\tfirst flag",
		"--first=\tfirst flag",
		"-f\tfirst flag",
		"--second\tsecond flag",
		"-s\tsecond flag",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that flag names are completed when a prefix is given
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "--f")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--first\tfirst flag",
		"--first=\tfirst flag",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that flag names are completed in a sub-cmd
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "childCmd", "-")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--second\tsecond flag",
		"-s\tsecond flag",
		"--subFlag\tsub flag",
		"--subFlag=\tsub flag",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}

func TestRequiredFlagNameCompletionInGo(t *testing.T) {
	rootCmd := &Command{
		Use:       "root",
		ValidArgs: []string{"realArg"},
		Run:       emptyRun,
	}
	childCmd := &Command{
		Use: "childCmd",
		ValidArgsFunction: func(cmd *Command, args []string, toComplete string) ([]string, ShellCompDirective) {
			return []string{"subArg"}, ShellCompDirectiveNoFileComp
		},
		Run: emptyRun,
	}
	rootCmd.AddCommand(childCmd)

	rootCmd.Flags().IntP("requiredFlag", "r", -1, "required flag")
	rootCmd.MarkFlagRequired("requiredFlag")
	requiredFlag := rootCmd.Flags().Lookup("requiredFlag")

	rootCmd.PersistentFlags().IntP("requiredPersistent", "p", -1, "required persistent")
	rootCmd.MarkPersistentFlagRequired("requiredPersistent")
	requiredPersistent := rootCmd.PersistentFlags().Lookup("requiredPersistent")

	rootCmd.Flags().StringP("release", "R", "", "Release name")

	childCmd.Flags().BoolP("subRequired", "s", false, "sub required flag")
	childCmd.MarkFlagRequired("subRequired")
	childCmd.Flags().BoolP("subNotRequired", "n", false, "sub not required flag")

	// Test that a required flag is suggested even without the - prefix
	output, err := executeCommand(rootCmd, ShellCompNoDescRequestCmd, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		"childCmd",
		"--requiredFlag",
		"--requiredFlag=",
		"-r",
		"--requiredPersistent",
		"--requiredPersistent=",
		"-p",
		"realArg",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that a required flag is suggested without other flags when using the '-' prefix
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "-")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--requiredFlag",
		"--requiredFlag=",
		"-r",
		"--requiredPersistent",
		"--requiredPersistent=",
		"-p",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that if no required flag matches, the normal flag is suggested
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "--relea")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--release",
		"--release=",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test required flags for sub-commands
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "childCmd", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--requiredPersistent",
		"--requiredPersistent=",
		"-p",
		"--subRequired",
		"-s",
		"subArg",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "childCmd", "-")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--requiredPersistent",
		"--requiredPersistent=",
		"-p",
		"--subRequired",
		"-s",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "childCmd", "--subNot")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--subNotRequired",
		":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that when a required flag is present, it is not suggested anymore
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "--requiredFlag", "1", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Reset the flag for the next command
	requiredFlag.Changed = false

	expected = strings.Join([]string{
		"childCmd",
		"--requiredPersistent",
		"--requiredPersistent=",
		"-p",
		"realArg",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that when a persistent required flag is present, it is not suggested anymore
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "--requiredPersistent", "1", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Reset the flag for the next command
	requiredPersistent.Changed = false

	expected = strings.Join([]string{
		"childCmd",
		"--requiredFlag",
		"--requiredFlag=",
		"-r",
		"realArg",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that when all required flags are present, normal completion is done
	output, err = executeCommand(rootCmd, ShellCompNoDescRequestCmd, "--requiredFlag", "1", "--requiredPersistent", "1", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Reset the flags for the next command
	requiredFlag.Changed = false
	requiredPersistent.Changed = false

	expected = strings.Join([]string{
		"childCmd",
		"realArg",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}
