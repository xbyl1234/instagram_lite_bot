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
 * A visitor that populates fields of an object with data from its equivalent
 * JSON representation
 *
 * @author Inderjeet Singh
 * @author Joel Leitch
 */
final class JsonObjectDeserializationVisitor<T> extends JsonDeserializationVisitor<T> {

  JsonObjectDeserializationVisitor(JsonElement json, Type type,
      ObjectNavigator objectNavigator, FieldNamingStrategy2 fieldNamingPolicy,
      ObjectConstructor objectConstructor,
      ParameterizedTypeHandlerMap<JsonDeserializer<?>> deserializers,
      JsonDeserializationContext context) {
    super(json, type, objectNavigator, fieldNamingPolicy, objectConstructor, deserializers, context);
  }

  @Override
  @SuppressWarnings("unchecked")
  protected T constructTarget() {
    return (T) objectConstructor.construct(targetType);
  }

  public void startVisitingObject(Object node) {
    // do nothing
  }

  public void visitArray(Object array, Type componentType) {
    // should not be called since this case should invoke JsonArrayDeserializationVisitor
    throw new JsonParseException("Expecting object but found array: " + array);
  }

  public void visitObjectField(FieldAttributes f, Type typeOfF, Object obj) {
    try {
      if (!json.isJsonObject()) {
        throw new JsonParseException("Expecting object found: " + json);
      }
      JsonObject jsonObject = json.getAsJsonObject();
      String fName = getFieldName(f);
      JsonElement jsonChild = jsonObject.get(fName);
      if (jsonChild != null) {
        Object child = visitChildAsObject(typeOfF, jsonChild);
        f.set(obj, child);
      } else {
        f.set(obj, null);
      }
    } catch (IllegalAccessException e) {
      throw new RuntimeException(e);
    }
  }

  public void visitArrayField(FieldAttributes f, Type typeOfF, Object obj) {
    try {
      if (!json.isJsonObject()) {
        throw new JsonParseException("Expecting object found: " + json);
      }
      JsonObject jsonObject = json.getAsJsonObject();
      String fName = getFieldName(f);
      JsonArray jsonChild = (JsonArray) jsonObject.get(fName);
      if (jsonChild != null) {
        Object array = visitChildAsArray(typeOfF, jsonChild);
        f.set(obj, array);
      } else {
        f.set(obj, null);
      }
    } catch (IllegalAccessException e) {
      throw new RuntimeException(e);
    }
  }

  private String getFieldName(FieldAttributes f) {
    return fieldNamingPolicy.translateName(f);
  }

  public boolean visitFieldUsingCustomHandler(FieldAttributes f, Type declaredTypeOfField, Object parent) {
    try {
      String fName = getFieldName(f);
      if (!json.isJsonObject()) {
        throw new JsonParseException("Expecting object found: " + json);
      }
      JsonElement child = json.getAsJsonObject().get(fName);
      boolean isPrimitive = Primitives.isPrimitive(declaredTypeOfField);
      if (child == null) { // Child will be null if the field wasn't present in Json
        return true;
      } else if (child.isJsonNull()) {
        if (!isPrimitive) {
          f.set(parent, null);
        }
        return true;
      }
      ObjectTypePair objTypePair = new ObjectTypePair(null, declaredTypeOfField, false);
      Pair<JsonDeserializer<?>, ObjectTypePair> pair = objTypePair.getMatchingHandler(deserializers);
      if (pair == null) {
        return false;
      }
      Object value = invokeCustomDeserializer(child, pair);
      if (value != null || !isPrimitive) {
        f.set(parent, value);
      }
      return true;
    } catch (IllegalAccessException e) {
      throw new RuntimeException();
    }
  }

  @SuppressWarnings("unchecked")
  public void visitPrimitive(Object primitive) {
    if (!json.isJsonPrimitive()) {
      throw new JsonParseException(
          "Type information is unavailable, and the target object is not a primitive: " + json);
    }
    JsonPrimitive prim = json.getAsJsonPrimitive();
    target = (T) prim.getAsObject();
  }
}
