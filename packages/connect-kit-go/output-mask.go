package connectkit

import (
	"maps"
	"strconv"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

func ApplyOutputMask(msg protoreflect.ProtoMessage, mask map[string]*structpb.Value) {
	if msg == nil || mask == nil {
		return
	}
	applyMessageMask(msg.ProtoReflect(), mask)
}

func applyMessageMask(m protoreflect.Message, fields map[string]*structpb.Value) {
	type rewrite struct {
		fd  protoreflect.FieldDescriptor
		val protoreflect.Value
	}
	var pending []rewrite

	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		name := string(fd.Name())
		mv, ok := fields[name]
		if !ok || mv == nil {
			m.Clear(fd)
			return true
		}

		switch kind := mv.GetKind().(type) {

		case *structpb.Value_BoolValue:
			if !kind.BoolValue {
				m.Clear(fd)
			}
			return true

		case *structpb.Value_StructValue:
			sub := kind.StructValue.GetFields()

			if fd.IsList() {
				if fd.Kind() != protoreflect.MessageKind {
					m.Clear(fd)
					return true
				}

				if len(sub) == 0 {
					m.Clear(fd)
					return true
				}

				indexed := true
				for k := range sub {
					if _, err := strconv.Atoi(k); err != nil {
						indexed = false
						break
					}
				}

				l := v.List()

				if indexed {
					kept := m.NewField(fd).List()
					for i := 0; i < l.Len(); i++ {
						ev, ok := sub[strconv.Itoa(i)]
						if !ok || ev == nil {
							continue
						}
						switch k := ev.GetKind().(type) {
						case *structpb.Value_BoolValue:
							if k.BoolValue {
								kept.Append(l.Get(i))
							}
						case *structpb.Value_StructValue:
							elem := l.Get(i)
							if em := elem.Message(); em.IsValid() {
								applyMessageMask(em, k.StructValue.GetFields())
							}
							kept.Append(elem)
						}
					}
					pending = append(pending, rewrite{fd, protoreflect.ValueOfList(kept)})
					return true
				}

				for i := 0; i < l.Len(); i++ {
					elem := l.Get(i).Message()
					if elem.IsValid() {
						applyMessageMask(elem, sub)
					}
				}
				return true
			}

			if fd.IsMap() {
				mvDesc := fd.MapValue()
				if mvDesc.Kind() != protoreflect.MessageKind {
					m.Clear(fd)
					return true
				}
				mp := v.Map()
				mp.Range(func(k protoreflect.MapKey, val protoreflect.Value) bool {
					elem := val.Message()
					if elem.IsValid() {
						applyMessageMask(elem, sub)
					}
					return true
				})
				return true
			}

			if fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
				if m.Has(fd) {
					applyMessageMask(v.Message(), sub)
				}
				return true
			}

			m.Clear(fd)
			return true

		case *structpb.Value_ListValue:
			if !fd.IsList() || fd.Kind() != protoreflect.MessageKind {
				m.Clear(fd)
				return true
			}

			vals := kind.ListValue.GetValues()
			if len(vals) == 0 {
				return true
			}

			for _, ev := range vals {
				if b, ok := ev.GetKind().(*structpb.Value_BoolValue); ok && b.BoolValue {
					return true
				}
			}

			merged := map[string]*structpb.Value{}
			for _, ev := range vals {
				if s, ok := ev.GetKind().(*structpb.Value_StructValue); ok && s.StructValue != nil {
					maps.Copy(merged, s.StructValue.GetFields())
				}
			}
			if len(merged) == 0 {
				m.Clear(fd)
				return true
			}

			l := v.List()
			for i := 0; i < l.Len(); i++ {
				elem := l.Get(i).Message()
				if elem.IsValid() {
					applyMessageMask(elem, merged)
				}
			}
			return true

		default:
			m.Clear(fd)
			return true
		}
	})

	for _, r := range pending {
		m.Set(r.fd, r.val)
	}
}
