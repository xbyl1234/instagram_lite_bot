#include "dobby_internal.h"

#include "logging/logging.h"

#include "Interceptor.h"
#include "InterceptRouting/InterceptRouting.h"

#include "function-wrapper.h"

PUBLIC int DobbyWrap(void *function_address, PreCallTy pre_call, PostCallTy post_call) {
  DLOG(0, "Initialize 'DobbyWrap' hook at %p", function_address);

  Interceptor *interceptor = Interceptor::SharedInstance();

  InterceptEntry *entry = new InterceptEntry();
  entry->id = interceptor->entries->getCount();
  entry->type = kFunctionWrapper;
  entry->function_address = function_address;

  FunctionWrapperRouting *routing = new FunctionWrapperRouting(entry);
  routing->DispatchRouting();
  interceptor->addHookEntry(entry);
  routing->Commit();

  DLOG(0, "Finalize %p", function_address);
  return RS_SUCCESS;
}
