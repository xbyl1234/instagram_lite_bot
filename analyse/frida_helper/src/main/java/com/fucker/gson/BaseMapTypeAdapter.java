/*
 * Copyright (C) 2011 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.fucker.gson;

import java.lang.reflect.Type;
import java.util.Map;

/**
 * Captures all the common/shared logic between the old, ({@link MapTypeAdapter}, and
 * the new, {@link MapAsArrayTypeAdapter}, map type adapters.
 *
 * @author Joel Leitch
 */
abstract class BaseMapTypeAdapter
    implements JsonSerializer<Map<?, ?>>, JsonDeserializer<Map<?, ?>> {

  protected static final JsonElement serialize(JsonSerializationContext context,
      Object src, Type srcType) {
    JsonSerializationContextDefault contextImpl = (JsonSerializationContextDefault) context;
    return contextImpl.serialize(src, srcType, false);
  }

  @SuppressWarnings("unchecked")
  protected static final Map<Object, Object> constructMapType(
      Type mapType, JsonDeserializationContext context) {
    JsonDeserializationContextDefault contextImpl = (JsonDeserializationContextDefault) context;
    ObjectConstructor objectConstructor = contextImpl.getObjectConstructor();
    return (Map<Object, Object>) objectConstructor.construct(mapType);
  }
}