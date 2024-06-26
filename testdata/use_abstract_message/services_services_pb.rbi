# Code generated by protoc-gen-rbi. DO NOT EDIT.
# source: services.proto
# typed: strict

module Testdata::SimpleMathematics
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
    end

    # Negates the input
    sig do
      params(
        request: Testdata::Subdir::IntegerMessage
      ).returns(Testdata::Subdir::IntegerMessage)
    end
    def negate(request)
    end

    # Report the median of a stream of integers
    sig do
      params(
        request: T::Enumerable[Testdata::Subdir::IntegerMessage]
      ).returns(Testdata::Subdir::IntegerMessage)
    end
    def median(request)
    end
  end
end

module Testdata::ComplexMathematics
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
    end

    # Stream the first N numbers in the Fibonacci sequence
    sig do
      params(
        request: Testdata::Subdir::IntegerMessage
      ).returns(T::Enumerable[Testdata::Subdir::IntegerMessage])
    end
    def fibonacci(request)
    end

    # Accept a stream of integers, and report whenever a new maximum is found
    sig do
      params(
        request: T::Enumerable[Testdata::Subdir::IntegerMessage]
      ).returns(T::Enumerable[Testdata::Subdir::IntegerMessage])
    end
    def running_max(request)
    end

    # Accept a stream of integers, and report the maximum every second
    sig do
      params(
        request: T::Enumerable[Testdata::Subdir::IntegerMessage]
      ).returns(T::Enumerable[Testdata::Subdir::IntegerMessage])
    end
    def periodic_max(request)
    end
  end
end
