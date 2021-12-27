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
	"regexp"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	path string
	refs bool
)

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

		diagnostics := make([]string, 0)
		diags := make(chan string)
		go func() {
			for d := range diags {
				diagnostics = append(diagnostics, d)
			}
		}()

		type intersphinxResult struct {
			domain string
			file   []byte
		}

		basepath, err := filepath.Abs(path)
		if err != nil {
			log.Panic(err)
		}
		snootyToml := utils.GetLocalFile(filepath.Join(basepath, "snooty.toml"))
		projectSnooty, err := sources.NewTomlConfig(snootyToml)
		if err != nil {
			log.Panic(err)
		}
		intersphinxes := make([]intersphinx.SphinxMap, len(projectSnooty.Intersphinx))
		var wgSetup sync.WaitGroup
		ixs := make(chan intersphinxResult, len(projectSnooty.Intersphinx))
		for _, intersphinx := range projectSnooty.Intersphinx {
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
		close(ixs)

		sphinxMap := intersphinx.JoinSphinxes(intersphinxes)
		files := collectors.GatherFiles(basepath)

		allShared := collectors.GatherSharedIncludes(files)

		sharedRefs := make(collectors.RstRoleMap)
		sharedLocals := make(collectors.RefTargetMap)

		for _, share := range allShared {
			sharedFile := utils.GetNetworkFile(projectSnooty.SharedPath + share.Path)
			sharedRefs.Union(collectors.GatherSharedRefs(sharedFile, *projectSnooty))
			sharedLocals.Union(collectors.GatherSharedLocalRefs(sharedFile, *projectSnooty))
		}

		allConstants := collectors.GatherConstants(files)
		allRoleTargets := collectors.GatherRoles(files)
		allHTTPLinks := collectors.GatherHTTPLinks(files)
		allLocalRefs := collectors.GatherLocalRefs(files).SSLToTLS()

		allRoleTargets.Union(sharedRefs)
		allLocalRefs.Union(sharedLocals)

		allRoleTargets = allRoleTargets.ConvertConstants(projectSnooty)

		for con, filename := range allConstants {
			if _, ok := projectSnooty.Constants[con.Name]; !ok {
				diags <- fmt.Sprintf("%s is not defined in config", con)
			}
			testCon := rst.RstConstant{Name: con.Name, Target: projectSnooty.Constants[filename] + con.Name}
			if testCon.IsHTTPLink() {
				allHTTPLinks[rst.RstHTTPLink(testCon.Target)] = filename
			}
		}

		checkedUrls := sync.Map{}
		var wgValidate sync.WaitGroup
		workStack := make([]func(), 0)
		rstSpecRoles := sources.NewRoleMap(utils.GetNetworkFile(utils.GetLatestSnootyParserTag()))

		// limit concurrency to 5
		semaphore := make(chan struct{}, 5)

		// have a max rate of 100/sec
		rate := make(chan struct{}, 100)
		for i := 0; i < cap(rate); i++ {
			rate <- struct{}{}
		}

		// leaky bucket
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			for range ticker.C {
				_, ok := <-rate
				if !ok {
					return
				}
			}
		}()

		for role, filename := range allRoleTargets {
			switch role.Name {
			case "ref":
				if refs {
					if _, ok := sphinxMap[role.Target]; !ok {
						if _, ok := allLocalRefs.Get(&role); !ok {
							diags <- fmt.Sprintf("in %s: %+v is not a valid ref", filename, role)
						}
					}
				}
			case "doc":
				found := false
				for _, f := range files {
					if matched, _ := regexp.Match(strings.TrimSuffix(role.Target, "/"), []byte(f)); matched {
						found = true
						break
					}
				}
				if !found {
					diags <- fmt.Sprintf("in %s: %s is not a valid file found in this docset", filename, role)
				}

			case "py:meth":
				if refs {
					if _, ok := sphinxMap[role.Target]; !ok {
						if _, ok := allLocalRefs.Get(&role); !ok {
							diags <- fmt.Sprintf("in %s: %+v is not a valid ref", filename, role)
						}
					}
				}
			case "py:class":
				if refs {
					if _, ok := sphinxMap[role.Target]; !ok {
						if _, ok := allLocalRefs.Get(&role); !ok {
							diags <- fmt.Sprintf("in %s: %+v is not a valid ref", filename, role)
						}
					}
				}
			default:
				if _, ok := rstSpecRoles.Roles[role.Name]; !ok {
					if _, ok := rstSpecRoles.RawRoles[role.Name]; !ok {
						if _, ok := rstSpecRoles.RstObjects[role.Name]; !ok {
							diags <- fmt.Sprintf("in %s: %s is not a valid role", filename, role)
						}
					}
					continue
				}
				url := fmt.Sprintf(rstSpecRoles.Roles[role.Name], role.Target)
				workFunc := func() {
					defer wgValidate.Done()
					rate <- struct{}{}
					semaphore <- struct{}{}
					defer func() {
						<-semaphore
					}()
					wgValidate.Add(1)
					if _, ok := checkedUrls.Load(url); !ok {
						checkedUrls.Store(url, true)
						if resp, ok := utils.IsReachable(url); !ok {
							errmsg := fmt.Sprintf("in %s: interpeted url %s from  %+v was not valid. Got response %+v", filename, url, role, resp)
							diags <- errmsg
						}
					}
				}
				workStack = append(workStack, workFunc)
			}
		}

		for link, filename := range allHTTPLinks {
			workStack = append(workStack, func() {
				defer wgValidate.Done()
				rate <- struct{}{}
				semaphore <- struct{}{}
				defer func() {
					<-semaphore
				}()
				wgValidate.Add(1)
				if _, ok := checkedUrls.Load(link); !ok {
					checkedUrls.Store(link, true)
					if resp, ok := utils.IsReachable(string(link)); !ok {
						errmsg := fmt.Sprintf("in %s: %s is not a valid http link. Got response %+v", filename, link, resp)
						diags <- errmsg
					}
				}
			})
		}
		for _, f := range workStack {
			go f()
		}
		wgValidate.Wait()
		for _, msg := range diagnostics {
			log.Error(msg)
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
	rootCmd.PersistentFlags().BoolVar(&refs, "refs", false, "check refs")
}
