# Java SDK

Java 11+ SDK for 2Finance services using `java.net.http`.

```java
SdkClient client = SdkClient.fromEnv();
String indicators = client.analytics.indicators();
```

Compile smoke test:

```bash
javac -d build/classes $(find src/main/java src/test/java -name '*.java')
java -cp build/classes com.twofinance.sdk.SdkClientTest
```

Build package metadata:

```bash
mvn package
```

The Maven project is dependency-free and publishes the package coordinates
`com.twofinance:sdk-client`.

See `examples/Quickstart.java` for a minimal analytics and planner flow.
See `examples/RequestOptionsExample.java` for per-call headers, idempotency,
query, pagination, timeout, and retry options.
See `examples/AuthClientCredentialsExample.java` for client credentials token
source configuration.
See `examples/ErrorHandlingExample.java` for catching `ServiceException`.
