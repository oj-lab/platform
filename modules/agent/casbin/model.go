package casbin_agent

// RBAC with domain model for Casbin, supports evaluation of extra information
//   - sub: user / role
//   - ext: extra information
//   - ext_rule: expression to evaluate extra information
//   - dom: domain
//   - obj: object, supports gin router key match (EX: /api/v1/user/:id/*any)
//   - act: action, supports regex match (EX: GET | POST)
//   - eft: effect
//
// Router key match will not work if keyMatchGin is not defined in the model.
//
// You should also prepare an extra information struct to evaluate the extra information
// and make sure the policy ext_rule is using the correct field
// (which is included in the extra information struct).
const ExtendedRBACWithDomainModelString = `
[request_definition]
r = sub, ext, dom, obj, act

[policy_definition]
p = sub, ext_rule, dom, obj, act, eft

[role_definition]
g = _, _，_

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && eval(p.ext_rule) && r.dom == p.dom && keyMatchGin(r.obj, p.obj) && regexMatch(r.act, p.act)
`

type ExtraInfo struct {
	IsVIP bool
}
