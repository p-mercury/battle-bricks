package services.catalogue.update_squad

# METADATA
# scope: document
# schemas:
#   - input.context: schema["catalogue.v1.UpdateSquadContext.jsonschema.strict.bundle"]
default authz := false

authz if {
	data.user.is_active
}

decision := {
	"authz": authz,
	"output_mask": {"squad": authz},
}
