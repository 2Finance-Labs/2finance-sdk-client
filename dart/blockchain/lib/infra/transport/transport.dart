abstract class FinanceNetworkTransport {
  Future<dynamic> sendRequest(String method, dynamic params, String replyTo);
}
