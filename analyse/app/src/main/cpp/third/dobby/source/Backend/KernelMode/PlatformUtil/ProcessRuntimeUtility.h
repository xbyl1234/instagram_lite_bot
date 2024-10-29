#pragma once

#include "PlatformUnifiedInterface/MemoryAllocator.h"

#include "UnifiedInterface/platform.h"

typedef struct _RuntimeModule {
  char path[1024];
  void *load_address;
} RuntimeModule;

struct MemRegion : MemRange {
  MemoryPermission permission;
  MemRegion(addr_t addr, size_t size, MemoryPermission perm): MemRange(addr, size), permission(perm) {

  }
};

class ProcessRuntimeUtility {
public:
  static const std::vector<MemRegion> &GetProcessMemoryLayout();

  static const std::vector<RuntimeModule> *GetProcessModuleMap();

  static RuntimeModule GetProcessModule(const char *name);
};