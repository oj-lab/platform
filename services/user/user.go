package user_service

import (
	"context"
	"fmt"

	judge_model "github.com/oj-lab/platform/models/judge"
	user_model "github.com/oj-lab/platform/models/user"
	casbin_agent "github.com/oj-lab/platform/modules/agent/casbin"
	gorm_agent "github.com/oj-lab/platform/modules/agent/gorm"
	auth_module "github.com/oj-lab/platform/modules/auth"
	log_module "github.com/oj-lab/platform/modules/log"
)

func CreateUser(ctx context.Context, request user_model.User) (*user_model.User, error) {
	db := gorm_agent.GetDefaultDB()

	user, err := user_model.CreateUser(db, request)
	if err != nil {
		return nil, err
	}

	_, err = judge_model.CreateJudgeRankCache(db, judge_model.NewJudgeRankCache(request.Account))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUser(ctx context.Context, account string) (*user_model.User, error) {
	db := gorm_agent.GetDefaultDB()
	user, err := user_model.GetUser(db, account)
	if err != nil {
		return nil, err
	}
	enforcer := casbin_agent.GetDefaultCasbinEnforcer()
	if enforcer != nil {
		user.Roles = enforcer.GetRolesForUserInDomain(casbin_agent.UserSubjectPrefix+account, "system")
		for i, role := range user.Roles {
			user.Roles[i] = role[len(casbin_agent.RoleSubjectPrefix):]
		}
	}

	return user, nil
}

// TODO: Tidy this function
func DeleteUser(ctx context.Context, account string) error {
	db := gorm_agent.GetDefaultDB()

	judges, err := judge_model.GetJudgeListByOptions(db, judge_model.GetJudgeOptions{
		UserAccount: account,
	})
	if err != nil {
		return err
	}
	for _, judge := range judges {
		err = judge_model.DeleteJudgeResultByJudgeUID(db, judge.UID)
		if err != nil {
			log_module.AppLogger().WithField("judge", judge).Errorf("delete judge result failed: %v", err)
		}
	}
	err = judge_model.DeleteJudgesByAccount(db, account)
	if err != nil {
		log_module.AppLogger().WithField("account", account).Errorf("delete judges failed: %v", err)
	}
	err = judge_model.DeleteJudgeRankCache(db, account)
	if err != nil {
		log_module.AppLogger().WithField("account", account).Errorf("delete judge rank cache failed: %v", err)
	}
	err = judge_model.DeleteJudgeScoreCacheByUserAccount(db, account)
	if err != nil {
		log_module.AppLogger().WithField("account", account).Errorf("delete judge score cache failed: %v", err)
	}

	err = user_model.DeleteUser(db, account)
	if err != nil {
		return err
	}

	return nil
}

func GetUserList(
	ctx context.Context, options user_model.GetUserOptions,
) ([]user_model.User, int64, error) {
	db := gorm_agent.GetDefaultDB()
	total, err := user_model.CountUsersByOptions(db, options)
	if err != nil {
		return nil, 0, err
	}
	users, err := user_model.GetUsersByOptions(db, options)
	if err != nil {
		return nil, 0, err
	}
	for i := range users {
		enforcer := casbin_agent.GetDefaultCasbinEnforcer()
		if enforcer != nil {
			users[i].Roles = enforcer.GetRolesForUserInDomain(
				casbin_agent.UserSubjectPrefix+users[i].Account, "system")
			for j, role := range users[i].Roles {
				users[i].Roles[j] = role[len(casbin_agent.RoleSubjectPrefix):]
			}
		}
	}

	return users, total, nil
}

func UpdateUser(ctx context.Context, user user_model.User) error {
	db := gorm_agent.GetDefaultDB()
	err := user_model.UpdateUser(db, user)
	if err != nil {
		return err
	}

	return auth_module.UpdateLoginSessionByAccount(ctx,
		user.Account,
		auth_module.LoginSessionData{})
}

func GrantUserRole(ctx context.Context, account, role, domain string) error {
	exist, err := CheckUserExist(ctx, account)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("user not exist")
	}

	enforcer := casbin_agent.GetDefaultCasbinEnforcer()
	account = casbin_agent.UserSubjectPrefix + account
	role = casbin_agent.RoleSubjectPrefix + role
	notDuplicated, err := enforcer.AddRoleForUserInDomain(account, role, domain)
	if err != nil {
		return err
	}
	if !notDuplicated {
		return fmt.Errorf("role already granted")
	}
	err = enforcer.SavePolicy()
	if err != nil {
		return err
	}

	return nil
}

func RevokeUserRole(ctx context.Context, account, role, domain string) error {
	exist, err := CheckUserExist(ctx, account)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("user not exist")
	}

	enforcer := casbin_agent.GetDefaultCasbinEnforcer()
	account = casbin_agent.UserSubjectPrefix + account
	role = casbin_agent.RoleSubjectPrefix + role
	notDuplicated, err := enforcer.DeleteRoleForUserInDomain(account, role, domain)
	if err != nil {
		return err
	}
	if !notDuplicated {
		return fmt.Errorf("role not granted")
	}
	err = enforcer.SavePolicy()
	if err != nil {
		return err
	}

	return nil
}

func CheckUserExist(ctx context.Context, account string) (bool, error) {
	getOptions := user_model.GetUserOptions{
		AccountQuery: account,
	}
	db := gorm_agent.GetDefaultDB()
	count, err := user_model.CountUsersByOptions(db, getOptions)
	if err != nil {
		return false, err
	}

	if count > 1 {
		log_module.AppLogger().
			WithField("account", account).
			WithField("count", count).
			Warn("user account is not unique")
	}

	return count > 0, nil
}

func StartLoginSession(ctx context.Context, account string) (*auth_module.LoginSession, error) {
	ls := auth_module.NewLoginSession(account, auth_module.LoginSessionData{})
	err := ls.SaveToRedis(ctx)
	if err != nil {
		return nil, err
	}

	return ls, nil
}
