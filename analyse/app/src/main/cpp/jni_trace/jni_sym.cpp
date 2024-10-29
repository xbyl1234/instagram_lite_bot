#include <jni.h>

#include "format.h"
#include "jni_sym.h"
#include "art_method_name.h"

using namespace format;

ExternHookStub(GetStringUTFLength, jsize, JNIEnv *env, jstring jstr);

ExternHookStub(GetStringUTFChars, const char*, JNIEnv *env, jstring java_string, jboolean *is_copy);

ExternHookStub(NewStringUTF, jstring, JNIEnv *env, const char *utf);

ExternHookStub(RegisterNatives, jint, JNIEnv *env, jclass java_class,
               const JNINativeMethod *methods, jint method_count);

ExternHookStub(InvokeVirtualOrInterfaceWithVarArgs, jvalue, ScopedObjectAccessAlreadyRunnable soa,
               jobject obj, jmethodID mid, va_list args);

ExternHookStub(InvokeWithVarArgs, jvalue, ScopedObjectAccessAlreadyRunnable soa, jobject obj,
               jmethodID mid, va_list args);

ExternHookStub(SetObjectArrayElement, void, JNIEnv *env, jobjectArray array, jsize index,
               jobject value);

ExternHookStub(GetObjectArrayElement, jobject, JNIEnv *env, jobjectArray array, jsize index);

#define ExternGetFieldHook(type) ExternHookStub(Get##type##Field, jobject, JNIEnv *env, jobject obj, jfieldID field)

#define ExternSetFieldHook(type) ExternHookStub(Set##type##Field, void, JNIEnv *env, jobject obj, jfieldID field, jvalue v);

ExternGetFieldHook(Object);

ExternGetFieldHook(Boolean);

ExternGetFieldHook(Byte);

ExternGetFieldHook(Char);

ExternGetFieldHook(Short);

ExternGetFieldHook(Int);

ExternGetFieldHook(Long);

ExternGetFieldHook(Float);

ExternGetFieldHook(Double);

ExternSetFieldHook(Object);

ExternSetFieldHook(Boolean);

ExternSetFieldHook(Byte);

ExternSetFieldHook(Char);

ExternSetFieldHook(Short);

ExternSetFieldHook(Int);

ExternSetFieldHook(Long);

ExternSetFieldHook(Float);

ExternSetFieldHook(Double);

#define AddSymbolInfo(symName)  SymbolInfo{.isReg=false, .sym=  #symName, .stub=(void*)Hook_##symName, .org=(void**) &pHook_##symName,.target=(void*)env->functions->symName}
#define AddSymbolInfoBySym(symName)  SymbolInfo{.isReg=true, .sym=  #symName, .stub=(void*)Hook_##symName, .org=(void**) &pHook_##symName}

bool jni_sym::init(fake_dlctx_ref_t handleLibArt, JNIEnv *env) {
    jniHooks = {
            AddSymbolInfoBySym(InvokeVirtualOrInterfaceWithVarArgs),
            AddSymbolInfoBySym(InvokeWithVarArgs),
            AddSymbolInfo(NewStringUTF),
            AddSymbolInfo(GetStringUTFChars),
            AddSymbolInfo(GetStringUTFLength),
            AddSymbolInfo(RegisterNatives),
            AddSymbolInfo(GetObjectField),
            AddSymbolInfo(GetBooleanField),
            AddSymbolInfo(GetByteField),
            AddSymbolInfo(GetCharField),
            AddSymbolInfo(GetShortField),
            AddSymbolInfo(GetIntField),
            AddSymbolInfo(GetLongField),
            AddSymbolInfo(GetFloatField),
            AddSymbolInfo(GetDoubleField),
            AddSymbolInfo(SetObjectField),
            AddSymbolInfo(SetBooleanField),
            AddSymbolInfo(SetByteField),
            AddSymbolInfo(SetCharField),
            AddSymbolInfo(SetShortField),
            AddSymbolInfo(SetIntField),
            AddSymbolInfo(SetLongField),
            AddSymbolInfo(SetFloatField),
            AddSymbolInfo(SetDoubleField),
            AddSymbolInfo(GetObjectArrayElement),
            AddSymbolInfo(SetObjectArrayElement),
    };

    auto names = getSynName();
    vector<SymbolInfo *> needResolve;
    for (int i = 0; i < jniHooks.size(); ++i) {
        if (!jniHooks[i].sym.empty()) {
            auto find = names.find(jniHooks[i].sym);
            if (find != names.end()) {
                jniHooks[i].sym = find->second;
                needResolve.push_back(&jniHooks[i]);
            }
        }
    }
    if (!resolve(handleLibArt, &needResolve)) {
        loge("SymbolInfo::resolve error!");
        return false;
    }

    for (int i = 0; i < jniHooks.size(); ++i) {
        logi("%s: %p", jniHooks[i].sym.c_str(), jniHooks[i].target);
    }
    return true;
}
