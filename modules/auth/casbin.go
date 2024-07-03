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

const RBACModelString = `
[request_definition]
r = sub, info, dom, obj, act

[policy_definition]
p = sub, info_rule, dom, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && eval(p.info_rule) && r.dom == p.dom && r.obj == p.obj && regexMatch(r.act, p.act)
`

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
		model, err := model.NewModelFromString(RBACModelString)
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
	_, err := enforcer.AddPolicy(`admin`, `true`, `system`, `testData`, http.MethodGet, "allow")
	if err != nil {
		return err
	}
	_, err = enforcer.AddPolicy(`admin`, `true`, `system`, `adminRequired`,
		strings.Join([]string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		}, "|"), "allow")
	if err != nil {
		return err
	}
	_, err = enforcer.AddGroupingPolicy(`test_user`, `admin`)
	if err != nil {
		return err
	}
	err = enforcer.SavePolicy()
	if err != nil {
		return err
	}
	return nil
}
