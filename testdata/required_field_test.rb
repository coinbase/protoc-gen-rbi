require_relative 'required_field_pb'

t = Thing.decode_json('{"foo": "foo", "bar": {}}')

puts "t.class: #{t.class}"
puts "t.has_foo?: #{!!t.has_foo?}"
puts "t.foo?: #{t.foo}"
puts "t.has_bar?: #{!!t.has_bar?}"
puts "t.bar: #{t.bar}"
puts "t.has_optional_value?: #{!!t.has_optional_value?}"
puts "t.optional_value: #{t.optional_value}"
puts "t.bar.has_another_optional_value?: #{!!t.bar.has_another_optional_value?}"
puts "t.bar.another_optional_value: #{t.bar.another_optional_value}"

puts "foo, required_field?: #{Google::Protobuf::DescriptorPool.generated_pool.lookup('required_field').get(Thing.descriptor.lookup('foo').options) == true}"

puts "bar, required_field?: #{Google::Protobuf::DescriptorPool.generated_pool.lookup('required_field').get(Thing.descriptor.lookup('bar').options) == true}"

puts "optional_value, required_field?: #{Google::Protobuf::DescriptorPool.generated_pool.lookup('required_field').get(Thing.descriptor.lookup('optional_value').options) == true}"

puts "another_optional_value, required_field?: #{Google::Protobuf::DescriptorPool.generated_pool.lookup('required_field').get(Thing::InnerThing.descriptor.lookup('another_optional_value').options) == true}"
