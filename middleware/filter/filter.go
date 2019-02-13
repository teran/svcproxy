package filter

import (
	"net"
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/teran/svcproxy/middleware/types"
)

var _ types.Middleware = (*Filter)(nil)

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
	Logic          string
	AllowFrom      []*net.IPNet
	DenyFrom       []*net.IPNet
	DenyUserAgents []*regexp.Regexp
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
			allowFrom, ok := x.(map[interface{}]interface{})["allowFrom"]
			if ok {
				for _, i := range allowFrom.([]interface{}) {
					_, ipnet, err := net.ParseCIDR(i.(string))
					if err != nil {
						log.Printf("Error parsing CIDR: %s. Rule skipped.", err)
						continue
					}
					rule.AllowFrom = append(rule.AllowFrom, ipnet)
				}
			}

			denyFrom, ok := x.(map[interface{}]interface{})["denyFrom"]
			if ok {
				for _, i := range denyFrom.([]interface{}) {
					_, ipnet, err := net.ParseCIDR(i.(string))
					if err != nil {
						log.Printf("Error parsing CIDR: %s. Rule skipped.", err)
						continue
					}
					rule.DenyFrom = append(rule.DenyFrom, ipnet)
				}
			}

			useragents, ok := x.(map[interface{}]interface{})["denyUserAgents"]
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
				rule.DenyUserAgents = uaList
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
		addrString, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Warnf("Error parsing remote addr: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		addr := net.ParseIP(addrString)

		if f.isUserAgentDenied(userAgent) || f.isIPDenied(addr) || !f.isIPAllowed(addr) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (f *Filter) isUserAgentDenied(userAgent string) bool {
	for _, rule := range f.options.Rules {
		for _, ua := range rule.DenyUserAgents {
			if ua.MatchString(userAgent) {
				return true
			}
		}
	}
	return false
}

func (f *Filter) isIPDenied(addr net.IP) bool {
	for _, rule := range f.options.Rules {
		for _, ip := range rule.DenyFrom {
			if ip.Contains(addr) {
				return true
			}
		}
	}
	return false
}

func (f *Filter) isIPAllowed(addr net.IP) bool {
	var totalRules int
	for _, rule := range f.options.Rules {
		for _, ip := range rule.AllowFrom {
			totalRules++
			if ip.Contains(addr) {
				return true
			}
		}
	}

	if totalRules == 0 {
		return true
	}
	return false
}
