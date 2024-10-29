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

/**
 * The new mechanism for providing custom field naming in Gson.  This allows the client code
 * to translate field names into a particular convention that is not supported as a normal
 * Java field declaration rules.  For example, Java does not support "-" characters in a
 * field name.
 *
 * @author Inderjeet Singh
 * @author Joel Leitch
 */
interface FieldNamingStrategy2 {

  /**
   * Translates the field name into its JSON field name representation.
   *
   * @param f the field that is being translated
   * @return the translated field name.
   */
  public String translateName(FieldAttributes f);
}
