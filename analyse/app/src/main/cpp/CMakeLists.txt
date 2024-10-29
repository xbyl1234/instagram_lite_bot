cmake_minimum_required(VERSION 3.22.1)
project("analyse")
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY "${CMAKE_CURRENT_SOURCE_DIR}/../../../libs/${ANDROID_ABI}")

find_library(log-lib log)

add_subdirectory(./base)
add_subdirectory(./global)
add_subdirectory(./log2file)
add_subdirectory(./third)
include_directories(./third/libunwind/include)

file(GLOB analyse_src
        "*.h"
        "*.hpp"
        "*.cpp"

        "jni_trace/*.cpp"
        "jni_trace/*.h"
        )


add_library(
        analyse
        SHARED
        ${analyse_src}
)

target_link_libraries(
        analyse
        dobby
        libbase
        libutils
        libbyopen
        libjson
        libbase64
        unwindstack
        global
        log2file
        ${log-lib}
)

add_subdirectory(sniffer)
add_subdirectory(sqllite)