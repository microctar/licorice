package acl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/microctar/licorice/app/utils"
	"golang.org/x/sync/errgroup"
)

var _ ACLReader = (*clashDiverter)(nil)

type clashDiverter struct {
	Offline                bool
	Ruleset                map[string][]string
	CustomProxyGroup       []map[string]any
	OverwriteOriginalRules bool
	EnableRuleGenerator    bool
}

func (d *clashDiverter) ReadFile(basedir, path string) error {
	// absPath :string => filter configuration file
	// join path
	absPath := fmt.Sprintf("%s/%s", basedir, path)

	content, readErr := utils.ReadAll(absPath)

	if readErr != nil {
		return readErr
	}

	// regexp
	rRuleset := regexp.MustCompile("ruleset=(.*)")
	rCpg := regexp.MustCompile("custom_proxy_group=(.*)")
	rErg := regexp.MustCompile("enable_rule_generator=(.*)")
	rOor := regexp.MustCompile("overwrite_original_rules=(.*)")
	rCommonRule := regexp.MustCompile("(?m:^(DOMAIN|DOMAIN-(KEYWORD|SUFFIX)|PROCESS-NAME)(.*?)$)")
	rUnsupportedRule := regexp.MustCompile("(?m:^(USER-AGENT|URL-REGEX)(.*?)$)")
	rNoresolve := regexp.MustCompile("(?im:^(.*?),no-resolve$)")
	rAllrule := regexp.MustCompile("(?im:^[^\\#\\n].*$)")
	rOnline := regexp.MustCompile("(?i:online)")
	rFinal := regexp.MustCompile("(?i:final)")

	ruleset := rRuleset.FindAllStringSubmatch(string(content), 64)
	cpg := rCpg.FindAllStringSubmatch(string(content), 64)
	erp := rErg.FindStringSubmatch(string(content))
	oor := rOor.FindStringSubmatch(string(content))

	if !rOnline.MatchString(absPath) {
		d.Offline = true
	}

	if ruleset != nil {
		// initialize
		d.Ruleset = make(map[string][]string)
		errgrp := new(errgroup.Group)
		MapMutex := sync.RWMutex{}

		// rule :[]string => ["matched_string", "rulename,location of rule file"]
		// rule :[]string => ["matched_string", "rulename,[]groupname"]
		for _, rule := range ruleset {

			rule := rule

			errgrp.Go(func() error {

				// kvpair :[]string => ["rulename", "location of rule file"]
				kvpair := strings.Split(rule[1], ",")

				if strings.Contains(rule[1], "[]") {

					if rFinal.MatchString(rule[1]) {
						rule[1] = strings.Replace(rule[1], "FINAL", "MATCH", -1)
					}

					kvpair = strings.Split(strings.Replace(rule[1], "[]", "", -1), ",")
					strategy := strings.Join(kvpair[1:], ",")

					MapMutex.Lock()
					d.Ruleset[kvpair[0]] = append(d.Ruleset[kvpair[0]], fmt.Sprintf("%s,%s", strategy, kvpair[0]))
					MapMutex.Unlock()

				} else {

					var ruleContent string
					var readRuleErr error

					if d.Offline {
						// ruleFilePath => absolute path of offline configuration file
						ruleFilePath := fmt.Sprintf("%s/%s", basedir, kvpair[1])
						ruleContent, readRuleErr = utils.ReadAll(ruleFilePath)

					} else {
						ruleContent, readRuleErr = utils.GetOnlineContent(kvpair[1])
					}

					if readRuleErr != nil {
						return readRuleErr
					}

					// remove unsupported rule
					ruleContent = rUnsupportedRule.ReplaceAllString(ruleContent, "\n")
					// "$0" => matched string
					ruleContent = rCommonRule.ReplaceAllString(ruleContent, fmt.Sprintf("$0,%s", kvpair[0]))
					// "$1" => matched substring
					ruleContent = rNoresolve.ReplaceAllString(ruleContent, fmt.Sprintf("$1,%s,no-resolve", kvpair[0]))

					allRule := rAllrule.FindAllString(ruleContent, 8192)

					MapMutex.Lock()
					d.Ruleset[kvpair[0]] = append(d.Ruleset[kvpair[0]], allRule...)
					MapMutex.Unlock()
				}

				return nil

			})
		}

		if gerr := errgrp.Wait(); gerr != nil {
			return gerr
		}
	}

	if cpg != nil {
		wg := new(sync.WaitGroup)
		MapMutex := sync.RWMutex{}

		for _, group := range cpg {
			wg.Add(1)
			group := group

			go func() {
				defer wg.Done()
				// initialize
				cpgrp := make(map[string]any)
				// grpinfo :[]string => ["group_name", "type", "group", ...]
				// grpinfo :[]string => ["group_name", "type", ".*", "url", "interval_time,,interval_time"]
				grpinfo := strings.Split(strings.Replace(group[1], "[]", "", -1), "`")
				cpgrp["name"] = grpinfo[0]
				cpgrp["type"] = grpinfo[1]

				// grpinfo[1] => type :string
				if grpinfo[1] != "url-test" {
					cpgrp["proxies"] = append([]string{}, grpinfo[2:]...)
				} else {
					// grpinfo[3] => "url"
					cpgrp["url"] = grpinfo[3]
					// grpinfo[4] => "interval,,xx"
					intervalStr := strings.Split(grpinfo[4], ",")[0]
					interval, _ := strconv.Atoi(intervalStr)
					cpgrp["interval"] = interval
					// grpinfo[2] => ".*"
					cpgrp["proxies"] = append([]string{}, grpinfo[2])
				}

				MapMutex.Lock()
				defer MapMutex.Unlock()
				d.CustomProxyGroup = append(d.CustomProxyGroup, cpgrp)
			}()
		}

		wg.Wait()
	}

	if erp != nil {
		// erp :[]string => ["matched_string", "boolean"]
		d.EnableRuleGenerator = utils.Str2Bool(erp[1])
	}

	if oor != nil {
		// oor :[]string => ["matched_string", "boolean"]
		d.OverwriteOriginalRules = utils.Str2Bool(oor[1])
	}

	return nil
}

func (d *clashDiverter) Expose() any {
	return (*ClashDiverter)(d)
}
