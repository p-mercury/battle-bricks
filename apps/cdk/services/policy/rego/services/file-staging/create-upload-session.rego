package services.file_staging.create_upload_session

# METADATA
# scope: document
# schemas:
#   - input.context: schema["filestaging.v1.CreateUploadSessionContext.jsonschema.strict.bundle"]
default authz := false

authz if {
	data.user.id != null
}

decision := {
	"authz": authz,
	"output_mask": {"session": authz},
}
