// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"code.cloudfoundry.org/consuladapter"
	"github.com/hashicorp/consul/api"
)

type FakeClient struct {
	AgentStub        func() consuladapter.Agent
	agentMutex       sync.RWMutex
	agentArgsForCall []struct{}
	agentReturns     struct {
		result1 consuladapter.Agent
	}
	SessionStub        func() consuladapter.Session
	sessionMutex       sync.RWMutex
	sessionArgsForCall []struct{}
	sessionReturns     struct {
		result1 consuladapter.Session
	}
	CatalogStub        func() consuladapter.Catalog
	catalogMutex       sync.RWMutex
	catalogArgsForCall []struct{}
	catalogReturns     struct {
		result1 consuladapter.Catalog
	}
	KVStub        func() consuladapter.KV
	kVMutex       sync.RWMutex
	kVArgsForCall []struct{}
	kVReturns     struct {
		result1 consuladapter.KV
	}
	StatusStub        func() consuladapter.Status
	statusMutex       sync.RWMutex
	statusArgsForCall []struct{}
	statusReturns     struct {
		result1 consuladapter.Status
	}
	LockOptsStub        func(opts *api.LockOptions) (consuladapter.Lock, error)
	lockOptsMutex       sync.RWMutex
	lockOptsArgsForCall []struct {
		opts *api.LockOptions
	}
	lockOptsReturns struct {
		result1 consuladapter.Lock
		result2 error
	}
}

func (fake *FakeClient) Agent() consuladapter.Agent {
	fake.agentMutex.Lock()
	fake.agentArgsForCall = append(fake.agentArgsForCall, struct{}{})
	fake.agentMutex.Unlock()
	if fake.AgentStub != nil {
		return fake.AgentStub()
	} else {
		return fake.agentReturns.result1
	}
}

func (fake *FakeClient) AgentCallCount() int {
	fake.agentMutex.RLock()
	defer fake.agentMutex.RUnlock()
	return len(fake.agentArgsForCall)
}

func (fake *FakeClient) AgentReturns(result1 consuladapter.Agent) {
	fake.AgentStub = nil
	fake.agentReturns = struct {
		result1 consuladapter.Agent
	}{result1}
}

func (fake *FakeClient) Session() consuladapter.Session {
	fake.sessionMutex.Lock()
	fake.sessionArgsForCall = append(fake.sessionArgsForCall, struct{}{})
	fake.sessionMutex.Unlock()
	if fake.SessionStub != nil {
		return fake.SessionStub()
	} else {
		return fake.sessionReturns.result1
	}
}

func (fake *FakeClient) SessionCallCount() int {
	fake.sessionMutex.RLock()
	defer fake.sessionMutex.RUnlock()
	return len(fake.sessionArgsForCall)
}

func (fake *FakeClient) SessionReturns(result1 consuladapter.Session) {
	fake.SessionStub = nil
	fake.sessionReturns = struct {
		result1 consuladapter.Session
	}{result1}
}

func (fake *FakeClient) Catalog() consuladapter.Catalog {
	fake.catalogMutex.Lock()
	fake.catalogArgsForCall = append(fake.catalogArgsForCall, struct{}{})
	fake.catalogMutex.Unlock()
	if fake.CatalogStub != nil {
		return fake.CatalogStub()
	} else {
		return fake.catalogReturns.result1
	}
}

func (fake *FakeClient) CatalogCallCount() int {
	fake.catalogMutex.RLock()
	defer fake.catalogMutex.RUnlock()
	return len(fake.catalogArgsForCall)
}

func (fake *FakeClient) CatalogReturns(result1 consuladapter.Catalog) {
	fake.CatalogStub = nil
	fake.catalogReturns = struct {
		result1 consuladapter.Catalog
	}{result1}
}

func (fake *FakeClient) KV() consuladapter.KV {
	fake.kVMutex.Lock()
	fake.kVArgsForCall = append(fake.kVArgsForCall, struct{}{})
	fake.kVMutex.Unlock()
	if fake.KVStub != nil {
		return fake.KVStub()
	} else {
		return fake.kVReturns.result1
	}
}

func (fake *FakeClient) KVCallCount() int {
	fake.kVMutex.RLock()
	defer fake.kVMutex.RUnlock()
	return len(fake.kVArgsForCall)
}

func (fake *FakeClient) KVReturns(result1 consuladapter.KV) {
	fake.KVStub = nil
	fake.kVReturns = struct {
		result1 consuladapter.KV
	}{result1}
}

func (fake *FakeClient) Status() consuladapter.Status {
	fake.statusMutex.Lock()
	fake.statusArgsForCall = append(fake.statusArgsForCall, struct{}{})
	fake.statusMutex.Unlock()
	if fake.StatusStub != nil {
		return fake.StatusStub()
	} else {
		return fake.statusReturns.result1
	}
}

func (fake *FakeClient) StatusCallCount() int {
	fake.statusMutex.RLock()
	defer fake.statusMutex.RUnlock()
	return len(fake.statusArgsForCall)
}

func (fake *FakeClient) StatusReturns(result1 consuladapter.Status) {
	fake.StatusStub = nil
	fake.statusReturns = struct {
		result1 consuladapter.Status
	}{result1}
}

func (fake *FakeClient) LockOpts(opts *api.LockOptions) (consuladapter.Lock, error) {
	fake.lockOptsMutex.Lock()
	fake.lockOptsArgsForCall = append(fake.lockOptsArgsForCall, struct {
		opts *api.LockOptions
	}{opts})
	fake.lockOptsMutex.Unlock()
	if fake.LockOptsStub != nil {
		return fake.LockOptsStub(opts)
	} else {
		return fake.lockOptsReturns.result1, fake.lockOptsReturns.result2
	}
}

func (fake *FakeClient) LockOptsCallCount() int {
	fake.lockOptsMutex.RLock()
	defer fake.lockOptsMutex.RUnlock()
	return len(fake.lockOptsArgsForCall)
}

func (fake *FakeClient) LockOptsArgsForCall(i int) *api.LockOptions {
	fake.lockOptsMutex.RLock()
	defer fake.lockOptsMutex.RUnlock()
	return fake.lockOptsArgsForCall[i].opts
}

func (fake *FakeClient) LockOptsReturns(result1 consuladapter.Lock, result2 error) {
	fake.LockOptsStub = nil
	fake.lockOptsReturns = struct {
		result1 consuladapter.Lock
		result2 error
	}{result1, result2}
}

var _ consuladapter.Client = new(FakeClient)
