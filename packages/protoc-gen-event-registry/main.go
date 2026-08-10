package main

import (
	"path"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	sdopts "protoc-gen-event-registry/dist/eventregistry/v1"
)

type Event struct {
	DetailType string
	Message    *protogen.Message
}

type Package struct {
	GoImport     protogen.GoImportPath
	GoPackage    protogen.GoPackageName
	OutDir       string
	SourceSuffix string
	Events       []Event
}

func main() {
	protogen.Options{}.Run(func(p *protogen.Plugin) error {
		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		packages := map[string]*Package{}

		for _, f := range p.Files {
			if !f.Generate {
				continue
			}

			key := string(f.GoImportPath) + "|" + string(f.GoPackageName)
			if packages[key] == nil {
				packages[key] = &Package{
					GoImport:     f.GoImportPath,
					GoPackage:    f.GoPackageName,
					OutDir:       path.Dir(f.GeneratedFilenamePrefix),
					SourceSuffix: "",
					Events:       []Event{},
				}
			}

			if fo, _ := f.Desc.Options().(*descriptorpb.FileOptions); fo != nil && proto.HasExtension(fo, sdopts.E_EventFile) {
				if ext, ok := proto.GetExtension(fo, sdopts.E_EventFile).(*sdopts.FileOptions); ok && ext != nil {
					packages[key].SourceSuffix = strings.TrimSpace(ext.GetSource())
				}
			}

			for _, message := range f.Messages {
				mo, _ := message.Desc.Options().(*descriptorpb.MessageOptions)
				if mo == nil || !proto.HasExtension(mo, sdopts.E_Event) {
					continue
				}

				e, ok := proto.GetExtension(mo, sdopts.E_Event).(*sdopts.EventOptions)
				if !ok || e == nil {
					continue
				}

				for _, detailType := range e.GetDetailTypes() {
					packages[key].Events = append(packages[key].Events, Event{
						DetailType: strings.TrimSpace(detailType),
						Message:    message,
					})
				}
			}
		}

		for _, pac := range packages {
			if err := generateGo(p, pac); err != nil {
				return err
			}

			if err := generateTs(p, pac); err != nil {
				return err
			}

			if err := generateSfnStepsTs(p, pac); err != nil {
				return err
			}
		}

		return nil
	})
}
