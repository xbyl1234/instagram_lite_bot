/*
 * Copyright (C) 2011 Google Inc.
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

import com.fucker.gson.annotations.Expose;

/**
 * Excludes fields that do not have the {@link Expose} annotation
 *
 * @author Joel Leitch
 */
final class ExposeAnnotationDeserializationExclusionStrategy implements ExclusionStrategy {
  public boolean shouldSkipClass(Class<?> clazz) {
    return false;
  }

  public boolean shouldSkipField(FieldAttributes f) {
    Expose annotation = f.getAnnotation(Expose.class);
    if (annotation == null) {
      return true;
    }
    return !annotation.deserialize();
  }
}
