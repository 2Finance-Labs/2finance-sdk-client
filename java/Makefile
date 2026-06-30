.PHONY: test clean

BUILD_DIR := build/classes
JAVA_SOURCES := $(shell find src/main/java src/test/java examples -name '*.java')

test:
	mkdir -p $(BUILD_DIR)
	javac -d $(BUILD_DIR) $(JAVA_SOURCES)
	java -cp $(BUILD_DIR) com.twofinance.sdk.SdkClientTest

clean:
	rm -rf build
