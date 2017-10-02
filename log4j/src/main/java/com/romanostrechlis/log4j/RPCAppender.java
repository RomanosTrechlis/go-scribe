package com.romanostrechlis.log4j;

import com.romanostrechlis.log4j.api.LogStream;
import com.romanostrechlis.log4j.api.LogStreamClient;

import org.apache.log4j.FileAppender;
import org.apache.log4j.Layout;
import org.apache.log4j.spi.LoggingEvent;

import java.io.IOException;
import java.text.SimpleDateFormat;
import java.util.Date;

public class RPCAppender extends FileAppender {

  private String host = "localhost";
  private int port = 8080;

  public RPCAppender() {
  }

  public RPCAppender(Layout layout,
                     String filename,
                     boolean append,
                     boolean bufferedIO,
                     int bufferSize) throws IOException {
    super(layout, filename, append, bufferedIO, bufferSize);
  }

  public RPCAppender(Layout layout,
                     String filename,
                     boolean append) throws IOException {
    super(layout, filename, append);
  }

  public RPCAppender(Layout layout,
                     String filename) throws IOException {
    super(layout, filename);
  }

  public String getHost() {
    return host;
  }

  public void setHost(String host) {
    this.host = host;
  }

  public int getPort() {
    return port;
  }

  public void setPort(int port) {
    this.port = port;
  }

  protected void subAppend(LoggingEvent event) {
    LogStream.LogRequest req = LogStream.LogRequest.newBuilder()
        .setFilename(this.getName())
        .setPath(this.getFile())
        .setLine(customFormat(event))
        .build();
    try {
      LogStreamClient client = new LogStreamClient(host, port);
      client.send(req);
    } catch (Exception e) {
    } finally {
    }
  }

  public String customFormat(LoggingEvent event) {
    SimpleDateFormat sdf = new SimpleDateFormat("dd-MM-yyyy hh:mm:ss.S");
    StringBuilder sbuf = new StringBuilder();
    sbuf.append(sdf.format(new Date()) + " ");
    sbuf.append(event.getLevel().toString());
    sbuf.append(" - ");
    sbuf.append(event.getRenderedMessage());
    sbuf.append(Layout.LINE_SEP);
    return sbuf.toString();
  }
}
