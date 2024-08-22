require_relative "required_field_pb"

t = Thing.decode_json('{"foo": "foo", "bar": {}}')

!!t.has_foo? == true
!!t.has_bar? == true
!!t.has_optional_value? == false
!!t.bar.has_another_optional_value? == false

Google::Protobuf::DescriptorPool.generated_pool.lookup("required_field").get(Thing.descriptor.lookup("foo").options) == true

Google::Protobuf::DescriptorPool.generated_pool.lookup("required_field").get(Thing.descriptor.lookup("bar").options) == true

Google::Protobuf::DescriptorPool.generated_pool.lookup("required_field").get(Thing.descriptor.lookup("optional_value").options) == false

Google::Protobuf::DescriptorPool.generated_pool.lookup("required_field").get(Thing::InnerThing.descriptor.lookup("another_optional_value").options) == false
