/**************************************************************************************************/
// File: wsusscn2cli.go
// Author: Jon Smith
// Copyright: Hash Authority, LLC 2018
// Description: Command-line tool using wsusscn2.cab API
// Versions:
// 0.1.0: Initial release with basic commands
// 0.1.1: Added update parameters
// 0.1.2: Added update boolean parameters
// 0.1.3: Added columns, limit, and offset
// 0.1.4: Added update_creation_date filters
/**************************************************************************************************/
package main

import (
	"encoding/json"     //api
	"errors"            //new error
	"fmt"               //printing
	"io/ioutil"         //writing to file
	"log"               //logging
	"net/http"          //http client
	"net/http/httputil" //http debug
	"os"                //testing existence of a file
	"path/filepath"     //splitting paths
	"regexp"            //include/exclude pattern matching
	"strconv"           //parsing boolean
	"strings"
	"time" //saving current time to marks.json

	"github.com/urfave/cli" //cli structure
)

/**************************************************************************************************/
/*                                                                                                */
/*                                           CONSTANTS                                            */
/*                                                                                                */
/**************************************************************************************************/

/**************************************************************************************************/
/*                                                                                                */
/*                                             TYPES                                              */
/*                                                                                                */
/**************************************************************************************************/
// Update: Structure for update records
type Update struct {
	Bundles             string `json:"bundles"`
	ClassificationTitle string `json:"classification_title"`
	CompanyTitle        string `json:"company_title"`
	Description         string `json:"description"`
	InstallBehavior     string `json:"install_behavior"`
	IsBeta              string `json:"is_beta"`
	IsBundled           string `json:"is_bundled"`
	IsPublic            string `json:"is_public"`
	IsSuperseded        string `json:"is_superseded"`
	Kb                  string `json:"kb"`
	Language            string `json:"language"`
	MoreInfoUrl         string `json:"more_info_url"`
	MsrcSeverity        string `json:"msrc_severity"`
	ProductFamilyTitle  string `json:"product_family_title"`
	ProductTitle        string `json:"product_title"`
	PublicationState    string `json:"publication_state"`
	Readiness           string `json:"readiness"`
	Supersedes          string `json:"supersedes"`
	SupportUrl          string `json:"support_url"`
	UninstallBehavior   string `json:"uninstall_behavior"`
	UninstallNotes      string `json:"uninstall_notes"`
	UpdateCreationDate  string `json:"update_creation_date"`
	UpdateRevision      string `json:"update_revision"`
	UpdateTitle         string `json:"update_title"`
	UpdateType          string `json:"update_type"`
	UpdateUid           string `json:"update_uid"`
}

// Classification: Structure for classification records
type Classification struct {
	ClassificationUid      string `json:"classification_uid"`
	ClassificationRevision string `json:"classification_revision"`
	ClassificationTitle    string `json:"classification_title"`
}

// Product: Structure for product records
type Product struct {
	ProductUid      string `json:"product_uid"`
	ProductRevision string `json:"product_revision"`
	ProductTitle    string `json:"product_title"`
}

// ProductFamily: Structure for product family records
type ProductFamily struct {
	ProductFamilyUid      string `json:"product_family_uid"`
	ProductFamilyRevision string `json:"product_family_revision"`
	ProductFamilyTitle    string `json:"product_family_title"`
}

type wConfig struct {
	ApiServer string `json:"api_server"`
	ApiPort   string `json:"api_port"`
	ApiKey    string `json:"api_key"`
}

/**************************************************************************************************/
/*                                                                                                */
/*                                           FUNCTIONS                                            */
/*                                                                                                */
/**************************************************************************************************/
// check logs and exits if an error occurs
func check(e error) {
	if e != nil {
		log.Fatalf("%s", e)
	}
}

func strToBool(str string) bool {
	b, err := strconv.ParseBool(str)
	check(err)
	return b
}

func strToMap(input string) map[string]int {
	m := make(map[string]int)
	for _, v := range strings.Split(input, ",") {
		trimmedV := strings.ToLower(strings.TrimSpace(v))
		m[trimmedV] = 1
	}
	return m
}

func strToSlice(input string) []string {
	inputSlice := strings.Split(input, ",")
	for i, v := range inputSlice {
		inputSlice[i] = strings.ToLower(strings.TrimSpace(v))
	}
	return inputSlice
}

// pathMatch returns true if any pattern matches
func pathMatch(matches []string, path string) bool {
	for _, v := range matches {
		matched, _ := regexp.MatchString(v, path)
		if matched {
			return true
		}
	}
	return false
}

func createNewHttpReq(url string, key string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	req.SetBasicAuth("u", key)
	return req
}

func getJson(c *http.Client, req *http.Request, debug bool, target interface{}) error {
	log.Println("GET " + req.URL.String())

	if debug {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump))
	}

	r, err := c.Do(req)

	if err != nil {
		return err
	}

	if debug {
		responseDump, err := httputil.DumpResponse(r, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(responseDump))
	}

	if r.StatusCode == http.StatusUnauthorized {
		return errors.New("Unauthorized request to service")
	}

	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// readConfig reads in configuration items
func readConfig(file string) wConfig {
	var c = wConfig{}
	_, err := os.Stat(file)
	if err == nil {
		f, _ := os.Open(file)
		defer f.Close()
		decoder := json.NewDecoder(f)
		err := decoder.Decode(&c)
		check(err)
	} else {
		if !os.IsNotExist(err) {
			log.Fatalf("Unable to read %s: %s\n", file, err)
		}
	}
	return c
}

/**************************************************************************************************/
/*                                                                                                */
/*                                             MAIN                                               */
/*                                                                                                */
/**************************************************************************************************/
func main() {
	var apiKey string //userinput, required
	var apiUrl string //url to wsusscn2 api
	//var execPath string //path to the go executable
	var debug bool     //debug logging
	var countOnly bool //only print number of records returned

	var productTitle []string
	var updateUid []string
	var updateTitle []string
	var updateKb []string
	var updateType []string
	var productFamilyTitle []string
	var classificationTitle []string
	var msrcSeverity []string
	var columns string

	var limit int
	var offset int
	var recordLimit int

	defaultLimit := 1000
	defaultOffset := 0
	defaultRecordLimit := 20000

	var defaultUpdateColumns = "update_uid, kb, update_title, update_creation_date, product_title, product_family_title, update_type, is_superseded, classification_title, company_title, description, install_behavior, is_beta, is_bundled, is_public, language, more_info_url, msrc_severity, publication_state, readiness, support_url, uninstall_behavior, uninstall_notes, update_revision"
	var defaultUpdateColumnsTitle = make(map[string]string)

	for _, v := range strings.Split(defaultUpdateColumns, ",") {
		c := strings.ToLower(strings.TrimSpace(v))
		words := strings.Split(c, "_")
		newTitle := ""
		for _, w := range words {
			newTitle += strings.Title(w)
		}
		defaultUpdateColumnsTitle[c] = newTitle
	}

	var isSuperseded string
	var isBundled string
	var isBeta string
	var isPublic string

	var updateCreationDateAfter string
	var updateCreationDateBefore string
	var updateCreationDateOn string

	apiUrl = "https://wsusscn2.cab:443"

	api := &http.Client{Timeout: 30 * time.Second}

	//setup
	ex, err := os.Executable()
	check(err)
	execPath := filepath.Dir(ex)

	config := readConfig(execPath + "/wsusscn2cli.json")

	app := cli.NewApp()
	app.Name = "wsusscn2cli"
	app.Version = "0.1.4"
	app.Usage = "wsusscn2.cab integration"
	app.Copyright = "(c) 2018 Hash Authority, LLC"
	app.Commands = []cli.Command{
		{
			Name:  "listclassification",
			Usage: "List all classifications",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "debug, d",
					Usage:       "Output debug level logging",
					Destination: &debug,
				},
				cli.StringFlag{
					Name:        "api_key, a",
					Usage:       "API key (required if not using config file)",
					Destination: &apiKey,
				},
			},
			Action: func(c *cli.Context) error {
				log.Println("List classification called")

				//Authentication setup
				if apiKey == "" && config.ApiKey == "" {
					log.Fatalf("Unable to find api key. use api_key or set one using wsusscn2cli setapikey ASDF")
				}

				if apiKey == "" {
					apiKey = config.ApiKey
				}

				var classification []Classification

				req := createNewHttpReq(apiUrl+"/classification", apiKey)

				err := getJson(api, req, debug, &classification)
				check(err)
				fmt.Println(`"ClassificationUid","ClassificationRevision","ClassificationTitle"`)
				for _, v := range classification {
					fmt.Printf("\"%s\",\"%s\",\"%s\"\n", v.ClassificationUid, v.ClassificationRevision, v.ClassificationTitle)
				}

				return nil
			},
		},
		{
			Name:  "listproduct",
			Usage: "List all products",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "debug, d",
					Usage:       "Output debug level logging",
					Destination: &debug,
				},
				cli.StringFlag{
					Name:        "api_key, a",
					Usage:       "API key (required if not using config file)",
					Destination: &apiKey,
				},
			},
			Action: func(c *cli.Context) error {
				log.Println("List product called")

				//Authentication setup
				if apiKey == "" && config.ApiKey == "" {
					log.Fatalf("Unable to find api key. use api_key or set one using wsusscn2cli setapikey ASDF")
				}

				if apiKey == "" {
					apiKey = config.ApiKey
				}

				var product []Product

				req := createNewHttpReq(apiUrl+"/product", apiKey)

				err := getJson(api, req, debug, &product)
				check(err)
				fmt.Println(`"ProductUid","ProductRevision","ProductTitle"`)
				for _, v := range product {
					fmt.Printf("\"%s\",\"%s\",\"%s\"\n", v.ProductUid, v.ProductRevision, v.ProductTitle)
				}

				return nil
			},
		},
		{
			Name:  "listproductfamily",
			Usage: "List all product families",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "debug, d",
					Usage:       "Output debug level logging",
					Destination: &debug,
				},
				cli.StringFlag{
					Name:        "api_key, a",
					Usage:       "API key (required if not using config file)",
					Destination: &apiKey,
				},
			},
			Action: func(c *cli.Context) error {
				log.Println("List productfamily called")

				//Authentication setup
				if apiKey == "" && config.ApiKey == "" {
					log.Fatalf("Unable to find api key. use api_key or set one using wsusscn2cli setapikey ASDF")
				}

				if apiKey == "" {
					apiKey = config.ApiKey
				}

				var productfamily []ProductFamily

				req := createNewHttpReq(apiUrl+"/productfamily", apiKey)

				err := getJson(api, req, debug, &productfamily)
				check(err)
				fmt.Println(`"ProductFamilyUid","ProductFamilyRevision","ProductFamilyTitle"`)
				for _, v := range productfamily {
					fmt.Printf("\"%s\",\"%s\",\"%s\"\n", v.ProductFamilyUid, v.ProductFamilyRevision, v.ProductFamilyTitle)
				}

				return nil
			},
		},
		{
			Name:  "listupdate",
			Usage: "List updates",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "debug, d",
					Usage:       "Output debug level logging",
					Destination: &debug,
				},
				cli.StringFlag{
					Name:        "api_key, a",
					Usage:       "API key (required if not using config file)",
					Destination: &apiKey,
				},
				cli.BoolFlag{
					Name:        "count_only",
					Usage:       "Only print number of records",
					Destination: &countOnly,
				},
				cli.StringSliceFlag{
					Name:  "product_title",
					Usage: "Name of product.",
				},
				cli.StringSliceFlag{
					Name:  "update_uid",
					Usage: "Update Uid.",
				},
				cli.StringSliceFlag{
					Name:  "update_title",
					Usage: "Update Title.",
				},
				cli.StringSliceFlag{
					Name:  "kb",
					Usage: "Update KB.",
				},
				cli.StringSliceFlag{
					Name:  "update_type",
					Usage: "Update Type.",
				},
				cli.StringSliceFlag{
					Name:  "product_family_title",
					Usage: "Product Family Title.",
				},
				cli.StringSliceFlag{
					Name:  "classification_title",
					Usage: "Classification Title.",
				},
				cli.StringSliceFlag{
					Name:  "msrc_severity",
					Usage: "MSRC Severity.",
				},
				cli.StringFlag{
					Name:        "is_superseded",
					Usage:       "Is Superseded.",
					Destination: &isSuperseded,
				},
				cli.StringFlag{
					Name:        "is_bundled",
					Usage:       "Is Bundled.",
					Destination: &isBundled,
				},
				cli.StringFlag{
					Name:        "is_public",
					Usage:       "Is Public.",
					Destination: &isPublic,
				},
				cli.StringFlag{
					Name:        "is_beta",
					Usage:       "Is Beta.",
					Destination: &isBeta,
				},
				cli.StringFlag{
					Name:        "update_creation_date_after",
					Usage:       "Updates created after this date [YYYY-MM-DD] (exclusive).",
					Destination: &updateCreationDateAfter,
				},
				cli.StringFlag{
					Name:        "update_creation_date_before",
					Usage:       "Updates created before this date [YYYY-MM-DD] (exclusive).",
					Destination: &updateCreationDateBefore,
				},
				cli.StringFlag{
					Name:        "update_creation_date_on",
					Usage:       "Updates created on this date [YYYY-MM-DD].",
					Destination: &updateCreationDateOn,
				},
				cli.StringFlag{
					Name:        "columns",
					Usage:       "Restrict output to listed columns.",
					Destination: &columns,
				},
				cli.IntFlag{
					Name:        "limit",
					Usage:       "Number of records per page.",
					Value:       defaultLimit,
					Destination: &limit,
				},
				cli.IntFlag{
					Name:        "offset",
					Usage:       "Number of records to skip.",
					Value:       defaultOffset,
					Destination: &offset,
				},
				cli.IntFlag{
					Name:        "record_limit",
					Usage:       "Max number of records to return.",
					Value:       defaultRecordLimit,
					Destination: &recordLimit,
				},
			},
			Action: func(c *cli.Context) error {
				log.Println("List update called")

				//Authentication setup
				if apiKey == "" && config.ApiKey == "" {
					log.Fatalf("Unable to find api key. use api_key or set one using wsusscn2cli setapikey ASDF")
				}

				if apiKey == "" {
					apiKey = config.ApiKey
				}

				productTitle = c.StringSlice("product_title")
				updateUid = c.StringSlice("update_uid")
				updateTitle = c.StringSlice("update_title")
				updateKb = c.StringSlice("kb")
				updateType = c.StringSlice("update_type")
				productFamilyTitle = c.StringSlice("product_family_title")
				classificationTitle = c.StringSlice("classification_title")
				msrcSeverity = c.StringSlice("msrc_severity")

				if columns == "" {
					columns = defaultUpdateColumns
				}

				columnFilter := strToSlice(columns)

				recordCnt := 0
				done := false

				if limit <= 0 {
					limit = defaultLimit
				}

				if offset <= 0 {
					offset = defaultOffset
				}

				if recordLimit <= 0 {
					recordLimit = defaultRecordLimit
				}

				if recordLimit < limit {
					limit = recordLimit
				}

				for recordCnt < recordLimit && !done {
					var update []Update

					req := createNewHttpReq(apiUrl+"/update", apiKey)

					q := req.URL.Query()
					q.Add("limit", strconv.Itoa(limit))
					q.Add("offset", strconv.Itoa(offset))

					if len(productTitle) > 0 {
						for _, p := range productTitle {
							q.Add("product_title", p)
						}
					}
					if len(updateUid) > 0 {
						for _, p := range updateUid {
							q.Add("uid", p)
						}
					}
					if len(updateTitle) > 0 {
						for _, p := range updateTitle {
							q.Add("title", p)
						}
					}
					if len(updateKb) > 0 {
						for _, p := range updateKb {
							q.Add("kb", p)
						}
					}
					if len(updateType) > 0 {
						for _, p := range updateType {
							q.Add("type", p)
						}
					}
					if len(productFamilyTitle) > 0 {
						for _, p := range productFamilyTitle {
							q.Add("product_family_title", p)
						}
					}
					if len(classificationTitle) > 0 {
						for _, p := range classificationTitle {
							q.Add("classification_title", p)
						}
					}
					if len(msrcSeverity) > 0 {
						for _, p := range msrcSeverity {
							q.Add("msrc_severity", p)
						}
					}

					if isSuperseded != "" {
						q.Add("is_superseded", isSuperseded)
					}

					if isBundled != "" {
						q.Add("is_bundled", isBundled)
					}

					if isPublic != "" {
						q.Add("is_public", isPublic)
					}

					if isBeta != "" {
						q.Add("is_beta", isBeta)
					}

					layout := "2006-01-02"
					if updateCreationDateAfter != "" {
						_, err := time.Parse(layout, updateCreationDateAfter)
						if err != nil {
							log.Fatal(err)
							log.Println("Unable to parse provided date. Expected: YYYY-MM-DD. Found %s", updateCreationDateAfter)
							return nil
						}
						q.Add("update_creation_date_after", updateCreationDateAfter)
					}

					if updateCreationDateBefore != "" {
						_, err := time.Parse(layout, updateCreationDateBefore)
						if err != nil {
							log.Fatal(err)
							log.Println("Unable to parse provided date. Expected: YYYY-MM-DD. Found %s", updateCreationDateBefore)
							return nil
						}
						q.Add("update_creation_date_before", updateCreationDateBefore)
					}

					if updateCreationDateOn != "" {
						if strings.ToLower(updateCreationDateOn) != "today" {
							_, err := time.Parse(layout, updateCreationDateBefore)
							if err != nil {
								log.Fatal(err)
								log.Println("Unable to parse provided date. Expected: YYYY-MM-DD. Found %s", updateCreationDateOn)
								return nil
							}
						}
						q.Add("update_creation_date_on", updateCreationDateOn)
					}

					req.URL.RawQuery = q.Encode()

					err := getJson(api, req, debug, &update)
					check(err)

					curRecordCnt := len(update)

					if countOnly {
						//do nothing here
					} else {
						firstCol := true
						for _, c := range columnFilter {
							if val, ok := defaultUpdateColumnsTitle[c]; ok {
								if firstCol {
									firstCol = false
								} else {
									fmt.Printf(",")
								}
								fmt.Printf("\"%s\"", val)
							}
						}
						fmt.Printf("\n")

						for _, v := range update {
							firstCol = true
							for _, c := range columnFilter {
								if firstCol {
									firstCol = false
								} else {
									fmt.Printf(",")
								}

								switch c {
								case "classification_title":
									fmt.Printf("\"%s\"", v.ClassificationTitle)
								case "company_title":
									fmt.Printf("\"%s\"", v.CompanyTitle)
								case "description":
									fmt.Printf("\"%s\"", v.Description)
								case "install_behavior":
									fmt.Printf("\"%s\"", v.InstallBehavior)
								case "is_beta":
									fmt.Printf("\"%s\"", v.IsBeta)
								case "is_bundled":
									fmt.Printf("\"%s\"", v.IsBundled)
								case "is_public":
									fmt.Printf("\"%s\"", v.IsPublic)
								case "is_superseded":
									fmt.Printf("\"%s\"", v.IsSuperseded)
								case "kb":
									fmt.Printf("\"%s\"", v.Kb)
								case "language":
									fmt.Printf("\"%s\"", v.Language)
								case "more_info_url":
									fmt.Printf("\"%s\"", v.MoreInfoUrl)
								case "msrc_severity":
									fmt.Printf("\"%s\"", v.MsrcSeverity)
								case "product_family_title":
									fmt.Printf("\"%s\"", v.ProductFamilyTitle)
								case "product_title":
									fmt.Printf("\"%s\"", v.ProductTitle)
								case "publication_state":
									fmt.Printf("\"%s\"", v.PublicationState)
								case "readiness":
									fmt.Printf("\"%s\"", v.Readiness)
								case "support_url":
									fmt.Printf("\"%s\"", v.SupportUrl)
								case "uninstall_behavior":
									fmt.Printf("\"%s\"", v.UninstallBehavior)
								case "uninstall_notes":
									fmt.Printf("\"%s\"", v.UninstallNotes)
								case "update_creation_date":
									fmt.Printf("\"%s\"", v.UpdateCreationDate)
								case "update_revision":
									fmt.Printf("\"%s\"", v.UpdateRevision)
								case "update_title":
									fmt.Printf("\"%s\"", v.UpdateTitle)
								case "update_type":
									fmt.Printf("\"%s\"", v.UpdateType)
								case "update_uid":
									fmt.Printf("\"%s\"", v.UpdateUid)
								}
							}
							fmt.Printf("\n")
						}
					}

					recordCnt += curRecordCnt
					offset += curRecordCnt

					if curRecordCnt == 0 {
						log.Println("No more records returned")
						done = true
					}
					if curRecordCnt < limit {
						log.Println("Last page of records reached")
						done = true
					}
				}

				if countOnly {
					fmt.Printf("Number of records: %d\n", recordCnt)
				}

				return nil
			},
		},
		{
			Name:  "setapikey",
			Usage: "Set API key for repeated usage",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "debug, d",
					Usage:       "Output debug level logging",
					Destination: &debug,
				},
				cli.StringFlag{
					Name:        "api_key, a",
					Usage:       "Authentication to API",
					Destination: &apiKey,
				},
			},
			Action: func(c *cli.Context) error {
				log.Println("Set API Key called")

				if apiKey == "" {
					log.Fatalf("--api_key argument is blank. Pass a valid api_key argument")
				}

				config = wConfig{}
				config.ApiKey = apiKey

				configJson, _ := json.Marshal(config)
				err = ioutil.WriteFile(execPath+"/wsusscn2cli.json", configJson, 0644)
				check(err)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
