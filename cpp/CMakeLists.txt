cmake_minimum_required(VERSION 3.16)

project(twofinance_sdk_client VERSION 0.1.0 LANGUAGES CXX)

add_library(twofinance_sdk_client INTERFACE)
add_library(twofinance::sdk_client ALIAS twofinance_sdk_client)

target_include_directories(
  twofinance_sdk_client
  INTERFACE
    $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
    $<INSTALL_INTERFACE:include>
)
target_compile_features(twofinance_sdk_client INTERFACE cxx_std_17)

include(GNUInstallDirs)

install(
  DIRECTORY include/
  DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}
)

install(
  TARGETS twofinance_sdk_client
  EXPORT twofinance_sdk_client_targets
)

install(
  EXPORT twofinance_sdk_client_targets
  NAMESPACE twofinance::
  FILE twofinance_sdk_client_targets.cmake
  DESTINATION ${CMAKE_INSTALL_LIBDIR}/cmake/twofinance_sdk_client
)

include(CTest)
if(BUILD_TESTING)
  add_executable(twofinance_sdk_client_test tests/sdk_client_test.cpp)
  target_link_libraries(twofinance_sdk_client_test PRIVATE twofinance::sdk_client)
  add_test(NAME twofinance_sdk_client_test COMMAND twofinance_sdk_client_test)
  set_tests_properties(
    twofinance_sdk_client_test
    PROPERTIES WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  )

  add_executable(twofinance_sdk_client_quickstart examples/quickstart.cpp)
  target_link_libraries(twofinance_sdk_client_quickstart PRIVATE twofinance::sdk_client)

  add_executable(twofinance_sdk_client_request_options examples/request_options.cpp)
  target_link_libraries(twofinance_sdk_client_request_options PRIVATE twofinance::sdk_client)

  add_executable(twofinance_sdk_client_auth_token_source examples/auth_token_source.cpp)
  target_link_libraries(twofinance_sdk_client_auth_token_source PRIVATE twofinance::sdk_client)

  add_executable(twofinance_sdk_client_error_handling examples/error_handling.cpp)
  target_link_libraries(twofinance_sdk_client_error_handling PRIVATE twofinance::sdk_client)
endif()
