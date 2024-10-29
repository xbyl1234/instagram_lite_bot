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

import java.lang.reflect.Array;
import java.lang.reflect.Type;

/**
 * This class contains a mapping of all the application specific
 * {@link InstanceCreator} instances.  Registering an {@link InstanceCreator}
 * with this class will override the default object creation that is defined
 * by the ObjectConstructor that this class is wrapping.  Using this class
 * with the JSON framework provides the application with "pluggable" modules
 * to customize framework to suit the application's needs.
 *
 * @author Joel Leitch
 */
final class MappedObjectConstructor implements ObjectConstructor {
  private static final UnsafeAllocator unsafeAllocator = UnsafeAllocator.create();
  private static final DefaultConstructorAllocator defaultConstructorAllocator =
      new DefaultConstructorAllocator(500);

  private final ParameterizedTypeHandlerMap<InstanceCreator<?>> instanceCreatorMap;

  public MappedObjectConstructor(
      ParameterizedTypeHandlerMap<InstanceCreator<?>> instanceCreators) {
    instanceCreatorMap = instanceCreators;
  }

  @SuppressWarnings("unchecked")
  public <T> T construct(Type typeOfT) {
    InstanceCreator<T> creator = (InstanceCreator<T>) instanceCreatorMap.getHandlerFor(typeOfT);
    if (creator != null) {
      return creator.createInstance(typeOfT);
    }
    return (T) constructWithAllocators(typeOfT);
  }

  public Object constructArray(Type type, int length) {
    return Array.newInstance($Gson$Types.getRawType(type), length);
  }

  @SuppressWarnings({"unchecked", "cast"})
  private <T> T constructWithAllocators(Type typeOfT) {
    try {
      Class<T> clazz = (Class<T>) $Gson$Types.getRawType(typeOfT);
      T obj = defaultConstructorAllocator.newInstance(clazz);
      return (obj == null)
          ? unsafeAllocator.newInstance(clazz)
          : obj;
    } catch (Exception e) {
      throw new RuntimeException(("Unable to invoke no-args constructor for " + typeOfT + ". "
          + "Register an InstanceCreator with Gson for this type may fix this problem."), e);
    }
  }

  @Override
  public String toString() {
    return instanceCreatorMap.toString();
  }
}
