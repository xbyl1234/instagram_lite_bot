/*
 * Copyright (C) 2008 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.fucker.gson;

import java.lang.reflect.Type;

/**
 * implementation of a deserialization context for Gson
 *
 * @author Inderjeet Singh
 */
final class JsonDeserializationContextDefault implements JsonDeserializationContext {

  private final ObjectNavigator objectNavigator;
  private final FieldNamingStrategy2 fieldNamingPolicy;
  private final ParameterizedTypeHandlerMap<JsonDeserializer<?>> deserializers;
  private final MappedObjectConstructor objectConstructor;

  JsonDeserializationContextDefault(ObjectNavigator objectNavigator,
      FieldNamingStrategy2 fieldNamingPolicy,
      ParameterizedTypeHandlerMap<JsonDeserializer<?>> deserializers,
      MappedObjectConstructor objectConstructor) {
    this.objectNavigator = objectNavigator;
    this.fieldNamingPolicy = fieldNamingPolicy;
    this.deserializers = deserializers;
    this.objectConstructor = objectConstructor;
  }

  ObjectConstructor getObjectConstructor() {
    return objectConstructor;
  }

  @SuppressWarnings("unchecked")
  public <T> T deserialize(JsonElement json, Type typeOfT) throws JsonParseException {
    if (json == null || json.isJsonNull()) {
      return null;
    } else if (json.isJsonArray()) {
      return (T) fromJsonArray(typeOfT, json.getAsJsonArray(), this);
    } else if (json.isJsonObject()) {
      return (T) fromJsonObject(typeOfT, json.getAsJsonObject(), this);
    } else if (json.isJsonPrimitive()) {
      return (T) fromJsonPrimitive(typeOfT, json.getAsJsonPrimitive(), this);
    } else {
      throw new JsonParseException("Failed parsing JSON source: " + json + " to Json");
    }
  }

  private <T> T fromJsonArray(Type arrayType, JsonArray jsonArray,
      JsonDeserializationContext context) throws JsonParseException {
    JsonArrayDeserializationVisitor<T> visitor = new JsonArrayDeserializationVisitor<T>(
        jsonArray, arrayType, objectNavigator, fieldNamingPolicy,
        objectConstructor, deserializers, context);
    objectNavigator.accept(new ObjectTypePair(null, arrayType, true), visitor);
    return visitor.getTarget();
  }

  private <T> T fromJsonObject(Type typeOfT, JsonObject jsonObject,
      JsonDeserializationContext context) throws JsonParseException {
    JsonObjectDeserializationVisitor<T> visitor = new JsonObjectDeserializationVisitor<T>(
        jsonObject, typeOfT, objectNavigator, fieldNamingPolicy,
        objectConstructor, deserializers, context);
    objectNavigator.accept(new ObjectTypePair(null, typeOfT, true), visitor);
    return visitor.getTarget();
  }

  @SuppressWarnings("unchecked")
  private <T> T fromJsonPrimitive(Type typeOfT, JsonPrimitive json,
      JsonDeserializationContext context) throws JsonParseException {
    JsonObjectDeserializationVisitor<T> visitor = new JsonObjectDeserializationVisitor<T>(
        json, typeOfT, objectNavigator, fieldNamingPolicy, objectConstructor, deserializers, context);
    objectNavigator.accept(new ObjectTypePair(json.getAsObject(), typeOfT, true), visitor);
    Object target = visitor.getTarget();
    return (T) target;
  }
}
