package services.catalogue.list_squads

# METADATA
# scope: document
# schemas:
#   - input.context: schema["catalogue.v1.ListSquadsContext.jsonschema.strict.bundle"]
default authz := false

authz if {
	data.user.is_active
}

decision := {
	"authz": authz,
	"output_mask": {"squads": authz},
}
