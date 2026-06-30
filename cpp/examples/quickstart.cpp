#include "twofinance/sdk_client.hpp"

#include <iostream>

int main() {
  twofinance::SdkConfig config = twofinance::config_from_environment();

  twofinance::Transport transport = [](const twofinance::HttpRequest& request) {
    std::cout << request.method << " " << request.url << "\n";
    return twofinance::HttpResponse{200, "{\"ok\":true}"};
  };

  twofinance::SdkClient client(config, transport);
  client.analytics.indicators();
  client.planner.trading_plan(
      "{\"goal\":\"prepare a BTC rebalancing plan\",\"useAnalytics\":true,\"useTrading\":true}",
      true,
      true);

  return 0;
}
