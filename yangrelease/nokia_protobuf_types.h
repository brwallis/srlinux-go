/**********************************************************************************************************************
 * Description:
 *   Helper functions that make it easier to use the nokia-extensions protobuf messages.
 *   For example, if your message has a field 'path' of type 'string', you set it like this:
 *      using nokia_extensions::CreateStringValue;
 *      message.set_allocated_path(CreateStringValue(application.path).release());
 *
 * Copyright (c) 2018 Nokia
 ***********************************************************************************************************************/

#ifndef SRLINUX_YANG_MODULES_NOKIA_PROTOBUF_TYPES_H_
#define SRLINUX_YANG_MODULES_NOKIA_PROTOBUF_TYPES_H_

// #include files here.
#include <memory>
#include <string>
#include <utility>
#include "nokia_extensions.pb.h"

namespace nokia_extensions
{
// Note: to see what arguments are accepted,
// take a look at the 'set_value' methods defined in nokia_extensions.pb.h

template <typename... T>
std::unique_ptr<BytesValue> CreateBytesValue(T... value);
template <typename... T>
std::unique_ptr<BoolValue> CreateBoolValue(T... value);
template <typename... T>
std::unique_ptr<Decimal64Value> CreateDecimal64Value(T... value);
template <typename... T>
std::unique_ptr<IntValue> CreateIntValue(T... value);
template <typename... T>
std::unique_ptr<StringValue> CreateStringValue(T... value);
template <typename... T>
std::unique_ptr<UintValue> CreateUintValue(T... value);

}  // namespace nokia_extensions

//-----------------------------------------------------------------------------
// Implementation
//-----------------------------------------------------------------------------

namespace nokia_extensions
{
template <typename... T>
std::unique_ptr<BytesValue> CreateBytesValue(T... value)
{
    auto result = std::make_unique<BytesValue>();
    result->set_value(std::forward<T>(value)...);
    return result;
}

template <typename... T>
std::unique_ptr<BoolValue> CreateBoolValue(T... value)
{
    auto result = std::make_unique<BoolValue>();
    result->set_value(std::forward<T>(value)...);
    return result;
}

template <typename... T>
std::unique_ptr<IntValue> CreateIntValue(T... value)
{
    auto result = std::make_unique<IntValue>();
    result->set_value(std::forward<T>(value)...);
    return result;
}

template <typename... T>
std::unique_ptr<StringValue> CreateStringValue(T... value)
{
    auto result = std::make_unique<StringValue>();
    result->set_value(std::forward<T>(value)...);
    return result;
}

template <typename... T>
std::unique_ptr<UintValue> CreateUintValue(T... value)
{
    auto result = std::make_unique<UintValue>();
    result->set_value(std::forward<T>(value)...);
    return result;
}

}  // namespace nokia_extensions

#endif  // SRLINUX_YANG_MODULES_NOKIA_PROTOBUF_TYPES_H_
