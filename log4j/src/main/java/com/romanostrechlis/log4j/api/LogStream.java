package com.romanostrechlis.log4j.api;

public final class LogStream {
  private LogStream() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  public interface LogRequestOrBuilder extends
      // @@protoc_insertion_point(interface_extends:api.LogRequest)
      com.google.protobuf.MessageOrBuilder {

    /**
     * <code>string filename = 1;</code>
     */
    String getFilename();
    /**
     * <code>string filename = 1;</code>
     */
    com.google.protobuf.ByteString
    getFilenameBytes();

    /**
     * <code>string path = 2;</code>
     */
    String getPath();
    /**
     * <code>string path = 2;</code>
     */
    com.google.protobuf.ByteString
    getPathBytes();

    /**
     * <code>string line = 3;</code>
     */
    String getLine();
    /**
     * <code>string line = 3;</code>
     */
    com.google.protobuf.ByteString
    getLineBytes();
  }
  /**
   * Protobuf type {@code api.LogRequest}
   */
  public  static final class LogRequest extends
      com.google.protobuf.GeneratedMessageV3 implements
      // @@protoc_insertion_point(message_implements:api.LogRequest)
      LogRequestOrBuilder {
    // Use LogRequest.newBuilder() to construct.
    private LogRequest(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
      super(builder);
    }
    private LogRequest() {
      filename_ = "";
      path_ = "";
      line_ = "";
    }

    @Override
    public final com.google.protobuf.UnknownFieldSet
    getUnknownFields() {
      return com.google.protobuf.UnknownFieldSet.getDefaultInstance();
    }
    private LogRequest(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      this();
      int mutable_bitField0_ = 0;
      try {
        boolean done = false;
        while (!done) {
          int tag = input.readTag();
          switch (tag) {
            case 0:
              done = true;
              break;
            default: {
              if (!input.skipField(tag)) {
                done = true;
              }
              break;
            }
            case 10: {
              String s = input.readStringRequireUtf8();

              filename_ = s;
              break;
            }
            case 18: {
              String s = input.readStringRequireUtf8();

              path_ = s;
              break;
            }
            case 26: {
              String s = input.readStringRequireUtf8();

              line_ = s;
              break;
            }
          }
        }
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        throw e.setUnfinishedMessage(this);
      } catch (java.io.IOException e) {
        throw new com.google.protobuf.InvalidProtocolBufferException(
            e).setUnfinishedMessage(this);
      } finally {
        makeExtensionsImmutable();
      }
    }
    public static final com.google.protobuf.Descriptors.Descriptor
    getDescriptor() {
      return LogStream.internal_static_api_LogRequest_descriptor;
    }

    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
    internalGetFieldAccessorTable() {
      return LogStream.internal_static_api_LogRequest_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              LogStream.LogRequest.class, LogStream.LogRequest.Builder.class);
    }

    public static final int FILENAME_FIELD_NUMBER = 1;
    private volatile Object filename_;
    /**
     * <code>string filename = 1;</code>
     */
    public String getFilename() {
      Object ref = filename_;
      if (ref instanceof String) {
        return (String) ref;
      } else {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        String s = bs.toStringUtf8();
        filename_ = s;
        return s;
      }
    }
    /**
     * <code>string filename = 1;</code>
     */
    public com.google.protobuf.ByteString
    getFilenameBytes() {
      Object ref = filename_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b =
            com.google.protobuf.ByteString.copyFromUtf8(
                (String) ref);
        filename_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }

    public static final int PATH_FIELD_NUMBER = 2;
    private volatile Object path_;
    /**
     * <code>string path = 2;</code>
     */
    public String getPath() {
      Object ref = path_;
      if (ref instanceof String) {
        return (String) ref;
      } else {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        String s = bs.toStringUtf8();
        path_ = s;
        return s;
      }
    }
    /**
     * <code>string path = 2;</code>
     */
    public com.google.protobuf.ByteString
    getPathBytes() {
      Object ref = path_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b =
            com.google.protobuf.ByteString.copyFromUtf8(
                (String) ref);
        path_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }

    public static final int LINE_FIELD_NUMBER = 3;
    private volatile Object line_;
    /**
     * <code>string line = 3;</code>
     */
    public String getLine() {
      Object ref = line_;
      if (ref instanceof String) {
        return (String) ref;
      } else {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        String s = bs.toStringUtf8();
        line_ = s;
        return s;
      }
    }
    /**
     * <code>string line = 3;</code>
     */
    public com.google.protobuf.ByteString
    getLineBytes() {
      Object ref = line_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b =
            com.google.protobuf.ByteString.copyFromUtf8(
                (String) ref);
        line_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }

    private byte memoizedIsInitialized = -1;
    public final boolean isInitialized() {
      byte isInitialized = memoizedIsInitialized;
      if (isInitialized == 1) return true;
      if (isInitialized == 0) return false;

      memoizedIsInitialized = 1;
      return true;
    }

    public void writeTo(com.google.protobuf.CodedOutputStream output)
        throws java.io.IOException {
      if (!getFilenameBytes().isEmpty()) {
        com.google.protobuf.GeneratedMessageV3.writeString(output, 1, filename_);
      }
      if (!getPathBytes().isEmpty()) {
        com.google.protobuf.GeneratedMessageV3.writeString(output, 2, path_);
      }
      if (!getLineBytes().isEmpty()) {
        com.google.protobuf.GeneratedMessageV3.writeString(output, 3, line_);
      }
    }

    public int getSerializedSize() {
      int size = memoizedSize;
      if (size != -1) return size;

      size = 0;
      if (!getFilenameBytes().isEmpty()) {
        size += com.google.protobuf.GeneratedMessageV3.computeStringSize(1, filename_);
      }
      if (!getPathBytes().isEmpty()) {
        size += com.google.protobuf.GeneratedMessageV3.computeStringSize(2, path_);
      }
      if (!getLineBytes().isEmpty()) {
        size += com.google.protobuf.GeneratedMessageV3.computeStringSize(3, line_);
      }
      memoizedSize = size;
      return size;
    }

    private static final long serialVersionUID = 0L;
    @Override
    public boolean equals(final Object obj) {
      if (obj == this) {
        return true;
      }
      if (!(obj instanceof LogStream.LogRequest)) {
        return super.equals(obj);
      }
      LogStream.LogRequest other = (LogStream.LogRequest) obj;

      boolean result = true;
      result = result && getFilename()
          .equals(other.getFilename());
      result = result && getPath()
          .equals(other.getPath());
      result = result && getLine()
          .equals(other.getLine());
      return result;
    }

    @Override
    public int hashCode() {
      if (memoizedHashCode != 0) {
        return memoizedHashCode;
      }
      int hash = 41;
      hash = (19 * hash) + getDescriptor().hashCode();
      hash = (37 * hash) + FILENAME_FIELD_NUMBER;
      hash = (53 * hash) + getFilename().hashCode();
      hash = (37 * hash) + PATH_FIELD_NUMBER;
      hash = (53 * hash) + getPath().hashCode();
      hash = (37 * hash) + LINE_FIELD_NUMBER;
      hash = (53 * hash) + getLine().hashCode();
      hash = (29 * hash) + unknownFields.hashCode();
      memoizedHashCode = hash;
      return hash;
    }

    public static LogStream.LogRequest parseFrom(
        java.nio.ByteBuffer data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static LogStream.LogRequest parseFrom(
        java.nio.ByteBuffer data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static LogStream.LogRequest parseFrom(
        com.google.protobuf.ByteString data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static LogStream.LogRequest parseFrom(
        com.google.protobuf.ByteString data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static LogStream.LogRequest parseFrom(byte[] data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static LogStream.LogRequest parseFrom(
        byte[] data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static LogStream.LogRequest parseFrom(java.io.InputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input);
    }
    public static LogStream.LogRequest parseFrom(
        java.io.InputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input, extensionRegistry);
    }
    public static LogStream.LogRequest parseDelimitedFrom(java.io.InputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseDelimitedWithIOException(PARSER, input);
    }
    public static LogStream.LogRequest parseDelimitedFrom(
        java.io.InputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
    }
    public static LogStream.LogRequest parseFrom(
        com.google.protobuf.CodedInputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input);
    }
    public static LogStream.LogRequest parseFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input, extensionRegistry);
    }

    public Builder newBuilderForType() { return newBuilder(); }
    public static Builder newBuilder() {
      return DEFAULT_INSTANCE.toBuilder();
    }
    public static Builder newBuilder(LogStream.LogRequest prototype) {
      return DEFAULT_INSTANCE.toBuilder().mergeFrom(prototype);
    }
    public Builder toBuilder() {
      return this == DEFAULT_INSTANCE
          ? new Builder() : new Builder().mergeFrom(this);
    }

    @Override
    protected Builder newBuilderForType(
        com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
      Builder builder = new Builder(parent);
      return builder;
    }
    /**
     * Protobuf type {@code api.LogRequest}
     */
    public static final class Builder extends
        com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
        // @@protoc_insertion_point(builder_implements:api.LogRequest)
        LogStream.LogRequestOrBuilder {
      public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
        return LogStream.internal_static_api_LogRequest_descriptor;
      }

      protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
        return LogStream.internal_static_api_LogRequest_fieldAccessorTable
            .ensureFieldAccessorsInitialized(
                LogStream.LogRequest.class, LogStream.LogRequest.Builder.class);
      }

      // Construct using api.LogStream.LogRequest.newBuilder()
      private Builder() {
        maybeForceBuilderInitialization();
      }

      private Builder(
          com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
        super(parent);
        maybeForceBuilderInitialization();
      }
      private void maybeForceBuilderInitialization() {
        if (com.google.protobuf.GeneratedMessageV3
            .alwaysUseFieldBuilders) {
        }
      }
      public Builder clear() {
        super.clear();
        filename_ = "";

        path_ = "";

        line_ = "";

        return this;
      }

      public com.google.protobuf.Descriptors.Descriptor
      getDescriptorForType() {
        return LogStream.internal_static_api_LogRequest_descriptor;
      }

      public LogStream.LogRequest getDefaultInstanceForType() {
        return LogStream.LogRequest.getDefaultInstance();
      }

      public LogStream.LogRequest build() {
        LogStream.LogRequest result = buildPartial();
        if (!result.isInitialized()) {
          throw newUninitializedMessageException(result);
        }
        return result;
      }

      public LogStream.LogRequest buildPartial() {
        LogStream.LogRequest result = new LogStream.LogRequest(this);
        result.filename_ = filename_;
        result.path_ = path_;
        result.line_ = line_;
        onBuilt();
        return result;
      }

      public Builder clone() {
        return (Builder) super.clone();
      }
      public Builder setField(
          com.google.protobuf.Descriptors.FieldDescriptor field,
          Object value) {
        return (Builder) super.setField(field, value);
      }
      public Builder clearField(
          com.google.protobuf.Descriptors.FieldDescriptor field) {
        return (Builder) super.clearField(field);
      }
      public Builder clearOneof(
          com.google.protobuf.Descriptors.OneofDescriptor oneof) {
        return (Builder) super.clearOneof(oneof);
      }
      public Builder setRepeatedField(
          com.google.protobuf.Descriptors.FieldDescriptor field,
          int index, Object value) {
        return (Builder) super.setRepeatedField(field, index, value);
      }
      public Builder addRepeatedField(
          com.google.protobuf.Descriptors.FieldDescriptor field,
          Object value) {
        return (Builder) super.addRepeatedField(field, value);
      }
      public Builder mergeFrom(com.google.protobuf.Message other) {
        if (other instanceof LogStream.LogRequest) {
          return mergeFrom((LogStream.LogRequest)other);
        } else {
          super.mergeFrom(other);
          return this;
        }
      }

      public Builder mergeFrom(LogStream.LogRequest other) {
        if (other == LogStream.LogRequest.getDefaultInstance()) return this;
        if (!other.getFilename().isEmpty()) {
          filename_ = other.filename_;
          onChanged();
        }
        if (!other.getPath().isEmpty()) {
          path_ = other.path_;
          onChanged();
        }
        if (!other.getLine().isEmpty()) {
          line_ = other.line_;
          onChanged();
        }
        onChanged();
        return this;
      }

      public final boolean isInitialized() {
        return true;
      }

      public Builder mergeFrom(
          com.google.protobuf.CodedInputStream input,
          com.google.protobuf.ExtensionRegistryLite extensionRegistry)
          throws java.io.IOException {
        LogStream.LogRequest parsedMessage = null;
        try {
          parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
        } catch (com.google.protobuf.InvalidProtocolBufferException e) {
          parsedMessage = (LogStream.LogRequest) e.getUnfinishedMessage();
          throw e.unwrapIOException();
        } finally {
          if (parsedMessage != null) {
            mergeFrom(parsedMessage);
          }
        }
        return this;
      }

      private Object filename_ = "";
      /**
       * <code>string filename = 1;</code>
       */
      public String getFilename() {
        Object ref = filename_;
        if (!(ref instanceof String)) {
          com.google.protobuf.ByteString bs =
              (com.google.protobuf.ByteString) ref;
          String s = bs.toStringUtf8();
          filename_ = s;
          return s;
        } else {
          return (String) ref;
        }
      }
      /**
       * <code>string filename = 1;</code>
       */
      public com.google.protobuf.ByteString
      getFilenameBytes() {
        Object ref = filename_;
        if (ref instanceof String) {
          com.google.protobuf.ByteString b =
              com.google.protobuf.ByteString.copyFromUtf8(
                  (String) ref);
          filename_ = b;
          return b;
        } else {
          return (com.google.protobuf.ByteString) ref;
        }
      }
      /**
       * <code>string filename = 1;</code>
       */
      public Builder setFilename(
          String value) {
        if (value == null) {
          throw new NullPointerException();
        }

        filename_ = value;
        onChanged();
        return this;
      }
      /**
       * <code>string filename = 1;</code>
       */
      public Builder clearFilename() {

        filename_ = getDefaultInstance().getFilename();
        onChanged();
        return this;
      }
      /**
       * <code>string filename = 1;</code>
       */
      public Builder setFilenameBytes(
          com.google.protobuf.ByteString value) {
        if (value == null) {
          throw new NullPointerException();
        }
        checkByteStringIsUtf8(value);

        filename_ = value;
        onChanged();
        return this;
      }

      private Object path_ = "";
      /**
       * <code>string path = 2;</code>
       */
      public String getPath() {
        Object ref = path_;
        if (!(ref instanceof String)) {
          com.google.protobuf.ByteString bs =
              (com.google.protobuf.ByteString) ref;
          String s = bs.toStringUtf8();
          path_ = s;
          return s;
        } else {
          return (String) ref;
        }
      }
      /**
       * <code>string path = 2;</code>
       */
      public com.google.protobuf.ByteString
      getPathBytes() {
        Object ref = path_;
        if (ref instanceof String) {
          com.google.protobuf.ByteString b =
              com.google.protobuf.ByteString.copyFromUtf8(
                  (String) ref);
          path_ = b;
          return b;
        } else {
          return (com.google.protobuf.ByteString) ref;
        }
      }
      /**
       * <code>string path = 2;</code>
       */
      public Builder setPath(
          String value) {
        if (value == null) {
          throw new NullPointerException();
        }

        path_ = value;
        onChanged();
        return this;
      }
      /**
       * <code>string path = 2;</code>
       */
      public Builder clearPath() {

        path_ = getDefaultInstance().getPath();
        onChanged();
        return this;
      }
      /**
       * <code>string path = 2;</code>
       */
      public Builder setPathBytes(
          com.google.protobuf.ByteString value) {
        if (value == null) {
          throw new NullPointerException();
        }
        checkByteStringIsUtf8(value);

        path_ = value;
        onChanged();
        return this;
      }

      private Object line_ = "";
      /**
       * <code>string line = 3;</code>
       */
      public String getLine() {
        Object ref = line_;
        if (!(ref instanceof String)) {
          com.google.protobuf.ByteString bs =
              (com.google.protobuf.ByteString) ref;
          String s = bs.toStringUtf8();
          line_ = s;
          return s;
        } else {
          return (String) ref;
        }
      }
      /**
       * <code>string line = 3;</code>
       */
      public com.google.protobuf.ByteString
      getLineBytes() {
        Object ref = line_;
        if (ref instanceof String) {
          com.google.protobuf.ByteString b =
              com.google.protobuf.ByteString.copyFromUtf8(
                  (String) ref);
          line_ = b;
          return b;
        } else {
          return (com.google.protobuf.ByteString) ref;
        }
      }
      /**
       * <code>string line = 3;</code>
       */
      public Builder setLine(
          String value) {
        if (value == null) {
          throw new NullPointerException();
        }

        line_ = value;
        onChanged();
        return this;
      }
      /**
       * <code>string line = 3;</code>
       */
      public Builder clearLine() {

        line_ = getDefaultInstance().getLine();
        onChanged();
        return this;
      }
      /**
       * <code>string line = 3;</code>
       */
      public Builder setLineBytes(
          com.google.protobuf.ByteString value) {
        if (value == null) {
          throw new NullPointerException();
        }
        checkByteStringIsUtf8(value);

        line_ = value;
        onChanged();
        return this;
      }
      public final Builder setUnknownFields(
          final com.google.protobuf.UnknownFieldSet unknownFields) {
        return this;
      }

      public final Builder mergeUnknownFields(
          final com.google.protobuf.UnknownFieldSet unknownFields) {
        return this;
      }


      // @@protoc_insertion_point(builder_scope:api.LogRequest)
    }

    // @@protoc_insertion_point(class_scope:api.LogRequest)
    private static final LogStream.LogRequest DEFAULT_INSTANCE;
    static {
      DEFAULT_INSTANCE = new LogStream.LogRequest();
    }

    public static LogStream.LogRequest getDefaultInstance() {
      return DEFAULT_INSTANCE;
    }

    private static final com.google.protobuf.Parser<LogRequest>
        PARSER = new com.google.protobuf.AbstractParser<LogRequest>() {
      public LogRequest parsePartialFrom(
          com.google.protobuf.CodedInputStream input,
          com.google.protobuf.ExtensionRegistryLite extensionRegistry)
          throws com.google.protobuf.InvalidProtocolBufferException {
        return new LogRequest(input, extensionRegistry);
      }
    };

    public static com.google.protobuf.Parser<LogRequest> parser() {
      return PARSER;
    }

    @Override
    public com.google.protobuf.Parser<LogRequest> getParserForType() {
      return PARSER;
    }

    public LogStream.LogRequest getDefaultInstanceForType() {
      return DEFAULT_INSTANCE;
    }

  }

  public interface LogResponseOrBuilder extends
      // @@protoc_insertion_point(interface_extends:api.LogResponse)
      com.google.protobuf.MessageOrBuilder {

    /**
     * <code>string res = 1;</code>
     */
    String getRes();
    /**
     * <code>string res = 1;</code>
     */
    com.google.protobuf.ByteString
    getResBytes();
  }
  /**
   * Protobuf type {@code api.LogResponse}
   */
  public  static final class LogResponse extends
      com.google.protobuf.GeneratedMessageV3 implements
      // @@protoc_insertion_point(message_implements:api.LogResponse)
      LogResponseOrBuilder {
    // Use LogResponse.newBuilder() to construct.
    private LogResponse(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
      super(builder);
    }
    private LogResponse() {
      res_ = "";
    }

    @Override
    public final com.google.protobuf.UnknownFieldSet
    getUnknownFields() {
      return com.google.protobuf.UnknownFieldSet.getDefaultInstance();
    }
    private LogResponse(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      this();
      int mutable_bitField0_ = 0;
      try {
        boolean done = false;
        while (!done) {
          int tag = input.readTag();
          switch (tag) {
            case 0:
              done = true;
              break;
            default: {
              if (!input.skipField(tag)) {
                done = true;
              }
              break;
            }
            case 10: {
              String s = input.readStringRequireUtf8();

              res_ = s;
              break;
            }
          }
        }
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        throw e.setUnfinishedMessage(this);
      } catch (java.io.IOException e) {
        throw new com.google.protobuf.InvalidProtocolBufferException(
            e).setUnfinishedMessage(this);
      } finally {
        makeExtensionsImmutable();
      }
    }
    public static final com.google.protobuf.Descriptors.Descriptor
    getDescriptor() {
      return LogStream.internal_static_api_LogResponse_descriptor;
    }

    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
    internalGetFieldAccessorTable() {
      return LogStream.internal_static_api_LogResponse_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              LogStream.LogResponse.class, LogStream.LogResponse.Builder.class);
    }

    public static final int RES_FIELD_NUMBER = 1;
    private volatile Object res_;
    /**
     * <code>string res = 1;</code>
     */
    public String getRes() {
      Object ref = res_;
      if (ref instanceof String) {
        return (String) ref;
      } else {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        String s = bs.toStringUtf8();
        res_ = s;
        return s;
      }
    }
    /**
     * <code>string res = 1;</code>
     */
    public com.google.protobuf.ByteString
    getResBytes() {
      Object ref = res_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b =
            com.google.protobuf.ByteString.copyFromUtf8(
                (String) ref);
        res_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }

    private byte memoizedIsInitialized = -1;
    public final boolean isInitialized() {
      byte isInitialized = memoizedIsInitialized;
      if (isInitialized == 1) return true;
      if (isInitialized == 0) return false;

      memoizedIsInitialized = 1;
      return true;
    }

    public void writeTo(com.google.protobuf.CodedOutputStream output)
        throws java.io.IOException {
      if (!getResBytes().isEmpty()) {
        com.google.protobuf.GeneratedMessageV3.writeString(output, 1, res_);
      }
    }

    public int getSerializedSize() {
      int size = memoizedSize;
      if (size != -1) return size;

      size = 0;
      if (!getResBytes().isEmpty()) {
        size += com.google.protobuf.GeneratedMessageV3.computeStringSize(1, res_);
      }
      memoizedSize = size;
      return size;
    }

    private static final long serialVersionUID = 0L;
    @Override
    public boolean equals(final Object obj) {
      if (obj == this) {
        return true;
      }
      if (!(obj instanceof LogStream.LogResponse)) {
        return super.equals(obj);
      }
      LogStream.LogResponse other = (LogStream.LogResponse) obj;

      boolean result = true;
      result = result && getRes()
          .equals(other.getRes());
      return result;
    }

    @Override
    public int hashCode() {
      if (memoizedHashCode != 0) {
        return memoizedHashCode;
      }
      int hash = 41;
      hash = (19 * hash) + getDescriptor().hashCode();
      hash = (37 * hash) + RES_FIELD_NUMBER;
      hash = (53 * hash) + getRes().hashCode();
      hash = (29 * hash) + unknownFields.hashCode();
      memoizedHashCode = hash;
      return hash;
    }

    public static LogStream.LogResponse parseFrom(
        java.nio.ByteBuffer data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static LogStream.LogResponse parseFrom(
        java.nio.ByteBuffer data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static LogStream.LogResponse parseFrom(
        com.google.protobuf.ByteString data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static LogStream.LogResponse parseFrom(
        com.google.protobuf.ByteString data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static LogStream.LogResponse parseFrom(byte[] data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static LogStream.LogResponse parseFrom(
        byte[] data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static LogStream.LogResponse parseFrom(java.io.InputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input);
    }
    public static LogStream.LogResponse parseFrom(
        java.io.InputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input, extensionRegistry);
    }
    public static LogStream.LogResponse parseDelimitedFrom(java.io.InputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseDelimitedWithIOException(PARSER, input);
    }
    public static LogStream.LogResponse parseDelimitedFrom(
        java.io.InputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
    }
    public static LogStream.LogResponse parseFrom(
        com.google.protobuf.CodedInputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input);
    }
    public static LogStream.LogResponse parseFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input, extensionRegistry);
    }

    public Builder newBuilderForType() { return newBuilder(); }
    public static Builder newBuilder() {
      return DEFAULT_INSTANCE.toBuilder();
    }
    public static Builder newBuilder(LogStream.LogResponse prototype) {
      return DEFAULT_INSTANCE.toBuilder().mergeFrom(prototype);
    }
    public Builder toBuilder() {
      return this == DEFAULT_INSTANCE
          ? new Builder() : new Builder().mergeFrom(this);
    }

    @Override
    protected Builder newBuilderForType(
        com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
      Builder builder = new Builder(parent);
      return builder;
    }
    /**
     * Protobuf type {@code api.LogResponse}
     */
    public static final class Builder extends
        com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
        // @@protoc_insertion_point(builder_implements:api.LogResponse)
        LogStream.LogResponseOrBuilder {
      public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
        return LogStream.internal_static_api_LogResponse_descriptor;
      }

      protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
        return LogStream.internal_static_api_LogResponse_fieldAccessorTable
            .ensureFieldAccessorsInitialized(
                LogStream.LogResponse.class, LogStream.LogResponse.Builder.class);
      }

      // Construct using api.LogStream.LogResponse.newBuilder()
      private Builder() {
        maybeForceBuilderInitialization();
      }

      private Builder(
          com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
        super(parent);
        maybeForceBuilderInitialization();
      }
      private void maybeForceBuilderInitialization() {
        if (com.google.protobuf.GeneratedMessageV3
            .alwaysUseFieldBuilders) {
        }
      }
      public Builder clear() {
        super.clear();
        res_ = "";

        return this;
      }

      public com.google.protobuf.Descriptors.Descriptor
      getDescriptorForType() {
        return LogStream.internal_static_api_LogResponse_descriptor;
      }

      public LogStream.LogResponse getDefaultInstanceForType() {
        return LogStream.LogResponse.getDefaultInstance();
      }

      public LogStream.LogResponse build() {
        LogStream.LogResponse result = buildPartial();
        if (!result.isInitialized()) {
          throw newUninitializedMessageException(result);
        }
        return result;
      }

      public LogStream.LogResponse buildPartial() {
        LogStream.LogResponse result = new LogStream.LogResponse(this);
        result.res_ = res_;
        onBuilt();
        return result;
      }

      public Builder clone() {
        return (Builder) super.clone();
      }
      public Builder setField(
          com.google.protobuf.Descriptors.FieldDescriptor field,
          Object value) {
        return (Builder) super.setField(field, value);
      }
      public Builder clearField(
          com.google.protobuf.Descriptors.FieldDescriptor field) {
        return (Builder) super.clearField(field);
      }
      public Builder clearOneof(
          com.google.protobuf.Descriptors.OneofDescriptor oneof) {
        return (Builder) super.clearOneof(oneof);
      }
      public Builder setRepeatedField(
          com.google.protobuf.Descriptors.FieldDescriptor field,
          int index, Object value) {
        return (Builder) super.setRepeatedField(field, index, value);
      }
      public Builder addRepeatedField(
          com.google.protobuf.Descriptors.FieldDescriptor field,
          Object value) {
        return (Builder) super.addRepeatedField(field, value);
      }
      public Builder mergeFrom(com.google.protobuf.Message other) {
        if (other instanceof LogStream.LogResponse) {
          return mergeFrom((LogStream.LogResponse)other);
        } else {
          super.mergeFrom(other);
          return this;
        }
      }

      public Builder mergeFrom(LogStream.LogResponse other) {
        if (other == LogStream.LogResponse.getDefaultInstance()) return this;
        if (!other.getRes().isEmpty()) {
          res_ = other.res_;
          onChanged();
        }
        onChanged();
        return this;
      }

      public final boolean isInitialized() {
        return true;
      }

      public Builder mergeFrom(
          com.google.protobuf.CodedInputStream input,
          com.google.protobuf.ExtensionRegistryLite extensionRegistry)
          throws java.io.IOException {
        LogStream.LogResponse parsedMessage = null;
        try {
          parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
        } catch (com.google.protobuf.InvalidProtocolBufferException e) {
          parsedMessage = (LogStream.LogResponse) e.getUnfinishedMessage();
          throw e.unwrapIOException();
        } finally {
          if (parsedMessage != null) {
            mergeFrom(parsedMessage);
          }
        }
        return this;
      }

      private Object res_ = "";
      /**
       * <code>string res = 1;</code>
       */
      public String getRes() {
        Object ref = res_;
        if (!(ref instanceof String)) {
          com.google.protobuf.ByteString bs =
              (com.google.protobuf.ByteString) ref;
          String s = bs.toStringUtf8();
          res_ = s;
          return s;
        } else {
          return (String) ref;
        }
      }
      /**
       * <code>string res = 1;</code>
       */
      public com.google.protobuf.ByteString
      getResBytes() {
        Object ref = res_;
        if (ref instanceof String) {
          com.google.protobuf.ByteString b =
              com.google.protobuf.ByteString.copyFromUtf8(
                  (String) ref);
          res_ = b;
          return b;
        } else {
          return (com.google.protobuf.ByteString) ref;
        }
      }
      /**
       * <code>string res = 1;</code>
       */
      public Builder setRes(
          String value) {
        if (value == null) {
          throw new NullPointerException();
        }

        res_ = value;
        onChanged();
        return this;
      }
      /**
       * <code>string res = 1;</code>
       */
      public Builder clearRes() {

        res_ = getDefaultInstance().getRes();
        onChanged();
        return this;
      }
      /**
       * <code>string res = 1;</code>
       */
      public Builder setResBytes(
          com.google.protobuf.ByteString value) {
        if (value == null) {
          throw new NullPointerException();
        }
        checkByteStringIsUtf8(value);

        res_ = value;
        onChanged();
        return this;
      }
      public final Builder setUnknownFields(
          final com.google.protobuf.UnknownFieldSet unknownFields) {
        return this;
      }

      public final Builder mergeUnknownFields(
          final com.google.protobuf.UnknownFieldSet unknownFields) {
        return this;
      }


      // @@protoc_insertion_point(builder_scope:api.LogResponse)
    }

    // @@protoc_insertion_point(class_scope:api.LogResponse)
    private static final LogStream.LogResponse DEFAULT_INSTANCE;
    static {
      DEFAULT_INSTANCE = new LogStream.LogResponse();
    }

    public static LogStream.LogResponse getDefaultInstance() {
      return DEFAULT_INSTANCE;
    }

    private static final com.google.protobuf.Parser<LogResponse>
        PARSER = new com.google.protobuf.AbstractParser<LogResponse>() {
      public LogResponse parsePartialFrom(
          com.google.protobuf.CodedInputStream input,
          com.google.protobuf.ExtensionRegistryLite extensionRegistry)
          throws com.google.protobuf.InvalidProtocolBufferException {
        return new LogResponse(input, extensionRegistry);
      }
    };

    public static com.google.protobuf.Parser<LogResponse> parser() {
      return PARSER;
    }

    @Override
    public com.google.protobuf.Parser<LogResponse> getParserForType() {
      return PARSER;
    }

    public LogStream.LogResponse getDefaultInstanceForType() {
      return DEFAULT_INSTANCE;
    }

  }

  private static final com.google.protobuf.Descriptors.Descriptor
      internal_static_api_LogRequest_descriptor;
  private static final
  com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_api_LogRequest_fieldAccessorTable;
  private static final com.google.protobuf.Descriptors.Descriptor
      internal_static_api_LogResponse_descriptor;
  private static final
  com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_api_LogResponse_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
  getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    String[] descriptorData = {
        "\n\017logStream.proto\022\003api\":\n\nLogRequest\022\020\n\010" +
            "filename\030\001 \001(\t\022\014\n\004path\030\002 \001(\t\022\014\n\004line\030\003 \001" +
            "(\t\"\032\n\013LogResponse\022\013\n\003res\030\001 \001(\t29\n\013LogStr" +
            "eamer\022*\n\003Log\022\017.api.LogRequest\032\020.api.LogR" +
            "esponse\"\000b\006proto3"
    };
    com.google.protobuf.Descriptors.FileDescriptor.InternalDescriptorAssigner assigner =
        new com.google.protobuf.Descriptors.FileDescriptor.    InternalDescriptorAssigner() {
          public com.google.protobuf.ExtensionRegistry assignDescriptors(
              com.google.protobuf.Descriptors.FileDescriptor root) {
            descriptor = root;
            return null;
          }
        };
    com.google.protobuf.Descriptors.FileDescriptor
        .internalBuildGeneratedFileFrom(descriptorData,
            new com.google.protobuf.Descriptors.FileDescriptor[] {
            }, assigner);
    internal_static_api_LogRequest_descriptor =
        getDescriptor().getMessageTypes().get(0);
    internal_static_api_LogRequest_fieldAccessorTable = new
        com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_api_LogRequest_descriptor,
        new String[] { "Filename", "Path", "Line", });
    internal_static_api_LogResponse_descriptor =
        getDescriptor().getMessageTypes().get(1);
    internal_static_api_LogResponse_fieldAccessorTable = new
        com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_api_LogResponse_descriptor,
        new String[] { "Res", });
  }

  // @@protoc_insertion_point(outer_class_scope)
}