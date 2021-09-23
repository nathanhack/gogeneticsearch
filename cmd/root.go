package cmd

import (
	"fmt"
	"github.com/nathanhack/gogeneticsearch/search"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os/exec"
)

var iteration, totalPerIter, mutatedPerIter int
var randomScript, historyScript, mutateScript, testScript, storeScript string

func init() {
	rootCmd.Flags().IntVarP(&iteration, "iter", "", 1_000_000, "Number of iterations to run")
	rootCmd.Flags().IntVarP(&totalPerIter, "samples", "", 5, "Total number of samples to run per iteration (random + mutated).")
	rootCmd.Flags().IntVarP(&mutatedPerIter, "mut", "", 5, "Number of mutated samples to try and generate per iteration.")

	rootCmd.Flags().StringVarP(&randomScript, "rand", "", "random.sh", "A script or program to generate random samples to test (return must be a utf-8 string).")
	rootCmd.Flags().StringVarP(&historyScript, "history", "", "history.sh", "A script or program to select a random sample from the TOP samples to later mutate (return must be a utf-8 string).")
	rootCmd.Flags().StringVarP(&mutateScript, "mutate", "", "mutate.sh", "A script or program that takes two commandline args (strings) to mutate returns the result (return must be a utf-8 string).")
	rootCmd.Flags().StringVarP(&testScript, "test", "", "test.sh", "A script or program that takes one commandline arg (string) to test (return must be a utf-8 string).")
	rootCmd.Flags().StringVarP(&storeScript, "store", "", "store.sh", "A script or program that takes two commandline args (sample and results - both strings) to be stored for later queries.")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gogeneticsearch",
	Short: "Simple tool for running genetic search algorithm",
	Long: `Simple tool for running genetic search algorithm. 
gogeneticsearch will orchestrate the search by calling specific 
programs/scripts for required actions. All input and output is 
expected to be in the form of a utf-8 string.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if totalPerIter < mutatedPerIter {
			return fmt.Errorf("total sample per iteration must be larger equal to or larger than mutated")
		}
		search.Run(iteration, totalPerIter-mutatedPerIter, mutatedPerIter, func() string {
			out, err := exec.Command(randomScript).Output()
			if err != nil {
				logrus.Fatal(err)
			}
			return string(out)
		}, func() string {
			out, err := exec.Command(historyScript).Output()
			if err != nil {
				logrus.Fatal(err)
			}
			return string(out)
		}, func(sample1, sample2 string) string {
			out, err := exec.Command(mutateScript, sample1, sample2).Output()
			if err != nil {
				logrus.Fatal(err)
			}
			return string(out)
		}, func(sample string) string {
			out, err := exec.Command(testScript, sample).Output()
			if err != nil {
				logrus.Fatal(err)
			}
			return string(out)
		}, func(sample, result string) {
			err := exec.Command(storeScript, sample, result).Run()
			if err != nil {
				logrus.Fatal(err)
			}
		})
		return nil
	},
}
