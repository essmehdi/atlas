package cmd

import (
	"atlas/compiler"
	"atlas/vm"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Executes specified bytecode file",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		byteCodeFilePath := args[0]

		data, err := os.ReadFile(byteCodeFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		compiler.RegisterObjectsToGob()

		var bytesBuffer bytes.Buffer
		bytesBuffer.Write(data) 
		decoder := gob.NewDecoder(&bytesBuffer)

		var byteCode compiler.ByteCode
		err = decoder.Decode(&byteCode)
		if err != nil {
			fmt.Println("Could not deserialize bytecode: ", err)
			return
		}

		vm := vm.New(byteCode)
		err = vm.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
}
