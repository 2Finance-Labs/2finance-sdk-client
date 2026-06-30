#include "twofinance/sdk_client.hpp"

#include <iostream>

int main() {
  twofinance::SdkConfig config;
  config.analytics_url = "https://analytics.example";

  twofinance::Transport transport = [](const twofinance::HttpRequest&) {
    return twofinance::HttpResponse{429, "{\"error\":\"rate_limited\"}"};
  };

  twofinance::SdkClient client(config, transport);
  try {
    client.analytics.indicators();
  } catch (const twofinance::ServiceError& error) {
    std::cout << "request failed with status " << error.status_code() << ": " << error.body() << "\n";
  }

  return 0;
}
