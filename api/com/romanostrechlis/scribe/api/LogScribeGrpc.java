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
 * <pre>
 * service defines a new service names LogScribe
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.7.0)",
    comments = "Source: logScribe.proto")
public final class LogScribeGrpc {

  private LogScribeGrpc() {}

  public static final String SERVICE_NAME = "com.romanostrechlis.scribe.api.LogScribe";

  // Static method descriptors that strictly reflect the proto.
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static final io.grpc.MethodDescriptor<com.romanostrechlis.scribe.api.LogScribeOuterClass.LogRequest,
      com.romanostrechlis.scribe.api.LogScribeOuterClass.LogResponse> METHOD_LOG =
      io.grpc.MethodDescriptor.<com.romanostrechlis.scribe.api.LogScribeOuterClass.LogRequest, com.romanostrechlis.scribe.api.LogScribeOuterClass.LogResponse>newBuilder()
          .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
          .setFullMethodName(generateFullMethodName(
              "com.romanostrechlis.scribe.api.LogScribe", "Log"))
          .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
              com.romanostrechlis.scribe.api.LogScribeOuterClass.LogRequest.getDefaultInstance()))
          .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
              com.romanostrechlis.scribe.api.LogScribeOuterClass.LogResponse.getDefaultInstance()))
          .setSchemaDescriptor(new LogScribeMethodDescriptorSupplier("Log"))
          .build();

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static LogScribeStub newStub(io.grpc.Channel channel) {
    return new LogScribeStub(channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static LogScribeBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    return new LogScribeBlockingStub(channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static LogScribeFutureStub newFutureStub(
      io.grpc.Channel channel) {
    return new LogScribeFutureStub(channel);
  }

  /**
   * <pre>
   * service defines a new service names LogScribe
   * </pre>
   */
  public static abstract class LogScribeImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * the LogScribe service sends a LogRequest
     * and recieves a LogResponse
     * </pre>
     */
    public void log(com.romanostrechlis.scribe.api.LogScribeOuterClass.LogRequest request,
        io.grpc.stub.StreamObserver<com.romanostrechlis.scribe.api.LogScribeOuterClass.LogResponse> responseObserver) {
      asyncUnimplementedUnaryCall(METHOD_LOG, responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            METHOD_LOG,
            asyncUnaryCall(
              new MethodHandlers<
                com.romanostrechlis.scribe.api.LogScribeOuterClass.LogRequest,
                com.romanostrechlis.scribe.api.LogScribeOuterClass.LogResponse>(
                  this, METHODID_LOG)))
          .build();
    }
  }

  /**
   * <pre>
   * service defines a new service names LogScribe
   * </pre>
   */
  public static final class LogScribeStub extends io.grpc.stub.AbstractStub<LogScribeStub> {
    private LogScribeStub(io.grpc.Channel channel) {
      super(channel);
    }

    private LogScribeStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected LogScribeStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new LogScribeStub(channel, callOptions);
    }

    /**
     * <pre>
     * the LogScribe service sends a LogRequest
     * and recieves a LogResponse
     * </pre>
     */
    public void log(com.romanostrechlis.scribe.api.LogScribeOuterClass.LogRequest request,
        io.grpc.stub.StreamObserver<com.romanostrechlis.scribe.api.LogScribeOuterClass.LogResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_LOG, getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * service defines a new service names LogScribe
   * </pre>
   */
  public static final class LogScribeBlockingStub extends io.grpc.stub.AbstractStub<LogScribeBlockingStub> {
    private LogScribeBlockingStub(io.grpc.Channel channel) {
      super(channel);
    }

    private LogScribeBlockingStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected LogScribeBlockingStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new LogScribeBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * the LogScribe service sends a LogRequest
     * and recieves a LogResponse
     * </pre>
     */
    public com.romanostrechlis.scribe.api.LogScribeOuterClass.LogResponse log(com.romanostrechlis.scribe.api.LogScribeOuterClass.LogRequest request) {
      return blockingUnaryCall(
          getChannel(), METHOD_LOG, getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * service defines a new service names LogScribe
   * </pre>
   */
  public static final class LogScribeFutureStub extends io.grpc.stub.AbstractStub<LogScribeFutureStub> {
    private LogScribeFutureStub(io.grpc.Channel channel) {
      super(channel);
    }

    private LogScribeFutureStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected LogScribeFutureStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new LogScribeFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * the LogScribe service sends a LogRequest
     * and recieves a LogResponse
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<com.romanostrechlis.scribe.api.LogScribeOuterClass.LogResponse> log(
        com.romanostrechlis.scribe.api.LogScribeOuterClass.LogRequest request) {
      return futureUnaryCall(
          getChannel().newCall(METHOD_LOG, getCallOptions()), request);
    }
  }

  private static final int METHODID_LOG = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final LogScribeImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(LogScribeImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_LOG:
          serviceImpl.log((com.romanostrechlis.scribe.api.LogScribeOuterClass.LogRequest) request,
              (io.grpc.stub.StreamObserver<com.romanostrechlis.scribe.api.LogScribeOuterClass.LogResponse>) responseObserver);
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

  private static abstract class LogScribeBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    LogScribeBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.romanostrechlis.scribe.api.LogScribeOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("LogScribe");
    }
  }

  private static final class LogScribeFileDescriptorSupplier
      extends LogScribeBaseDescriptorSupplier {
    LogScribeFileDescriptorSupplier() {}
  }

  private static final class LogScribeMethodDescriptorSupplier
      extends LogScribeBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    LogScribeMethodDescriptorSupplier(String methodName) {
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
      synchronized (LogScribeGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new LogScribeFileDescriptorSupplier())
              .addMethod(METHOD_LOG)
              .build();
        }
      }
    }
    return result;
  }
}
