# Code generated by protoc-gen-rbi. DO NOT EDIT.
# source: example.proto
# typed: strict

# some description for request message
class Example::Request < ::Google::Protobuf::AbstractMessage
  sig do
    params(
      name: T.nilable(String)
    ).void
  end
  def initialize(
    name: ""
  )
  end

  # some description for name field
  sig { returns(String) }
  def name
  end

  # some description for name field
  sig { params(value: String).void }
  def name=(value)
  end

  # some description for name field
  sig { void }
  def clear_name
  end
end

# some description for responsee message that is multi line and has a # in it
# that needs to be escaped
class Example::Response < ::Google::Protobuf::AbstractMessage
  sig do
    params(
      greeting: T.nilable(String)
    ).void
  end
  def initialize(
    greeting: ""
  )
  end

  # some description for greeting field
  sig { returns(String) }
  def greeting
  end

  # some description for greeting field
  sig { params(value: String).void }
  def greeting=(value)
  end

  # some description for greeting field
  sig { void }
  def clear_greeting
  end
end
