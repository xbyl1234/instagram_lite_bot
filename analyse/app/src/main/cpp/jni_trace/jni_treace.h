#pragma once

#include <string>
#include <set>
#include <jni.h>
#include <dlfcn.h>

#include "../third/byopen/hack_dlopen.h"
#include "../third/utils/utils.h"
#include "../third/utils/log.h"
#include "../base/hook.h"

#include "jni_helper.h"
#include "jni_sym.h"

using namespace std;


class JniTrace {
public:
    explicit JniTrace() {}

    bool
    Init(jclass frida_helper, fake_dlctx_ref_t handleLibArt,
         const vector<string> &targetModule, const vector<string> &passJavaMethod) {
        JNIEnv *env = jniHelper.GetEnv();
        if (!this->sym.init(handleLibArt, env)) {
            return false;
        }
        this->targetModule = targetModule;
        for (const auto &item: passJavaMethod) {
            this->passJavaMethod.insert(item);
        }
        this->frida_helper = frida_helper;
        object_2_string = env->GetStaticMethodID(frida_helper, "object_2_string",
                                                 "(Ljava/lang/Object;)Ljava/lang/String;");
        init = true;
        return true;
    }

    bool Hook() {
        return hookAll(&sym.jniHooks);
    }

    bool CheckPassJavaMethod(jmethodID method) {
        return passJavaMethodId.contains(method);
    }

    bool CheckPassJavaMethod(jmethodID methodId, const string &method) {
        bool ret = passJavaMethod.contains(method);
        if (ret) {
            passJavaMethodId.insert(methodId);
        }
        return ret;
    }

    bool CheckTargetModule(const vector<Stack> &frame) {
        bool find = false;
        for (const auto &item: frame) {
            for (const auto &target: targetModule) {
                if (item.name.find(target) != -1) {
                    find = true;
                    break;
                }
            }
            if (find) {
                break;
            }
        }
        return find;
    }

public:
    bool init = false;
    jclass frida_helper;
    jmethodID object_2_string;
    vector<string> targetModule;
    set<string> passJavaMethod;
    set<jmethodID> passJavaMethodId;
    jni_sym sym;
};

extern JniTrace jniTrace;
extern __thread bool passJniTrace;
extern __thread bool passCallMethod;