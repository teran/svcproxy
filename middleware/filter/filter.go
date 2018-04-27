package filter

import (
	"net"
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/teran/svcproxy/middleware/types"
)

var _ types.Middleware = &Filter{}

// Filter middleware type
type Filter struct {
	options *Options
}

// Options type to care about filter middleware options without
// using inerface{} types
type Options struct {
	Name  string
	Rules []Rule
}

// Rule type
type Rule struct {
	Logic      string
	IPs        []string
	UserAgents []*regexp.Regexp
}

// NewMiddleware returns new Middleware instance
func NewMiddleware() *Filter {
	return &Filter{}
}

// SetOptions sets passed options for middleware at startup time(i.e. Chaining procedure)
func (f *Filter) SetOptions(options map[string]interface{}) {
	var rules []Rule

	r, ok := options["rules"]
	if ok {
		for _, x := range r.([]interface{}) {
			rule := Rule{}
			ips, ok := x.(map[interface{}]interface{})["ips"]
			if ok {
				var ipList []string
				for _, i := range ips.([]interface{}) {
					ipList = append(ipList, i.(string))
				}
				rule.IPs = ipList
			}

			useragents, ok := x.(map[interface{}]interface{})["userAgents"]
			if ok {
				var uaList []*regexp.Regexp
				for _, ua := range useragents.([]interface{}) {
					uaStr := ua.(string)
					pattern, err := regexp.Compile(uaStr)
					if err != nil {
						log.Fatalf("Error compiling regexp: %s", uaStr)
					}
					uaList = append(uaList, pattern)
				}
				rule.UserAgents = uaList
			}

			rules = append(rules, rule)
		}
	}

	f.options = &Options{
		Name:  options["name"].(string),
		Rules: rules,
	}
}

// Middleware implements middleware filter middleware
func (f *Filter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.UserAgent()
		addr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Warnf("Error parsing remote addr: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		for _, x := range f.options.Rules {
			if isIPBlacklisted(addr, x.IPs) {
				http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
				return
			}
			if isUserAgentBlacklisted(userAgent, x.UserAgents) {
				http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

//
func isIPBlacklisted(sourceIP string, filterList []string) bool {
	for _, ip := range filterList {
		if sourceIP == ip {
			return true
		}
	}
	return false
}

func isUserAgentBlacklisted(sourceUserAgent string, userAgentFilterList []*regexp.Regexp) bool {
	for _, userAgent := range userAgentFilterList {
		if userAgent.MatchString(sourceUserAgent) {
			return true
		}
	}
	return false
}
