/*
Copyright © 2021 Nathan Leniz <terakilobyte@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"checker/internal/collectors"
	"checker/internal/parsers/intersphinx"
	"checker/internal/parsers/rst"
	"checker/internal/sources"
	"checker/internal/utils"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var path string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "checker",
	Short: "Checks refs, roles, and links in a docs project",
	Long: `Checker is a tool for checking refs, roles, and links in a docs project.
It will check refs against locally found refs and those found in intersphinx targets,
and checks roles against the latest RELEASE of rstspec.toml. Once they are validated,
all links are checked for validity.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		type intersphinxResult struct {
			domain string
			file   []byte
		}

		// rstSpec := utils.GetFile(utils.GetLatestSnootyParserTag())
		basepath, err := filepath.Abs(path)
		if err != nil {
			log.Panic(err)
		}
		snootyToml := utils.GetLocalFile(filepath.Join(basepath, "snooty.toml"))
		projectCfg, err := sources.NewTomlConfig(snootyToml)
		if err != nil {
			log.Panic(err)
		}
		intersphinxes := make([]intersphinx.SphinxMap, len(projectCfg.Intersphinx))
		var wgSetup sync.WaitGroup
		ixs := make(chan intersphinxResult, len(projectCfg.Intersphinx))
		for _, intersphinx := range projectCfg.Intersphinx {
			wgSetup.Add(1)
			go func(phx string) {
				domain := strings.Split(phx, "objects.inv")[0]
				file := utils.GetNetworkFile(phx)
				ixs <- intersphinxResult{domain: domain, file: file}
			}(intersphinx)
		}
		go func() {
			for res := range ixs {
				intersphinxes = append(intersphinxes, intersphinx.Intersphinx(res.file, res.domain))
				wgSetup.Done()
			}
		}()
		wgSetup.Wait()

		files := collectors.GatherFiles(basepath)

		allConstants := collectors.GatherConstants(files)
		allRoleTargets := collectors.GatherRoles(files)
		allHTTPLinks := collectors.GatherHTTPLinks(files)
		// allLocalRefs := collectors.GatherLocalRefs(files)

		for con, filename := range allConstants {
			if _, ok := projectCfg.Constants[con.Name]; !ok {
				log.Errorf("%s is not defined in the config", con)
			}
			testCon := rst.RstConstant{Name: con.Name, Target: projectCfg.Constants[filename] + con.Name}
			if testCon.IsHTTPLink() {
				allHTTPLinks[rst.RstHTTPLink(con.Target)] = filename
			}
		}
		checkedUrls := sync.Map{}
		var wgValidate sync.WaitGroup
		workStack := make([]func(), 0)
		rstSpecRoles := sources.NewRoleMap(utils.GetNetworkFile(utils.GetLatestSnootyParserTag()))
		for role, filename := range allRoleTargets {
			if role.RoleType == "role" && role.Name == "manual" {
				if _, ok := rstSpecRoles[role.Name]; !ok {
					log.Errorf("%s is not defined in the rstspec", role)
				}
				url := fmt.Sprintf(rstSpecRoles[role.Name], role.Target)
				errmsg := fmt.Sprintf("in %s: interpeted url %s from  %+v was not valid", filename, url, role)
				workFunc := func() {
					defer wgValidate.Done()
					wgValidate.Add(1)
					if _, ok := checkedUrls.Load(url); !ok {
						checkedUrls.Store(url, true)
						if !utils.IsReachable(url) {
							log.Error(errmsg)
						}
					}
				}
				workStack = append(workStack, workFunc)
			}
			for _, f := range workStack {
				f()
			}
			wgValidate.Wait()

			// if role.RoleType == "ref" {
			// 	for k, v := range allRoleTargets {

			// 	}
			// }
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&path, "path", "", "path to the project")
	if err := rootCmd.MarkPersistentFlagRequired("path"); err != nil {
		log.Panic(err)
	}
}
