package services.catalogue.list_units

# METADATA
# scope: document
# schemas:
#   - input.context: schema["catalogue.v1.ListUnitsContext.jsonschema.strict.bundle"]
default authz := false

authz if {
	data.user.is_active
}

decision := {
	"authz": authz,
	"output_mask": {"units": authz},
}
