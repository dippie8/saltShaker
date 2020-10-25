/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	utils "saltShaker/utils"
	"strings"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test Salt states from the selected directory",
	Long: `Test Salt states from the selected directory.
States are applied on a container with Centos.
Use -s flag if you want to apply your state in a "virgin" container.`,
	Run: func(cmd *cobra.Command, args []string) {

		// parsing parameters

		var absSaltRoot string

		if len(args) != 1 {
			fmt.Println("[ERROR] You must insert exactly one argument")
			return
		}

		saltRoot := args[0]
		pwd, err := os.Getwd()

		if err != nil {
			fmt.Println("[ERROR] Please retry...")
			return
		}

		if saltRoot[0] == '.' {
			absSaltRoot = pwd + saltRoot[1:]
		} else if saltRoot[0] == '/' {
			absSaltRoot = saltRoot
		} else {
			fmt.Println("[ERROR] please insert a correct path")
		}

		if absSaltRoot[len(absSaltRoot)-1] == '/' {
			absSaltRoot = absSaltRoot[:len(absSaltRoot)-1]
		}


		// se il container non esiste, crealo TODO condizione
		err = utils.BuildSaltshakerImage()
		if err != nil {
			panic(err)
		}

		// se esiste ed è spento, accendilo TODO condizione
		// err = utils.RunSaltshakerContainer()
		// if err != nil {
		// 	panic(err)
		// }

		// se esiste ed è acceso, vai avanti

		containerID := ""
		containers := utils.GetRunningContainers()
		for _, container := range containers {
			if container.Names[0] == "saltshaker" || container.Names[0] =="/saltshaker" {
				containerID = container.ID
			}
		}
		deleteSrvCmd := []string{"rm", "-rf", "/srv/salt"}
		_, err = utils.RunCommand(containerID, deleteSrvCmd)
		if err != nil {
			panic(err)
		}

		createSrvCmd := []string{"mkdir", "/srv/salt"}
		_, err = utils.RunCommand(containerID, createSrvCmd)
		if err != nil {
			panic(err)
		}

		err = utils.CopyToContainer(containerID, absSaltRoot)
		if err != nil {
			panic(err)
		}

		// se non è passato il parametro con lo stato da applicare, passa l'init della cartella principale
		splittedRoot := strings.Split(absSaltRoot, "/")
		module := splittedRoot[len(splittedRoot)-1]
		resp, err := utils.ApplyState(containerID, module)
		if err != nil {
			panic(err)
		}

		fmt.Println(resp.StdOut)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	testCmd.Flags().BoolP("scratch", "s", false, "start test on container from scratch")

}
