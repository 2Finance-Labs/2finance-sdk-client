#include "twofinance/sdk_client.hpp"

#include <iostream>

int main() {
  twofinance::SdkConfig config;
  config.analytics_url = "https://analytics.example";

  twofinance::Transport transport = [](const twofinance::HttpRequest& request) {
    std::cout << request.headers.at("Authorization") << "\n";
    return twofinance::HttpResponse{200, "{\"ok\":true}"};
  };

  twofinance::SdkClient client(
      config,
      transport,
      twofinance::static_token_source("example-client-credentials-token"));
  client.analytics.indicators();
  return 0;
}
