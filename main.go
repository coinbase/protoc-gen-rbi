package main

import (
	"log"
	"regexp"
	"strings"
	"text/template"

	"google.golang.org/protobuf/types/pluginpb"

	"github.com/vaidehi-figma/protoc-gen-rbi/ruby_types"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"

	pgs "github.com/vaidehi-figma/protoc-gen-star"
	pgsgo "github.com/vaidehi-figma/protoc-gen-star/lang/go"
)

var (
	validRubyField = regexp.MustCompile(`\A[a-z][A-Za-z0-9_]*\z`)
)

var SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL) | uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
var MaximumEdition = int32(descriptorpb.Edition_EDITION_2023)
var MinimumEdition = int32(descriptorpb.Edition_EDITION_2023)

type rbiModule struct {
	*pgs.ModuleBase
	ctx                pgsgo.Context
	tpl                *template.Template
	serviceTpl         *template.Template
	hideCommonMethods  bool
	useAbstractMessage bool
	hasEditions2023    bool
}

func (m *rbiModule) HideCommonMethods() bool {
	return m.hideCommonMethods
}

func (m *rbiModule) UseAbstractMessage() bool {
	return m.useAbstractMessage
}

func (m *rbiModule) HasEditions2023() bool {
	return m.hasEditions2023
}

func RBI() *rbiModule { return &rbiModule{ModuleBase: &pgs.ModuleBase{}} }

func (m *rbiModule) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.ctx = pgsgo.InitContext(c.Parameters())

	hideCommonMethods, err := m.ctx.Params().BoolDefault("hide_common_methods", false)
	if err != nil {
		log.Panicf("Bad parameter: hide_common_methods\n")
	}
	m.hideCommonMethods = hideCommonMethods

	useAbstractMessage, err := m.ctx.Params().BoolDefault("use_abstract_message", false)
	if err != nil {
		log.Panicf("Bad parameter: use_abstract_message\n")
	}
	m.useAbstractMessage = useAbstractMessage

	funcs := map[string]interface{}{
		"increment":                m.increment,
		"optional":                 m.optional,
		"optionalOneOf":            m.optionalOneOf,
		"willGenerateInvalidRuby":  m.willGenerateInvalidRuby,
		"rubyPackage":              ruby_types.RubyPackage,
		"rubyMessageType":          ruby_types.RubyMessageType,
		"rubyMessageTypeComment":   ruby_types.RubyMessageTypeComment,
		"rubyFieldTypeComment":     ruby_types.RubyFieldTypeComment,
		"rubyGetterFieldType":      ruby_types.RubyGetterFieldType,
		"rubySetterFieldType":      ruby_types.RubySetterFieldType,
		"rubyInitializerFieldType": ruby_types.RubyInitializerFieldType,
		"rubyFieldValue":           ruby_types.RubyFieldValue,
		"rubyMethodTypeComment":    ruby_types.RubyMethodTypeComment,
		"rubyMethodParamType":      ruby_types.RubyMethodParamType,
		"rubyMethodReturnType":     ruby_types.RubyMethodReturnType,
		"hideCommonMethods":        m.HideCommonMethods,
		"useAbstractMessage":       m.UseAbstractMessage,
		"edition2023":              m.HasEditions2023,
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
	m.hasEditions2023 = f.Descriptor().GetEdition() == descriptorpb.Edition_EDITION_2023
	m.AddGeneratorTemplateFile(op, m.tpl, f)
}

func (m *rbiModule) generateServices(f pgs.File) {
	op := strings.TrimSuffix(f.InputPath().String(), ".proto") + "_services_pb.rbi"
	m.AddGeneratorTemplateFile(op, m.serviceTpl, f)
}

func (m *rbiModule) increment(i int) int {
	return i + 1
}

func (m *rbiModule) optional(field pgs.Field) bool {
	return field.Descriptor().GetProto3Optional()
}

func (m *rbiModule) optionalOneOf(oneOf pgs.OneOf) bool {
	return len(oneOf.Fields()) == 1 && oneOf.Fields()[0].Descriptor().GetProto3Optional()
}

func (m *rbiModule) willGenerateInvalidRuby(fields []pgs.Field) bool {
	for _, field := range fields {
		if !validRubyField.MatchString(string(field.Name())) {
			return true
		}
	}
	return false
}

func main() {
	pgs.Init(
		pgs.DebugEnv("DEBUG"),
		pgs.SupportedFeatures(&SupportedFeatures),
		pgs.MaximumEdition(&MaximumEdition),
		pgs.MinimumEdition(&MinimumEdition),
	).RegisterModule(
		RBI(),
	).RegisterPostProcessor(
		pgsgo.GoFmt(),
	).Render()
}

const tpl = `# Code generated by figma's fork of protoc-gen-rbi! DO NOT EDIT.
# source: {{ .InputPath }}
# typed: strict
{{ range .AllMessages }}{{ if rubyMessageTypeComment . }}
# {{ rubyMessageTypeComment . }}{{ end }}
class {{ rubyMessageType . }}{{ if useAbstractMessage }} < ::Google::Protobuf::AbstractMessage{{ else }}
  include ::Google::Protobuf::MessageExts
  extend ::Google::Protobuf::MessageExts::ClassMethods
{{ end }}{{ if hideCommonMethods }}{{ else }}
  sig { params(str: String).returns({{ rubyMessageType . }}) }
  def self.decode(str)
  end

  sig { params(msg: {{ rubyMessageType . }}).returns(String) }
  def self.encode(msg)
  end

  sig { params(str: String, kw: T.untyped).returns({{ rubyMessageType . }}) }
  def self.decode_json(str, **kw)
  end

  sig { params(msg: {{ rubyMessageType . }}, kw: T.untyped).returns(String) }
  def self.encode_json(msg, **kw)
  end

  sig { returns(::Google::Protobuf::Descriptor) }
  def self.descriptor
  end
{{ end }}{{ if willGenerateInvalidRuby .Fields }}
  # Constants of the form Constant_1 are invalid. We've declined to type this as a result, taking a hash instead.
  sig { params(args: T::Hash[T.untyped, T.untyped]).void }
  def initialize(args); end
{{ else if gt (len .Fields) 0 }}
  sig do
    params({{ $index := 0 }}{{ range .Fields }}{{ if gt $index 0 }},{{ end }}{{ $index = increment $index }}
      {{ .Name }}: {{ rubyInitializerFieldType . }}{{ end }}
    ).void
  end
  def initialize({{ $index := 0 }}{{ range .Fields }}{{ if gt $index 0 }},{{ end }}{{ $index = increment $index }}
    {{ .Name }}: {{ rubyFieldValue . }}{{ end }}
  )
  end
{{ else }}
  sig {void}
  def initialize; end
{{ end }}{{ range .Fields }}{{ if rubyFieldTypeComment . }}
  # {{ rubyFieldTypeComment . }}{{ end }}
  sig { returns({{ rubyGetterFieldType . }}) }
  def {{ .Name }}
  end
{{ if rubyFieldTypeComment . }}
  # {{ rubyFieldTypeComment . }}{{ end }}
  sig { params(value: {{ rubySetterFieldType . }}).void }
  def {{ .Name }}=(value)
  end
{{ if rubyFieldTypeComment . }}
  # {{ rubyFieldTypeComment . }}{{ end }}
  sig { void }
  def clear_{{ .Name }}
  end
{{ if or (edition2023) (optional .) }}
  sig { returns(T::Boolean) }
  def has_{{ .Name }}?
  end
{{ end }}{{ end }}{{ range .OneOfs }}{{ if not (optionalOneOf .) }}
  sig { returns(T.nilable(Symbol)) }
  def {{ .Name }}
  end
{{ end }}{{ end }}{{ if hideCommonMethods }}{{ else }}
  sig { params(field: String).returns(T.untyped) }
  def [](field)
  end

  sig { params(field: String, value: T.untyped).void }
  def []=(field, value)
  end

  sig { returns(T::Hash[Symbol, T.untyped]) }
  def to_h
  end
{{ end }}end
{{ end }}{{ range .AllEnums }}
module {{ rubyMessageType . }}{{ range .Values }}
  self::{{ .Name }} = T.let({{ .Value }}, Integer){{ end }}

  sig { params(value: Integer).returns(T.nilable(Symbol)) }
  def self.lookup(value)
  end

  sig { params(value: Symbol).returns(T.nilable(Integer)) }
  def self.resolve(value)
  end

  sig { returns(::Google::Protobuf::EnumDescriptor) }
  def self.descriptor
  end
end
{{ end }}`

const serviceTpl = `# Code generated by protoc-gen-rbi. DO NOT EDIT.
# source: {{ .InputPath }}
# typed: strict
{{ range .Services }}
module {{ rubyPackage .File }}::{{ .Name }}
  class Service
    include ::GRPC::GenericService
  end

  class Stub < ::GRPC::ClientStub
    sig do
      params(
        host: String,
        creds: T.any(::GRPC::Core::ChannelCredentials, Symbol),
        kw: T.untyped,
      ).void
    end
    def initialize(host, creds, **kw)
    end{{ range .Methods }}
{{ if rubyMethodTypeComment . }}
    # {{ rubyMethodTypeComment . }}{{ end }}
    sig do
      params(
        request: {{ rubyMethodParamType . }}
      ).returns({{ rubyMethodReturnType . }})
    end
    def {{ .Name.LowerSnakeCase }}(request)
    end{{ end }}
  end
end
{{ end }}`
