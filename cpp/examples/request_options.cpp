#include "twofinance/sdk_client.hpp"

#include <iostream>

int main() {
  twofinance::SdkConfig config;
  config.analytics_url = "https://analytics.example";

  twofinance::Transport transport = [](const twofinance::HttpRequest& request) {
    std::cout << request.method << " " << request.url << "\n";
    std::cout << request.headers.at("Idempotency-Key") << "\n";
    return twofinance::HttpResponse{200, "{\"ok\":true}"};
  };

  twofinance::RequestOptions options;
  options.headers["X-Trace-ID"] = "trace-1";
  options.idempotency_key = "candles-upsert-001";
  options.query["source"] = "sdk-example";
  options.page = 1;
  options.limit = 25;
  options.timeout_ms = 5000;
  options.max_retries = 1;

  twofinance::SdkClient client(config, transport);
  client.analytics.post("/analytics/candles:upsert", "{\"symbol\":\"BTC-USDT\"}", options);
  return 0;
}
