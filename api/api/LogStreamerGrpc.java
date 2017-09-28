package api;

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
 * service defines a new service names LogStreamer
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.6.1)",
    comments = "Source: logStream.proto")
public final class LogStreamerGrpc {

  private LogStreamerGrpc() {}

  public static final String SERVICE_NAME = "api.LogStreamer";

  // Static method descriptors that strictly reflect the proto.
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static final io.grpc.MethodDescriptor<api.LogStream.LogRequest,
      api.LogStream.LogResponse> METHOD_LOG =
      io.grpc.MethodDescriptor.<api.LogStream.LogRequest, api.LogStream.LogResponse>newBuilder()
          .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
          .setFullMethodName(generateFullMethodName(
              "api.LogStreamer", "Log"))
          .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
              api.LogStream.LogRequest.getDefaultInstance()))
          .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
              api.LogStream.LogResponse.getDefaultInstance()))
          .build();

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static LogStreamerStub newStub(io.grpc.Channel channel) {
    return new LogStreamerStub(channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static LogStreamerBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    return new LogStreamerBlockingStub(channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static LogStreamerFutureStub newFutureStub(
      io.grpc.Channel channel) {
    return new LogStreamerFutureStub(channel);
  }

  /**
   * <pre>
   * service defines a new service names LogStreamer
   * </pre>
   */
  public static abstract class LogStreamerImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * the LogStreamer service sends a LogRequest
     * and recieves a LogResponse
     * </pre>
     */
    public void log(api.LogStream.LogRequest request,
        io.grpc.stub.StreamObserver<api.LogStream.LogResponse> responseObserver) {
      asyncUnimplementedUnaryCall(METHOD_LOG, responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            METHOD_LOG,
            asyncUnaryCall(
              new MethodHandlers<
                api.LogStream.LogRequest,
                api.LogStream.LogResponse>(
                  this, METHODID_LOG)))
          .build();
    }
  }

  /**
   * <pre>
   * service defines a new service names LogStreamer
   * </pre>
   */
  public static final class LogStreamerStub extends io.grpc.stub.AbstractStub<LogStreamerStub> {
    private LogStreamerStub(io.grpc.Channel channel) {
      super(channel);
    }

    private LogStreamerStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected LogStreamerStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new LogStreamerStub(channel, callOptions);
    }

    /**
     * <pre>
     * the LogStreamer service sends a LogRequest
     * and recieves a LogResponse
     * </pre>
     */
    public void log(api.LogStream.LogRequest request,
        io.grpc.stub.StreamObserver<api.LogStream.LogResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_LOG, getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * service defines a new service names LogStreamer
   * </pre>
   */
  public static final class LogStreamerBlockingStub extends io.grpc.stub.AbstractStub<LogStreamerBlockingStub> {
    private LogStreamerBlockingStub(io.grpc.Channel channel) {
      super(channel);
    }

    private LogStreamerBlockingStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected LogStreamerBlockingStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new LogStreamerBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * the LogStreamer service sends a LogRequest
     * and recieves a LogResponse
     * </pre>
     */
    public api.LogStream.LogResponse log(api.LogStream.LogRequest request) {
      return blockingUnaryCall(
          getChannel(), METHOD_LOG, getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * service defines a new service names LogStreamer
   * </pre>
   */
  public static final class LogStreamerFutureStub extends io.grpc.stub.AbstractStub<LogStreamerFutureStub> {
    private LogStreamerFutureStub(io.grpc.Channel channel) {
      super(channel);
    }

    private LogStreamerFutureStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected LogStreamerFutureStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new LogStreamerFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * the LogStreamer service sends a LogRequest
     * and recieves a LogResponse
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<api.LogStream.LogResponse> log(
        api.LogStream.LogRequest request) {
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
    private final LogStreamerImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(LogStreamerImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_LOG:
          serviceImpl.log((api.LogStream.LogRequest) request,
              (io.grpc.stub.StreamObserver<api.LogStream.LogResponse>) responseObserver);
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

  private static final class LogStreamerDescriptorSupplier implements io.grpc.protobuf.ProtoFileDescriptorSupplier {
    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return api.LogStream.getDescriptor();
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (LogStreamerGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new LogStreamerDescriptorSupplier())
              .addMethod(METHOD_LOG)
              .build();
        }
      }
    }
    return result;
  }
}