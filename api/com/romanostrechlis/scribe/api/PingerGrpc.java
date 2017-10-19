package com.romanostrechlis.scribe.api;

import static io.grpc.stub.ClientCalls.asyncUnaryCall;
import static io.grpc.stub.ClientCalls.asyncServerStreamingCall;
import static io.grpc.stub.ClientCalls.asyncClientStreamingCall;
import static io.grpc.stub.ClientCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ClientCalls.blockingUnaryCall;
import static io.grpc.stub.ClientCalls.blockingServerStreamingCall;
import static io.grpc.stub.ClientCalls.futureUnaryCall;
import static io.grpc.MethodDescriptor.generateFullMethodName;
import static io.grpc.stub.ServerCalls.asyncUnaryCall;
import static io.grpc.stub.ServerCalls.asyncServerStreamingCall;
import static io.grpc.stub.ServerCalls.asyncClientStreamingCall;
import static io.grpc.stub.ServerCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.7.0)",
    comments = "Source: logScribe.proto")
public final class PingerGrpc {

  private PingerGrpc() {}

  public static final String SERVICE_NAME = "com.romanostrechlis.scribe.api.Pinger";

  // Static method descriptors that strictly reflect the proto.
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static final io.grpc.MethodDescriptor<com.romanostrechlis.scribe.api.LogScribeOuterClass.PingRequest,
      com.romanostrechlis.scribe.api.LogScribeOuterClass.PingResponse> METHOD_PING =
      io.grpc.MethodDescriptor.<com.romanostrechlis.scribe.api.LogScribeOuterClass.PingRequest, com.romanostrechlis.scribe.api.LogScribeOuterClass.PingResponse>newBuilder()
          .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
          .setFullMethodName(generateFullMethodName(
              "com.romanostrechlis.scribe.api.Pinger", "Ping"))
          .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
              com.romanostrechlis.scribe.api.LogScribeOuterClass.PingRequest.getDefaultInstance()))
          .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
              com.romanostrechlis.scribe.api.LogScribeOuterClass.PingResponse.getDefaultInstance()))
          .setSchemaDescriptor(new PingerMethodDescriptorSupplier("Ping"))
          .build();

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static PingerStub newStub(io.grpc.Channel channel) {
    return new PingerStub(channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static PingerBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    return new PingerBlockingStub(channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static PingerFutureStub newFutureStub(
      io.grpc.Channel channel) {
    return new PingerFutureStub(channel);
  }

  /**
   */
  public static abstract class PingerImplBase implements io.grpc.BindableService {

    /**
     */
    public void ping(com.romanostrechlis.scribe.api.LogScribeOuterClass.PingRequest request,
        io.grpc.stub.StreamObserver<com.romanostrechlis.scribe.api.LogScribeOuterClass.PingResponse> responseObserver) {
      asyncUnimplementedUnaryCall(METHOD_PING, responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            METHOD_PING,
            asyncUnaryCall(
              new MethodHandlers<
                com.romanostrechlis.scribe.api.LogScribeOuterClass.PingRequest,
                com.romanostrechlis.scribe.api.LogScribeOuterClass.PingResponse>(
                  this, METHODID_PING)))
          .build();
    }
  }

  /**
   */
  public static final class PingerStub extends io.grpc.stub.AbstractStub<PingerStub> {
    private PingerStub(io.grpc.Channel channel) {
      super(channel);
    }

    private PingerStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PingerStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new PingerStub(channel, callOptions);
    }

    /**
     */
    public void ping(com.romanostrechlis.scribe.api.LogScribeOuterClass.PingRequest request,
        io.grpc.stub.StreamObserver<com.romanostrechlis.scribe.api.LogScribeOuterClass.PingResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_PING, getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class PingerBlockingStub extends io.grpc.stub.AbstractStub<PingerBlockingStub> {
    private PingerBlockingStub(io.grpc.Channel channel) {
      super(channel);
    }

    private PingerBlockingStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PingerBlockingStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new PingerBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.romanostrechlis.scribe.api.LogScribeOuterClass.PingResponse ping(com.romanostrechlis.scribe.api.LogScribeOuterClass.PingRequest request) {
      return blockingUnaryCall(
          getChannel(), METHOD_PING, getCallOptions(), request);
    }
  }

  /**
   */
  public static final class PingerFutureStub extends io.grpc.stub.AbstractStub<PingerFutureStub> {
    private PingerFutureStub(io.grpc.Channel channel) {
      super(channel);
    }

    private PingerFutureStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PingerFutureStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new PingerFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.romanostrechlis.scribe.api.LogScribeOuterClass.PingResponse> ping(
        com.romanostrechlis.scribe.api.LogScribeOuterClass.PingRequest request) {
      return futureUnaryCall(
          getChannel().newCall(METHOD_PING, getCallOptions()), request);
    }
  }

  private static final int METHODID_PING = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final PingerImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(PingerImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_PING:
          serviceImpl.ping((com.romanostrechlis.scribe.api.LogScribeOuterClass.PingRequest) request,
              (io.grpc.stub.StreamObserver<com.romanostrechlis.scribe.api.LogScribeOuterClass.PingResponse>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class PingerBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    PingerBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.romanostrechlis.scribe.api.LogScribeOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Pinger");
    }
  }

  private static final class PingerFileDescriptorSupplier
      extends PingerBaseDescriptorSupplier {
    PingerFileDescriptorSupplier() {}
  }

  private static final class PingerMethodDescriptorSupplier
      extends PingerBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    PingerMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (PingerGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new PingerFileDescriptorSupplier())
              .addMethod(METHOD_PING)
              .build();
        }
      }
    }
    return result;
  }
}
