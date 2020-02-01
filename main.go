package main

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

// intersection between pgs.FieldType and pgs.FieldTypeElem
type FieldType interface {
	ProtoType() pgs.ProtoType
	IsEmbed() bool
	IsEnum() bool
	Imports() []pgs.File
	Enum() pgs.Enum
	Embed() pgs.Message
}

// intersection between pgs.Message and pgs.Enum
type EntityWithParent interface {
	pgs.Entity
	Parent() pgs.ParentEntity
}

type rbiModule struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
	tpl *template.Template
	serviceTpl *template.Template
}

func RBI() *rbiModule { return &rbiModule{ModuleBase: &pgs.ModuleBase{}} }

func (m *rbiModule) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.ctx = pgsgo.InitContext(c.Parameters())

	funcs := map[string]interface{}{
		"modules": m.modules,
		"rubyPackage": m.rubyPackage,
		"rubyMessageType": m.rubyMessageType,
		"rubyFieldType": m.rubyFieldType,
		"rubyFieldValue": m.rubyFieldValue,
		"rubyMethodType": m.rubyMethodType,
		"increment": m.increment,
	}

	m.tpl = template.Must(template.New("rbi").Funcs(funcs).Parse(tpl))
	m.serviceTpl = template.Must(template.New("rbiService").Funcs(funcs).Parse(serviceTpl))
}

func (m *rbiModule) Name() string { return "rbi" }

func (m *rbiModule) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, t := range targets {
		m.generate(t)

		grpc, err := m.ctx.Params().BoolDefault("grpc", true)
		if err != nil {
			log.Panicf("Bad parameter: grpc\n")
		}

		if len(t.Services()) > 0 && grpc {
			m.generateServices(t)
		}
	}
	return m.Artifacts()
}

func (m *rbiModule) generate(f pgs.File) {
	op := strings.TrimSuffix(f.InputPath().String(), ".proto") + "_pb.rbi"
	m.AddGeneratorTemplateFile(op, m.tpl, f)
}

func (m *rbiModule) generateServices(f pgs.File) {
	op := strings.TrimSuffix(f.InputPath().String(), ".proto") + "_services_pb.rbi"
	m.AddGeneratorTemplateFile(op, m.serviceTpl, f)
}

func (m *rbiModule) modules(file pgs.File) []string {
	p := m.rubyPackage(file)
	split := strings.Split(p, "::")
	modules := make([]string, 0)
	for i := 0; i < len(split); i++ {
		modules = append(modules, strings.Join(split[0:i+1], "::"))
	}
	return modules
}

func (m *rbiModule) rubyPackage(file pgs.File) string {
	pkg := file.Descriptor().GetOptions().GetRubyPackage()
	if pkg == "" {
		pkg = file.Descriptor().GetPackage()
	}
	pkg = strings.Replace(pkg, ".", "::", -1)
	// right now the ruby_out doesn't camelcase the ruby_package, but this results in invalid classes, so do it:
	return pgs.Name(pkg).UpperCamelCase().String()
}

func (m *rbiModule) rubyMessageType(entity EntityWithParent) string {
	names := make([]string, 0)
	outer := entity
	ok := true
	for ok {
		names = append([]string{outer.Name().String()}, names...)
		outer, ok = outer.Parent().(pgs.Message)
	}
	return fmt.Sprintf("%s::%s", m.rubyPackage(entity.File()), strings.Join(names, "::"))
}

func (m *rbiModule) rubyFieldType(field pgs.Field, setter bool) string {
	t := field.Type()
	if t.IsMap() {
		if setter {
			return "Google::Protobuf::Map"
		}
		key := m.rubyProtoTypeElem(field, t.Key())
		value := m.rubyProtoTypeElem(field, t.Element())
		return fmt.Sprintf("T::Hash[%s, %s]", key, value)
	} else if t.IsRepeated() {
		value := m.rubyProtoTypeElem(field, t.Element())
		return fmt.Sprintf("T::Array[%s]", value)
	}
	return m.rubyProtoTypeElem(field, t)
}

func (m *rbiModule) rubyFieldValue(field pgs.Field) string {
	t := field.Type()
	if t.IsMap() {
		key := m.rubyMapType(t.Key())
		if t.Element().ProtoType() == pgs.MessageT {
			value := m.rubyMessageType(t.Element().Embed())
			return fmt.Sprintf("Google::Protobuf::Map.new(%s, :message, %s)", key, value)
		}
		value := m.rubyMapType(t.Element())
		return fmt.Sprintf("Google::Protobuf::Map.new(%s, %s)", key, value)
	} else if t.IsRepeated() {
		return "[]"
	}
	return m.rubyProtoTypeValue(field, t)
}

func (m *rbiModule) rubyProtoTypeElem(field pgs.Field, ft FieldType) string {
	pt := ft.ProtoType()
	if pt.IsInt() {
		return "Integer"
	}
	if pt.IsNumeric() {
		return "Float"
	}
	if pt == pgs.StringT || pt == pgs.BytesT {
		return "String"
	}
	if pt == pgs.BoolT {
		return "T::Boolean"
	}
	if pt == pgs.EnumT {
		return "Symbol"
	}
	if pt == pgs.MessageT {
		return fmt.Sprintf("T.nilable(%s)", m.rubyMessageType(ft.Embed()))
	}
	log.Panicf("Unsupported field type for field: %v\n", field.Name().String())
	return ""
}

func (m *rbiModule) rubyProtoTypeValue(field pgs.Field, ft FieldType) string {
	pt := ft.ProtoType()
	if pt.IsInt() {
		return "0"
	}
	if pt.IsNumeric() {
		return "0.0"
	}
	if pt == pgs.StringT || pt == pgs.BytesT {
		return "\"\""
	}
	if pt == pgs.BoolT {
		return "false"
	}
	if pt == pgs.EnumT {
		return fmt.Sprintf(":%s", ft.Enum().Values()[0].Name().String())
	}
	if pt == pgs.MessageT {
		return "nil"
	}
	log.Panicf("Unsupported field type for field: %v\n", field.Name().String())
	return ""
}

func (m *rbiModule) rubyMapType(ft FieldType) string {
	switch ft.ProtoType() {
	case pgs.DoubleT:
		return ":double"
	case pgs.FloatT:
		return ":float"
	case pgs.Int64T:
		return ":int64"
	case pgs.UInt64T:
		return ":uint64"
	case pgs.Int32T:
		return ":int32"
	case pgs.Fixed64T:
		return ":fixed64"
	case pgs.Fixed32T:
		return ":fixed32"
	case pgs.BoolT:
		return ":bool"
	case pgs.StringT:
		return ":string"
	case pgs.BytesT:
		return ":bytes"
	case pgs.UInt32T:
		return ":uint32"
	case pgs.EnumT:
		return ":enum"
	case pgs.SFixed32:
		return ":sfixed32"
	case pgs.SFixed64:
		return ":sfixed64"
	case pgs.SInt32:
		return ":sint32"
	case pgs.SInt64:
		return ":sint64"
	}
	log.Panicf("Unsupported map field type\n")
	return ""
}

func (m *rbiModule) rubyMethodType(method pgs.Method, input bool) string {
	var streaming bool
	var message pgs.Message
	if input {
		streaming = method.ClientStreaming()
		message = method.Input()
	} else {
		streaming = method.ServerStreaming()
		message = method.Output()
	}
	t := m.rubyMessageType(message)
	if streaming {
		return fmt.Sprintf("T::Enumerable[%s]", t)
	}
	return t
}

func (m *rbiModule) increment(i int) int {
	return i + 1
}

func main() {
	pgs.Init(
		pgs.DebugEnv("DEBUG"),
	).RegisterModule(
		RBI(),
	).RegisterPostProcessor(
		pgsgo.GoFmt(),
	).Render()
}

const tpl = `# Code generated by protoc-gen-rbi. DO NOT EDIT.
# source: {{ .InputPath }}
# typed: strict
{{ range modules . }}
module {{ . }}; end{{ end }}
{{ range .AllMessages }}
class {{ rubyMessageType . }}
  include Google::Protobuf
  include Google::Protobuf::MessageExts
{{ if gt (len .Fields) 0 }}
  sig do
    params({{ $index := 0 }}{{ range .Fields }}{{ if gt $index 0 }},{{ end }}{{ $index = increment $index }}
      {{ .Name }}: {{ rubyFieldType . true }}{{ end }}
    ).void
  end
  def initialize({{ $index := 0 }}{{ range .Fields }}{{ if gt $index 0 }},{{ end }}{{ $index = increment $index }}
    {{ .Name }}: {{ rubyFieldValue . }}{{ end }}
  )
  end
{{ end }}{{ range .Fields }}
  sig { returns({{ rubyFieldType . false }}) }
  def {{ .Name }}
  end

  sig { params(value: {{ rubyFieldType . true }}).void }
  def {{ .Name }}=(value)
  end
{{ end }}end
{{ end }}{{ range .AllEnums }}
module {{ rubyMessageType . }}{{ range .Values }}
  {{ .Name }} = T.let({{ .Value }}, Integer){{ end }}
end
{{ end }}`

const serviceTpl = `# Code generated by protoc-gen-rbi. DO NOT EDIT.
# source: {{ .InputPath }}
# typed: strict
{{ range .Services }}
module {{ rubyPackage .File }}::{{ .Name }}
  class Service
    include GRPC::GenericService
  end

  class Stub
    sig do
      params(
        host: String,
        creds: T.any(GRPC::Core::ChannelCredentials, Symbol),
        kw: T::Hash[Symbol, T.untyped]
      ).void
    end
    def initialize(host, creds, **kw)
    end{{ range .Methods }}

    sig do
      params(
        request: {{ rubyMethodType . true }}
      ).returns({{ rubyMethodType . false }})
    end
    def {{ .Name.LowerSnakeCase }}(request)
    end{{ end }}
  end
end
{{ end }}`
