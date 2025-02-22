// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	context "context"

	storage "github.com/nebula-contrib/nebula-sirius/nebula/storage"
	mock "github.com/stretchr/testify/mock"
)

// StorageAdminService is an autogenerated mock type for the StorageAdminService type
type StorageAdminService struct {
	mock.Mock
}

// AddAdminTask provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) AddAdminTask(ctx context.Context, req *storage.AddTaskRequest) (*storage.AddTaskResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for AddAdminTask")
	}

	var r0 *storage.AddTaskResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.AddTaskRequest) (*storage.AddTaskResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.AddTaskRequest) *storage.AddTaskResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.AddTaskResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.AddTaskRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddLearner provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) AddLearner(ctx context.Context, req *storage.AddLearnerReq) (*storage.AdminExecResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for AddLearner")
	}

	var r0 *storage.AdminExecResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.AddLearnerReq) (*storage.AdminExecResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.AddLearnerReq) *storage.AdminExecResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.AdminExecResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.AddLearnerReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddPart provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) AddPart(ctx context.Context, req *storage.AddPartReq) (*storage.AdminExecResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for AddPart")
	}

	var r0 *storage.AdminExecResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.AddPartReq) (*storage.AdminExecResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.AddPartReq) *storage.AdminExecResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.AdminExecResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.AddPartReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BlockingWrites provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) BlockingWrites(ctx context.Context, req *storage.BlockingSignRequest) (*storage.BlockingSignResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for BlockingWrites")
	}

	var r0 *storage.BlockingSignResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.BlockingSignRequest) (*storage.BlockingSignResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.BlockingSignRequest) *storage.BlockingSignResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.BlockingSignResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.BlockingSignRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckPeers provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) CheckPeers(ctx context.Context, req *storage.CheckPeersReq) (*storage.AdminExecResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for CheckPeers")
	}

	var r0 *storage.AdminExecResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.CheckPeersReq) (*storage.AdminExecResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.CheckPeersReq) *storage.AdminExecResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.AdminExecResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.CheckPeersReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClearSpace provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) ClearSpace(ctx context.Context, req *storage.ClearSpaceReq) (*storage.ClearSpaceResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for ClearSpace")
	}

	var r0 *storage.ClearSpaceResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.ClearSpaceReq) (*storage.ClearSpaceResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.ClearSpaceReq) *storage.ClearSpaceResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.ClearSpaceResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.ClearSpaceReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateCheckpoint provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) CreateCheckpoint(ctx context.Context, req *storage.CreateCPRequest) (*storage.CreateCPResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for CreateCheckpoint")
	}

	var r0 *storage.CreateCPResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.CreateCPRequest) (*storage.CreateCPResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.CreateCPRequest) *storage.CreateCPResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.CreateCPResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.CreateCPRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DropCheckpoint provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) DropCheckpoint(ctx context.Context, req *storage.DropCPRequest) (*storage.DropCPResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for DropCheckpoint")
	}

	var r0 *storage.DropCPResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.DropCPRequest) (*storage.DropCPResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.DropCPRequest) *storage.DropCPResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.DropCPResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.DropCPRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLeaderParts provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) GetLeaderParts(ctx context.Context, req *storage.GetLeaderReq) (*storage.GetLeaderPartsResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetLeaderParts")
	}

	var r0 *storage.GetLeaderPartsResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.GetLeaderReq) (*storage.GetLeaderPartsResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.GetLeaderReq) *storage.GetLeaderPartsResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.GetLeaderPartsResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.GetLeaderReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MemberChange provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) MemberChange(ctx context.Context, req *storage.MemberChangeReq) (*storage.AdminExecResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for MemberChange")
	}

	var r0 *storage.AdminExecResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.MemberChangeReq) (*storage.AdminExecResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.MemberChangeReq) *storage.AdminExecResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.AdminExecResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.MemberChangeReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemovePart provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) RemovePart(ctx context.Context, req *storage.RemovePartReq) (*storage.AdminExecResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for RemovePart")
	}

	var r0 *storage.AdminExecResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.RemovePartReq) (*storage.AdminExecResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.RemovePartReq) *storage.AdminExecResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.AdminExecResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.RemovePartReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StopAdminTask provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) StopAdminTask(ctx context.Context, req *storage.StopTaskRequest) (*storage.StopTaskResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for StopAdminTask")
	}

	var r0 *storage.StopTaskResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.StopTaskRequest) (*storage.StopTaskResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.StopTaskRequest) *storage.StopTaskResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.StopTaskResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.StopTaskRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransLeader provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) TransLeader(ctx context.Context, req *storage.TransLeaderReq) (*storage.AdminExecResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for TransLeader")
	}

	var r0 *storage.AdminExecResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.TransLeaderReq) (*storage.AdminExecResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.TransLeaderReq) *storage.AdminExecResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.AdminExecResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.TransLeaderReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WaitingForCatchUpData provides a mock function with given fields: ctx, req
func (_m *StorageAdminService) WaitingForCatchUpData(ctx context.Context, req *storage.CatchUpDataReq) (*storage.AdminExecResp, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for WaitingForCatchUpData")
	}

	var r0 *storage.AdminExecResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.CatchUpDataReq) (*storage.AdminExecResp, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.CatchUpDataReq) *storage.AdminExecResp); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.AdminExecResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.CatchUpDataReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStorageAdminService creates a new instance of StorageAdminService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorageAdminService(t interface {
	mock.TestingT
	Cleanup(func())
}) *StorageAdminService {
	mock := &StorageAdminService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
