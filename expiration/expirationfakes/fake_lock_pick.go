// Code generated by counterfeiter. DO NOT EDIT.
package expirationfakes

import (
	"sync"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/locket/db"
	"code.cloudfoundry.org/locket/expiration"
)

type FakeLockPick struct {
	ExpirationCountsStub        func() (uint32, uint32)
	expirationCountsMutex       sync.RWMutex
	expirationCountsArgsForCall []struct {
	}
	expirationCountsReturns struct {
		result1 uint32
		result2 uint32
	}
	expirationCountsReturnsOnCall map[int]struct {
		result1 uint32
		result2 uint32
	}
	RegisterTTLStub        func(lager.Logger, *db.Lock)
	registerTTLMutex       sync.RWMutex
	registerTTLArgsForCall []struct {
		arg1 lager.Logger
		arg2 *db.Lock
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLockPick) ExpirationCounts() (uint32, uint32) {
	fake.expirationCountsMutex.Lock()
	ret, specificReturn := fake.expirationCountsReturnsOnCall[len(fake.expirationCountsArgsForCall)]
	fake.expirationCountsArgsForCall = append(fake.expirationCountsArgsForCall, struct {
	}{})
	fake.recordInvocation("ExpirationCounts", []interface{}{})
	expirationCountsStubCopy := fake.ExpirationCountsStub
	fake.expirationCountsMutex.Unlock()
	if expirationCountsStubCopy != nil {
		return expirationCountsStubCopy()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.expirationCountsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeLockPick) ExpirationCountsCallCount() int {
	fake.expirationCountsMutex.RLock()
	defer fake.expirationCountsMutex.RUnlock()
	return len(fake.expirationCountsArgsForCall)
}

func (fake *FakeLockPick) ExpirationCountsCalls(stub func() (uint32, uint32)) {
	fake.expirationCountsMutex.Lock()
	defer fake.expirationCountsMutex.Unlock()
	fake.ExpirationCountsStub = stub
}

func (fake *FakeLockPick) ExpirationCountsReturns(result1 uint32, result2 uint32) {
	fake.expirationCountsMutex.Lock()
	defer fake.expirationCountsMutex.Unlock()
	fake.ExpirationCountsStub = nil
	fake.expirationCountsReturns = struct {
		result1 uint32
		result2 uint32
	}{result1, result2}
}

func (fake *FakeLockPick) ExpirationCountsReturnsOnCall(i int, result1 uint32, result2 uint32) {
	fake.expirationCountsMutex.Lock()
	defer fake.expirationCountsMutex.Unlock()
	fake.ExpirationCountsStub = nil
	if fake.expirationCountsReturnsOnCall == nil {
		fake.expirationCountsReturnsOnCall = make(map[int]struct {
			result1 uint32
			result2 uint32
		})
	}
	fake.expirationCountsReturnsOnCall[i] = struct {
		result1 uint32
		result2 uint32
	}{result1, result2}
}

func (fake *FakeLockPick) RegisterTTL(arg1 lager.Logger, arg2 *db.Lock) {
	fake.registerTTLMutex.Lock()
	fake.registerTTLArgsForCall = append(fake.registerTTLArgsForCall, struct {
		arg1 lager.Logger
		arg2 *db.Lock
	}{arg1, arg2})
	fake.recordInvocation("RegisterTTL", []interface{}{arg1, arg2})
	registerTTLStubCopy := fake.RegisterTTLStub
	fake.registerTTLMutex.Unlock()
	if registerTTLStubCopy != nil {
		registerTTLStubCopy(arg1, arg2)
	}
}

func (fake *FakeLockPick) RegisterTTLCallCount() int {
	fake.registerTTLMutex.RLock()
	defer fake.registerTTLMutex.RUnlock()
	return len(fake.registerTTLArgsForCall)
}

func (fake *FakeLockPick) RegisterTTLCalls(stub func(lager.Logger, *db.Lock)) {
	fake.registerTTLMutex.Lock()
	defer fake.registerTTLMutex.Unlock()
	fake.RegisterTTLStub = stub
}

func (fake *FakeLockPick) RegisterTTLArgsForCall(i int) (lager.Logger, *db.Lock) {
	fake.registerTTLMutex.RLock()
	defer fake.registerTTLMutex.RUnlock()
	argsForCall := fake.registerTTLArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeLockPick) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.expirationCountsMutex.RLock()
	defer fake.expirationCountsMutex.RUnlock()
	fake.registerTTLMutex.RLock()
	defer fake.registerTTLMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeLockPick) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ expiration.LockPick = new(FakeLockPick)
