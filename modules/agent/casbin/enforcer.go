package casbin_agent

import (
	"strings"

	"github.com/casbin/casbin/persist"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	gorm_agent "github.com/oj-lab/platform/modules/agent/gorm"
	redis_agent "github.com/oj-lab/platform/modules/agent/redis"
	log_module "github.com/oj-lab/platform/modules/log"
)

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
		model, err := model.NewModelFromString(ExtendedRBACWithDomainModelString)
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
			log_module.AppLogger().Info("Casbin enforcer watcher initialized")
		}
		casbinEnforcer.AddFunction("keyMatchGin", keyMatchGinFunc)
		log_module.AppLogger().Info("Casbin enforcer initialized")
	}

	return casbinEnforcer
}
