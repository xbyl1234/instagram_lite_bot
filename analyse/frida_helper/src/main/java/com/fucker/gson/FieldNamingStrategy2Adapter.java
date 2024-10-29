/*
 * Copyright (C) 2010 Google Inc.
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

/**
 * Adapts the old FieldNamingStrategy to the new {@link FieldNamingStrategy2}
 * type.
 *
 * @author Inderjeet Singh
 * @author Joel Leitch
 */
final class FieldNamingStrategy2Adapter implements FieldNamingStrategy2 {
  private final FieldNamingStrategy adaptee;

  FieldNamingStrategy2Adapter(FieldNamingStrategy adaptee) {
    this.adaptee = $Gson$Preconditions.checkNotNull(adaptee);
  }

  public String translateName(FieldAttributes f) {
    return adaptee.translateName(f.getFieldObject());
  }
}
