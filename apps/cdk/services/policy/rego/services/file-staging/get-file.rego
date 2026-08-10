package services.file_staging.get_file

# METADATA
# scope: document
# schemas:
#   - input.context: schema["filestaging.v1.GetFileContext.jsonschema.strict.bundle"]
default authz := false

authz if {
	input.context.subject.userId == data.user.id
}

decision := {
	"authz": authz,
	"output_mask": {"file": authz},
}
