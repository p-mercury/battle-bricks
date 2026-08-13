package services.catalogue.list_loadouts

# METADATA
# scope: document
# schemas:
#   - input.context: schema["catalogue.v1.ListLoadoutsContext.jsonschema.strict.bundle"]
default authz := false

authz if {
	data.user.is_active
}

decision := {
	"authz": authz,
	"output_mask": {"loadouts": authz},
}
