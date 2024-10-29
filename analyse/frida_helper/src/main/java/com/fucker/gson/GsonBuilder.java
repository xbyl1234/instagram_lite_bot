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

import com.fucker.gson.annotations.Expose;
import com.fucker.gson.DefaultTypeAdapters.DefaultDateTypeAdapter;
import com.fucker.gson.internal.$Gson$Preconditions;

import java.lang.reflect.Type;
import java.sql.Timestamp;
import java.text.DateFormat;
import java.util.Arrays;
import java.util.Date;
import java.util.HashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * <p>Use this builder to construct a {@link Gson} instance when you need to set configuration
 * options other than the default. For {@link Gson} with default configuration, it is simpler to
 * use {@code new Gson()}. {@code GsonBuilder} is best used by creating it, and then invoking its
 * various configuration methods, and finally calling create.</p>
 *
 * <p>The following is an example shows how to use the {@code GsonBuilder} to construct a Gson
 * instance:
 *
 * <pre>
 * Gson gson = new GsonBuilder()
 *     .registerTypeAdapter(Id.class, new IdTypeAdapter())
 *     .enableComplexMapKeySerialization()
 *     .serializeNulls()
 *     .setDateFormat(DateFormat.LONG)
 *     .setFieldNamingPolicy(FieldNamingPolicy.UPPER_CAMEL_CASE)
 *     .setPrettyPrinting()
 *     .setVersion(1.0)
 *     .create();
 * </pre></p>
 *
 * <p>NOTES:
 * <ul>
 * <li> the order of invocation of configuration methods does not matter.</li>
 * <li> The default serialization of {@link Date} and its subclasses in Gson does
 *  not contain time-zone information. So, if you are using date/time instances,
 *  use {@code GsonBuilder} and its {@code setDateFormat} methods.</li>
 *  </ul>
 * </p>
 *
 * @author Inderjeet Singh
 * @author Joel Leitch
 */
public final class GsonBuilder {
  private static final MapAsArrayTypeAdapter COMPLEX_KEY_MAP_TYPE_ADAPTER =
      new MapAsArrayTypeAdapter();
  private static final InnerClassExclusionStrategy innerClassExclusionStrategy =
      new InnerClassExclusionStrategy();
  private static final ExposeAnnotationDeserializationExclusionStrategy
  exposeAnnotationDeserializationExclusionStrategy =
      new ExposeAnnotationDeserializationExclusionStrategy();
  private static final ExposeAnnotationSerializationExclusionStrategy
  exposeAnnotationSerializationExclusionStrategy =
      new ExposeAnnotationSerializationExclusionStrategy();

  private final Set<ExclusionStrategy> serializeExclusionStrategies =
      new HashSet<ExclusionStrategy>();
  private final Set<ExclusionStrategy> deserializeExclusionStrategies =
      new HashSet<ExclusionStrategy>();

  private double ignoreVersionsAfter;
  private ModifierBasedExclusionStrategy modifierBasedExclusionStrategy;
  private boolean serializeInnerClasses;
  private boolean excludeFieldsWithoutExposeAnnotation;
  private LongSerializationPolicy longSerializationPolicy;
  private FieldNamingStrategy2 fieldNamingPolicy;
  private final ParameterizedTypeHandlerMap<InstanceCreator<?>> instanceCreators;
  private final ParameterizedTypeHandlerMap<JsonSerializer<?>> serializers;
  private final ParameterizedTypeHandlerMap<JsonDeserializer<?>> deserializers;
  private boolean serializeNulls;
  private String datePattern;
  private int dateStyle;
  private int timeStyle;
  private boolean serializeSpecialFloatingPointValues;
  private boolean escapeHtmlChars;
  private boolean prettyPrinting;
  private boolean generateNonExecutableJson;

  /**
   * Creates a GsonBuilder instance that can be used to build Gson with various configuration
   * settings. GsonBuilder follows the builder pattern, and it is typically used by first
   * invoking various configuration methods to set desired options, and finally calling
   * {@link #create()}.
   */
  public GsonBuilder() {
    // add default exclusion strategies
    deserializeExclusionStrategies.add(Gson.DEFAULT_ANON_LOCAL_CLASS_EXCLUSION_STRATEGY);
    deserializeExclusionStrategies.add(Gson.DEFAULT_SYNTHETIC_FIELD_EXCLUSION_STRATEGY);
    serializeExclusionStrategies.add(Gson.DEFAULT_ANON_LOCAL_CLASS_EXCLUSION_STRATEGY);
    serializeExclusionStrategies.add(Gson.DEFAULT_SYNTHETIC_FIELD_EXCLUSION_STRATEGY);

    // setup default values
    ignoreVersionsAfter = VersionConstants.IGNORE_VERSIONS;
    serializeInnerClasses = true;
    prettyPrinting = false;
    escapeHtmlChars = true;
    modifierBasedExclusionStrategy = Gson.DEFAULT_MODIFIER_BASED_EXCLUSION_STRATEGY;
    excludeFieldsWithoutExposeAnnotation = false;
    longSerializationPolicy = LongSerializationPolicy.DEFAULT;
    fieldNamingPolicy = Gson.DEFAULT_NAMING_POLICY;
    instanceCreators = new ParameterizedTypeHandlerMap<InstanceCreator<?>>();
    serializers = new ParameterizedTypeHandlerMap<JsonSerializer<?>>();
    deserializers = new ParameterizedTypeHandlerMap<JsonDeserializer<?>>();
    serializeNulls = false;
    dateStyle = DateFormat.DEFAULT;
    timeStyle = DateFormat.DEFAULT;
    serializeSpecialFloatingPointValues = false;
    generateNonExecutableJson = false;
  }

  /**
   * Configures Gson to enable versioning support.
   *
   * @param ignoreVersionsAfter any field or type marked with a version higher than this value
   * are ignored during serialization or deserialization.
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   */
  public GsonBuilder setVersion(double ignoreVersionsAfter) {
    this.ignoreVersionsAfter = ignoreVersionsAfter;
    return this;
  }

  /**
   * Configures Gson to excludes all class fields that have the specified modifiers. By default,
   * Gson will exclude all fields marked transient or static. This method will override that
   * behavior.
   *
   * @param modifiers the field modifiers. You must use the modifiers specified in the
   * {@link java.lang.reflect.Modifier} class. For example,
   * {@link java.lang.reflect.Modifier#TRANSIENT},
   * {@link java.lang.reflect.Modifier#STATIC}.
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   */
  public GsonBuilder excludeFieldsWithModifiers(int... modifiers) {
    modifierBasedExclusionStrategy = new ModifierBasedExclusionStrategy(modifiers);
    return this;
  }

  /**
   * Makes the output JSON non-executable in Javascript by prefixing the generated JSON with some
   * special text. This prevents attacks from third-party sites through script sourcing. See
   * <a href="http://code.google.com/p/google-gson/issues/detail?id=42">Gson Issue 42</a>
   * for details.
   *
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.3
   */
  public GsonBuilder generateNonExecutableJson() {
    this.generateNonExecutableJson = true;
    return this;
  }

  /**
   * Configures Gson to exclude all fields from consideration for serialization or deserialization
   * that do not have the {@link Expose} annotation.
   *
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   */
  public GsonBuilder excludeFieldsWithoutExposeAnnotation() {
    excludeFieldsWithoutExposeAnnotation = true;
    return this;
  }

  /**
   * Configure Gson to serialize null fields. By default, Gson omits all fields that are null
   * during serialization.
   *
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.2
   */
  public GsonBuilder serializeNulls() {
    this.serializeNulls = true;
    return this;
  }

  /**
   * Enabling this feature will only change the serialized form if the map key is
   * a complex type (i.e. non-primitive) in its <strong>serialized</strong> JSON
   * form. The default implementation of map serialization uses {@code toString()}
   * on the key; however, when this is called then one of the following cases
   * apply:
   *
   * <h3>Maps as JSON objects</h3>
   * For this case, assume that a type adapter is registered to serialize and
   * deserialize some {@code Point} class, which contains an x and y coordinate,
   * to/from the JSON Primitive string value {@code "(x,y)"}. The Java map would
   * then be serialized as a {@link JsonObject}.
   *
   * <p>Below is an example:
   * <pre>  {@code
   *   Gson gson = new GsonBuilder()
   *       .register(Point.class, new MyPointTypeAdapter())
   *       .enableComplexMapKeySerialization()
   *       .create();
   *
   *   Map<Point, String> original = new LinkedHashMap<Point, String>();
   *   original.put(new Point(5, 6), "a");
   *   original.put(new Point(8, 8), "b");
   *   System.out.println(gson.toJson(original, type));
   * }</pre>
   * The above code prints this JSON object:<pre>  {@code
   *   {
   *     "(5,6)": "a",
   *     "(8,8)": "b"
   *   }
   * }</pre>
   *
   * <h3>Maps as JSON arrays</h3>
   * For this case, assume that a type adapter was NOT registered for some
   * {@code Point} class, but rather the default Gson serialization is applied.
   * In this case, some {@code new Point(2,3)} would serialize as {@code
   * {"x":2,"y":5}}.
   *
   * <p>Given the assumption above, a {@code Map<Point, String>} will be
   * serialize as an array of arrays (can be viewed as an entry set of pairs).
   *
   * <p>Below is an example of serializing complex types as JSON arrays:
   * <pre> {@code
   *   Gson gson = new GsonBuilder()
   *       .enableComplexMapKeySerialization()
   *       .create();
   *
   *   Map<Point, String> original = new LinkedHashMap<Point, String>();
   *   original.put(new Point(5, 6), "a");
   *   original.put(new Point(8, 8), "b");
   *   System.out.println(gson.toJson(original, type));
   * }
   *
   * The JSON output would look as follows:
   * <pre>   {@code
   *   [
   *     [
   *       {
   *         "x": 5,
   *         "y": 6
   *       },
   *       "a"
   *     ],
   *     [
   *       {
   *         "x": 8,
   *         "y": 8
   *       },
   *       "b"
   *     ]
   *   ]
   * }</pre>
   *
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.7
   */
  public GsonBuilder enableComplexMapKeySerialization() {
    registerTypeHierarchyAdapter(Map.class, COMPLEX_KEY_MAP_TYPE_ADAPTER);
    return this;
  }

  /**
   * Configures Gson to exclude inner classes during serialization.
   *
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.3
   */
  public GsonBuilder disableInnerClassSerialization() {
    serializeInnerClasses = false;
    return this;
  }

  /**
   * Configures Gson to apply a specific serialization policy for {@code Long} and {@code long}
   * objects.
   *
   * @param serializationPolicy the particular policy to use for serializing longs.
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.3
   */
  public GsonBuilder setLongSerializationPolicy(LongSerializationPolicy serializationPolicy) {
    this.longSerializationPolicy = serializationPolicy;
    return this;
  }

  /**
   * Configures Gson to apply a specific naming policy to an object's field during serialization
   * and deserialization.
   *
   * @param namingConvention the JSON field naming convention to use for serialization and
   * deserialization.
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   */
  public GsonBuilder setFieldNamingPolicy(FieldNamingPolicy namingConvention) {
    return setFieldNamingStrategy(namingConvention.getFieldNamingPolicy());
  }

  /**
   * Configures Gson to apply a specific naming policy strategy to an object's field during
   * serialization and deserialization.
   *
   * @param fieldNamingStrategy the actual naming strategy to apply to the fields
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.3
   */
  public GsonBuilder setFieldNamingStrategy(FieldNamingStrategy fieldNamingStrategy) {
    return setFieldNamingStrategy(new FieldNamingStrategy2Adapter(fieldNamingStrategy));
  }

  /**
   * Configures Gson to apply a specific naming policy strategy to an object's field during
   * serialization and deserialization.
   *
   * @param fieldNamingStrategy the actual naming strategy to apply to the fields
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   */
  GsonBuilder setFieldNamingStrategy(FieldNamingStrategy2 fieldNamingStrategy) {
    this.fieldNamingPolicy =
        new SerializedNameAnnotationInterceptingNamingPolicy(fieldNamingStrategy);
    return this;
  }

  /**
   * Configures Gson to apply a set of exclusion strategies during both serialization and
   * deserialization. Each of the {@code strategies} will be applied as a disjunction rule.
   * This means that if one of the {@code strategies} suggests that a field (or class) should be
   * skipped then that field (or object) is skipped during serializaiton/deserialization.
   *
   * @param strategies the set of strategy object to apply during object (de)serialization.
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.4
   */
  public GsonBuilder setExclusionStrategies(ExclusionStrategy... strategies) {
    List<ExclusionStrategy> strategyList = Arrays.asList(strategies);
    serializeExclusionStrategies.addAll(strategyList);
    deserializeExclusionStrategies.addAll(strategyList);
    return this;
  }

  /**
   * Configures Gson to apply the passed in exclusion strategy during serialization.
   * If this method is invoked numerous times with different exclusion strategy objects
   * then the exclusion strategies that were added will be applied as a disjunction rule.
   * This means that if one of the added exclusion strategies suggests that a field (or
   * class) should be skipped then that field (or object) is skipped during its
   * serialization.
   *
   * @param strategy an exclusion strategy to apply during serialization.
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.7
   */
  public GsonBuilder addSerializationExclusionStrategy(ExclusionStrategy strategy) {
    serializeExclusionStrategies.add(strategy);
    return this;
  }

  /**
   * Configures Gson to apply the passed in exclusion strategy during deserialization.
   * If this method is invoked numerous times with different exclusion strategy objects
   * then the exclusion strategies that were added will be applied as a disjunction rule.
   * This means that if one of the added exclusion strategies suggests that a field (or
   * class) should be skipped then that field (or object) is skipped during its
   * deserialization.
   *
   * @param strategy an exclusion strategy to apply during deserialization.
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.7
   */
  public GsonBuilder addDeserializationExclusionStrategy(ExclusionStrategy strategy) {
    deserializeExclusionStrategies.add(strategy);
    return this;
  }
  /**
   * Configures Gson to output Json that fits in a page for pretty printing. This option only
   * affects Json serialization.
   *
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   */
  public GsonBuilder setPrettyPrinting() {
    prettyPrinting = true;
    return this;
  }

  /**
   * By default, Gson escapes HTML characters such as &lt; &gt; etc. Use this option to configure
   * Gson to pass-through HTML characters as is.
   *
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.3
   */
  public GsonBuilder disableHtmlEscaping() {
    this.escapeHtmlChars = false;
    return this;
  }

  /**
   * Configures Gson to serialize {@code Date} objects according to the pattern provided. You can
   * call this method or {@link #setDateFormat(int)} multiple times, but only the last invocation
   * will be used to decide the serialization format.
   *
   * <p>The date format will be used to serialize and deserialize {@link java.util.Date}, {@link
   * java.sql.Timestamp} and {@link java.sql.Date}.
   *
   * <p>Note that this pattern must abide by the convention provided by {@code SimpleDateFormat}
   * class. See the documentation in {@link java.text.SimpleDateFormat} for more information on
   * valid date and time patterns.</p>
   *
   * @param pattern the pattern that dates will be serialized/deserialized to/from
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.2
   */
  public GsonBuilder setDateFormat(String pattern) {
    // TODO(Joel): Make this fail fast if it is an invalid date format
    this.datePattern = pattern;
    return this;
  }

  /**
   * Configures Gson to to serialize {@code Date} objects according to the style value provided.
   * You can call this method or {@link #setDateFormat(String)} multiple times, but only the last
   * invocation will be used to decide the serialization format.
   *
   * <p>Note that this style value should be one of the predefined constants in the
   * {@code DateFormat} class. See the documentation in {@link java.text.DateFormat} for more
   * information on the valid style constants.</p>
   *
   * @param style the predefined date style that date objects will be serialized/deserialized
   * to/from
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.2
   */
  public GsonBuilder setDateFormat(int style) {
    this.dateStyle = style;
    this.datePattern = null;
    return this;
  }

  /**
   * Configures Gson to to serialize {@code Date} objects according to the style value provided.
   * You can call this method or {@link #setDateFormat(String)} multiple times, but only the last
   * invocation will be used to decide the serialization format.
   *
   * <p>Note that this style value should be one of the predefined constants in the
   * {@code DateFormat} class. See the documentation in {@link java.text.DateFormat} for more
   * information on the valid style constants.</p>
   *
   * @param dateStyle the predefined date style that date objects will be serialized/deserialized
   * to/from
   * @param timeStyle the predefined style for the time portion of the date objects
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.2
   */
  public GsonBuilder setDateFormat(int dateStyle, int timeStyle) {
    this.dateStyle = dateStyle;
    this.timeStyle = timeStyle;
    this.datePattern = null;
    return this;
  }

  /**
   * Configures Gson for custom serialization or deserialization. This method combines the
   * registration of an {@link InstanceCreator}, {@link JsonSerializer}, and a
   * {@link JsonDeserializer}. It is best used when a single object {@code typeAdapter} implements
   * all the required interfaces for custom serialization with Gson. If an instance creator,
   * serializer or deserializer was previously registered for the specified {@code type}, it is
   * overwritten.
   *
   * @param type the type definition for the type adapter being registered
   * @param typeAdapter This object must implement at least one of the {@link InstanceCreator},
   * {@link JsonSerializer}, and a {@link JsonDeserializer} interfaces.
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   */
  public GsonBuilder registerTypeAdapter(Type type, Object typeAdapter) {
    $Gson$Preconditions.checkArgument(typeAdapter instanceof JsonSerializer<?>
            || typeAdapter instanceof JsonDeserializer<?> || typeAdapter instanceof InstanceCreator<?>);
    if (typeAdapter instanceof InstanceCreator<?>) {
      registerInstanceCreator(type, (InstanceCreator<?>) typeAdapter);
    }
    if (typeAdapter instanceof JsonSerializer<?>) {
      registerSerializer(type, (JsonSerializer<?>) typeAdapter);
    }
    if (typeAdapter instanceof JsonDeserializer<?>) {
      registerDeserializer(type, (JsonDeserializer<?>) typeAdapter);
    }
    return this;
  }

  /**
   * Configures Gson to use a custom {@link InstanceCreator} for the specified type. If an instance
   * creator was previously registered for the specified class, it is overwritten. Since this method
   * takes a type instead of a Class object, it can be used to register a specific handler for a
   * generic type corresponding to a raw type.
   *
   * @param <T> the type for which instance creator is being registered
   * @param typeOfT The Type definition for T
   * @param instanceCreator the instance creator for T
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   */
  private <T> GsonBuilder registerInstanceCreator(Type typeOfT,
      InstanceCreator<? extends T> instanceCreator) {
    instanceCreators.register(typeOfT, instanceCreator);
    return this;
  }

  /**
   * Configures Gson to use a custom JSON serializer for the specified type. You should use this
   * method if you want to register different serializers for different generic types corresponding
   * to a raw type.
   *
   * @param <T> the type for which the serializer is being registered
   * @param typeOfT The type definition for T
   * @param serializer the custom serializer
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   */
  private <T> GsonBuilder registerSerializer(Type typeOfT, final JsonSerializer<T> serializer) {
    serializers.register(typeOfT, serializer);
    return this;
  }

  /**
   * Configures Gson to use a custom JSON deserializer for the specified type. You should use this
   * method if you want to register different deserializers for different generic types
   * corresponding to a raw type.
   *
   * @param <T> the type for which the deserializer is being registered
   * @param typeOfT The type definition for T
   * @param deserializer the custom deserializer
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   */
  private <T> GsonBuilder registerDeserializer(Type typeOfT, JsonDeserializer<T> deserializer) {
    deserializers.register(typeOfT, new JsonDeserializerExceptionWrapper<T>(deserializer));
    return this;
  }

  /**
   * Configures Gson for custom serialization or deserialization for an inheritance type hierarchy.
   * This method combines the registration of an {@link InstanceCreator}, {@link JsonSerializer},
   * and a {@link JsonDeserializer}. It is best used when a single object {@code typeAdapter}
   * implements all the required interfaces for custom serialization with Gson.
   * If an instance creator, serializer or deserializer was previously registered for the specified
   * type hierarchy, it is overwritten. If an instance creator, serializer or deserializer is
   * registered for a specific type in the type hierarchy, it will be invoked instead of the one
   * registered for the type hierarchy.
   *
   * @param baseType the class definition for the type adapter being registered for the base class
   *        or interface
   * @param typeAdapter This object must implement at least one of the {@link InstanceCreator},
   * {@link JsonSerializer}, and a {@link JsonDeserializer} interfaces.
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.7
   */
  public GsonBuilder registerTypeHierarchyAdapter(Class<?> baseType, Object typeAdapter) {
    $Gson$Preconditions.checkArgument(typeAdapter instanceof JsonSerializer<?>
            || typeAdapter instanceof JsonDeserializer<?> || typeAdapter instanceof InstanceCreator<?>);
    if (typeAdapter instanceof InstanceCreator<?>) {
      registerInstanceCreatorForTypeHierarchy(baseType, (InstanceCreator<?>) typeAdapter);
    }
    if (typeAdapter instanceof JsonSerializer<?>) {
      registerSerializerForTypeHierarchy(baseType, (JsonSerializer<?>) typeAdapter);
    }
    if (typeAdapter instanceof JsonDeserializer<?>) {
      registerDeserializerForTypeHierarchy(baseType, (JsonDeserializer<?>) typeAdapter);
    }
    return this;
  }

  private <T> GsonBuilder registerInstanceCreatorForTypeHierarchy(Class<?> classOfT,
      InstanceCreator<? extends T> instanceCreator) {
    instanceCreators.registerForTypeHierarchy(classOfT, instanceCreator);
    return this;
  }

  private <T> GsonBuilder registerSerializerForTypeHierarchy(Class<?> classOfT,
      final JsonSerializer<T> serializer) {
    serializers.registerForTypeHierarchy(classOfT, serializer);
    return this;
  }

  private <T> GsonBuilder registerDeserializerForTypeHierarchy(Class<?> classOfT,
      JsonDeserializer<T> deserializer) {
    deserializers.registerForTypeHierarchy(classOfT,
        new JsonDeserializerExceptionWrapper<T>(deserializer));
    return this;
  }

  /**
   * Section 2.4 of <a href="http://www.ietf.org/rfc/rfc4627.txt">JSON specification</a> disallows
   * special double values (NaN, Infinity, -Infinity). However,
   * <a href="http://www.ecma-international.org/publications/files/ECMA-ST/Ecma-262.pdf">Javascript
   * specification</a> (see section 4.3.20, 4.3.22, 4.3.23) allows these values as valid Javascript
   * values. Moreover, most JavaScript engines will accept these special values in JSON without
   * problem. So, at a practical level, it makes sense to accept these values as valid JSON even
   * though JSON specification disallows them.
   *
   * <p>Gson always accepts these special values during deserialization. However, it outputs
   * strictly compliant JSON. Hence, if it encounters a float value {@link Float#NaN},
   * {@link Float#POSITIVE_INFINITY}, {@link Float#NEGATIVE_INFINITY}, or a double value
   * {@link Double#NaN}, {@link Double#POSITIVE_INFINITY}, {@link Double#NEGATIVE_INFINITY}, it
   * will throw an {@link IllegalArgumentException}. This method provides a way to override the
   * default behavior when you know that the JSON receiver will be able to handle these special
   * values.
   *
   * @return a reference to this {@code GsonBuilder} object to fulfill the "Builder" pattern
   * @since 1.3
   */
  public GsonBuilder serializeSpecialFloatingPointValues() {
    this.serializeSpecialFloatingPointValues = true;
    return this;
  }

  /**
   * Creates a {@link Gson} instance based on the current configuration. This method is free of
   * side-effects to this {@code GsonBuilder} instance and hence can be called multiple times.
   *
   * @return an instance of Gson configured with the options currently set in this builder
   */
  public Gson create() {
    List<ExclusionStrategy> deserializationStrategies =
        new LinkedList<ExclusionStrategy>(deserializeExclusionStrategies);
    List<ExclusionStrategy> serializationStrategies =
        new LinkedList<ExclusionStrategy>(serializeExclusionStrategies);
    deserializationStrategies.add(modifierBasedExclusionStrategy);
    serializationStrategies.add(modifierBasedExclusionStrategy);

    if (!serializeInnerClasses) {
      deserializationStrategies.add(innerClassExclusionStrategy);
      serializationStrategies.add(innerClassExclusionStrategy);
    }
    if (ignoreVersionsAfter != VersionConstants.IGNORE_VERSIONS) {
      VersionExclusionStrategy versionExclusionStrategy =
          new VersionExclusionStrategy(ignoreVersionsAfter);
      deserializationStrategies.add(versionExclusionStrategy);
      serializationStrategies.add(versionExclusionStrategy);
    }
    if (excludeFieldsWithoutExposeAnnotation) {
      deserializationStrategies.add(exposeAnnotationDeserializationExclusionStrategy);
      serializationStrategies.add(exposeAnnotationSerializationExclusionStrategy);
    }

    ParameterizedTypeHandlerMap<JsonSerializer<?>> customSerializers =
        DefaultTypeAdapters.DEFAULT_HIERARCHY_SERIALIZERS.copyOf();
    customSerializers.register(serializers.copyOf());
    ParameterizedTypeHandlerMap<JsonDeserializer<?>> customDeserializers =
        DefaultTypeAdapters.DEFAULT_HIERARCHY_DESERIALIZERS.copyOf();
    customDeserializers.register(deserializers.copyOf());
    addTypeAdaptersForDate(datePattern, dateStyle, timeStyle, customSerializers,
        customDeserializers);

    customSerializers.registerIfAbsent(DefaultTypeAdapters.getDefaultSerializers(
        serializeSpecialFloatingPointValues, longSerializationPolicy));

    customDeserializers.registerIfAbsent(DefaultTypeAdapters.getDefaultDeserializers());

    ParameterizedTypeHandlerMap<InstanceCreator<?>> customInstanceCreators =
        instanceCreators.copyOf();
    customInstanceCreators.registerIfAbsent(DefaultTypeAdapters.getDefaultInstanceCreators());

    customSerializers.makeUnmodifiable();
    customDeserializers.makeUnmodifiable();
    instanceCreators.makeUnmodifiable();

    MappedObjectConstructor objConstructor = new MappedObjectConstructor(customInstanceCreators);

    Gson gson = new Gson(new DisjunctionExclusionStrategy(deserializationStrategies),
        new DisjunctionExclusionStrategy(serializationStrategies),
        fieldNamingPolicy, objConstructor, serializeNulls,
        customSerializers, customDeserializers, generateNonExecutableJson, escapeHtmlChars,
        prettyPrinting);
    return gson;
  }

  private static void addTypeAdaptersForDate(String datePattern, int dateStyle, int timeStyle,
      ParameterizedTypeHandlerMap<JsonSerializer<?>> serializers,
      ParameterizedTypeHandlerMap<JsonDeserializer<?>> deserializers) {
    DefaultDateTypeAdapter dateTypeAdapter = null;
    if (datePattern != null && !"".equals(datePattern.trim())) {
      dateTypeAdapter = new DefaultDateTypeAdapter(datePattern);
    } else if (dateStyle != DateFormat.DEFAULT && timeStyle != DateFormat.DEFAULT) {
      dateTypeAdapter = new DefaultDateTypeAdapter(dateStyle, timeStyle);
    }

    if (dateTypeAdapter != null) {
      registerIfAbsent(Date.class, serializers, dateTypeAdapter);
      registerIfAbsent(Date.class, deserializers, dateTypeAdapter);
      registerIfAbsent(Timestamp.class, serializers, dateTypeAdapter);
      registerIfAbsent(Timestamp.class, deserializers, dateTypeAdapter);
      registerIfAbsent(java.sql.Date.class, serializers, dateTypeAdapter);
      registerIfAbsent(java.sql.Date.class, deserializers, dateTypeAdapter);
    }
  }

  private static <T> void registerIfAbsent(Class<?> type,
      ParameterizedTypeHandlerMap<T> adapters, T adapter) {
    if (!adapters.hasSpecificHandlerFor(type)) {
      adapters.register(type, adapter);
    }
  }
}
