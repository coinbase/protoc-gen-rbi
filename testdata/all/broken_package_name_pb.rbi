# Code generated by protoc-gen-rbi. DO NOT EDIT.
# source: broken_package_name.proto
# typed: strict

class Package2test::Message2test < ::Google::Protobuf::AbstractMessage
  sig do
    params(
      field2test: T.nilable(String)
    ).void
  end
  def initialize(
    field2test: ""
  )
  end

  sig { returns(String) }
  def field2test
  end

  sig { params(value: String).void }
  def field2test=(value)
  end

  sig { void }
  def clear_field2test
  end
end