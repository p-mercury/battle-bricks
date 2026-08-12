package services.catalogue.list_items

# METADATA
# scope: document
# schemas:
#   - input.context: schema["catalogue.v1.ListItemsContext.jsonschema.strict.bundle"]
default authz := false

authz if {
	data.user.is_active
}

decision := {
	"authz": authz,
	"output_mask": {"items": authz},
}
