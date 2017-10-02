package com.romanostrechlis.log4j.api;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.stub.StreamObserver;

public class LogStreamClient {

  private ManagedChannel channel;
  private LogStreamerGrpc.LogStreamerStub asyncStub;

  public LogStreamClient(String host, int port) {
    this(ManagedChannelBuilder.forAddress(host, port).usePlaintext(true));
  }

  public LogStreamClient(ManagedChannelBuilder<?> channelBuilder) {
    this.channel = channelBuilder.build();
    //blockingStub = RouteGuideGrpc.newBlockingStub(channel);
    this.asyncStub = LogStreamerGrpc.newStub(channel);

  }

  public void send(LogStream.LogRequest req) {
    LogStream.LogResponse res = LogStream.LogResponse.newBuilder().build();

    asyncStub.log(req, new StreamObserver<LogStream.LogResponse>() {
      public void onNext(LogStream.LogResponse logResponse) {
        // System.out.println("onNext");
      }
      public void onError(Throwable throwable) {
        System.out.println(throwable.getMessage());
      }

      public void onCompleted() {
        // System.out.println("onComplete");
      }
    });
  }
}
