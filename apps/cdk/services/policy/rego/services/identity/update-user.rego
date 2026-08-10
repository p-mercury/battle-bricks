package services.identity.update_user

# METADATA
# scope: document
# schemas:
#   - input.context: schema["identity.v1.UpdateUserContext.jsonschema.strict.bundle"]
default authz := false

authz if {
	data.user.is_active

	not "correlations" in input.context.request.updateMask
}

decision := {
	"authz": authz,
	"output_mask": {"user": {
		"id": authz,
		"correlations": authz,
		"version": authz,
		"created_time": authz,
		"modified_time": authz,
		"organisation_id": authz,
		"email_address": authz,
		"status": authz,
		"name": authz,
		"job_title": authz,
		"image": authz,
		"language": authz,
		"notification_settings": authz,
	}},
}
