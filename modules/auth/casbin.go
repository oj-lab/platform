package auth

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/persist"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
	redis_agent "github.com/oj-lab/oj-lab-platform/modules/agent/redis"
	"github.com/oj-lab/oj-lab-platform/modules/log"
)

const ABACModelString = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub_rule, obj, act, eft

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = eval(p.sub_rule) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`

type CasbinSubject struct {
	Age  int
	Role string
}

var casbinEnforcer *casbin.SyncedCachedEnforcer

func GetDefaultCasbinEnforcer() *casbin.SyncedCachedEnforcer {
	if casbinEnforcer == nil {
		var err error
		var watcher persist.Watcher
		if len(redis_agent.RedisHosts) == 1 {
			watcher, err = rediswatcher.NewWatcher(redis_agent.RedisHosts[0], rediswatcher.WatcherOptions{})
			if err != nil {
				panic(err)
			}
		} else if len(redis_agent.RedisHosts) > 1 {
			addrs := strings.Join(redis_agent.RedisHosts, ",")
			watcher, err = rediswatcher.NewWatcherWithCluster(addrs, rediswatcher.WatcherOptions{})
			if err != nil {
				panic(err)
			}
		}

		adapter, err := gormadapter.NewAdapterByDB(gorm_agent.GetDefaultDB())
		if err != nil && adapter == nil {
			panic(err)
		}
		model, err := model.NewModelFromString(ABACModelString)
		if err != nil {
			panic(err)
		}
		casbinEnforcer, err = casbin.NewSyncedCachedEnforcer(model, adapter)
		if err != nil {
			panic(err)
		}
		if watcher != nil {
			err := casbinEnforcer.SetWatcher(watcher)
			if err != nil {
				panic(err)
			}
			err = watcher.SetUpdateCallback(rediswatcher.DefaultUpdateCallback(casbinEnforcer))
			if err != nil {
				panic(err)
			}
			log.AppLogger().Info("Casbin enforcer watcher initialized")
		}
		log.AppLogger().Info("Casbin enforcer initialized")
	}

	return casbinEnforcer
}

func LoadDefaultCasbinPolicies() error {
	enforcer := GetDefaultCasbinEnforcer()
	_, err := enforcer.AddPolicy(`r.sub.Age > 18 && r.sub.Age < 60`, `testData`, http.MethodGet, "allow")
	if err != nil {
		return err
	}
	_, err = enforcer.AddPolicy(`r.sub.Role == 'admin'`, `adminRequired`,
		strings.Join([]string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		}, "|"), "allow")
	if err != nil {
		return err
	}
	err = enforcer.SavePolicy()
	if err != nil {
		return err
	}
	return nil
}
