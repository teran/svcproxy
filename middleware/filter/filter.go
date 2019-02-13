package filter

import (
	"errors"
	"net"
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/teran/svcproxy/middleware/types"
)

var _ types.Middleware = (*Filter)(nil)

// InputRuleSet type
type InputRuleSet = []InputRule

// InputRule type
type InputRule = map[string][]string

// Config type
type Config struct {
	Name  string
	Rules []Rule
}

// Unpack unpacks input configuration
func (fc *Config) Unpack(options map[string]interface{}) error {
	var rules []Rule

	r, ok := options["rules"]
	if !ok {
		log.Printf("no rules defined. Skipping configuration")
		return nil
	}

	rulesInput, ok := r.(InputRuleSet)
	if !ok {
		return errors.New("improper configuration: input rule set doesn't look so. Should be []map[string][]string")
	}

	for _, x := range rulesInput {
		rule := Rule{}
		allowFrom, ok := x["allowFrom"]
		if ok {
			for _, i := range allowFrom {
				_, ipnet, err := net.ParseCIDR(i)
				if err != nil {
					log.Printf("Error parsing CIDR: %s. Rule skipped.", err)
					continue
				}
				rule.AllowFrom = append(rule.AllowFrom, ipnet)
			}
		}

		denyFrom, ok := x["denyFrom"]
		if ok {
			for _, i := range denyFrom {
				_, ipnet, err := net.ParseCIDR(i)
				if err != nil {
					log.Printf("Error parsing CIDR: %s. Rule skipped.", err)
					continue
				}
				rule.DenyFrom = append(rule.DenyFrom, ipnet)
			}
		}

		useragents, ok := x["denyUserAgents"]
		if ok {
			var uaList []*regexp.Regexp
			for _, ua := range useragents {
				uaStr := ua
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

	fc.Name = options["name"].(string)
	fc.Rules = rules

	return nil
}

// Filter middleware type
type Filter struct {
	config *Config
}

// Rule type
type Rule struct {
	Logic          string
	AllowFrom      []*net.IPNet
	DenyFrom       []*net.IPNet
	DenyUserAgents []*regexp.Regexp
}

// NewMiddleware returns new Middleware instance
func NewMiddleware() types.Middleware {
	return &Filter{}
}

// SetConfig applies config to the middleware
func (f *Filter) SetConfig(c types.MiddlewareConfig) error {
	var ok bool
	f.config, ok = c.(*Config)
	if !ok {
		return errors.New("the map passed doesn't implement FilterConfig")
	}
	return nil
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
	for _, rule := range f.config.Rules {
		for _, ua := range rule.DenyUserAgents {
			if ua.MatchString(userAgent) {
				return true
			}
		}
	}
	return false
}

func (f *Filter) isIPDenied(addr net.IP) bool {
	for _, rule := range f.config.Rules {
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
	for _, rule := range f.config.Rules {
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
