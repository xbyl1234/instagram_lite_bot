#include <jni.h>
#include "third/utils/log.h"
#include "third/utils/jni_helper.hpp"
#include "third/utils/linux_helper.h"
#include "third/byopen/hack_dlopen.h"
#include "jni_trace/parse_java_sig.h"
#include "unwindstack/LocalUnwinder.h"
#include "dump_so.h"
#include "global/global.h"

void InitLibcHook();

void TestWhenHook();

bool dump_so_when_decode();

bool hook_sp_when_decode();

void watch_value();

void magisk_check();

bool check_memory_readable(void *addr);

void *get_start(const string &name);

bool hook_log_when_init();

void hook_pangle_log();

bool dump_pangle_so();

void testJniTrace();

void hook_aqy();

void hook_unity();

void HookOpengl();

void hook_pangle_pkg();

//    const-string v0, "analyse"
//    invoke-static {v0}, Ljava/lang/System;->loadLibrary(Ljava/lang/String;)V

void hook_pangle_log_pkg();

void hook_pangle_decode();

JNIEXPORT jint JNI_OnLoad(JavaVM *vm, void *reserved) {
//    hook_pangle_decode();
//    hook_pangle_log_pkg();
//    hook_pangle_pkg();
//    hook_unity();
//    HookOpengl();
//    hook_unity();
//    dump_so_delay("libmsaoaidsec.so", 5);
//    dump_so_delay("libmsaoaidauth.so", 5);
//    InitLibcHook();
//    hook_aqy();
//    dump_so_delay("libnms.so", 60);
//    dump_pangle_so();
//    hook_pangle_log();
//    dump_so_com2();
//    TestWhenHook();
//    dump_so_when_decode();
//    hook_sp_when_decode();
//    watch_value();
//    magisk_check();
//    LOGI("%d, %d, %d %d",
//    check_memory_readable(nullptr),
//    check_memory_readable((void *) 0x1515616500),
//    check_memory_readable((void *) check_memory_readable),
//    check_memory_readable(new char[1]));
//    hook_log_when_init();
//    LOGI("analyse libc: %p", get_start("libc.so"));
    return JNI_VERSION_1_6;
}

#include "third/dobby/include/dobby.h"

extern "C"
JNIEXPORT void JNICALL
Java_com_android_analyse_hook_Native_initNative(JNIEnv *env, jclass clazz, jstring pkg_name) {
    setPkgName(lsplant::JUTFString(env, pkg_name).get());
    LOGI("analyse inject pid: %d, %s", getpid(), getPkgName().c_str());
}