// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"k8s.io/client-go/tools/cache"
)

type FakeController struct {
	HasSyncedStub        func() bool
	hasSyncedMutex       sync.RWMutex
	hasSyncedArgsForCall []struct {
	}
	hasSyncedReturns struct {
		result1 bool
	}
	hasSyncedReturnsOnCall map[int]struct {
		result1 bool
	}
	LastSyncResourceVersionStub        func() string
	lastSyncResourceVersionMutex       sync.RWMutex
	lastSyncResourceVersionArgsForCall []struct {
	}
	lastSyncResourceVersionReturns struct {
		result1 string
	}
	lastSyncResourceVersionReturnsOnCall map[int]struct {
		result1 string
	}
	RunStub        func(<-chan struct{})
	runMutex       sync.RWMutex
	runArgsForCall []struct {
		arg1 <-chan struct{}
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeController) HasSynced() bool {
	fake.hasSyncedMutex.Lock()
	ret, specificReturn := fake.hasSyncedReturnsOnCall[len(fake.hasSyncedArgsForCall)]
	fake.hasSyncedArgsForCall = append(fake.hasSyncedArgsForCall, struct {
	}{})
	fake.recordInvocation("HasSynced", []interface{}{})
	fake.hasSyncedMutex.Unlock()
	if fake.HasSyncedStub != nil {
		return fake.HasSyncedStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.hasSyncedReturns
	return fakeReturns.result1
}

func (fake *FakeController) HasSyncedCallCount() int {
	fake.hasSyncedMutex.RLock()
	defer fake.hasSyncedMutex.RUnlock()
	return len(fake.hasSyncedArgsForCall)
}

func (fake *FakeController) HasSyncedCalls(stub func() bool) {
	fake.hasSyncedMutex.Lock()
	defer fake.hasSyncedMutex.Unlock()
	fake.HasSyncedStub = stub
}

func (fake *FakeController) HasSyncedReturns(result1 bool) {
	fake.hasSyncedMutex.Lock()
	defer fake.hasSyncedMutex.Unlock()
	fake.HasSyncedStub = nil
	fake.hasSyncedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeController) HasSyncedReturnsOnCall(i int, result1 bool) {
	fake.hasSyncedMutex.Lock()
	defer fake.hasSyncedMutex.Unlock()
	fake.HasSyncedStub = nil
	if fake.hasSyncedReturnsOnCall == nil {
		fake.hasSyncedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.hasSyncedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeController) LastSyncResourceVersion() string {
	fake.lastSyncResourceVersionMutex.Lock()
	ret, specificReturn := fake.lastSyncResourceVersionReturnsOnCall[len(fake.lastSyncResourceVersionArgsForCall)]
	fake.lastSyncResourceVersionArgsForCall = append(fake.lastSyncResourceVersionArgsForCall, struct {
	}{})
	fake.recordInvocation("LastSyncResourceVersion", []interface{}{})
	fake.lastSyncResourceVersionMutex.Unlock()
	if fake.LastSyncResourceVersionStub != nil {
		return fake.LastSyncResourceVersionStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.lastSyncResourceVersionReturns
	return fakeReturns.result1
}

func (fake *FakeController) LastSyncResourceVersionCallCount() int {
	fake.lastSyncResourceVersionMutex.RLock()
	defer fake.lastSyncResourceVersionMutex.RUnlock()
	return len(fake.lastSyncResourceVersionArgsForCall)
}

func (fake *FakeController) LastSyncResourceVersionCalls(stub func() string) {
	fake.lastSyncResourceVersionMutex.Lock()
	defer fake.lastSyncResourceVersionMutex.Unlock()
	fake.LastSyncResourceVersionStub = stub
}

func (fake *FakeController) LastSyncResourceVersionReturns(result1 string) {
	fake.lastSyncResourceVersionMutex.Lock()
	defer fake.lastSyncResourceVersionMutex.Unlock()
	fake.LastSyncResourceVersionStub = nil
	fake.lastSyncResourceVersionReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeController) LastSyncResourceVersionReturnsOnCall(i int, result1 string) {
	fake.lastSyncResourceVersionMutex.Lock()
	defer fake.lastSyncResourceVersionMutex.Unlock()
	fake.LastSyncResourceVersionStub = nil
	if fake.lastSyncResourceVersionReturnsOnCall == nil {
		fake.lastSyncResourceVersionReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.lastSyncResourceVersionReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeController) Run(arg1 <-chan struct{}) {
	fake.runMutex.Lock()
	fake.runArgsForCall = append(fake.runArgsForCall, struct {
		arg1 <-chan struct{}
	}{arg1})
	fake.recordInvocation("Run", []interface{}{arg1})
	fake.runMutex.Unlock()
	if fake.RunStub != nil {
		fake.RunStub(arg1)
	}
}

func (fake *FakeController) RunCallCount() int {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return len(fake.runArgsForCall)
}

func (fake *FakeController) RunCalls(stub func(<-chan struct{})) {
	fake.runMutex.Lock()
	defer fake.runMutex.Unlock()
	fake.RunStub = stub
}

func (fake *FakeController) RunArgsForCall(i int) <-chan struct{} {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	argsForCall := fake.runArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeController) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.hasSyncedMutex.RLock()
	defer fake.hasSyncedMutex.RUnlock()
	fake.lastSyncResourceVersionMutex.RLock()
	defer fake.lastSyncResourceVersionMutex.RUnlock()
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeController) recordInvocation(key string, args []interface{}) {
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

var _ cache.Controller = new(FakeController)