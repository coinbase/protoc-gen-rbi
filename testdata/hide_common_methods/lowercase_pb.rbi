# Code generated by protoc-gen-rbi. DO NOT EDIT.
# source: lowercase.proto
# typed: strict

class Example::Lowercase
  include ::Google::Protobuf::MessageExts
  extend ::Google::Protobuf::MessageExts::ClassMethods

  sig do
    params(
      example_proto_field: T.nilable(String)
    ).void
  end
  def initialize(
    example_proto_field: ""
  )
  end

  sig { returns(String) }
  def example_proto_field
  end

  sig { params(value: String).void }
  def example_proto_field=(value)
  end

  sig { void }
  def clear_example_proto_field
  end
end

class Example::Lowercase_with_underscores
  include ::Google::Protobuf::MessageExts
  extend ::Google::Protobuf::MessageExts::ClassMethods

  sig do
    params(
      example_proto_field: T.nilable(String)
    ).void
  end
  def initialize(
    example_proto_field: ""
  )
  end

  sig { returns(String) }
  def example_proto_field
  end

  sig { params(value: String).void }
  def example_proto_field=(value)
  end

  sig { void }
  def clear_example_proto_field
  end
end
