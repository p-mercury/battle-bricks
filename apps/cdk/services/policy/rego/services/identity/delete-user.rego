package services.identity.delete_user

# METADATA
# scope: document
# schemas:
#   - input.context: schema["identity.v1.DeleteUserContext.jsonschema.strict.bundle"]
default authz := false

authz if {
	data.user.is_active
}

authz if {
	data.user.id == input.context.subject.id
}

decision := {
	"authz": authz,
	"output_mask": {
		"id": authz,
		"correlations": authz,
		"version": authz,
		"created_time": authz,
		"modified_time": authz,
	},
}
