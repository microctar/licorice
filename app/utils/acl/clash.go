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

var _ ACLReader = (*ClashDiverter)(nil)

type ClashDiverter struct {
	Offline                bool
	Ruleset                map[string][]string
	CustomProxyGroup       []map[string]any
	OverwriteOriginalRules bool
	EnableRuleGenerator    bool
}

func (d *ClashDiverter) ReadFile(basedir, path string) error {
	// absolute_path :string => filter configuration file
	// join path
	absolute_path := fmt.Sprintf("%s/%s", basedir, path)

	content, read_err := utils.ReadAll(absolute_path)

	if read_err != nil {
		return read_err
	}

	// regexp
	r_ruleset := regexp.MustCompile("ruleset=(.*)")
	r_cpg := regexp.MustCompile("custom_proxy_group=(.*)")
	r_erg := regexp.MustCompile("enable_rule_generator=(.*)")
	r_oor := regexp.MustCompile("overwrite_original_rules=(.*)")
	r_common_rule := regexp.MustCompile("(?m:^(DOMAIN|DOMAIN-(KEYWORD|SUFFIX)|PROCESS-NAME)(.*?)$)")
	r_unsupported_rule := regexp.MustCompile("(?m:^(USER-AGENT|URL-REGEX)(.*?)$)")
	r_noresolve := regexp.MustCompile("(?im:^(.*?),no-resolve$)")
	r_allrule := regexp.MustCompile("(?im:^[^\\#\\n].*$)")
	r_online := regexp.MustCompile("(?i:online)")
	r_final := regexp.MustCompile("(?i:final)")

	ruleset := r_ruleset.FindAllStringSubmatch(string(content), 64)
	cpg := r_cpg.FindAllStringSubmatch(string(content), 64)
	erp := r_erg.FindStringSubmatch(string(content))
	oor := r_oor.FindStringSubmatch(string(content))

	if !r_online.MatchString(absolute_path) {
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

					if r_final.MatchString(rule[1]) {
						rule[1] = strings.Replace(rule[1], "FINAL", "MATCH", -1)
					}

					kvpair = strings.Split(strings.Replace(rule[1], "[]", "", -1), ",")
					strategy := strings.Join(kvpair[1:], ",")

					MapMutex.Lock()
					d.Ruleset[kvpair[0]] = append(d.Ruleset[kvpair[0]], fmt.Sprintf("%s,%s", strategy, kvpair[0]))
					MapMutex.Unlock()

				} else {

					var rule_content string
					var get_rule_error error

					if d.Offline {
						// rule_file_path => absolute path of offline configuration file
						rule_file_path := fmt.Sprintf("%s/%s", basedir, kvpair[1])
						rule_content, get_rule_error = utils.ReadAll(rule_file_path)

					} else {
						rule_content, get_rule_error = utils.GetOnlineContent(kvpair[1])
					}

					if get_rule_error != nil {
						return get_rule_error
					}

					// remove unsupported rule
					rule_content = r_unsupported_rule.ReplaceAllString(rule_content, "\n")
					// "$0" => matched string
					rule_content = r_common_rule.ReplaceAllString(rule_content, fmt.Sprintf("$0,%s", kvpair[0]))
					// "$1" => matched substring
					rule_content = r_noresolve.ReplaceAllString(rule_content, fmt.Sprintf("$1,%s,no-resolve", kvpair[0]))

					all_rule := r_allrule.FindAllString(rule_content, 8192)

					MapMutex.Lock()
					d.Ruleset[kvpair[0]] = append(d.Ruleset[kvpair[0]], all_rule...)
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
					interval_str := strings.Split(grpinfo[4], ",")[0]
					interval, _ := strconv.Atoi(interval_str)
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

func (d *ClashDiverter) Expose() any {
	return d
}
