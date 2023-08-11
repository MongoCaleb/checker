package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/MongoCaleb/checker/internal/collectors"
	"github.com/MongoCaleb/checker/internal/parsers/intersphinx"
	"github.com/MongoCaleb/checker/internal/parsers/rst"
	"github.com/MongoCaleb/checker/internal/sources"
	"github.com/MongoCaleb/checker/internal/utils"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

var (
	path     string
	refs     bool
	docs     bool
	changes  []string
	progress bool
	workers  int
	throttle int
	loglevel int
)

type bypassJson struct {
	Exclude string `json:"exclude"`
	Reason  string `json:"reason"`
}

var BypassList []bypassJson

func loadBypassList(bypassPath string) {
	jsonFile, err := os.Open(bypassPath + "/config/link_checker_bypass_list.json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var list []bypassJson
	json.Unmarshal(byteValue, &list)

	for i := 0; i < len(list); i++ {
		BypassList = append(BypassList, list[i])
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "checker",
	Version: "0.2.0",
	Short:   "Checks links, and optionally :ref:s, :doc:s, and other :role:s in a docs project.",

	Run: func(cmd *cobra.Command, args []string) {

		if val, ok := os.LookupEnv("CHECKER_WORKERS"); ok {
			v, err := strconv.Atoi(val)
			if err != nil {
				log.Panicf("couldn't convert %s to an int: %v", val, err)
			}
			workers = v
		}

		if val, ok := os.LookupEnv("CHECKER_THROTTLE"); ok {
			v, err := strconv.Atoi(val)
			if err != nil {
				log.Panicf("couldn't convert %s to an int: %v", val, err)
			}
			throttle = v
		}

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

		loadBypassList(path)
		basepath, err := filepath.Abs(path)
		checkErr(err)
		snootyToml := utils.GetLocalFile(filepath.Join(path, "snooty.toml"))
		projectSnooty, err := sources.NewTomlConfig(snootyToml)
		checkErr(err)
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
			testCon := rst.RstConstant{Name: con.Name, Target: projectSnooty.Constants[con.Name] + con.Target}
			if !isBlocked(testCon.Target) && testCon.IsHTTPLink() {
				allHTTPLinks[rst.RstHTTPLink(testCon.Target)] = filename
			}
		}

		checkedUrls := sync.Map{}
		workStack := make([]func(), 0)
		rstSpecRoles := sources.NewRoleMap(utils.GetNetworkFile(utils.GetLatestSnootyParserTag()))

		if len(changes) == 0 {
			changes = files
		}
		for role, filename := range allRoleTargets {
			if !contains(changes, strings.TrimPrefix(filename, "/")) {
				continue
			}

			switch role.Name {
			case "guilabel":
				break
			case "ref":
				if refs {
					if _, ok := sphinxMap[role.Target]; !ok {
						if _, ok := allLocalRefs.Get(&role); !ok {
							diags <- fmt.Sprintf("in %s: %+v is not a valid ref", filename, role)
						}
					}
					break
				}
			case "doc":
				if docs {
					if !contains(files, filename) {
						diags <- fmt.Sprintf("in %s: %s is not a valid file found in this docset", filename, role)
					}
					break
				}

			case "py:meth": // this is a fancy magic ref
				if refs {
					if _, ok := sphinxMap[role.Target]; !ok {
						if _, ok := allLocalRefs.Get(&role); !ok {
							diags <- fmt.Sprintf("in %s: %+v is not a valid ref", filename, role)
						}
					}
					break
				}
			case "py:class": // this is a fancy magic ref
				if refs {
					if _, ok := sphinxMap[role.Target]; !ok {
						if _, ok := allLocalRefs.Get(&role); !ok {
							diags <- fmt.Sprintf("in %s: %+v is not a valid ref", filename, role)
						}
					}
					break
				}
			default:
				if isBlocked(rstSpecRoles.Roles[role.Name]) {
					break
				}
				if _, ok := rstSpecRoles.Roles[role.Name]; !ok {
					if _, ok := rstSpecRoles.RawRoles[role.Name]; !ok {
						if _, ok := rstSpecRoles.RstObjects[role.Name]; !ok {
							diags <- fmt.Sprintf("in %s: %s is not a valid role", filename, role)
						}
					}
					break
				}
				workFunc := func(role rst.RstRole, filename string) func() {
					url := fmt.Sprintf(rstSpecRoles.Roles[role.Name], role.Target)

					if _, ok := checkedUrls.Load(url); !ok {
						return func() {
							checkedUrls.Store(url, true)
							if resp, ok := utils.IsReachable(url); !ok {
								errmsg := fmt.Sprintf("in %s: interpreted url %s from  %+v was not valid. Got response %s", filename, url, role, resp)
								diags <- errmsg
							}
						}
					} else {
						return func() {}
					}

				}

				i := isBlocked(role.Target)
				if !i {
					workStack = append(workStack, workFunc(role, filename))
				} else {
					log.Error("roletarget_excluded: ", role.Target)
				}
			}
		}

		for link, filename := range allHTTPLinks {
			if !contains(changes, strings.TrimPrefix(filename, "/")) {
				continue
			}
			workFunc := func(link rst.RstHTTPLink, filename string) func() {
				if _, ok := checkedUrls.Load(link); !ok {
					return func() {
						checkedUrls.Store(link, true)
						if resp, ok := utils.IsReachable(string(link)); !ok {
							errmsg := fmt.Sprintf("%s | %s", filename, resp)
							diags <- errmsg
						}
					}
				} else {
					return func() {}
				}
			}

			i := isBlocked(string(link))
			if !i {
				workStack = append(workStack, workFunc(link, filename))
			}
		}

		jobChannel := make(chan func())
		doneChannel := make(chan struct{})

		var wgValidate sync.WaitGroup
		wgValidate.Add(workers)
		for i := 0; i < workers; i++ {
			go worker(&wgValidate, jobChannel, doneChannel)
		}

		bar := pb.StartNew(len(workStack)).SetMaxWidth(120)
		if loglevel > 0 {
			log.Info(fmt.Sprintf("Checking %d links", len(workStack)))
		}
		if progress && loglevel > 1 {
			log.Info(progress)
			bar.SetWriter(os.Stdout)
		} else {
			bar.SetWriter(ioutil.Discard)
		}
		go func() {
			for range doneChannel {
				bar.Increment()
			}
		}()

		for _, f := range workStack {
			jobChannel <- f
		}

		close(jobChannel)
		wgValidate.Wait()
		bar.Finish()
		for _, msg := range diagnostics {
			if loglevel > 0 {
				log.Error(msg)
			}
		}

		if len(diagnostics) > 0 {
			log.Fatal(len(diagnostics), " errors found.\n")
		} else {
			{
				log.Info("No errors found.\n")
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.SetVersionTemplate("checker {{.Version}}\n")
	rootCmd.PersistentFlags().IntVarP(&loglevel, "loglevel", "l", 2, "0=silence all, 1=results only, 2=info and results")
	rootCmd.PersistentFlags().StringVar(&path, "path", ".", "path to the project")
	rootCmd.PersistentFlags().BoolVarP(&refs, "refs", "r", true, "check :refs:")
	rootCmd.PersistentFlags().BoolVarP(&docs, "docs", "d", true, "check :docs:")
	rootCmd.PersistentFlags().StringSliceVar(&changes, "changes", []string{}, "The list of files to check")
	rootCmd.PersistentFlags().BoolVarP(&progress, "progress", "p", true, "show progress bar")
	rootCmd.PersistentFlags().IntVarP(&workers, "workers", "w", 100, "The number of workers to spawn to do work.")
	rootCmd.PersistentFlags().IntVarP(&throttle, "throttle", "t", 100, "The throttle factor. Each worker will process at most (1e9 / (throttle / workers)) jobs per second.")
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
func isBlocked(input string) bool {
	for _, a := range BypassList {
		if !strings.Contains(input, a.Exclude) {
			continue
		} else {
			if loglevel >= 2 {
				log.Printf("Excluded: %s - Reason: %s %s", input, a.Exclude, a.Reason)
			}
			return true
		}
	}
	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {

		if strings.Contains(a, e) {
			return true
		}
	}
	return false
}

func worker(wg *sync.WaitGroup, jobChannel <-chan func(), doneChannel chan<- struct{}) {
	defer wg.Done()
	lastExecutionTime := time.Now()
	minimumTimeBetweenEachExecution := time.Duration(math.Ceil(1e9 / (float64(throttle) / float64(workers))))
	for job := range jobChannel {
		timeUntilNextExecution := -(time.Since(lastExecutionTime) - minimumTimeBetweenEachExecution)
		if timeUntilNextExecution > 0 {
			time.Sleep(timeUntilNextExecution)
		}
		lastExecutionTime = time.Now()
		job()
		doneChannel <- struct{}{}
	}
}
