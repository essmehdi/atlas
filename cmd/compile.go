/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"atlas/compiler"
	"atlas/parser"
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compiles Atlas code to atlas bytecode",
	Long: `Takes the input code from stdin and compiles to bytecode.

Example: cat code.atl | atlas compile
	
The output file is the bytecode saved in compiled.atlb.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		outputFile, _ := cmd.Flags().GetString("output")

		var codeBuffer strings.Builder

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			codeBuffer.WriteString(line)
		}
		
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
			return
		}

		code := codeBuffer.String()

		parser := parser.New(&code)
		comp := compiler.New()

		program := parser.Parse()

		if len(parser.Errors) > 0 {
			fmt.Println("Parsing failed")
			for _, err := range parser.Errors {
				fmt.Println(err)
			}
			return
		}

		err := comp.Compile(&program)
		if err != nil {
			fmt.Println(err)
			return
		}

		byteCode := comp.ByteCode()

		compiler.RegisterObjectsToGob()

		var bytesBuffer bytes.Buffer
		var encoder = gob.NewEncoder(&bytesBuffer)
		encoder.Encode(byteCode)

		file, _ := os.Create(outputFile)
		file.Write(bytesBuffer.Bytes())
		file.Close()
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
	compileCmd.Flags().String("output", "compiled.atlb", "Output file of the compiled bytecode")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// compileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// compileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}