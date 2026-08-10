# METADATA
# scope: package
# schemas:
#   - input.user: schema["policy.v1.User.jsonschema.strict.bundle"]
package user

default is_active := false

is_active if {
	input.user.status == "USER_STATUS_ACTIVE"
}

default id := null

id := input.user.id
