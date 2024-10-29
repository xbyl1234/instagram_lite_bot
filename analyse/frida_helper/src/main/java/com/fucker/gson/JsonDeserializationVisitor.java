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

import com.fucker.gson.internal.$Gson$Preconditions;

import java.lang.reflect.Type;

/**
 * Abstract data value container for the {@link ObjectNavigator.Visitor}
 * implementations.  This class exposes the {@link #getTarget()} method
 * which returns the class that was visited by this object.
 *
 * @author Inderjeet Singh
 * @author Joel Leitch
 */
abstract class JsonDeserializationVisitor<T> implements ObjectNavigator.Visitor {

  protected final ObjectNavigator objectNavigator;
  protected final FieldNamingStrategy2 fieldNamingPolicy;
  protected final ObjectConstructor objectConstructor;
  protected final ParameterizedTypeHandlerMap<JsonDeserializer<?>> deserializers;
  protected T target;
  protected final JsonElement json;
  protected final Type targetType;
  protected final JsonDeserializationContext context;
  protected boolean constructed;

  JsonDeserializationVisitor(JsonElement json, Type targetType,
      ObjectNavigator objectNavigator, FieldNamingStrategy2 fieldNamingPolicy,
      ObjectConstructor objectConstructor,
      ParameterizedTypeHandlerMap<JsonDeserializer<?>> deserializers,
      JsonDeserializationContext context) {
    this.targetType = targetType;
    this.objectNavigator = objectNavigator;
    this.fieldNamingPolicy = fieldNamingPolicy;
    this.objectConstructor = objectConstructor;
    this.deserializers = deserializers;
    this.json = $Gson$Preconditions.checkNotNull(json);
    this.context = context;
    this.constructed = false;
  }

  public T getTarget() {
    if (!constructed) {
      target = constructTarget();
      constructed = true;
    }
    return target;
  }

  protected abstract T constructTarget();

  public void start(ObjectTypePair node) {
  }

  public void end(ObjectTypePair node) {
  }

  @SuppressWarnings("unchecked")
  public final boolean visitUsingCustomHandler(ObjectTypePair objTypePair) {
    Pair<JsonDeserializer<?>, ObjectTypePair> pair = objTypePair.getMatchingHandler(deserializers);
    if (pair == null) {
      return false;
    }
    Object value = invokeCustomDeserializer(json, pair);
    target = (T) value;
    constructed = true;
    return true;
  }

  protected Object invokeCustomDeserializer(JsonElement element,
      Pair<JsonDeserializer<?>, ObjectTypePair> pair) {
    if (element == null || element.isJsonNull()) {
      return null;
    }
    Type objType = pair.second.type;
    return (pair.first).deserialize(element, objType, context);
  }

  final Object visitChildAsObject(Type childType, JsonElement jsonChild) {
    JsonDeserializationVisitor<?> childVisitor =
        new JsonObjectDeserializationVisitor<Object>(jsonChild, childType,
            objectNavigator, fieldNamingPolicy, objectConstructor, deserializers, context);
    return visitChild(childType, childVisitor);
  }

  final Object visitChildAsArray(Type childType, JsonArray jsonChild) {
    JsonDeserializationVisitor<?> childVisitor =
        new JsonArrayDeserializationVisitor<Object>(jsonChild.getAsJsonArray(), childType,
            objectNavigator, fieldNamingPolicy, objectConstructor, deserializers, context);
    return visitChild(childType, childVisitor);
  }

  private Object visitChild(Type type, JsonDeserializationVisitor<?> childVisitor) {
    objectNavigator.accept(new ObjectTypePair(null, type, false), childVisitor);
    // the underlying object may have changed during the construction phase
    // This happens primarily because of custom deserializers
    return childVisitor.getTarget();
  }
}
