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

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compiles Atlas code to Atlas bytecode",
	Long: `Compiles code from provided file to Atlas bytecode. If a file is not provided, the code compiled is the content of stdin.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputFile, _ := cmd.Flags().GetString("output")

		var pars *parser.Parser

		if len(args) == 1 {
			newParser, err := parser.NewFromFile(args[0])
			if err != nil {
				fmt.Println("Could not create parser: ", err)
				return
			}

			pars = newParser
		} else {
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
	
			pars = parser.New(&code)
		}

		comp := compiler.New()

		program := pars.Parse()

		if len(pars.Errors) > 0 {
			fmt.Println("Parsing failed")
			for _, err := range pars.Errors {
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
	compileCmd.Flags().StringP("output", "o", "compiled.atlb", "Output file of the compiled bytecode")
}