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

import com.fucker.gson.internal.$Gson$Types;

import java.lang.reflect.Type;

/**
 * Provides ability to apply a visitor to an object and all of its fields
 * recursively.
 *
 * @author Inderjeet Singh
 * @author Joel Leitch
 */
final class ObjectNavigator {

  public interface Visitor {
    public void start(ObjectTypePair node);

    public void end(ObjectTypePair node);

    /**
     * This is called before the object navigator starts visiting the current
     * object
     */
    void startVisitingObject(Object node);

    /**
     * This is called to visit the current object if it is an array
     */
    void visitArray(Object array, Type componentType);

    /**
     * This is called to visit an object field of the current object
     */
    void visitObjectField(FieldAttributes f, Type typeOfF, Object obj);

    /**
     * This is called to visit an array field of the current object
     */
    void visitArrayField(FieldAttributes f, Type typeOfF, Object obj);

    /**
     * This is called to visit an object using a custom handler
     *
     * @return true if a custom handler exists, false otherwise
     */
    public boolean visitUsingCustomHandler(ObjectTypePair objTypePair);

    /**
     * This is called to visit a field of the current object using a custom
     * handler
     */
    public boolean visitFieldUsingCustomHandler(FieldAttributes f, Type actualTypeOfField,
        Object parent);

    void visitPrimitive(Object primitive);

    /**
     * Retrieve the current target
     */
    Object getTarget();
  }

  private final ExclusionStrategy exclusionStrategy;
  private final ReflectingFieldNavigator reflectingFieldNavigator;

  /**
   * @param strategy the concrete exclusion strategy object to be used to filter out fields of an
   *          object.
   */
  ObjectNavigator(ExclusionStrategy strategy) {
    this.exclusionStrategy = strategy == null ? new NullExclusionStrategy() : strategy;
    this.reflectingFieldNavigator = new ReflectingFieldNavigator(exclusionStrategy);
  }

  /**
   * Navigate all the fields of the specified object. If a field is null, it
   * does not get visited.
   * @param objTypePair The object,type (fully genericized) being navigated
   */
  public void accept(ObjectTypePair objTypePair, Visitor visitor) {
    if (exclusionStrategy.shouldSkipClass($Gson$Types.getRawType(objTypePair.type))) {
      return;
    }
    boolean visitedWithCustomHandler = visitor.visitUsingCustomHandler(objTypePair);
    if (!visitedWithCustomHandler) {
      Object obj = objTypePair.getObject();
      Object objectToVisit = (obj == null) ? visitor.getTarget() : obj;
      if (objectToVisit == null) {
        return;
      }
      objTypePair.setObject(objectToVisit);
      visitor.start(objTypePair);
      try {
        if ($Gson$Types.isArray(objTypePair.type)) {
          visitor.visitArray(objectToVisit, objTypePair.type);
        } else if (objTypePair.type == Object.class && isPrimitiveOrString(objectToVisit)) {
          // TODO(Joel): this is only used for deserialization of "primitives"
          // we should rethink this!!!
          visitor.visitPrimitive(objectToVisit);
          visitor.getTarget();
        } else {
          visitor.startVisitingObject(objectToVisit);
          reflectingFieldNavigator.visitFieldsReflectively(objTypePair, visitor);
        }
      } finally {
        visitor.end(objTypePair);
      }
    }
  }

  private static boolean isPrimitiveOrString(Object objectToVisit) {
    Class<?> realClazz = objectToVisit.getClass();
    return realClazz == Object.class || realClazz == String.class
        || Primitives.unwrap(realClazz).isPrimitive();
  }
}
