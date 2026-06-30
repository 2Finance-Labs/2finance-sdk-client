package com.twofinance.sdk;

@FunctionalInterface
public interface TokenSource {
    String token() throws Exception;
}
